package server

import (
	"errors"
	"strings"

	"github.com/PNCommand/dstm/utils"
	"github.com/PNCommand/dstm/utils/shell"
	"github.com/spf13/viper"
)

func StopShard(clusterName, shardName string) error {
	sessionName := generateSessionName(clusterName, shardName)
	if !isShardRunning(sessionName) {
		return errors.New("世界 " + sessionName + " 处于关闭状态！")
	}

	shell.SendCmdToTmuxSession(sessionName, "c_shutdown(true)")
	return nil
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

func RestartShard(clusterName, shardName string) {
	sessionName := generateSessionName(clusterName, shardName)
	if isShardRunning(sessionName) {
		StopShard(clusterName, shardName)
	}
	StartShard(clusterName, shardName)
}
