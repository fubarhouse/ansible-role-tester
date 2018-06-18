package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/fubarhouse/ansible-role-tester/util"
)

// fullCmd represents the full command
var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "Complete end-to-end test process.",
	Long: `Runs a complete end-to-end process which performs the following:
  - creates a container
  - installs a requirements file
  - test the role syntax
  - runs the role
  - tests for idempotence
  - removes the container
You should be able to dockerRun all of this from the role folder on
the local file system. If you encounter errors, there's a lot
of flexibility in configuration, just change the defaults as
required.
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := util.AnsibleConfig{
			source,
			destination,
			requirements,
			playbook,
			verbose,
		}

		dist, e := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil {
			log.Errorln("Incompatible disribution was inputted, attempting autofix.")
		}

		dist.CID = containerID
		dist.DockerRun(&config)
		dist.RoleInstall(&config)
		dist.RoleSyntaxCheck(&config)
		dist.RoleTest(&config)
		dist.IdempotenceTest(&config)
		dist.DockerKill()
	},
}

func init() {
	rootCmd.AddCommand(fullCmd)

	pwd, _ := os.Getwd()
	fullCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Name of the container")
	fullCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	fullCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	fullCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	fullCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
	fullCmd.Flags().BoolVarP(&noOutput, "no-output", "o", false, "Hide output from all Docker commands")
	fullCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")

	fullCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	fullCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	fullCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")
}