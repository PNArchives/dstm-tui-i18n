package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/PNCommand/dstm/utils/shell"
	"github.com/spf13/viper"

	l10n "github.com/PNCommand/dstm/localization"
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
		return errors.New(l10n.String4Data("_shard_is_running", map[string]interface{}{"sessionName": sessionName}))
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

	fmt.Println(l10n.String4Data("_start_shard", map[string]interface{}{"sessionName": sessionName}))
	fmt.Println(l10n.String("_mod_auto_update_off"))

	for begin := time.Now(); time.Since(begin) < 30*time.Second; {
		result, err := shell.GrepTmuxSessionOutput(sessionName, "Sim paused")
		if err == nil && len(result) > 0 {
			fmt.Println(l10n.String4Data("_start_shard_succeeded", map[string]interface{}{"sessionName": sessionName}))
			return nil
		}
		time.Sleep(time.Second)
	}

	return errors.New(l10n.String4Data("_start_shard_failed", map[string]interface{}{"sessionName": sessionName}))
}
