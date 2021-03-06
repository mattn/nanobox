//
package vagrant

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nanobox-io/nanobox/config"
)

//
var err error

// Exists ensure vagrant is installed
func Exists() (exists bool) {
	if err := exec.Command("vagrant", "-v").Run(); err == nil {
		exists = true
	}
	return
}

// run runs a vagrant command
func run(cmd *exec.Cmd) error {

	//
	handleCMDout(cmd)

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output above
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// runInContext runs a command in the context of a Vagrantfile (from the same dir)
func runInContext(cmd *exec.Cmd) error {

	// run the command from ~/.nanobox/apps/<config.App>. Running the command from
	// the directory that contains the Vagratfile ensure that the command can
	// atleast run (especially in cases like 'create' where a VM hadn't been created
	// yet, and a UUID isn't available)
	setContext(config.AppDir)

	//
	handleCMDout(cmd)

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output above
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		return err
	}

	// switch back to project dir
	setContext(config.CWDir)

	return nil
}

// setContext changes the working directory to the designated context
func setContext(context string) {
	if err := os.Chdir(context); err != nil {
		fmt.Printf("No app found at %s, exiting...\n", config.AppDir)
		os.Exit(1)
	}
}

func customScanner(data []byte, atEOF bool) (advance int, token []byte, err error) {

	//
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}

	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// handleCMDout
func handleCMDout(cmd *exec.Cmd) {

	// start a goroutine that will act as an 'outputer' allowing us to add 'dots'
	// to the end of each line (as these lines are a reduced version of the actual
	// output there will be some delay between output)
	output := make(chan string)
	go func() {

		tick := time.Second

		// block until any one message outputs
		msg, ok := <-output

		// print initial message to 'get the ball rolling' on our 'outputer'
		fmt.Printf("   - %s", msg)

		// begin a loop to read off the channel until it's closed
		for {
			select {

			// print any messages and reset ticker
			case msg, ok = <-output:

				// once the channel closes print the final newline and close the goroutine
				if !ok {
					fmt.Println("")
					return
				}

				fmt.Printf("\n   - %s", msg)

				tick = time.Second

			// after every tick print a '.' until we get another message one the channel
			// (at which point ticker is reset and it starts all over again)
			case <-time.After(tick):
				fmt.Print(".")

				// increase the wait time by half of the total previous time
				tick += tick / 2
			}
		}
	}()

	// create a stderr pipe that will write any error messages to the log
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Fatal("[util/vagrant/vagrant] cmd.StderrPipe() failed", err.Error())
	}

	// log any command errors to the log
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			Error("A vagrant error occured", stderrScanner.Text())
		}
	}()

	// create a stdout pipe that will allow for scanning the output line-by-line;
	// if needed a stderr pipe could also be created at some point
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Fatal("[util/vagrant/vagrant] cmd.StdoutPipe() failed", err.Error())
	}

	// scan the command output intercepting only 'important' lines of vagrant output'
	// and tailoring their message so as to not flood the output.
	// styled according to: http://nanodocs.gopagoda.io/engines/style-guide
	stdoutScanner := bufio.NewScanner(stdout)
	stdoutScanner.Split(customScanner)
	go func() {
		for stdoutScanner.Scan() {

			txt := strings.TrimSpace(stdoutScanner.Text())
			app := config.Nanofile.Name

			// log all vagrant output (might as well)
			Log.Info(txt)

			// handle generic cases
			switch {

			// show the progress bar when trying to download nanobox/boot2docker
			case strings.Contains(txt, "box: Progress:"):
				subMatch := regexp.MustCompile(`box: Progress: (\d{1,3})% \(Rate: (.*), Estimated time remaining: (\d*:\d*:\d*)`).FindStringSubmatch(txt)

				// ensure we have all the submatches needed before using them
				if len(subMatch) >= 4 {
					i, err := strconv.Atoi(subMatch[1])
					if err != nil {

					}

					// show download progress: [*** progress *** 0.0%] 00:00:00 remaining
					fmt.Printf("\r   [%-41s %s%%] %s (%s remaining)", strings.Repeat("*", int(float64(i)/2.5)), subMatch[1], subMatch[2], subMatch[3])
				}
			}

			// handle specific cases
			switch txt {

			// nanobox vm has not yet been created
			case fmt.Sprintf("==> %v: VM not created. Moving on...", app):
				output <- "Nanobox not yet created, use 'nanobox dev' or 'nanobox run' to create it."

			// nanobox is already running
			case fmt.Sprintf("==> %v: VirtualBox VM is already running.", app):
				continue

			case fmt.Sprintf("==> %v: Importing base box 'nanobox/boot2docker'...", app):
				output <- "Importing nanobox base image"
			case fmt.Sprintf("==> %v: Booting VM...", app):
				output <- "Booting virtual machine"
			case fmt.Sprintf("==> %v: Configuring and enabling network interfaces...", app):
				output <- "Configuring virtual network"
			case fmt.Sprintf("==> %v: Mounting shared folders...", app):
				output <- fmt.Sprintf("Mounting source code (%s)", config.CWDir)
			case fmt.Sprintf("==> %v: Waiting for nanobox server...", app):
				output <- "Starting nanobox server"
			case fmt.Sprintf("==> %v: Attempting graceful shutdown of VM...", app):
				output <- "Shutting down virtual machine"
			case fmt.Sprintf("==> %v: Destroying VM and associated drives...", app):
				// output <- "Destroying virtual machine"
			case fmt.Sprintf("==> %v: Forcing shutdown of VM...", app):
				output <- "Shutting down virtual machine"
			case fmt.Sprintf("==> %v: Saving VM state and suspending execution...", app):
				output <- "Saving virtual machine"
			case fmt.Sprintf("==> %v: Resuming suspended VM...", app):
				// output <- "Resuming virtual machine"
			}
		}

		// close the output channel once all lines of command output have been read
		close(output)
	}()
}
