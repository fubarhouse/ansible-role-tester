package util

import (
	"os/exec"
	log "github.com/Sirupsen/logrus"
		)

var (

	// docker is simply the path to the Docker binary.
	// this will be located using exec.LookPath().
	// this must validate to a working binary named
	// docker in $PATH, otherwise the program will
	// report a fatal error.
	docker string

	// docker_found is a simple boolean which is set
	// to false by default, it's used in addition to
	// the docker string above to identify if the
	// docker binary is found to simplify flow control.
	docker_found = false

)

type AnsibleConfig struct {

	// HostPath is the path to the directory containing the role
	// on the host machine, which could be anywhere.
	HostPath string

	// RemotePath is the path to the roles folder on the container
	// which should represent the roles folder (ie /etc/ansible/roles)
	RemotePath string

	// The path to the requirements file relative to HostPath.
	// Requirements will not attempt installation if the field
	// does not have a value (when value == "")
	RequirementsFile string

	// The path to the playbook located in the tests file relative to
	// HostPath (ie HostPath/tests/playbook.yml)
	PlaybookFile string

	// verbose
	Verbose bool

}

// Container is an interface which allows
// a user from plugging in a Distribution
// to use these functions to dockerRun Ansible tests.
// Details on
type Container interface {
	dockerRun(config *AnsibleConfig)
	dockerKill()
	roleInstall(config *AnsibleConfig)
	roleTest(config *AnsibleConfig)
}

func init() {
	d, e := exec.LookPath("docker")
	if e != nil {
		log.Errorln("executable 'docker' was not found in $PATH.")
	}
	docker = d
	docker_found = true
	if !docker_found {
		log.Fatalln("you cannot use this application without having docker installed")
	}
}
