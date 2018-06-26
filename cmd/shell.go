package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Shells into a container",
	Long:  `Shell into a container after creation.`,
	Run: func(cmd *cobra.Command, args []string) {
		arguments := []string{
			"exec",
			"-it",
			containerID,
			"bash",
		}

		dist, _ := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		dist.CID = containerID

		if dist.DockerCheck() {
			util.DockerExec(arguments, true)
		} else {
			log.Warnf("Container %v is not currently running", dist.CID)
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	shellCmd.MarkFlagRequired("name")
}
