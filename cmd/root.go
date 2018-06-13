package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
	"time"
)

var (

	// The container ID.
	containerID = fmt.Sprint(time.Now().Unix())

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

	// image is the Distribution identifier by name.
	// These are declared as a Distribution and this will
	// attempt to find a matching Distribution for use.
	// example: fubarhouse/docker-ansible:bionic
	image string

	// source is the location of the role to test.
	source string

	// destination for the role after mounting.
	// default: /etc/ansible/roles/role_under_test
	destination string

	// Path to the requirements file relative to source.
	requirements string

	// playbook is the path to the playbook to execute inside of
	// the 'tests' folder.
	playbook string

	// user is the optional argument which specifies the
	// user associated to the selected distribution, which
	// will be used to locate a Distribution with the same user.
	user   string

	// distribution is the optional argument which specifies the
	// distribution field associated to a Distribution, which is used
	// with the user variable to identify a given Distribution.
	distro string

	// noOutput is a boolean indicating the output of docker commands
	// should remain completely hidden from Stdout.
	noOutput = false

	// verbose is a boolean indicating all Ansible commands should
	// be run with the --verbose flag.
	verbose = false

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "ansible-test",
		Short: "Run an Ansible role for testing purposes in an isolated environment.",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
