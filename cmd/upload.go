package cmd

import (
	"fmt"
	"gdrivecli/utils"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [file path]",
	Short: "Upload a file to Google Drive",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		err := utils.UploadFile(filePath)
		if err != nil {
			fmt.Println("❌ Upload failed:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
