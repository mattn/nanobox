package commands

import (
	"fmt"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// ServiceListCommand satisfies the Command interface for listing an app's services
type ServiceListCommand struct{}

// Help prints detailed help text for the service list command
func (c *ServiceListCommand) Help() {
	ui.CPrintln(`
Description:
  Lists an app's services

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

Usage:
  pagoda service:list [-a app-name]

  ex. pagoda service:list -a app-name

Options:
  -a, --app [app-name]
    The name of the app
  `)
}

// Run disaplys select information about all of an app's services
func (c *ServiceListCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	// get services
	services, err := api.GetAppServices(fApp)
	if err != nil {
		fmt.Println("There was a problem getting '%v's' services. See ~/.pagodabox/log.txt for details", fApp)
		ui.Error("pagoda service:list", err)
	}

	//
	fmt.Println(`
state name - type (uid)
-------------------------
  `)

	var serviceColor string

	for _, service := range services {

		serviceColor = helpers.DetermineServiceStatus(service.State)

		switch service.State {

		//
		case "initialized", "active":
			ui.CPrint(serviceColor + "\u2022[reset] " + service.Name + " - (" + service.UID + ")")

		//
		case "inactive":
			ui.CPrint(serviceColor + "x[reset] " + service.Name + " - (" + service.UID + ")")

		//
		case "defunct":
			ui.CPrint(serviceColor + "![reset] " + service.Name + " - (" + service.UID + ")")
		}

		fmt.Println("")
	}

	fmt.Println("")
}