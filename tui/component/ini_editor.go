package component

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	l10n "github.com/PNCommand/dstm/localization"
	tui "github.com/PNCommand/dstm/tui"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	//go:embed json/cluster-ini.json
	clusterIniJson []byte
	//go:embed json/shard-ini.json
	shardIniJson []byte
)

type IniItem struct {
	Key        string   `json:"key"`
	Value      string   `json:"value"`
	Type       string   `json:"input-type"`
	Validation string   `json:"validation"`
	Options    []string `json:"options"`
	index      *int     `json:"-"`
}

func (ii *IniItem) init() {
	if ii.Type != "selector" || ii.index != nil {
		return
	}
	for index, value := range ii.Options {
		if ii.Value == value {
			ii.index = &index
			return
		}
	}
}

func validateIp(value string) error {
	if net.ParseIP(value) == nil {
		return errors.New(l10n.String("_invalid_ip"))
	}
	return nil
}

func validatePort(value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if port < 2000 || 65535 < port {
		return errors.New(l10n.String("_invalid_port"))
	}
	return nil
}

func validateNotEmpty(value string) error {
	if len(value) == 0 {
		return errors.New(l10n.String("_cannot_be_empty"))
	}
	return nil
}

type IniGroup struct {
	Name  string    `json:"name"`
	Items []IniItem `json:"items"`
}

type IniEditor struct {
	title string

	IniGroups []IniGroup
	groupIdx  int
	row       int

	paginator paginator.Model
	isEditing bool
	textField textinput.Model
	err       error
}

func NewIniEditor(isCluster bool, title string) (*IniEditor, error) {
	var editor IniEditor
	if isCluster {
		json.Unmarshal(clusterIniJson, &editor.IniGroups)
	} else {
		json.Unmarshal(shardIniJson, &editor.IniGroups)
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 5
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	editor.paginator = p

	t := textinput.New()
	t.CursorStyle = tui.FocusedStyle
	t.SetCursorMode(textinput.CursorBlink)
	t.CharLimit = 50
	t.Placeholder = "..."
	t.Prompt = ""
	t.Focus()
	editor.textField = t

	editor.title = title

	return &editor, nil
}

func (ie *IniEditor) currentGroup() *IniGroup {
	return &ie.IniGroups[ie.groupIdx]
}

func (ie *IniEditor) currentGroupItemsLength() int {
	return len(ie.currentGroup().Items)
}

func (ie *IniEditor) itemsOnPage() int {
	return ie.paginator.ItemsOnPage(ie.currentGroupItemsLength())
}

func (ie *IniEditor) updateTotalPages() {
	ie.paginator.SetTotalPages(ie.currentGroupItemsLength())
}

func (ie *IniEditor) selectedItem() *IniItem {
	index := ie.itemsOnPage()*ie.paginator.Page + ie.row
	return &ie.currentGroup().Items[index]
}

func (ie *IniEditor) validate(selected *IniItem) {
	if selected.Type == "selector" {
		ie.err = nil
		return
	}
	switch selected.Validation {
	case "ip":
		ie.err = validateIp(ie.textField.Value())
	case "port":
		ie.err = validatePort(ie.textField.Value())
	case "allow-empty":
		ie.err = nil
	default:
		ie.err = validateNotEmpty(ie.textField.Value())
	}
}

func (ie *IniEditor) Init() tea.Cmd { return nil }

func (ie *IniEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println()
			return ie, tea.Quit
		case " ", "enter":
			if ie.row == ie.itemsOnPage() {
				fmt.Println()
				return ie, tea.Quit
			}
			if ie.isEditing {
				selected := ie.selectedItem()
				if selected.Type == "selector" {
					selected.Value = selected.Options[*selected.index]
					ie.isEditing = false
				} else if msg.String() != " " {
					ie.validate(selected)
					if ie.err == nil {
						selected.Value = ie.textField.Value()
						ie.isEditing = false
					}
				}
			} else if ie.row < ie.itemsOnPage() {
				ie.isEditing = true
				selected := ie.selectedItem()
				if selected.Type == "textfield" {
					ie.textField.SetValue(selected.Value)
					ie.textField.CursorEnd()
				} else {
					selected.init()
				}
			}
		case "up":
			ie.row -= 1
			if ie.row < 0 {
				ie.row = ie.itemsOnPage()
			}
		case "down":
			ie.row += 1
			if ie.row > ie.itemsOnPage() {
				ie.row = 0
			}
		case "left":
			if ie.isEditing {
				selected := ie.selectedItem()
				if selected.Type == "selector" && *selected.index > 0 {
					*selected.index -= 1
				}
			} else if ie.paginator.Page == 0 {
				ie.groupIdx -= 1
				if ie.groupIdx < 0 {
					ie.groupIdx = 0
				} else {
					ie.updateTotalPages()
					ie.paginator.Page = ie.paginator.TotalPages - 1
				}
				ie.row = 0
			} else {
				ie.paginator.PrevPage()
				ie.row = 0
			}
		case "right":
			if ie.isEditing {
				selected := ie.selectedItem()
				if selected.Type == "selector" && *selected.index < len(selected.Options)-1 {
					*selected.index += 1
				}
			} else if ie.paginator.OnLastPage() {
				ie.groupIdx += 1
				if ie.groupIdx >= len(ie.IniGroups) {
					ie.groupIdx = len(ie.IniGroups) - 1
				} else {
					ie.updateTotalPages()
					ie.paginator.Page = 0
				}
				ie.row = 0
			} else {
				ie.paginator.NextPage()
				ie.row = 0
			}
		}
	}
	if ie.isEditing {
		t, cmd := ie.textField.Update(msg)
		ie.textField = t
		return ie, cmd
	}
	return ie, nil
}

func (ie *IniEditor) View() string {
	var b strings.Builder

	header := tui.CenterStyle.Render(ie.title)
	b.WriteString(header + "\n")

	if ie.groupIdx == 0 {
		b.WriteString(tui.CenterStyle.Render("   "+ie.currentGroup().Name+" ->") + "\n")
	} else if ie.groupIdx == len(ie.IniGroups)-1 {
		b.WriteString(tui.CenterStyle.Render("<- "+ie.currentGroup().Name+"   ") + "\n")
	} else {
		b.WriteString(tui.CenterStyle.Render("<- "+ie.currentGroup().Name+" ->") + "\n")
	}
	b.WriteString(tui.CenterStyle.Render(ie.paginator.View()) + "\n")

	items := ie.currentGroup().Items
	start, end := ie.paginator.GetSliceBounds(len(items))

	var bb strings.Builder
	for i, item := range items[start:end] {
		if !ie.isEditing && i == ie.row {
			bb.WriteString(tui.OptBlockStyle.Copy().Inherit(tui.FocusedStyle).Render(item.Key))
		} else {
			bb.WriteString(tui.OptBlockStyle.Render(item.Key))
		}
		bb.WriteString(tui.Gap)

		if ie.isEditing && i == ie.row {
			if item.Type == "textfield" {
				bb.WriteString(ie.textField.View())
				if ie.err != nil {
					bb.WriteString("      " + ie.err.Error())
				}
			} else {
				if *item.index == 0 {
					bb.WriteString("   ")
				} else {
					bb.WriteString("<- ")
				}
				bb.WriteString(tui.FocusedStyle.Render(item.Options[*item.index]))
				if *item.index == len(item.Options)-1 {
					bb.WriteString("   ")
				} else {
					bb.WriteString(" ->")
				}
			}
		} else {
			bb.WriteString(item.Value)
		}
		bb.WriteRune('\n')
	}
	b.WriteString(tui.EditorBlockStyle.Render(bb.String()) + "\n")
	if ie.row == ie.itemsOnPage() {
		b.WriteString(tui.CenterStyle.Copy().Inherit(tui.FocusedStyle).Render("OK"))
	} else {
		b.WriteString(tui.CenterStyle.Render("OK"))
	}

	return tui.BorderStyle.Render(b.String())
}
