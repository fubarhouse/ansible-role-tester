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
			source,
			destination,
			"",
			"",
			verbose,
		}

		dist, e := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil {
			log.Errorln("Incompatible disribution was inputted, attempting autofix.")
		}

		dist.CID = containerID
		dist.DockerRun(&config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	pwd, _ := os.Getwd()
	runCmd.Flags().StringVarP(&containerID, "id", "n", containerID, "Container ID")
	runCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	runCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	runCmd.Flags().BoolVarP(&noOutput, "no-output", "o", false, "Hide output from all Docker commands")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")

	runCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	runCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	runCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")
}
