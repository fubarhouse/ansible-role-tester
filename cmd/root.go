package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (

	// The container ID.
	containerID string

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
	user string

	// distribution is the optional argument which specifies the
	// distribution field associated to a Distribution, which is used
	// with the user variable to identify a given Distribution.
	distro string

	// noOutput is a boolean indicating the output of docker commands
	// should remain completely hidden from Stdout.
	noOutput = false

	// verbose is a boolean indicating all Ansible commands should
	// be dockerRun with the --verbose flag.
	verbose = false

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "ansible-test",
		Short: "Run an Ansible role for testing purposes in an isolated environment.",
		Long:  ``,
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
