package util

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

// DockerExec will execute a command to the docker binary
// and use the input args as arguments for that process.
// You can request output be printed using the bool stdout.
func DockerExec(args []string, stdout bool) (string, error) {

	// Generate the command, based on input.
	cmd := exec.Cmd{}
	cmd.Path = docker
	cmd.Args = []string{docker}

	// Add our arguments to the command.
	cmd.Args = append(cmd.Args, args...)

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
	var wg sync.WaitGroup
	wg.Add(1)
	if err := cmd.Run(); err != nil {
		log.Errorln(err)
		return out.String(), err
	}
	wg.Done()

	// Return out output as a string.
	return out.String(), nil
}

// DockerCheck checks if the specified container is running.
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

// buildDockerArgs returns a list of arguments for the docker daemon. Note that the order
// matters here, and beware of trailing whitespaces.
func buildDockerArgs(dist *Distribution, config *AnsibleConfig, report *AnsibleReport) []string {
	dockerArgs := []string{
		"run",
		"--detach",
		fmt.Sprintf("--name=%v", dist.CID),
	}

	// Basic volumes, assumed default.
	report.Docker.Volumes = append(report.Docker.Volumes, dist.Family.Volume)
	report.Docker.Volumes = append(report.Docker.Volumes, fmt.Sprintf("%v:%v", config.HostPath, config.RemotePath))

	// If we're dealing with commands inside the container directly,
	// it would be practical to mount into the proper namespace
	// so we can call the roles by name in the playbook, rather
	// than the role_under_test convention.
	if !config.Remote {
		pwd, _ := os.Getwd()
		pwds := strings.Split(pwd, string(os.PathSeparator))
		report.Docker.Volumes = append(report.Docker.Volumes, fmt.Sprintf("%s:/etc/ansible/roles/%v", config.HostPath, pwds[len(pwds)-1]))
	}

	if config.ExtraRolesPath != "" {
		report.Docker.Volumes = append(report.Docker.Volumes, fmt.Sprintf("%s:%v", config.ExtraRolesPath, "/root/.ansible/roles"))
	}

	if config.LibraryPath != "" {
		report.Docker.Volumes = append(report.Docker.Volumes, fmt.Sprintf("%s:%v", config.LibraryPath, "/root/.ansible/plugins/modules"))
	}

	// Mount the volumes!
	for _, Volume := range report.Docker.Volumes {
		dockerArgs = append(dockerArgs, fmt.Sprintf("--volume=%v", Volume))
	}

	if dist.Privileged {
		dockerArgs = append(dockerArgs, fmt.Sprint("--privileged"))
	}
	dockerArgs = append(dockerArgs, []string{
		dist.Container,
		dist.Family.Initialise,
	}...)
	return dockerArgs
}

// DockerRun will launch a new container (containerID) using
// the fields in a AnsibleConfig struct.
func (dist *Distribution) DockerRun(config *AnsibleConfig, report *AnsibleReport) bool {

	if dist.CID == "" {
		dist.CID = fmt.Sprint(time.Now().Unix())
	}

	if !dist.DockerCheck() {
		if !config.Quiet {
			log.Printf("Running %v", dist.CID)
		}

		if _, err := DockerExec(buildDockerArgs(dist, config, report), true); err != nil {
			log.Errorln(err)
		}

	} else {
		if !config.Quiet {
			log.Warnf("container %v is already running, skipping the dockerRun stage", dist.CID)
		}
	}

	return dist.DockerCheck()

}

// DockerKill will stop the container and remove it.
func (dist *Distribution) DockerKill(quiet bool) bool {

	if dist.CID != "" {

		if dist.DockerCheck() {

			if !quiet {
				log.Printf("Stopping %v\n", dist.CID)
			}
			if _, err := DockerExec([]string{
				"stop",
				dist.CID,
			}, false); err != nil {
				log.Errorln(err)
			}

			if !quiet {
				log.Printf("Removing %v\n", dist.CID)
			}
			if _, err := DockerExec([]string{
				"rm",
				dist.CID,
			}, false); err != nil {
				log.Errorln(err)
			}
		} else {
			if !quiet {
				log.Errorf("container %v is not running\n", dist.CID)
			}
		}

	} else {
		if !quiet {
			log.Errorln("container name was not specified")
		}
	}

	return dist.DockerCheck()

}
