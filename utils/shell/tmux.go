package shell

import (
	"os/exec"
	"strings"
)

func CreateNewTmuxSession(sessionName string) {
	if len(sessionName) > 0 {
		MakeCmdUseStdIO("tmux", "-u", "new", "-d", "-s", sessionName).Run()
	} else {
		MakeCmdUseStdIO("tmux", "-u", "new", "-d").Run()
	}
}

func KillTmuxSession(sessionName string) {
	MakeCmdUseStdIO("tmux", "kill-session", "-t", sessionName).Run()
}

func SendCmdToTmuxSession(sessionName, cmd string) {
	MakeCmdUseStdIO("tmux", "send-keys", "-t", sessionName, cmd, "ENTER").Run()
}

func ForceShutdownSession(sessionName string) {
	MakeCmdUseStdIO("tmux", "send-keys", "-t", sessionName, "C-c").Run()
}

func GetTmuxSessionList() []string {
	bytes, err := exec.Command("bash", "-c", "tmux ls | awk -F: '{print $1}'").Output()
	if err != nil {
		panic(err)
	}
	output := strings.Trim(string(bytes), "\n")
	return strings.Split(output, "\n")
}

func GetTmuxSessionOutput(sessionName string) string {
	bytes, err := exec.Command("tmux", "capture-pane", "-t", sessionName, "-p", "-S", "-").Output()
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(bytes), "\n")
}
