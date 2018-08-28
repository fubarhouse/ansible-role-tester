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

	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
)

// runCmd represents the dockerRun command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts a container",
	Long: `Start a container from a specified image.

Volume mount locations image and id are all configurable.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := util.AnsibleConfig{
			HostPath:         source,
			RemotePath:       destination,
			ExtraRolesPath:   extraRoles,
			RequirementsFile: "",
			PlaybookFile:     "",
			Verbose:          verbose,
			Quiet:            quiet,
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
			util.CustomDistributionValueSet(&dist, "Container", fmt.Sprintf("%s/%s:%s", user, container, tag))
			util.CustomDistributionValueSet(&dist, "User", user)
			util.CustomDistributionValueSet(&dist, "Distro", image)
			util.CustomFamilyValueSet(&dist.Family, "Initialise", initialise)
			util.CustomFamilyValueSet(&dist.Family, "Volume", volume)
		}

		dist.CID = containerID

		if !config.IsAnsibleRole() && !quiet {
			log.Fatalf("Path %v is not recognized as an Ansible role.", config.HostPath)
		}
		if !dist.DockerCheck() {
			dist.DockerRun(&config)
		} else {
			if !quiet {
				log.Warnf("Container %v is already running", dist.CID)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	pwd, _ := os.Getwd()
	runCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	runCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	runCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	runCmd.Flags().StringVarP(&extraRoles, "extra-roles", "e", "", "Path to roles folder with dependencies.")
	runCmd.Flags().BoolVarP(&custom, "custom", "c", false, "Provide my own custom distribution.")
	runCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Enable quiet mode")

	runCmd.Flags().StringVarP(&initialise, "initialise", "a", "/bin/systemd", "The initialise command for the image")
	runCmd.Flags().StringVarP(&volume, "volume", "l", "/sys/fs/cgroup:/sys/fs/cgroup:ro", "The volume argument for the image")

	runCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	runCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	runCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")

	runCmd.MarkFlagRequired("name")
}
