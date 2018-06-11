package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"bytes"
	"errors"
	log "github.com/Sirupsen/logrus"
	"io"
	"strings"
	"time"
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

// getDistribution will get the distribution object to allow dynamic
// loading of different distributions. A suitable struct will be compiled
// from the inputs and returned with an error if the specified container
// cannot be found.
func getDistribution(container, target, init, volume string) (error, Distribution) {

	// We will search for the exact container.
	for _, dist := range Distributions {
		if dist.Container == container {
			return nil, dist
		}
	}

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		os.Exit(1)
	}

	c := exec.Cmd{}
	c.Path = docker
	c.Args = []string{
		docker,
		"images",
		//"-q",
		container,
	}

	d, _ := c.Output()
	if !strings.Contains(string(d), container) {
		log.Errorf("no valid image was found for '%v'\n", container)
		os.Exit(1)
	} else {

	}
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Run()
	c.Wait()

	return errors.New("could not find matching distribution, returned a compatible data structure"), Distribution{
		init,
		target,
		true,
		volume,
		container,
	}
}

// run will launch a new container (containerID) using
// the fields in a AnsibleConfig struct.
func (dist *Distribution) run(config *AnsibleConfig) {

	if containerID == "" {
		containerID = fmt.Sprint(time.Now().Unix())
	}

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
func kill() {

	docker, e := exec.LookPath("docker")
	if e != nil {
		log.Errorf("executable 'docker' was not found in $PATH.")
		return
	}

	if containerID != "" {

		if is_running() {

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
		} else {
			log.Errorf("container %v is not running\n")
		}

	} else {
		log.Errorln("container name was not specified")
	}

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

	out, _ := i.Output()
	fmt.Println(string(out))

	idempotence := idempotence_result(string(out))
	if idempotence {
		log.Infoln("Idempotence test: PASS")
	} else {
		log.Errorln("Idempotence test: FAIL")
	}
}

func idempotence_result(output string) bool {

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
