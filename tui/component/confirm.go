package component

import (
	"fmt"
	"strings"

	l10n "github.com/PNCommand/dstm/localization"
	tui "github.com/PNCommand/dstm/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type Confirm struct {
	message string
	answer  bool
}

func (c Confirm) GetAns() bool { return c.answer }

func (c Confirm) Init() tea.Cmd { return nil }

func (c Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down", "left", "right":
			c.answer = !c.answer
		case "enter", " ":
			fmt.Println()
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c Confirm) View() string {
	var build strings.Builder

	build.WriteString(tui.CenterStyle.Render(c.message))
	build.WriteRune('\n')

	noOption := " " + l10n.String("_no")
	yesOption := " " + l10n.String("_yes")
	if c.answer {
		noOption = tui.NormalCursor + noOption
	} else {
		noOption = tui.FocusedCursor + noOption
	}
	if c.answer {
		yesOption = tui.FocusedCursor + yesOption
	} else {
		yesOption = tui.NormalCursor + yesOption
	}
	build.WriteString(tui.CenterStyle.Render(noOption + tui.Gap + yesOption))

	return tui.BorderStyle.Render(build.String())
}
