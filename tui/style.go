package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	focusedColor = lipgloss.Color("205")

	borderColor = lipgloss.Color("63")

	grayColor = lipgloss.Color("0")
)

var (
	ScreenWidth, ScreenHeight, _ = term.GetSize(int(os.Stdout.Fd()))

	BorderedScreenWidth, BorderScreenHeight = ScreenWidth - 2, ScreenHeight - 2

	FullWidthStyle = lipgloss.NewStyle().Width(BorderedScreenWidth).MaxWidth(BorderedScreenWidth)

	FocusedStyle = lipgloss.NewStyle().Foreground(focusedColor)

	DisabledStyle = lipgloss.NewStyle().Foreground(grayColor)

	CenterStyle = FullWidthStyle.Copy().Align(lipgloss.Center)

	RightStyle = FullWidthStyle.Copy().Align(lipgloss.Right)

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)

	OptBlockStyle = lipgloss.NewStyle().Width(20).MaxWidth(20).Padding(0, 3).Inline(true)

	EditorBlockStyle = lipgloss.NewStyle().Height(5)
)

var (
	Gap = "     "

	FocusedCursor = FocusedStyle.Render("[x]")

	NormalCursor = "[ ]"
)
