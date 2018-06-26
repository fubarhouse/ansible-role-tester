package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys a container with a specified ID",
	Long: `Destroys a container with a specified ID
`,
	Run: func(cmd *cobra.Command, args []string) {
		dist, _ := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		dist.CID = containerID
		if dist.DockerCheck() {
			dist.DockerKill()
		} else {
			log.Warnf("Container %v is not currently running", dist.CID)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().StringVarP(&containerID, "name", "n", "", "Container ID")
	destroyCmd.MarkFlagRequired("name")
}
