package env

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	osName string = ""
	osVer  string = ""
	cpuBit string = ""

	userName string = ""
	isSudoer bool   = false
)

func CheckSystem() {
	if runtime.GOOS != "linux" {
		fmt.Println("Only support linux!!!")
		os.Exit(1)
	}

	checkOSInfo()
	checkCpuArch()
	checkCurrentuser()
	checkSudoPerm()

	switch osName {
	case "ubuntu", "debian", "centos", "rhel":
		fmt.Println("OS:", osName, osVer)
	default:
		fmt.Println("Unsupported OS:", osName)
		os.Exit(1)
	}
	if cpuBit != "64" {
		fmt.Println("Unsupported CPU architecture:", cpuBit, "bit")
		os.Exit(1)
	}
	fmt.Println("CPU:", cpuBit, "bit")

	roleList := []string{}
	if userName == "root" {
		roleList = append(roleList, "root")
	}
	if isSudoer {
		roleList = append(roleList, "sudoer")
	}

	fmt.Println("User:", userName)
	fmt.Println("Role:", strings.Join(roleList, " "))

	if ok := isAllDepsInstalled(); !ok {
		if userName != "root" && !(isSudoer && isInstalled("sudo")) {
			fmt.Println("Permission requires")
			os.Exit(1)
		}
		installDeps()
	}

	installSteam()
	downloadDST(false)
}
