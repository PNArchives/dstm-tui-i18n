package interaction

import (
	"bufio"
	"fmt"
	"os"
)

func Confirm(message string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(message)
		fmt.Println("(1) Yes  (2) No")
		fmt.Println("请输入选项数字> ")
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "1":
			return true
		case "2":
			return false
		default:
			fmt.Println("请输入正确的数字！")
		}
	}
}
