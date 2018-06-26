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
	"github.com/fubarhouse/ansible-role-tester/util"
	log "github.com/Sirupsen/logrus"
)

// testCmd represents the test command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Run installation tasks for the mounted role",
	Long: `Run installation tasks for the mounted role (--name $NAME)`,
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

		if !config.IsAnsibleRole() {
			log.Fatalf("Path %v is not recognized as an Ansible role.", config.HostPath)
		}
		if dist.DockerCheck() {
			dist.RoleInstall(&config)
		} else {
			log.Warnf("Container %v is not currently running", dist.CID)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	installCmd.MarkFlagRequired("name")
}
