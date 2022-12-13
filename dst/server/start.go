package server

import (
	"errors"
	"fmt"
	"time"

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

	shell.CreateNewTmuxSessionExecCmd(sessionName, cmd)

	fmt.Println("正在启动世界 " + sessionName)
	fmt.Println("启动需要时间，请等待 90 秒")
	fmt.Println("本脚本启动世界时禁止自动更新mod, 有需要请在Mod管理面板更新")
	fmt.Println("如果启动失败，可以再尝试一两次")

	for begin := time.Now(); time.Since(begin) < 30*time.Second; {
		result, err := shell.GrepTmuxSessionOutput(sessionName, "Sim paused")
		if err == nil && len(result) > 0 {
			fmt.Println("成功启动世界 " + sessionName)
			return nil
		}
		time.Sleep(time.Second)
	}

	return errors.New("世界 " + sessionName + " 启动失败！")
}
