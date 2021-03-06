//
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/server"
)

//
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about the nanobox and your app",
	Long:  ``,

	Run: info,
}

// info runs 'vagrant status'
func info(ccmd *cobra.Command, args []string) {

	// determine status
	status := Vagrant.Status()

	fmt.Printf(`
Local Domain   : %s
App Name       : %s
Nanobox state  : %s
Nanobox Files  : %s

`, config.Nanofile.Domain, config.Nanofile.Name, status, config.AppDir)

	if status != "running" {
		fmt.Println("NOTE: To view 'services' information start nanobox with 'nanobox dev' or 'nanobox run'")
		return
	}

	//
	var services []server.Service
	res, err := server.Get("/services", &services)
	if err != nil {
		server.Fatal("[commands/info] failed", err.Error())
	}
	defer res.Body.Close()

	//
	if len(services) >= 1 {
		fmt.Printf("////////// SERVICES //////////\n")

		//
		for _, service := range services {

			fmt.Printf("\n%s :\n", service.UID)

			if service.Name != "" {
				fmt.Printf("   name  : %s\n", service.Name)
			}

			fmt.Printf("   host  : %s\n", config.Nanofile.Domain)
			fmt.Printf("   ports : %v\n", service.Ports)

			//
			if service.Username != "" {
				fmt.Printf("   username : %s\n", service.Username)
			}

			//
			if service.Password != "" {
				fmt.Printf("   password : %s\n", service.Password)
			}

			// show any environment variables
			if len(service.EnvVars) >= 1 {
				fmt.Printf("   evars :\n")

				for k, v := range service.EnvVars {
					fmt.Printf("      %s : %s\n", k, v)
				}
			}
		}
	}
}
