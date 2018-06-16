package util

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// Install will roleInstall the requirements if the file is configured.
func (dist *Distribution) RoleInstall(config *AnsibleConfig) {

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

		// Add verbose if configured
		if config.Verbose {
			args = append(args, "-vvvv")
		}

		_, err := DockerExec(args, true)
		if err != nil {
			log.Errorln(err)
		}

	} else {
		log.Warnln("Requirements file is not configured (empty/null), skipping...")
	}
}

func (dist *Distribution) RoleSyntaxCheck(config *AnsibleConfig) {

	// Ansible syntax check.
	log.Infoln("Checking role syntax...")

	args := []string{
		"exec",
		"--tty",
		dist.CID,
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
		"--syntax-check",
	}

	// Add verbose if configured
	if config.Verbose {
		args = append(args, "-vvvv")
	}

	_, err := DockerExec(args, true)
	if err != nil {
		log.Errorln("Syntax check: FAIL")
	} else {
		log.Infoln("Syntax check: PASS")
	}
}
func (dist *Distribution) RoleTest(config *AnsibleConfig) {

	// Test role.
	log.Infoln("Running the role...")

	args := []string{
		"exec",
		"--tty",
		dist.CID,
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
	}

	// Add verbose if configured
	if config.Verbose {
		args = append(args, "-vvvv")
	}

	DockerExec(args, true)
}
