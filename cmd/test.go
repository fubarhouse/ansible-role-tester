	package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
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
			HostPath:         source,
			RemotePath:       destination,
			RequirementsFile: requirements,
			PlaybookFile:     playbook,
			Verbose:          verbose,
		}

		dist, _ := util.GetDistribution(image, image, "/sbin/init", "/sys/fs/cgroup:/sys/fs/cgroup:ro", user, distro)

		dist.CID = containerID

		if dist.DockerCheck() {
			dist.RoleSyntaxCheck(&config)
			dist.RoleTest(&config)
			dist.IdempotenceTest(&config)
		} else {
			log.Warnf("Container %v is not currently running", dist.CID)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&containerID, "name", "n", containerID, "Container ID")
	testCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for Ansible commands.")
	testCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")

	testCmd.MarkFlagRequired("name")
}
