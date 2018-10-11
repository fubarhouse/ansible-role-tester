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
	"io/ioutil"
	"os"

	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
)

func newFullCmd() *cobra.Command {
	// Shared state between Run and PostRun
	var config util.AnsibleConfig
	var report util.AnsibleReport

	return &cobra.Command{
		Use:   "full",
		Short: "Complete end-to-end test process.",
		Long: `Runs a complete end-to-end process which performs the following:
  - creates a container
  - installs a requirements file
  - test the role syntax
  - runs the role
  - tests for idempotence
  - removes the container
You should be able to dockerRun all of this from the role folder on
the local file system. If you encounter errors, there's a lot
of flexibility in configuration, just change the defaults as
required.
`,
		Run: func(cmd *cobra.Command, args []string) {
			config = util.AnsibleConfig{
				HostPath:         source,
				Inventory:        inventory,
				RemotePath:       destination,
				ExtraRolesPath:   extraRoles,
				LibraryPath:      libraryPath,
				RequirementsFile: requirements,
				PlaybookFile:     playbook,
				Verbose:          verbose,
				Remote:           remote,
				Quiet:            quiet,
			}

			log := logrus.New()
			if quiet {
				log.Out = ioutil.Discard
			}

			var dist util.Distribution

			if !custom {
				var e error
				dist, e = util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
				if e != nil && !quiet {
					log.Fatalln("Incompatible distribution was inputted.")
				}
			} else {
				dist = *util.NewCustomDistribution()
				user := strings.Split(image, "/")[0]
				container := strings.Split(image, ":")[0]
				container = strings.Split(container, "/")[1]
				tag := strings.Split(image, ":")[1]

				dist.Privileged = true
				util.CustomDistributionValueSet(&dist, "Name", containerID)
				//util.CustomValueSet(&dist, "Privileged", "true")
				util.CustomDistributionValueSet(&dist, "Container", fmt.Sprintf("%s/%s:%s", user, container, tag))
				util.CustomDistributionValueSet(&dist, "User", user)
				util.CustomDistributionValueSet(&dist, "Distro", image)
				util.CustomFamilyValueSet(&dist.Family, "Initialise", initialise)
				util.CustomFamilyValueSet(&dist.Family, "Volume", volume)
			}

			dist.CID = containerID

			if !config.IsAnsibleRole() {
				if !quiet {
					log.Fatalf("Path %v is not recognized as an Ansible role.", config.HostPath)
				}
				os.Exit(util.NotARoleCode)
			}

			util.MapInventory(dist.CID, &config)
			util.MapRequirements(&config)
			util.MapPlaybook(&config)

			report = util.NewReport(&config)
			report.Meta.ReportFile = reportFilename
			report.Ansible.Distribution = dist

			if !dist.DockerCheck() {
				dist.DockerRun(&config, &report)
				report.Docker.Run = dist.DockerCheck()
			}
			hosts, _ := dist.AnsibleHosts(&config, &report)
			report.Ansible.Hosts = hosts
			if remote {
				for _, host := range hosts {
					if host == "localhost" {
						log.Errorln("remote runs should be run directly, not through this tool")
						dist.DockerKill(quiet)
					}
				}
			}

			report.Ansible.Requirements, _ = dist.RoleInstall(&config)
			if !remote {
				report.Ansible.Syntax, _ = dist.RoleSyntaxCheck(&config)
				report.Ansible.Run.Result, report.Ansible.Run.Time, _ = dist.RoleTest(&config)
				report.Ansible.Idempotence.Result, report.Ansible.Idempotence.Time = dist.IdempotenceTest(&config)
			} else {
				report.Ansible.Syntax, _ = dist.RoleSyntaxCheckRemote(&config)
				report.Ansible.Run.Result, report.Ansible.Run.Time, _ = dist.RoleTestRemote(&config)
				report.Ansible.Idempotence.Result, report.Ansible.Idempotence.Time = dist.IdempotenceTestRemote(&config)
			}

			dist.DockerKill(quiet)
			if !dist.DockerCheck() {
				report.Docker.Kill = true
			}

			if reportProvided {
				report.Ansible.Config = config
				report.Printf()
			}
		},
		// Analyze report and return the proper exit code.
		PostRun: func(cmd *cobra.Command, args []string) {
			if !report.Docker.Run {
				os.Exit(util.DockerRunCode)
			} else if !report.Ansible.Syntax {
				os.Exit(util.AnsibleSyntaxCode)
			} else if !report.Ansible.Run.Result {
				os.Exit(util.AnsibleRunCode)
			} else if !report.Ansible.Idempotence.Result {
				os.Exit(util.AnsibleIdempotenceCode)
			} else {
				os.Exit(util.OKCode)
			}
		},
	}
}

func addFullFlags(fullCmd *cobra.Command, dir string) {
	fullCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Name of the container")
	fullCmd.Flags().StringVarP(&source, "source", "s", dir, "Location of the role to test")
	fullCmd.Flags().StringVarP(&destination, "destination", "d", "", "Location which the role will be mounted to")
	fullCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	fullCmd.Flags().StringVarP(&extraRoles, "extra-roles", "x", "", "Path to roles folder with dependencies.")
	fullCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
	fullCmd.Flags().BoolVarP(&noOutput, "no-output", "o", false, "Hide output from all Docker commands")
	fullCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Enable quiet mode")
	fullCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")
	fullCmd.Flags().BoolVarP(&custom, "custom", "c", false, "Provide my own custom distribution.")
	fullCmd.Flags().StringVarP(&inventory, "inventory", "e", "", "Inventory file")
	fullCmd.Flags().BoolVarP(&remote, "remote", "m", false, "Run the test remotely to the container")
	fullCmd.Flags().BoolVarP(&reportProvided, "report", "f", false, "Provide a report after completion")
	fullCmd.Flags().StringVarP(&reportFilename, "report-output", "b", "report.yml", "Filename in current working directory to write a report to")
	fullCmd.Flags().StringVarP(&libraryPath, "library", "", "", "Path to library folder with modules.")

	fullCmd.Flags().StringVarP(&initialise, "initialise", "a", "/bin/systemd", "The initialise command for the image")
	fullCmd.Flags().StringVarP(&volume, "volume", "l", "/sys/fs/cgroup:/sys/fs/cgroup:ro", "The volume argument for the image")

	fullCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	fullCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	fullCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")
}

func init() {
	fullCmd := newFullCmd()
	pwd, _ := os.Getwd()
	addFullFlags(fullCmd, pwd)
	rootCmd.AddCommand(fullCmd)
}

func InitFullCmdForTest(dir string) *cobra.Command {
	fullCmd := newFullCmd()
	addFullFlags(fullCmd, dir)
	return fullCmd
}
