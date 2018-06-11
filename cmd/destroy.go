package cmd

import (
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys a container with a specified ID",
	Long: `Destroys a container with a specified ID
`,
	Run: func(cmd *cobra.Command, args []string) {
		kill()
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().StringVarP(&containerID, "id", "i", "", "Container ID")
	destroyCmd.MarkFlagRequired("containerID")
}
