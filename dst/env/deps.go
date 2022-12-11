package env

import (
	"fmt"
	"os/exec"
	"strings"
)

func isInstalled(pkg string) bool {
	switch osName {
	case "debian", "ubuntu":
		bashCmd := "dpkg-query -l | awk '{print $2}' | grep -s " + pkg
		result, _ := exec.Command("bash", "-c", bashCmd).Output()
		return len(result) > 0
	case "rhel", "centos":
		bashCmd := "yum list installed | grep -s " + pkg
		result, _ := exec.Command("bash", "-c", bashCmd).Output()
		return len(result) > 0
	}
	return false
}

func getDeps() []string {
	switch osName {
	case "ubuntu":
		if strings.HasPrefix(osVer, "18") {
			return []string{"lib32gcc1", "lua5.3", "tmux", "wget", "curl"}
		} else {
			return []string{"lib32gcc-s1", "lua5.3", "tmux", "wget", "curl"}
		}
	case "debian":
		return []string{"lib32gcc", "lua5.3", "tmux", "wget", "curl"}
	case "rhel", "centos":
		return []string{"glibc.i686", "libstdc++", "lua", "tmux", "wget", "curl"}
	default:
		return []string{"Unsupported OS"}
	}
}

func isAllDepsInstalled() bool {
	deps := getDeps()
	isAllInstalled := true

	for _, dep := range deps {
		if isInstalled(dep) {
			fmt.Println("[OK]", dep)
		} else {
			fmt.Println("[NG]", dep)
			isAllInstalled = false
		}
	}
	return isAllInstalled
}
