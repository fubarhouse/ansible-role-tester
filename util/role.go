package util

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
func (dist *Distribution) RoleInstall(config *AnsibleConfig) bool {

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
				log.Errorln(err)
				return false
			}
		} else {
			_, err := DockerExec(args, false)
			if err != nil {
				log.Errorln(err)
				return false
			}
		}

	} else {
		if !config.Quiet {
			log.Warnln("Requirements file is not configured (empty/null), skipping...")
			return false
		}
	}
	return true
}

// RoleSyntaxCheck will run a syntax check of the mounted volume inside
// of the active container. This helps with pure isolation of the syntax
// to separate it from other potential Ansible versions.
func (dist *Distribution) RoleSyntaxCheck(config *AnsibleConfig) bool {

	// Ansible syntax check.
	if !config.Quiet {
		log.Infoln("Checking role syntax...")
	}

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
			log.Errorln("Syntax check: FAIL")
			return false
		} else {
			log.Infoln("Syntax check: PASS")
			return true
		}
	} else {
		_, err := DockerExec(args, false)
		if err != nil {
			log.Errorln(err)
			return false
		}
	}
	return true
}

// RoleTest will execute the specified playbook inside
// the container once. It will assemble a request to
// pass into the Docker execution function DockerRun.
func (dist *Distribution) RoleTest(config *AnsibleConfig) (bool, time.Duration) {

	// Test role.
	if !config.Quiet {
		log.Infoln("Running the role...")
	}

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
			log.Errorln(err)
			return false, time.Since(now)
		}
	} else {
		if _, err := DockerExec(args, false); err != nil {
			log.Errorln(err)
			return false, time.Since(now)
		}
	}
	if !config.Quiet {
		log.Infof("Role ran in %v", time.Since(now))
	}
	return true, time.Since(now)
}
