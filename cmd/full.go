package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
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
You should be able to run all of this from the role folder on
the local file system. If you encounter errors, there's a lot
of flexibility in configuration, just change the defaults as
required.

usage: ansible-role-tester full -i fubarhouse/docker-ansible:bionic -p test.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := AnsibleConfig{
			source,
			destination,
			requirements,
			playbook,
		}

		e, dist := getDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro")
		if e != nil {
			log.Errorln("Incompatible disribution was inputted, attempting autofix.")
		}
		dist.run(&config)
		dist.install(&config)
		dist.test_syntax(&config)
		dist.test_role(&config)
		dist.test_idempotence(&config)
		kill()
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
	fullCmd.Flags().StringVarP(&image, "image", "i", "fubarhouse/docker-ansible:bionic", "The image reference to use.")
}
