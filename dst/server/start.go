package server

import (
	"errors"

	"github.com/PNCommand/dstm/utils/shell"
	"github.com/spf13/viper"
)

func generateSessionName(clusterName, shardName string) string {
	return clusterName + viper.GetString("separator") + shardName
}

func isShardRunning(sessionName string) bool {
	sessionList := shell.GetTmuxSessionList()
	for _, name := range sessionList {
		if name == sessionName {
			return true
		}
	}
	return false
}

func StartShard(clusterName, shardName string) error {
	sessionName := generateSessionName(clusterName, shardName)
	if isShardRunning(sessionName) {
		return errors.New("世界 " + sessionName + " 处于启动状态！")
	}

	dstRootDir := viper.GetString("dstRootDir")
	ugcDir := viper.GetString("ugcDir")
	kleiRootDir := viper.GetString("kleiRootDir")
	worldsDirName := viper.GetString("worldsDirName")

	cmd := "cd " + dstRootDir + "/bin64;"
	cmd += " ./dontstarve_dedicated_server_nullrenderer_x64"
	cmd += " -skip_update_server_mods"
	cmd += " -ugc_directory " + ugcDir
	cmd += " -persistent_storage_root " + kleiRootDir
	cmd += " -conf_dir " + worldsDirName
	cmd += " -cluster ttt -shard Main"

	shell.CreateNewTmuxSession(sessionName)
	shell.SendCmdToTmuxSession(sessionName, cmd)
	return nil
}
