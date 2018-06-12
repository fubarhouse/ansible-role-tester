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

	docker string

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

	// Path t
	playbook string

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
		log.Errorf("executable 'docker' was not found in $PATH.\n")
	}
	docker = d
	docker_found = true
	if !docker_found {
		log.Fatalln("you cannot use this application without having docker installed")
	}
}
