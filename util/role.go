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
		DockerExec([]string{
			"exec",
			"--tty",
			dist.CID,
			"ansible-galaxy",
			"roleInstall",
			fmt.Sprintf("-r %v", req),
		}, true)
	} else {
		log.Warnln("Requirements file is not configured (empty/null), skipping...")
	}
}

func (dist *Distribution) RoleSyntaxCheck(config *AnsibleConfig) {

	// Ansible syntax check.
	log.Infoln("Checking role syntax...")
	DockerExec([]string{
		"exec",
		"--tty",
		dist.CID,
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
		"--syntax-check",
	}, true)

	log.Infoln("PASS")
}
func (dist *Distribution) RoleTest(config *AnsibleConfig) {

	// Test role.
	log.Infoln("Running the role...")
	DockerExec([]string{
		"exec",
		"--tty",
		dist.CID,
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
	}, true)
}
