package component

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Scanner struct {
	message   string
	textfield textinput.Model
	err       error
}

func NewScanner(message string) *Scanner {
	t := textinput.New()
	t.Focus()
	t.CharLimit = 50

	scanner := Scanner{
		message:   message,
		textfield: t,
		err:       nil,
	}
	return &scanner
}

func (s *Scanner) GetInput() string {
	return s.textfield.Value()
}

func (s *Scanner) Init() tea.Cmd {
	return textinput.Blink
}

func (s *Scanner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			s.err = validateNotEmpty(s.textfield.Value())
			if s.err == nil {
				fmt.Println()
				return s, tea.Quit
			}
		}
	}

	s.textfield, cmd = s.textfield.Update(msg)
	return s, cmd
}

func (s *Scanner) View() string {
	var b strings.Builder

	b.WriteString(s.message)
	if s.err == nil {
		b.WriteString("\n")
	} else {
		b.WriteString("   " + s.err.Error() + "\n")
	}
	b.WriteString(s.textfield.View())

	return b.String()
}
