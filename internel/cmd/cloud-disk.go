package cmd

import (
	cloud_disk "github.com/YFR718/cmd-tool/internel/cloud-disk"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "file",
	Short: "云文件系统",
	Long:  "云文件系统",
	Run: func(cmd *cobra.Command, args []string) {
		cloud_disk.Run()
	},
}
