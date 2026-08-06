package cmd

import "github.com/spf13/cobra"

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "A TUI application for viewing docker containers, images, and volumes",
	Long:  `A TUI application for viewing docker containers, images, and volumes.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
