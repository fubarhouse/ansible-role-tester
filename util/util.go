package util

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

var (

	// ansibleplaybook is the path to ansible-playbook.
	// this will only be used on targeted remotes which
	// replaces the inventory path when calling
	// ansible-playbook from the host machine.
	ansibleplaybook string

	// docker is simply the path to the Docker binary.
	// this will be located using exec.LookPath().
	// this must validate to a working binary named
	// docker in $PATH, otherwise the program will
	// report a fatal error.
	docker string

	// dockerFound is a simple boolean which is set
	// to false by default, it's used in addition to
	// the docker string above to identify if the
	// docker binary is found to simplify flow control.
	dockerFound = false
)

// AnsibleConfig represents a series of configuration options
// for an ansible command to be executed.
type AnsibleConfig struct {

	// HostPath is the path to the directory containing the role
	// on the host machine, which could be anywhere.
	HostPath string

	// Inventory is the inventory argument in Ansible commands.
	// in this case we are focusing on inventory files or inline dicts.
	// Example: 'container_name,' or './tests/inventory.
	Inventory string

	// RemotePath is the path to the roles folder on the container
	// which should represent the roles folder (ie /etc/ansible/roles)
	RemotePath string

	// ExtraRolesPath is the path to the roles folder on the host which will
	// be mounted on the container to "/root/.ansible/roles" and available to the playbook
	// as dependencies. This is a useful workaround for CI/CD environments where the roles
	// are already downloaded on the host, or if the roles are in private git repos.
	ExtraRolesPath string

	// The path to the requirements file relative to HostPath.
	// Requirements will not attempt installation if the field
	// does not have a value (when value == "")
	RequirementsFile string

	// PlaybookFile is the path to the playbook located in the
	// tests file relative to HostPath (ie HostPath/tests/playbook.yml)
	PlaybookFile string

	// Remote indicates the playbook will be run on a remote host
	// likely which is inputted to the inventory field.
	Remote bool

	// verbose
	Verbose bool

	// Quiet will determine if all reporting mechanisms are hidden.
	Quiet bool
}

// Container is an interface which allows
// a user from plugging in a Distribution
// to use these functions to dockerRun Ansible tests.
// Details on
type Container interface {
	DockerRun(config *AnsibleConfig)
	DockerKill()
	RoleInstall(config *AnsibleConfig)
	RoleTest(config *AnsibleConfig)
}

func init() {
	d, e := exec.LookPath("docker")
	if e != nil {
		log.Errorln("executable 'docker' was not found in $PATH.")
		os.Exit(1)
	}
	docker = d
	dockerFound = true

	a, e := exec.LookPath("ansible-playbook")
	if e != nil {
		log.Errorln("executable 'ansible-playbook' was not found in $PATH.")
	}
	ansibleplaybook = a

}
