// Copyright © 2018 Karl Hepworth karl.hepworth@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	util "github.com/fubarhouse/ansible-role-tester/util"
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
			HostPath: source,
			RemotePath: destination,
			RequirementsFile: "",
			PlaybookFile: "",
			Verbose: verbose,
		}

		dist, e := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil {
			log.Fatalln("Incompatible distribution was inputted.")
		}

		dist.CID = containerID
		dist.DockerRun(&config)

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	pwd, _ := os.Getwd()
	runCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	runCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	runCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")

	runCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	runCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	runCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")

	runCmd.MarkFlagRequired("name")
}
