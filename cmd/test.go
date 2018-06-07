package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")
		config := AnsibleConfig{
			"/Users/karl/.ansible/roles/fubarhouse.golang",
			"/etc/ansible/roles/role_under_test",
			"",
			"test.yml",
		}
		Ubuntu1604.run(&config)
		Ubuntu1604.install(&config)
		Ubuntu1604.test_syntax(&config)
		Ubuntu1604.test_role(&config)
		Ubuntu1604.test_idempotence(&config)
		Ubuntu1604.kill()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	pwd, _ := os.Getwd()
	testCmd.Flags().StringVarP(&containerID, "name", "n", "", "Name of the container")
	testCmd.Flags().StringVarP(&source, "source", "s", pwd, "Location of the role to test")
	testCmd.Flags().StringVarP(&destination, "destination", "d", "/etc/ansible/roles/role_under_test", "Location which the role will be mounted to")
	testCmd.Flags().StringVarP(&requirements, "requirements", "r", "", "Path to requirements file.")
	testCmd.Flags().StringVarP(&playbook, "playbook", "p", "playbook.yml", "The filename of the playbook")
}
