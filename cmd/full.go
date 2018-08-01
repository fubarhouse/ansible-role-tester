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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
	"fmt"
	"strings"
)

// fullCmd represents the full command
var fullCmd = &cobra.Command{
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
		config := util.AnsibleConfig{
			HostPath:         source,
			RemotePath:       destination,
			RequirementsFile: requirements,
			PlaybookFile:     playbook,
			Verbose:          verbose,
		}

		var dist util.Distribution

		if !custom {
			var e error
			dist, e = util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
			if e != nil {
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
			util.CustomDistributionValueSet(&dist, "Container", fmt.Sprintf("%s/%s:%s", user, image, tag))
			util.CustomDistributionValueSet(&dist, "User", user)
			util.CustomDistributionValueSet(&dist, "Distro", image)
			util.CustomFamilyValueSet(&dist.Family, "Initialise", initialise)
			util.CustomFamilyValueSet(&dist.Family, "Volume", volume)
		}

		dist.CID = containerID

		if !config.IsAnsibleRole() {
			log.Fatalf("Path %v is not recognized as an Ansible role.", config.HostPath)
		}
		fp := fmt.Sprintf(source + "/tests/" + playbook)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			log.Fatalf("Specified playbook file %v does not exist.", fp)
		}
		if requirements != "" {
			fr := fmt.Sprintf(source + "/" + requirements)
			if _, err := os.Stat(fr); os.IsNotExist(err) {
				log.Fatalf("Specified requirements file %v does not exist.", fr)
			}
		}
		if !dist.DockerCheck() {
			dist.DockerRun(&config)
		}

		dist.RoleInstall(&config)
		dist.RoleSyntaxCheck(&config)
		dist.RoleTest(&config)
		dist.IdempotenceTest(&config)
		dist.DockerKill()
	},
}

func init() {
	rootCmd.AddCommand(fullCmd)

	pwd, _ := os.Getwd()
	fullCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Name of the container")
	fullCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	fullCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	fullCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	fullCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
	fullCmd.Flags().BoolVarP(&noOutput, "no-output", "o", false, "Hide output from all Docker commands")
	fullCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")
	fullCmd.Flags().BoolVarP(&custom, "custom", "c", false, "Provide my own custom distribution.")

	fullCmd.Flags().StringVarP(&initialise, "initialise", "a", "/bin/systemd", "The initialise command for the image")
	fullCmd.Flags().StringVarP(&volume, "volume", "l", "/sys/fs/cgroup:/sys/fs/cgroup:ro", "The volume argument for the image")

	fullCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	fullCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	fullCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")
}
