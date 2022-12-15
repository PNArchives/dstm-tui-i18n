package server

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PNCommand/dstm/utils"
	"github.com/PNCommand/dstm/utils/shell"
	"github.com/spf13/viper"

	l10n "github.com/PNCommand/dstm/localization"
)

func StopShard(clusterName, shardName string) error {
	sessionName := generateSessionName(clusterName, shardName)
	if !isShardRunning(sessionName) {
		return errors.New(l10n.String4Data("_shard_is_closed", map[string]interface{}{"sessionName": sessionName}))
	}

	shell.SendCmdToTmuxSession(sessionName, "c_shutdown(true)")
	fmt.Println(l10n.String4Data("_close_shard", map[string]interface{}{"sessionName": sessionName}))

	for begin := time.Now(); time.Since(begin) < 30*time.Second; {
		if !isShardRunning(sessionName) {
			fmt.Println(l10n.String4Data("_close_shard_succeeded", map[string]interface{}{"sessionName": sessionName}))
			return nil
		}
		time.Sleep(time.Second)
	}

	return errors.New(l10n.String4Data("_close_shard_failed", map[string]interface{}{"sessionName": sessionName}))
}

func StopAllShardsInCluster(clusterName string) {
	shardList := utils.GetShardListInCluster(clusterName)
	for _, shardname := range shardList {
		StopShard(clusterName, shardname)
	}
}

func StopAllShards() {
	sessionList := shell.GetTmuxSessionList()
	for _, session := range sessionList {
		names := strings.Split(session, viper.GetString("separator"))
		StopShard(names[0], names[1])
	}
}

func RestartShard(clusterName, shardName string) error {
	sessionName := generateSessionName(clusterName, shardName)
	if isShardRunning(sessionName) {
		StopShard(clusterName, shardName)
	}
	return StartShard(clusterName, shardName)
}
