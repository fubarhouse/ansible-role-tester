package cmd

import (
	"fmt"
	"os/exec"
	"os"

	log "github.com/Sirupsen/logrus"
	"bytes"
	"io"
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
}

// A Distribution declares the options to
// pass to Docker to run and test the container.
type Distribution struct {

	// Initialise is the initialisation command
	Initialise string

	// Name is the identifying name of the distribution
	Name       string

	// Privileged is a boolean to indicate to use privileged
	Privileged bool

	// The volume parameter when running a container.
	Volume     string

	// The fully qualified container name in the format:
	// name/image:version - ie fubarhouse/docker-ansible:bionic
	Container string
}

// CentOS 6
var CentOS6 = Distribution{
	"/sbin/init",
	"centos6",
	true,
	"",
	"geerlingguy/docker-centos6-ansible:latest",
}

// CentOS 7
var CentOS7 = Distribution{
	"/usr/lib/systemd/systemd",
	"centos7",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-centos7-ansible:latest",
}

// Ubuntu 14.04
var Ubuntu1404 = Distribution{
	"/sbin/init",
	"ubuntu1404",
	true,
	"",
	"geerlingguy/docker-ubuntu1404-ansible:latest",
}

// Ubuntu 16.04
var Ubuntu1604 = Distribution{
	"/lib/systemd/systemd",
	"ubuntu1604",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-ubuntu1604-ansible:latest",
}

// Ubuntu 18.04
var Ubuntu1804 = Distribution{
	"/lib/systemd/systemd",
	"ubuntu1804",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-ubuntu1804-ansible:latest",
}

// A slice of distributions.
var Distributions = []Distribution{
	CentOS6,
	CentOS7,
	Ubuntu1404,
	Ubuntu1604,
	Ubuntu1804,
}

// Container is an interface which allows
// a user from plugging in a Distribution
// to use these functions to run Ansible tests.
// Details on
type Container interface {
	run(config *AnsibleConfig)
	install(config *AnsibleConfig)
	kill()
	test(config *AnsibleConfig)
}

// Checks if the specified container is running.
func is_running() bool {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return false
	}

	// Users should not be able to re-run containers with the same name...
	d := exec.Cmd{}
	d.Path = docker
	d.Args = []string{
		docker,
		"ps",
		"-f",
		"status=running",
		"--format",
		"'{{.Names}}'",
	}

	var out bytes.Buffer
	multi := io.MultiWriter(&out)
	d.Stdout = multi

	if err := d.Run(); err != nil {
		log.Fatalln(err)
	}
	d.Wait()

	fmt.Printf("\n*** FULL OUTPUT *** %s\n", out.String())

	return true
}

// run will launch a new container (containerID) using
// the fields in a AnsibleConfig struct.
func (dist *Distribution) run(config *AnsibleConfig) {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return
	}

	log.Printf("Running %v\n", containerID)

	var run_options string
	if dist.Privileged {
		run_options += fmt.Sprintf("--privileged")
	}

	c := exec.Cmd{}
	c.Path = docker
	c.Args = []string{
		docker,
		"run",
		"--detach",
		fmt.Sprintf("--name=%v", containerID),
		fmt.Sprintf("--volume=%v", dist.Volume),
		fmt.Sprintf("--volume=%v:%v", config.HostPath, config.RemotePath),
		run_options,
		dist.Container,
		dist.Initialise,
	}

	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Run()
	c.Wait()
}

// Install will install the requirements if the file is configured.
func (dist *Distribution) install(config *AnsibleConfig) {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return
	}

	if config.RequirementsFile != "" {

		req := fmt.Sprintf("%v/%v", config.RemotePath, config.RequirementsFile)
		log.Printf("Installing requirements from %v\n", req)
		r := exec.Cmd{}
		r.Path = docker
		r.Args = []string{
			docker,
			"exec",
			"--tty",
			containerID,
			"ansible-galaxy",
			"install",
			fmt.Sprintf("-r %v", req),
		}
		r.Stderr = os.Stderr
		r.Stdin = os.Stdin
		r.Stdout = os.Stdout
		r.Run()
		r.Wait()
	} else {
		log.Warnln("Requirements file is not configured (empty/null), skipping...")
	}
}

// Kill will stop the container and remove it.
func (dist *Distribution) kill() {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return
	}

	log.Printf("Stopping %v\n", containerID)

	r := exec.Cmd{}
	r.Path = docker
	r.Args = []string{
		docker,
		"stop",
		containerID,
	}
	r.Stderr = os.Stderr
	r.Stdin = os.Stdin
	r.Stdout = os.Stdout
	r.Run()
	r.Wait()

}

func (dist *Distribution) test_syntax(config *AnsibleConfig) {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return
	}

	// Ansible syntax check.
	log.Infoln("Checking role syntax...")

	c := exec.Cmd{}
	c.Path = docker
	c.Args = []string{
		docker,
		"exec",
		"--tty",
		containerID,
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
		"--syntax-check",
	}
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Run()
	c.Wait()

	log.Infoln("PASS")
}
func (dist *Distribution) test_role(config *AnsibleConfig) {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return
	}

	// Test role.
	log.Infoln("Running the role...")

	r := exec.Cmd{}
	r.Path = docker
	r.Args = []string{
		docker,
		"exec",
		"--tty",
		containerID,
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
	}
	r.Stderr = os.Stderr
	r.Stdin = os.Stdin
	r.Stdout = os.Stdout
	r.Run()
	r.Wait()
}

func (dist *Distribution) test_idempotence(config *AnsibleConfig) {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return
	}

	// Test role idempotence.
	// TODO: Actually make this test for idempotence.
	log.Infoln("Testing role idempotence...")
	i := exec.Cmd{}
	i.Path = docker
	i.Args = []string{
		docker,
		"exec",
		"--tty",
		string(containerID),
		"ansible-playbook",
		fmt.Sprintf("%v/tests/%v", config.RemotePath, config.PlaybookFile),
	}
	i.Stderr = os.Stderr
	i.Stdin = os.Stdin
	i.Stdout = os.Stdout
	i.Run()
	i.Wait()

	//idempotence_log=$(mktemp)
	//docker exec --tty "$CONTAINER_ID" ansible-playbook "${ROLE_PATH_IN_CONTAINER}/tests/${PLAYBOOK}" | tee -a $idempotence_log
	//
	//tail $idempotence_log | grep -q 'changed=0.*failed=0' \
	//&& (printf "\n${GREEN}Idempotence test: PASS${NEUTRAL}\n\n" && exit 0) \
	//|| (printf "\n${GREEN}Idempotence test: ${RED}FAIL${NEUTRAL}\n\n" && exit 1)
}