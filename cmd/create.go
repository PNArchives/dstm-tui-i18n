/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/PNCommand/dstm/dst/cluster"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := cmd.Flags().Lookup("cluster").Value.String()
		shardName := cmd.Flags().Lookup("shard").Value.String()
		if len(shardName) > 0 {
			fmt.Println(clusterName, shardName)
		} else {
			cluster.NewCluster(clusterName)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("cluster", "c", "", "specify a cluster")
	createCmd.Flags().StringP("shard", "s", "", "specify a shard")
	createCmd.MarkFlagRequired("cluster")
}
