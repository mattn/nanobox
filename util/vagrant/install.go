//
package vagrant

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/config"
)

// Install downloads the nanobox vagrant and adds it to the list of vagrant boxes
func Install() error {

	// ensure nanobox/boot2docker has not already been installed
	if _, err := os.Stat(config.Home + "/.vagrant.d/boxes/nanobox-VAGRANTSLASH-boot2docker"); err == nil {
		// fmt.Printf(stylish.Bullet("nanobox/boot2docker already installed"))
		return nil
	}

	fmt.Printf(stylish.Bullet("Installing nanobox/boot2docker..."))
	return run(exec.Command("vagrant", "box", "add", "--force", "nanobox/boot2docker"))
}
