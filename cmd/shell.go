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
	)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Shells into a container",
	Long: `Shell into a container after creation.`,
	Run: func(cmd *cobra.Command, args []string) {
		arguments := []string{
			"exec",
			"-it",
			containerID,
			"bash",
		}
		util.DockerExec(arguments, true)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	shellCmd.MarkFlagRequired("name")
}
