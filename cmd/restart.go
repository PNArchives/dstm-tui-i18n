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

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use: "restart",
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := cmd.Flags().Lookup("cluster").Value.String()
		shardName := cmd.Flags().Lookup("shard").Value.String()
		if err := server.RestartShard(clusterName, shardName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)

	restartCmd.Flags().StringP("cluster", "c", "", "specify a cluster")
	restartCmd.Flags().StringP("shard", "s", "", "specify a shard")
	restartCmd.MarkFlagsRequiredTogether("cluster", "shard")
}
