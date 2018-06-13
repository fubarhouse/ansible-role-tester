package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"

	"github.com/fubarhouse/ansible-role-tester/util"
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
			source,
			destination,
			requirements,
			playbook,
			verbose,
		}

		dist, e := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)
		if e != nil {
			logrus.Errorln("Incompatible disribution was inputted, attempting autofix.")
		}

		dist.CID = containerID
		dist.RoleInstall(&config)
		dist.RoleSyntaxCheck(&config)
		dist.RoleTest(&config)
		dist.IdempotenceTest(&config)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	pwd, _ := os.Getwd()
	testCmd.Flags().StringVarP(&containerID, "id", "n", containerID, "Container ID")
	testCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	testCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	testCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	testCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
	testCmd.Flags().BoolVarP(&noOutput, "no-output", "o", false, "Hide output from all Docker commands")
	testCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")

	testCmd.Flags().StringVarP(&image, "image", "i", "", "The image reference to use.")
	testCmd.Flags().StringVarP(&user, "user", "u", "fubarhouse", "Selectively choose a compatible docker image from a specified user.")
	testCmd.Flags().StringVarP(&distro, "distribution", "t", "ubuntu1804", "Selectively choose a compatible docker image of a specified distribution.")
}
