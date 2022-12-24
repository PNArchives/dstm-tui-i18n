package component

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	l10n "github.com/PNCommand/dstm/localization"
	tui "github.com/PNCommand/dstm/tui"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfigItem struct {
	Key          string   `json:"key"`
	Display      string   `json:"display"`
	DefaultIndex int      `json:"index"`
	Index        *int     `json:"-"`
	Options      []string `json:"options"`
	OptsDisplay  []string `json:"opts-display"`
}

type ConfigGroup struct {
	Name    string       `json:"name"`
	Display string       `json:"display"`
	Items   []ConfigItem `json:"items"`
}

type WorldEditor struct {
	title string

	ConfigGroups []ConfigGroup
	groupIdx     int
	row          int

	paginator paginator.Model
	isEditing bool
}

func NewWorldEditor(filePath string, isForest, isGen bool) (*WorldEditor, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var editor WorldEditor
	json.NewDecoder(file).Decode(&editor.ConfigGroups)

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 5
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	editor.paginator = p

	if isForest {
		if isGen {
			editor.title = l10n.String("_forest_gen")
		} else {
			editor.title = l10n.String("_forest_set")
		}
	} else {
		if isGen {
			editor.title = l10n.String("_cave_gen")
		} else {
			editor.title = l10n.String("_cave_set")
		}
	}

	return &editor, nil
}

func (we *WorldEditor) currentGroup() *ConfigGroup {
	return &we.ConfigGroups[we.groupIdx]
}

func (we *WorldEditor) currentGroupItemsLength() int {
	return len(we.currentGroup().Items)
}

func (we *WorldEditor) itemsOnPage() int {
	return we.paginator.ItemsOnPage(we.currentGroupItemsLength())
}

func (we *WorldEditor) updateTotalPages() {
	we.paginator.SetTotalPages(we.currentGroupItemsLength())
}

func (we *WorldEditor) selectedItem() *ConfigItem {
	index := we.itemsOnPage()*we.paginator.Page + we.row
	return &we.ConfigGroups[we.groupIdx].Items[index]
}

func (we *WorldEditor) Init() tea.Cmd { return nil }

func (we *WorldEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println()
			return we, tea.Quit
		case " ", "enter":
			if we.row == we.itemsOnPage() {
				fmt.Println()
				return we, tea.Quit
			}
			if we.isEditing {
				we.isEditing = false
			} else if we.row < we.itemsOnPage() {
				we.isEditing = true
				selected := we.selectedItem()
				if selected.Index == nil {
					newIndex := selected.DefaultIndex
					selected.Index = &newIndex
				}
			}
		case "up":
			we.row -= 1
			if we.row < 0 {
				we.row = we.itemsOnPage()
			}
		case "down":
			we.row += 1
			if we.row > we.itemsOnPage() {
				we.row = 0
			}
		case "left":
			if we.isEditing {
				selected := we.selectedItem()
				if *selected.Index > 0 {
					*selected.Index -= 1
				}
			} else if we.paginator.Page == 0 {
				we.groupIdx -= 1
				if we.groupIdx < 0 {
					we.groupIdx = 0
				} else {
					we.updateTotalPages()
					we.paginator.Page = we.paginator.TotalPages - 1
				}
				we.row = 0
			} else {
				we.paginator.PrevPage()
				we.row = 0
			}
		case "right":
			if we.isEditing {
				selected := we.selectedItem()
				if *selected.Index < len(selected.Options)-1 {
					*selected.Index += 1
				}
			} else if we.paginator.OnLastPage() {
				we.groupIdx += 1
				if we.groupIdx >= len(we.ConfigGroups) {
					we.groupIdx = len(we.ConfigGroups) - 1
				} else {
					we.updateTotalPages()
					we.paginator.Page = 0
				}
				we.row = 0
			} else {
				we.paginator.NextPage()
				we.row = 0
			}
		}
	}
	return we, nil
}

func (we *WorldEditor) View() string {
	var b strings.Builder

	header := tui.CenterStyle.Render(we.title)
	b.WriteString(header + "\n")

	if we.groupIdx == 0 {
		b.WriteString(tui.CenterStyle.Render("   "+we.currentGroup().Display+" ->") + "\n")
	} else if we.groupIdx == len(we.ConfigGroups)-1 {
		b.WriteString(tui.CenterStyle.Render("<- "+we.currentGroup().Display+"   ") + "\n")
	} else {
		b.WriteString(tui.CenterStyle.Render("<- "+we.currentGroup().Display+" ->") + "\n")
	}
	b.WriteString(tui.CenterStyle.Render(we.paginator.View()) + "\n")

	items := we.currentGroup().Items
	start, end := we.paginator.GetSliceBounds(len(items))

	var bb strings.Builder
	for i, item := range items[start:end] {
		if !we.isEditing && i == we.row {
			bb.WriteString(tui.OptBlockStyle.Copy().Inherit(tui.FocusedStyle).Render(item.Display))
		} else {
			bb.WriteString(tui.OptBlockStyle.Render(item.Display))
		}
		bb.WriteString(tui.Gap)

		if we.isEditing && i == we.row {
			if *item.Index == 0 {
				bb.WriteString("   ")
			} else {
				bb.WriteString("<- ")
			}
			bb.WriteString(tui.FocusedStyle.Render(item.OptsDisplay[*item.Index]))
			if *item.Index == len(item.Options)-1 {
				bb.WriteString("   ")
			} else {
				bb.WriteString(" ->")
			}
		} else if item.Index == nil {
			bb.WriteString(item.OptsDisplay[item.DefaultIndex])
		} else {
			bb.WriteString(item.OptsDisplay[*item.Index])
		}
		bb.WriteRune('\n')
	}
	b.WriteString(tui.EditorBlockStyle.Render(bb.String()) + "\n")
	if we.row == we.itemsOnPage() {
		b.WriteString(tui.CenterStyle.Copy().Inherit(tui.FocusedStyle).Render("OK"))
	} else {
		b.WriteString(tui.CenterStyle.Render("OK"))
	}

	return tui.BorderStyle.Render(b.String())
}
