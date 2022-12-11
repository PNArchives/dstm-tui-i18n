package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/PNCommand/dstm/utils"
	"github.com/PNCommand/dstm/utils/interaction"
	"github.com/PNCommand/dstm/utils/shell"
	"github.com/spf13/viper"
)

func getPkgManager() string {
	switch osName {
	case "debian", "ubuntu":
		return "apt-get"
	case "rhel", "centos":
		return "yum"
	default:
		return ""
	}
}

func installDeps() {
	deps := getDeps()
	pkgManager := getPkgManager()
	useSudo := ""
	if userName != "root" {
		useSudo = "sudo "
	}

	fmt.Println("Next commands will be executed:")
	fmt.Println("    " + useSudo + pkgManager + " update")
	fmt.Println("    " + useSudo + pkgManager + " upgrade -y")
	fmt.Println("    " + useSudo + pkgManager + " install -y " + strings.Join(deps, " "))

	ans := interaction.Confirm("Ready?")
	if !ans {
		fmt.Println("Canceled")
		return
	}
	cmd := shell.MakeCmdUseStdIO("bash", "-c", useSudo+pkgManager+" update")
	cmd.Run()
	cmd = shell.MakeCmdUseStdIO("bash", "-c", useSudo+pkgManager+" upgrade -y")
	cmd.Run()
	cmd = shell.MakeCmdUseStdIO("bash", "-c", useSudo+pkgManager+" install -y "+strings.Join(deps, " "))
	cmd.Run()
	if !isAllDepsInstalled() {
		fmt.Println("Emmmm")
		os.Exit(1)
	}
}

// if [[ $1 == 'CentOS' ]] && [[ ! -e /usr/lib64/libcurl-gnutls.so.4 ]]; then
//     # To fix: libcurl-gnutls.so.4: cannot open shared object file: No such file or directory
//     if [[ $(whoami) == 'root' ]]; then
//         ln -s /usr/lib64/libcurl.so.4 /usr/lib64/libcurl-gnutls.so.4
//     else
//         sudo ln -s /usr/lib64/libcurl.so.4 /usr/lib64/libcurl-gnutls.so.4
//     fi
// fi

func installSteam() {
	steamDir := os.ExpandEnv("$HOME/Steam")
	scriptPath := steamDir + "/steamcmd.sh"

	if !utils.DirExists(steamDir) {
		if err := os.MkdirAll(steamDir, 0755); err != nil {
			fmt.Println("Can not create directory:", steamDir)
			os.Exit(1)
		}
	}
	if utils.FileExists(scriptPath) {
		fmt.Println("[OK]", scriptPath)
		return
	}

	// 该用哪个网站？
	// wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz
	// wget http://media.steampowered.com/installer/steamcmd_linux.tar.gz
	url := "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
	tarFilePath := steamDir + "/steamcmd_linux.tar.gz"

	cmd := shell.MakeCmdUseStdIO("wget", "--output-document", tarFilePath, url)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd = shell.MakeCmdUseStdIO("tar", "-xvzf", tarFilePath, "--directory", steamDir)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Tmp() {
	deleteDST()
}

func deleteDST() {
	dstRootDir := viper.GetString("dstRootDir")
	binPath := dstRootDir + "/bin64/dontstarve_dedicated_server_nullrenderer_x64"

	if !utils.FileExists(binPath) {
		return
	}

	fmt.Println("删除旧的DST服务端...")
	if err := os.RemoveAll(dstRootDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("削除しました...")
}

func downloadDST(isUpdate bool) {
	steamDir := os.ExpandEnv("$HOME/Steam")
	scriptPath := steamDir + "/steamcmd.sh"
	dstRootDir := viper.GetString("dstRootDir")
	binPath := dstRootDir + "/bin64/dontstarve_dedicated_server_nullrenderer_x64"

	if isUpdate {
		fmt.Println("开始升级服务端...")
	} else {
		if utils.FileExists(binPath) {
			fmt.Println("饥荒服务端已经下载好啦～")
			return
		}
		fmt.Println("路径 " + dstRootDir + " 未找到饥荒服务端，即将开始下载...")
	}

	fmt.Println("根据网络状况，可能会很耗时间，更新完成为止请勿息屏")
	// ~/Steam/steamcmd.sh +force_install_dir $1 +login anonymous +app_update 343050 validate +quit
	cmd := shell.MakeCmdUseStdIO(scriptPath, "+force_install_dir", dstRootDir, "+login", "anonymous", "+app_update", "343050", "validate", "+quit")
	if err := cmd.Run(); err != nil {
		fmt.Println("似乎出现了什么错误(国内主机的话, 一般是网络问题)...")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("饥荒服务端下载更新完成!")
	fmt.Println("途中可能会出现Failed to init SDL priority manager: SDL not found警告")
	fmt.Println("不用担心, 这个不影响下载/更新DST")
	fmt.Println("虽然可以解决, 但这需要下载一堆依赖包, 有可能会对其他运行中的服务造成影响, 所以无视它吧～")
}
