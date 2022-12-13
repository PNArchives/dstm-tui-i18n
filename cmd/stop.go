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

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use: "stop",
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := cmd.Flags().Lookup("cluster").Value.String()
		shardName := cmd.Flags().Lookup("shard").Value.String()
		if err := server.StopShard(clusterName, shardName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().StringP("cluster", "c", "", "specify a cluster")
	stopCmd.Flags().StringP("shard", "s", "", "specify a shard")
	stopCmd.MarkFlagsRequiredTogether("cluster", "shard")
}
