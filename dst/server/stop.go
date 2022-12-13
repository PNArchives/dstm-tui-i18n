package server

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
	fmt.Println("正在关闭世界 " + sessionName)

	for begin := time.Now(); time.Since(begin) < 30*time.Second; {
		if !isShardRunning(sessionName) {
			fmt.Println("成功关闭世界 " + sessionName)
			return nil
		}
		time.Sleep(time.Second)
	}

	return errors.New("世界 " + sessionName + " 关闭失败！")
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
