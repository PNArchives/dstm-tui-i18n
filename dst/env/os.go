package env

import (
	"io"
	"os"
	"os/exec"
	"strings"

	l10n "github.com/PNCommand/dstm/localization"
)

// ubuntu, debian, rhel(RedHat), centos, ol(Oracle), amzn(Amazon)
func checkOSInfo() {
	filePath := "/usr/lib/os-release"

	if _, err := os.Stat(filePath); err != nil {
		panic(l10n.String4Data("_file_not_found", map[string]interface{}{"path": filePath}))
	}

	f, err := os.Open(filePath)
	if err != nil {
		panic(l10n.String4Data("_cannot_open_file", map[string]interface{}{"path": filePath}))
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(l10n.String4Data("_cannot_read_file", map[string]interface{}{"path": filePath}))
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			osName = strings.ReplaceAll(strings.Trim(line[3:], " "), "\"", "")
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			osVer = strings.ReplaceAll(strings.Trim(line[11:], " "), "\"", "")
		}
	}
}

func checkCpuArch() {
	bit, err := exec.Command("getconf", "LONG_BIT").Output()
	if err != nil {
		panic(l10n.String("_get_cpu_bit_failed"))
	}
	cpuBit = strings.TrimRight(string(bit), "\n")
}

func checkCurrentuser() {
	result, _ := exec.Command("whoami").Output()
	userName = strings.TrimRight(string(result), "\n")
}

func checkSudoPerm() {
	sudoGroupName := " sudo"
	switch osName {
	case "rhel", "centos":
		sudoGroupName = " wheel"
	}
	result, _ := exec.Command("groups").Output()
	groups := string(result)
	isSudoer = strings.Contains(groups, sudoGroupName)
}
