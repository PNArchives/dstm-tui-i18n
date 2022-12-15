package env

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	l10n "github.com/PNCommand/dstm/localization"
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
		fmt.Println(l10n.String("_not_supported_os"))
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
		fmt.Println(l10n.String4Data("_not_supported_distribution", map[string]interface{}{"osName": osName}))
		os.Exit(1)
	}
	if cpuBit != "64" {
		fmt.Println(cpuBit, l10n.String4Data("_not_supported_bit", map[string]interface{}{"bit": cpuBit}))
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
			fmt.Println(l10n.String("_need_permission"))
			os.Exit(1)
		}
		installDeps()
	}

	installSteam()
	downloadDST(false)
}
