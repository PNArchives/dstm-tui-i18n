/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/PNCommand/dstm/dst/env"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
		env.DownloadDST(true)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
