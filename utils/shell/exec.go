package shell

import (
	"os"
	"os/exec"
)

// func Sample() {
// 	// cmd := MakeCmdUseStdIO("bash", "-c", "./test.sh")
// 	cmd := MakeCmdUseStdIO("ping", "-c", "6", "google.com")
// 	if err := cmd.Start(); err != nil {
// 		panic("error")
// 	}
// 	time.Sleep(time.Second * 3)
// 	os.Stdin.Write([]byte("他の処理...\n"))
// 	cmd.Wait()
// }

func MakeCmdUseStdIO(command string, arg ...string) *exec.Cmd {
	cmd := exec.Command(command, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
