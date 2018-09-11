// Copyright © 2018 Karl Hepworth Karl.Hepworth@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

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

	// inventory is the input path to the inventory file.
	inventory string

	// source is the location of the role to test.
	source string

	// destination for the role after mounting.
	// default: /etc/ansible/roles/role_under_test
	destination string

	// initialise is the initialisation command for custom distributions
	initialise string

	// Path to the requirements file relative to source.
	requirements string

	// extraRoles is an optional argument for binding a
	// host folder with roles into the container
	extraRoles string

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

	// quiet simply implies that no output should be presented
	// through ansible, docker or reporting mechanisms.
	quiet = true

	// remote indicates the execution of ansible-playbook should
	// happen from the localhost to a remote destination - namely
	// the name parameter which should refer to the container name being tested.
	remote = false

	// reportProvided indicates a report should be provided on the
	// completion of all tasks being run.
	reportProvided = false

	// verbose is a boolean indicating all Ansible commands should
	// be dockerRun with the --verbose flag.
	verbose = false

	// volume is the initialisation command for custom distributions
	volume string

	// custom is a boolean to indicate a custom distribution should be used.
	custom = false

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
	}
}
