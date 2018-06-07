package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy called")
		fmt.Println(source)
		fmt.Println(destination)
		fmt.Println(requirements)
		fmt.Println(playbook)
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	pwd, _ := os.Getwd()
	destroyCmd.Flags().StringVarP(&containerID, "name", "n", "", "Name of the container")
	destroyCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	destroyCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	destroyCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	destroyCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
}
