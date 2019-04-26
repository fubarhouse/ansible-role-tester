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

	"github.com/fubarhouse/ansible-role-tester/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// installCmd represents the test command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Run installation tasks for the mounted role",
	Long:  `Run installation tasks for the mounted role (--name $NAME)`,
	Run: func(cmd *cobra.Command, args []string) {
		config := util.AnsibleConfig{
			HostPath:         source,
			Inventory:        inventory,
			RemotePath:       destination,
			RequirementsFile: requirements,
			PlaybookFile:     playbook,
			Verbose:          verbose,
			Remote:           remote,
			Quiet:            quiet,
		}

		dist, e := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil && !quiet {
			log.Fatalln("Incompatible distribution was inputted.")
		}

		dist.CID = containerID

		if !config.IsAnsibleRole() && !quiet {
			log.Fatalf("Path %v is not recognized as an Ansible role.", config.HostPath)
		}
		if dist.DockerCheck() {

			util.MapInventory(dist.CID, &config)
			util.MapRequirements(&config)

			dist.RoleInstall(&config)

		} else {
			if !quiet {
				log.Warnf("Container %v is not currently running", dist.CID)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	pwd, _ := os.Getwd()
	installCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	installCmd.Flags().StringVarP(&inventory, "inventory", "e", "", "Inventory file")
	installCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	installCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Enable quiet mode")
	installCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	installCmd.MarkFlagRequired("name")
}
