package util

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"bytes"
	"io"
	"sync"

	log "github.com/Sirupsen/logrus"
	"strings"
)

// IdempotenceTestRemote will run an Ansible playbook once and check the
// output for any changed or failed tasks as reported by Ansible.
func (dist *Distribution) IdempotenceTestRemote(config *AnsibleConfig) {

	// Test role idempotence.
	if !config.Quiet {
		log.Infoln("Testing role idempotence...")
	}

	// Adjust the playbook path.
	if strings.HasPrefix(config.PlaybookFile, "/") {
		config.PlaybookFile = fmt.Sprintf("/%v", config.PlaybookFile)
	} else {
		strings.Replace(config.PlaybookFile, config.RemotePath, "./", -1)
	}

	args := []string{
		config.PlaybookFile,
		"-i",
		dist.CID + ",",
		"--connection",
		"docker",
	}

	// Add verbose if configured
	if config.Verbose {
		args = append(args, "-vvvv")
	}

	now := time.Now()
	var idempotence = false
	if !config.Quiet {
		out, _ := AnsiblePlaybook(args, true)
		idempotence = IdempotenceResult(out)
	} else {
		out, _ := AnsiblePlaybook(args, false)
		idempotence = IdempotenceResult(out)
	}

	if !config.Quiet {
		log.Infof("Idempotence was checked in %v", time.Since(now))
		if idempotence {
			log.Infoln("Idempotence test: PASS")
		} else {
			log.Errorln("Idempotence test: FAIL")
			os.Exit(1)
		}
	}
}

// RoleTestRemote will execute the specified playbook outside the
// container once. It will assemble a request to  pass into the
// Docker execution function DockerRun.
func (dist *Distribution) RoleTestRemote(config *AnsibleConfig) {

	// Test role.
	if !config.Quiet {
		log.Infoln("Running the role...")
	}

	// Adjust the playbook path.
	if strings.HasPrefix(config.PlaybookFile, "/") {
		config.PlaybookFile = fmt.Sprintf("/%v", config.PlaybookFile)
	} else {
		strings.Replace(config.PlaybookFile, config.RemotePath, "./", -1)
		//config.PlaybookFile = fmt.Sprintf("./%v", config.PlaybookFile)
	}

	args := []string{
		config.PlaybookFile,
		"-i",
		dist.CID + ",",
		"--connection",
		"docker",
	}

	// Add verbose if configured
	if config.Verbose {
		args = append(args, "-vvvv")
	}

	now := time.Now()
	if !config.Quiet {
		if _, err := AnsiblePlaybook(args, true); err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
	} else {
		if _, err := AnsiblePlaybook(args, false); err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
	}
	if !config.Quiet {
		log.Infof("Role ran in %v", time.Since(now))
	}
}


// AnsiblePlaybook will execute a command to the ansible-playbook
// binary and use the input args as arguments for that process.
// You can request output be printed using the bool stdout.
func AnsiblePlaybook(args []string, stdout bool) (string, error) {

	// Generate the command, based on input.
	cmd := exec.Cmd{}
	cmd.Path = ansibleplaybook
	cmd.Args = []string{ansibleplaybook}

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
		os.Exit(1)
		return out.String(), err
	}
	wg.Done()

	// Return out output as a string.
	return out.String(), nil
}

// RoleSyntaxCheckRemote will run a syntax check of the specified container.
// This helps with pure isolation of the syntax to separate it from other
// potential Ansible versions.
func (dist *Distribution) RoleSyntaxCheckRemote(config *AnsibleConfig) {

	// Ansible syntax check.
	if !config.Quiet {
		log.Infoln("Checking role syntax...")
	}

	args := []string{
		fmt.Sprintf("-i '%v,'", dist.CID),
		"-c docker",
		"--syntax-check",
	}

	// Add playbook path
	if strings.HasPrefix(config.PlaybookFile, "./") {
		pwd, _ := os.Getwd()
		config.PlaybookFile = strings.Replace(config.PlaybookFile, "./", "", -1)
		args = append(args, fmt.Sprintf("%v/%v", pwd, config.PlaybookFile))
	} else if strings.HasPrefix(config.PlaybookFile, "/") {
		args = append(args, fmt.Sprintf("%v", config.PlaybookFile))
	} else {
		args = append(args, fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile))
	}

	// Add verbose if configured
	if config.Verbose {
		args = append(args, "-vvvv")
	}

	if !config.Quiet {
		_, err := AnsiblePlaybook(args, true)
		if err != nil {
			log.Errorln("Syntax check: FAIL")
			os.Exit(1)
		} else {
			log.Infoln("Syntax check: PASS")
		}
	} else {
		_, err := AnsiblePlaybook(args, false)
		if err != nil {
			os.Exit(1)
		}
	}
}
