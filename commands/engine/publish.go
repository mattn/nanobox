//
package engine

import (
	"archive/tar"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	api "github.com/nanobox-io/nanobox-api-client"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	engineutil "github.com/nanobox-io/nanobox/util/engine"
	fileutil "github.com/nanobox-io/nanobox/util/file"
	s3util "github.com/nanobox-io/nanobox/util/s3"
)

var tw *tar.Writer

//
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes an engine to nanobox.io",
	Long:  ``,

	Run: publish,
}

// publish
func publish(ccmd *cobra.Command, args []string) {
	stylish.Header("publishing engine")

	//
	api.UserSlug, api.AuthToken = Auth.Authenticate()

	// create a new release
	fmt.Printf(stylish.Bullet("Creating release..."))
	release := &api.EngineRelease{}

	// create an annonymous struct to hold data that doesn't relate to a release but
	// is needed as part of the publish process
	opts := &struct {
		Generic  bool     `json:"generic"`
		Language string   `json:"language"`
		Overlays []string `json:"overlays"`
	}{}

	// ensure there is an Enginefile
	if _, err := os.Stat("./Enginefile"); err != nil {
		fmt.Println("Enginefile not found. Be sure to publish from a project directory. Exiting... ")
		os.Exit(1)
	}

	// parse the ./Enginefile into the new release
	if err := Config.ParseConfig("./Enginefile", release); err != nil {
		fmt.Printf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.\n")
		os.Exit(1)
	}

	// parse the ./Enginefile again to get the remaining fields
	if err := Config.ParseConfig("./Enginefile", opts); err != nil {
		fmt.Printf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.\n")
		os.Exit(1)
	}

	fmt.Printf(stylish.Bullet("Verifying engine is publishable..."))

	// determine if any required fields (name, version, language, summary) are missing,
	// if any are found to be missing exit 1
	// NOTE: I do this using fallthrough for asthetics onlye. The message is generic
	// enough that all cases will return the same message, and this looks better than
	// a single giant case (var == "" || var == "" || ...)
	switch {
	case opts.Language == "":
		fallthrough
	case release.Name == "":
		fallthrough
	case release.Summary == "":
		fallthrough
	case release.Version == "":
		fmt.Printf(stylish.Error("required fields missing", `Your Enginefile is missing one or more of the following required fields for publishing:

  name:      # the name of your project
  version:   # the current version of the project
  language:  # the lanauge (ruby, golang, etc.) of the engine
  summary:   # a 140 character summary of the project

Please ensure all required fields are provided and try again.`))

		os.Exit(1)
	}

	// attempt to read a README.md file and add it to the release...
	b, err := ioutil.ReadFile("./README.md")
	if err != nil {

		// this only fails if the file is not found, EOF is not an error. If no Readme
		// is found exit 1
		fmt.Printf(stylish.Error("missing readme", "Your engine is missing a README.md file. This file is required for publishing, as it is the only way for you to communicate how to use your engine. Please add a README.md and try again."))
		os.Exit(1)
	}

	//
	release.Readme = string(b)

	// check to see if the engine already exists on nanobox.io
	fmt.Printf(stylish.Bullet("Checking for existing engine on nanobox.io"))
	engine, err := api.GetEngine(api.UserSlug, release.Name)

	// if no engine exists, create a new one
	if err != nil {

		// if no engine is found create one
		if apiErr, _ := err.(api.APIError); apiErr.Code == 404 {

			fmt.Printf(stylish.SubTaskStart("Creating new engine on nanobox.io"))

			//
			engine = &api.Engine{
				Generic:      opts.Generic,
				LanguageName: opts.Language,
				Name:         release.Name,
			}

			//
			if _, err := api.CreateEngine(engine); err != nil {
				fmt.Printf(stylish.ErrBullet("Unable to create engine (%v).", err))
				os.Exit(1)
			}

			// wait until engine has been successfuly created before uploading to s3
			for {
				fmt.Print(".")

				p, err := api.GetEngine(api.UserSlug, release.Name)
				if err != nil {
					Config.Fatal("[commands/publish] api.GetEngine failed", err.Error())
				}

				// once the engine is "active", break
				if p.State == "active" {
					break
				}

				//
				time.Sleep(1000 * time.Millisecond)
			}

			// generically handle any other errors
		} else {
			Config.Fatal("[commands/publish] api.GetEngine failed", err.Error())
		}

		stylish.Success()
	}

	// create a meta.json file where we can add any extra data we might need; since
	// this is only used for internal purposes the file is removed once we're done
	// with it
	meta, err := os.Create("./meta.json")
	if err != nil {
		Config.Fatal("[commands/publish] os.Create() failed", err.Error())
	}
	defer meta.Close()
	defer os.Remove(meta.Name())

	// add any custom info to the metafile
	meta.WriteString(fmt.Sprintf(`{"engine_id": "%s"}`, engine.ID))

	// this is our predefined list of everything that gets archived as part of the
	// engine being published
	files := map[string][]string{
		"required": []string{"./bin", "./Enginefile", "./meta.json"},
		"optional": []string{"./lib", "./templates", "./files"},
	}

	// check to ensure no required files are missing
	for k, v := range files {
		if k == "required" {
			for _, f := range v {
				if _, err := os.Stat(f); err != nil {
					fmt.Printf(stylish.Error("required files missing", "Your Engine is missing one or more required files for publishing. Please read the following documentation to ensure all required files are included and try again.:\n\ndocs.nanobox.io/engines/project-creation/#example-engine-file-structure\n"))
					os.Exit(1)
				}
			}
		}
	}

	// create the temp engines folder for building the tarball
	tarPath := filepath.Join(config.EnginesDir, release.Name)
	if err := os.Mkdir(tarPath, 0755); err != nil {
		Config.Fatal("[commands/engine/publish] os.Create() failed", err.Error())
	}

	// remove tarDir once published
	defer func() {
		if err := os.RemoveAll(tarPath); err != nil {
			os.Stderr.WriteString(stylish.ErrBullet("Faild to remove '%v'...", tarPath))
		}
	}()

	// parse the ./Enginefile again to get the overlays
	if err := Config.ParseConfig("./Enginefile", opts); err != nil {
		fmt.Printf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.\n")
		os.Exit(1)
	}

	// iterate through each overlay fetching it and untaring to the tar path
	for _, overlay := range opts.Overlays {
		engineutil.GetOverlay(overlay, tarPath)
	}

	// range over each file from each file type, building the final list of files
	// to be tarballed
	for _, v := range files {
		for _, f := range v {

			// not handling error here because an error simply means the file doesn't
			// exist and therefor wont be copied to the final tarball
			fileutil.Copy(f, tarPath)
		}
	}

	// create an empty buffer for writing the file contents to for the subsequent
	// upload
	archive := bytes.NewBuffer(nil)

	//
	h := md5.New()

	//
	if err := fileutil.Tar(tarPath, archive, h); err != nil {
		Config.Fatal("[commands/engine/publish] file.Tar() failed", err.Error())
	}

	// add the checksum for the new release once its finished being archived
	release.Checksum = fmt.Sprintf("%x", h.Sum(nil))

	//
	// attempt to upload the release to S3
	fmt.Printf(stylish.Bullet("Uploading release to s3..."))

	v := url.Values{}
	v.Add("user_slug", api.UserSlug)
	v.Add("auth_token", api.AuthToken)
	v.Add("version", release.Version)

	//
	s3url, err := s3util.RequestURL(fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/request_upload?%v", release.Name, v.Encode()))
	if err != nil {
		Config.Fatal("[commands/publish] s3.RequestURL() failed", err.Error())
	}

	//
	if err := s3util.Upload(s3url, archive); err != nil {
		Config.Fatal("[commands/publish] s3.Upload() failed", err.Error())
	}

	//
	// if the release uploaded successfully to s3, created one on odin
	fmt.Printf(stylish.Bullet("Uploading release to nanobox.io"))
	if _, err := api.CreateEngineRelease(release.Name, release); err != nil {
		fmt.Printf(stylish.ErrBullet("Unable to publish release (%v).", err))
		os.Exit(1)
	}
}
