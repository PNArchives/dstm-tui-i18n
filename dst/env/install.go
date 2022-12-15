package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/PNCommand/dstm/utils"
	"github.com/PNCommand/dstm/utils/interaction"
	"github.com/PNCommand/dstm/utils/shell"
	"github.com/spf13/viper"

	l10n "github.com/PNCommand/dstm/localization"
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

	fmt.Println(l10n.String("_before_install_libs"))
	fmt.Println("    " + useSudo + pkgManager + " update")
	fmt.Println("    " + useSudo + pkgManager + " upgrade -y")
	fmt.Println("    " + useSudo + pkgManager + " install -y " + strings.Join(deps, " "))

	ans := interaction.Confirm(l10n.String("_confirm_install_libs"))
	if !ans {
		fmt.Println(l10n.String("_cancel_install_libs"))
		os.Exit(0)
	}
	cmd := shell.MakeCmdUseStdIO("bash", "-c", useSudo+pkgManager+" update")
	cmd.Run()
	cmd = shell.MakeCmdUseStdIO("bash", "-c", useSudo+pkgManager+" upgrade -y")
	cmd.Run()
	cmd = shell.MakeCmdUseStdIO("bash", "-c", useSudo+pkgManager+" install -y "+strings.Join(deps, " "))
	cmd.Run()
	if !isAllDepsInstalled() {
		fmt.Println(l10n.String("_install_libs_failed"))
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
			fmt.Println(l10n.String("_mkdir_failed"), steamDir)
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

func DeleteDST() {
	dstRootDir := viper.GetString("dstRootDir")
	binPath := dstRootDir + "/bin64/dontstarve_dedicated_server_nullrenderer_x64"

	if !utils.FileExists(binPath) {
		return
	}

	fmt.Println(l10n.String("_delete_dst"))
	if err := os.RemoveAll(dstRootDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(l10n.String("_dst_deleted"))
}

func downloadDST(isUpdate bool) {
	steamDir := os.ExpandEnv("$HOME/Steam")
	scriptPath := steamDir + "/steamcmd.sh"
	dstRootDir := viper.GetString("dstRootDir")
	binPath := dstRootDir + "/bin64/dontstarve_dedicated_server_nullrenderer_x64"

	if isUpdate {
		fmt.Println(l10n.String("_update_dst"))
	} else {
		if utils.FileExists(binPath) {
			fmt.Println(l10n.String("_dst_ready"))
			return
		}
		fmt.Println(l10n.String4Data("_download_dst", map[string]interface{}{"path": dstRootDir}))
	}

	fmt.Println(l10n.String("_download_need_time"))
	// ~/Steam/steamcmd.sh +force_install_dir $1 +login anonymous +app_update 343050 validate +quit
	cmd := shell.MakeCmdUseStdIO(scriptPath, "+force_install_dir", dstRootDir, "+login", "anonymous", "+app_update", "343050", "validate", "+quit")
	if err := cmd.Run(); err != nil {
		fmt.Println(l10n.String("_download_err"))
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(l10n.String("_finish_download"))
}
