package utils

import (
	"os"

	"github.com/spf13/viper"
)

func IsShard(dirPath string) bool {
	if !DirExists(dirPath) {
		return false
	}
	if !FileExists(dirPath + "server.ini") {
		return false
	}
	return true
}

func IsCluster(dirPath string) bool {
	if !DirExists(dirPath) {
		return false
	}
	if !FileExists(dirPath + "/cluster.ini") {
		return false
	}
	if !FileExists(dirPath + "/cluster_token.txt") {
		return false
	}
	return true
}

func IsClusterExists(clusterName string) bool {
	for _, cluster := range GetAllClusters() {
		if cluster == clusterName {
			return true
		}
	}
	return false
}

func GetAllClusters() []string {
	worldsDirPath := viper.GetString("kleiRootDir") + "/" + viper.GetString("worldsDirName")
	files, err := os.ReadDir(worldsDirPath)
	if err != nil {
		panic(err)
	}
	list := []string{}
	for _, dir := range files {
		if !dir.IsDir() {
			continue
		}
		if IsCluster(worldsDirPath + "/" + dir.Name()) {
			list = append(list, dir.Name())
		}
	}
	return list
}

func GetAllShards(clusterList []string) []string {
	list := []string{}
	for _, clusterName := range clusterList {
		shardList := GetShardListInCluster(clusterName)
		list = append(list, shardList...)
	}
	return list
}

func GetShardListInCluster(clusterName string) []string {
	clusterDirPath := viper.GetString("kleiRootDir") +
		"/" + viper.GetString("worldsDirName") + "/" + clusterName
	files, err := os.ReadDir(clusterDirPath)
	if err != nil {
		panic(err)
	}
	list := []string{}
	for _, dir := range files {
		if !dir.IsDir() {
			continue
		}
		if IsShard(clusterDirPath + "/" + dir.Name()) {
			list = append(list, dir.Name())
		}
	}
	return list
}
