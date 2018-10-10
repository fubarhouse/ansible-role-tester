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
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Execute tests against an existing container",
	Long: `Execute tests against an existing container

If container does not exist it will be created, however
containers won't be removed after completion.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if quiet {
			logrus := log.New()
			logrus.Out = ioutil.Discard
		}
	},
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

		dist, _ := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		report := util.NewReport(&config)

		dist.CID = containerID

		if dist.DockerCheck() {

			util.MapPlaybook(&config)
			util.MapInventory(dist.CID, &config)
			util.MapRequirements(&config)

			if !remote {
				dist.RoleSyntaxCheck(&config)
				dist.RoleTest(&config)
				dist.IdempotenceTest(&config)
			} else {
				dist.RoleSyntaxCheckRemote(&config)
				dist.RoleTestRemote(&config)
				dist.IdempotenceTestRemote(&config)
				hosts, _ := dist.AnsibleHosts(&config, &report)
				for _, host := range hosts {
					if host == "localhost" {
						log.Errorln("remote runs should be run directly, not through this tool")
						os.Exit(1)
					}
				}
			}
		} else {
			if !quiet {
				log.Warnf("Container %v is not currently running", dist.CID)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	pwd, _ := os.Getwd()
	testCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	testCmd.Flags().StringVarP(&inventory, "inventory", "e", "", "Inventory file")
	testCmd.Flags().StringVarP(&destination, "destination", "d", "", "Location which the role was mounted to")
	testCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")
	testCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
	testCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Enable quiet mode")
	testCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	testCmd.Flags().BoolVarP(&remote, "remote", "m", false, "Run the test remotely to the container")

	testCmd.MarkFlagRequired("name")
}
