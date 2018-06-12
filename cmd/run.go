package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts a container",
	Long: `Start a container from a specified image.

Volume mount locations image and id are all configurable.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := AnsibleConfig{
			source,
			destination,
			"",
			"",
		}

		e, dist := getDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil {
			log.Errorln("Incompatible disribution was inputted, attempting autofix.")
		}

		dist.run(&config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	pwd, _ := os.Getwd()
	runCmd.Flags().StringVarP(&containerID, "id", "n", containerID, "Container ID")
	runCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	runCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")

	runCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	runCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	runCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")
}
