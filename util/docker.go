package util

import (
	"os/exec"
	"strings"
	"os"
	"bytes"
	"io"

	log "github.com/Sirupsen/logrus"
	"fmt"
	"time"
)

// docer_exec will execute a command to the docker binary
// and use the input args as arguments for that process.
// You can request output be printed using the bool stdout.
func DockerExec(args []string, stdout bool) (string, error) {

	// Generate the command, based on input.
	cmd := exec.Cmd{}
	cmd.Path = docker
	cmd.Args = []string{docker}

	// Add our arguments to the command.
	for _, arg := range args {
		cmd.Args = append(cmd.Args, arg)
	}

	// If configured, print to os.Stdout.
	if stdout {
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
	}

	// Create a buffer for the output.
	var out bytes.Buffer
	multi := io.MultiWriter(&out)

	//if stdout && !noOutput {
	if stdout {
		multi = io.MultiWriter(&out, os.Stdout)
	}

	// Assign the output to the writer.
	cmd.Stdout = multi

	// Check the errors, return as needed.
	if err := cmd.Run(); err != nil {
		log.Errorln(err)
		return out.String(), err
	}
	cmd.Wait()

	// Return out output as a string.
	return out.String(), nil
}

// Checks if the specified container is running.
func (dist *Distribution) DockerCheck() bool {
	// Users should not be able to re-dockerRun containers with the same name...
	if dist.CID != "" {
		out, err := DockerExec([]string{
			"ps",
			"-a",
			"-f",
			"status=running",
			"--format",
			"'{{.Names}}'",
		}, false)

		if err != nil {
			return false
		}

		if strings.Contains(out, dist.CID) {
			return true
		}
	}

	return false
}

// dockerRun will launch a new container (containerID) using
// the fields in a AnsibleConfig struct.
func (dist *Distribution) DockerRun(config *AnsibleConfig) {

	if dist.CID == "" {
		dist.CID = fmt.Sprint(time.Now().Unix())
	}

	if !dist.DockerCheck() {
		log.Printf("Running %v", dist.CID)

		var run_options string
		if dist.Privileged {
			run_options += fmt.Sprint("--privileged")
		}

		DockerExec([]string{
			"run",
			"--detach",
			fmt.Sprintf("--name=%v", dist.CID),
			fmt.Sprintf("--volume=%v", dist.Family.Volume),
			fmt.Sprintf("--volume=%v:%v", config.HostPath, config.RemotePath),
			run_options,
			dist.Container,
			dist.Family.Initialise,
		}, true)
	} else {
		log.Warnf("container %v is already running, skipping the dockerRun stage", dist.CID)
	}
}

// Kill will stop the container and remove it.
func (dist *Distribution) DockerKill() {

	if dist.CID != "" {

		if dist.DockerCheck() {

			log.Printf("Stopping %v\n", dist.CID)
			DockerExec([]string{
				"stop",
				dist.CID,
			}, false)

			log.Printf("Removing %v\n", dist.CID)
			DockerExec([]string{
				"rm",
				dist.CID,
			}, false)
		} else {
			log.Errorf("container %v is not running\n", dist.CID)
		}

	} else {
		log.Errorln("container name was not specified")
	}

}
