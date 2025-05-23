package cmd

import (
	"fmt"
	"gdrivecli/utils"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Authorize and store Google Drive token",
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.Authorize()
		if err != nil {
			fmt.Println("❌ Config failed:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
