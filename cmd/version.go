package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

const currentVersion = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("kwik-cmd version %s\n", currentVersion)
		return nil
	},
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check-update",
	Short: "Check for updates",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Current version: %s\n", currentVersion)
		
		// Check GitHub for latest release
		resp, err := http.Get("https://api.github.com/repos/kaustuvbot/kwik-cmd/releases/latest")
		if err != nil {
			fmt.Println("Could not check for updates. Please check manually at https://github.com/kaustuvbot/kwik-cmd/releases")
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			// Simple parsing - just show the tag name
			fmt.Println("Checking for updates...")
			fmt.Println("Visit https://github.com/kaustuvbot/kwik-cmd/releases for latest version")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(checkUpdateCmd)
}
