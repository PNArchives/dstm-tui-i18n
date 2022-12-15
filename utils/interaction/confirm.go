package interaction

import (
	"bufio"
	"fmt"
	"os"

	l10n "github.com/PNCommand/dstm/localization"
)

func Confirm(message string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(message)
		fmt.Println("(1) Yes  (2) No")
		fmt.Println(l10n.String("_enter_number"))
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "1":
			return true
		case "2":
			return false
		default:
			fmt.Println(l10n.String("_enter_correct_number"))
		}
	}
}
