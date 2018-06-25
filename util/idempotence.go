package util

import (
		"strings"

	log "github.com/Sirupsen/logrus"
	"fmt"
)

// IdempotenceTest will run an Ansible playbook once and check the
// output for any changed or failed tasks as reported by Ansible.
func (dist *Distribution) IdempotenceTest(config *AnsibleConfig) {

	// Test role idempotence.
	log.Infoln("Testing role idempotence...")

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

	out, _ := DockerExec(args, true)

	idempotence := IdempotenceResult(out)
	if idempotence {
		log.Infoln("Idempotence test: PASS")
	} else {
		log.Errorln("Idempotence test: FAIL")
	}
}

// IdempotenceResult will get the result of an idempotence test
// which is the full output of a role, and it will identify each
// of the applicable checks for idempotence. In this case, we
// simply need the values of changed and failed and some basic logic.
func IdempotenceResult(output string) bool {

	lines := strings.Split(output, "\n")

	changed := ""
	failed := ""

	for _, line := range lines {
		if strings.Contains(line, "=") {
			f := strings.Split(line, "=")
			if strings.Contains(line, "changed=") {
				changed = strings.Split(f[2], " ")[0]
			}
			if strings.Contains(line, "failed=") {
				failed = strings.Split(f[4], " ")[0]
			}
		}
	}

	if failed != "0" {
		return false
	}

	if changed != "0" {
		return false
	}

	return true
}