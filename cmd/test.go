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
		"github.com/spf13/cobra"
	"os"

	"github.com/fubarhouse/ansible-role-tester/util"
	"log"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Execute tests against an existing container",
	Long: `Execute tests against an existing container

If container does not exist it will be created, however
containers won't be removed after completion.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := util.AnsibleConfig{
			HostPath: source,
			RemotePath: destination,
			RequirementsFile: requirements,
			PlaybookFile: playbook,
			Verbose: verbose,
		}

		dist, e := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil {
			log.Fatalln("Incompatible distribution was inputted.")
		}

		dist.CID = containerID
		dist.RoleTest(&config)
		dist.IdempotenceTest(&config)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	pwd, _ := os.Getwd()
	testCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	testCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	testCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	testCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	testCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
	testCmd.Flags().BoolVarP(&noOutput, "no-output", "o", false, "Hide output from all Docker commands")
	testCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")

	testCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	testCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	testCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")

	testCmd.MarkFlagRequired("name")
}
