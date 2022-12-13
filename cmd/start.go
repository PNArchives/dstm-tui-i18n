/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/PNCommand/dstm/dst/server"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := cmd.Flags().Lookup("cluster").Value.String()
		shardName := cmd.Flags().Lookup("shard").Value.String()
		if err := server.StartShard(clusterName, shardName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("cluster", "c", "", "specify a cluster")
	startCmd.Flags().StringP("shard", "s", "", "specify a shard")
	startCmd.MarkFlagsRequiredTogether("cluster", "shard")
}
