//
package engine

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

//
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generates a new engine inside the current working directory",
	Long:  ``,

	Run: new,
}

// new
func new(ccmd *cobra.Command, args []string) {

	// this will be removed once the new command is more fleshed out
	fmt.Println(`
This area of the CLI is under construction. For details on how to create an engine
please see the documentation provided here: https://docs.nanobox.io/engines/`)
	os.Exit(0)

	if len(args) < 1 {
		fmt.Println(stylish.Error("app name required", "Please provide a name when generating a new engine, (run 'nanobox new -h' for details)"))
		os.Exit(1)
	}

	name := fmt.Sprintf("nanobox-%s", args[0])
	version := "0.0.1"

	fmt.Println(stylish.Header("initializing new app: %s", name))

	// create a new project by the name, unless it already exists
	if _, err := os.Stat(name); err != nil {

		stylish.Bullet("Creating engine at")

		if err := os.Mkdir(name, 0755); err != nil {
			Config.Fatal("[commands/new] os.Mkdir() failed", err.Error())
		}

		entry := fmt.Sprintf(`
name: %-18s     # the name of your project (required)
version: %-18s  # the current version of the project (required)
language:                    # the lanauge (ruby, golang, etc.) of the engine (required)
summary:                     # a 140 character short summary of the project (required)
stability: alpha             # the current stability of the project (alpha/beta/stable)
generic: false               # is the engine generic enough to encompass multiple frameworks
                             # within the given language
license: MIT                 # the license to be applied to this software

# engine overlays - any overlays that your engine requires to run
overlays:
 -

# a list of authors/contributors
authors:
  -
`, name, version)

		enginefilePath := name + "/Enginefile"

		f, err := os.Create(enginefilePath)
		if err != nil {
			Config.Fatal("[commands/new] os.Create() failed", err.Error())
		}
		defer f.Close()

		// write the Enginefile
		if err := ioutil.WriteFile(enginefilePath, []byte(entry), 0644); err != nil {
			Config.Fatal("[commands/new] ioutil.WriteFile() failed", err.Error())
		}

	} else {
		fmt.Printf("A project by the name '%s' already exists at this location...\n", name)
	}

	stylish.Bullet("Default Enginefile created at")
}
