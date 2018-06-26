package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/ansible-role-tester/util"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Run installation tasks for the mounted role",
	Long:  `Run installation tasks for the mounted role (--name $NAME)`,
	Run: func(cmd *cobra.Command, args []string) {
		config := util.AnsibleConfig{
			HostPath:         source,
			RemotePath:       destination,
			RequirementsFile: requirements,
			PlaybookFile:     playbook,
			Verbose:          verbose,
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
	installCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	installCmd.MarkFlagRequired("name")
}
