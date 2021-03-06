//
package vagrant

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/config"
)

//
var (
	Console *lumber.ConsoleLogger
	Log     *lumber.FileLogger
	logFile string
)

// create a console and default file logger
func init() {

	// create a default console logger
	Console = config.Console

	// create a default file logger
	Log = config.Log
	logFile = config.LogFile
}

// NewLogger sets the vagrant logger to the given path
func NewLogger(path string) {

	var err error

	// create a file logger (append if already exists)
	if Log, err = lumber.NewAppendLogger(path); err != nil {
		config.Fatal("[util/vagrant/log] lumber.NewAppendLogger() failed", err.Error())
	}

	logFile = path

	fmt.Printf(stylish.Bullet("Created %s", path))
}

// Debug
func Debug(msg string, debug bool) {
	if debug {
		fmt.Printf(msg)
	}
}

// Error
func Error(msg, err string) {
	fmt.Printf("%s (See %s for details)\n", msg, logFile)
	Log.Error(err)
}

// Fatal
func Fatal(msg, err string) {
	fmt.Printf("A fatal Vagrant error occurred (See %s for details). Exiting...", logFile)
	Log.Fatal(fmt.Sprintf("%s - %s", msg, err))
	Log.Close()
	os.Exit(1)
}
