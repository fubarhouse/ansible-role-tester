package util

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"os"
	"time"
)

// IsAnsibleRole will identify if the mounted directory is an Ansible role.
func (config *AnsibleConfig) IsAnsibleRole() bool {
	if _, err := os.Stat(config.HostPath + string(os.PathSeparator) + "tasks" + string(os.PathSeparator) + "main.yml"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(config.HostPath + string(os.PathSeparator) + "meta" + string(os.PathSeparator) + "main.yml"); os.IsNotExist(err) {
		return false
	}
	return true
}

// RoleInstall will install the requirements if the file is configured.
func (dist *Distribution) RoleInstall(config *AnsibleConfig) (bool, error) {

	if config.RequirementsFile != "" {
		req := fmt.Sprintf("%v/%v", config.RemotePath, config.RequirementsFile)
		log.Printf("Installing requirements from %v\n", req)
		args := []string{
			"exec",
			"--tty",
			dist.CID,
			"ansible-galaxy",
			"install",
			"-r",
			req,
		}

		// Add inventory file if configured
		if config.Inventory != "" {
			args = append(args, fmt.Sprintf("-i=%v", config.Inventory))
		}

		// Add verbose if configured
		if config.Verbose {
			args = append(args, "-vvvv")
		}

		if !config.Quiet {
			_, err := DockerExec(args, true)
			if err != nil {
				return false, err
			}
		} else {
			_, err := DockerExec(args, false)
			if err != nil {
				return false, err
			}
		}

	} else {
		if !config.Quiet {
			log.Warnln("Requirements file is not configured (empty/null), skipping...")
			return false, nil
		}
	}
	return true, nil
}

// RoleSyntaxCheck will run a syntax check of the mounted volume inside
// of the active container. This helps with pure isolation of the syntax
// to separate it from other potential Ansible versions.
func (dist *Distribution) RoleSyntaxCheck(config *AnsibleConfig) (bool, error) {

	// Ansible syntax check.

	args := []string{
		"exec",
		"--tty",
		dist.CID,
		"ansible-playbook",
		"--syntax-check",
		fmt.Sprintf("%v/%v", config.RemotePath, config.PlaybookFile),
	}

	// Add inventory file if configured
	if config.Inventory != "" {
		args = append(args, fmt.Sprintf("-i=%v", config.Inventory))
	}

	// Add verbose if configured
	if config.Verbose {
		args = append(args, "-vvvv")
	}

	if !config.Quiet {
		_, err := DockerExec(args, true)
		if err != nil {
			return false, errors.New("syntax check failed")
		} else {
			return true, nil
		}
	} else {
		_, err := DockerExec(args, false)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// RoleTest will execute the specified playbook inside
// the container once. It will assemble a request to
// pass into the Docker execution function DockerRun.
func (dist *Distribution) RoleTest(config *AnsibleConfig) (bool, time.Duration, error) {

	// Test role.

	args := []string{
		"exec",
		"--tty",
		dist.CID,
		"ansible-playbook",
		fmt.Sprintf("%v/%v", config.RemotePath, config.PlaybookFile),
	}

	// Add inventory file if configured
	if config.Inventory != "" {
		args = append(args, fmt.Sprintf("-i=%v", config.Inventory))
	}

	// Add verbose if configured
	if config.Verbose {
		args = append(args, "-vvvv")
	}

	now := time.Now()
	if !config.Quiet {
		if _, err := DockerExec(args, true); err != nil {
			return false, time.Since(now), err
		}
	} else {
		if _, err := DockerExec(args, false); err != nil {
			return false, time.Since(now), err
		}
	}
	if !config.Quiet {
		log.Infof("Role ran in %v", time.Since(now))
	}
	return true, time.Since(now), nil
}
