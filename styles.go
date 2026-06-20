package main

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const listHeight = 14

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	inputBox     lipgloss.Style
}

var inputBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("51")).
	Padding(0, 1)

func newStyles(darkBG bool) styles {
	var s styles

	// Muted dark base, cyan selection, amber highlight
	s.title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("51")). // bright cyan
		MarginLeft(1)

	s.item = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")). // dim gray
		PaddingLeft(3)

	s.selectedItem = lipgloss.NewStyle().
		PaddingLeft(1).
		Bold(true).
		Foreground(lipgloss.Color("214")). // amber
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("214"))

	s.pagination = list.DefaultStyles(darkBG).PaginationStyle.PaddingLeft(3)

	return s
}

type item string

type itemDelegate struct {
	styles *styles
}

func (i item) FilterValue() string { return "" }

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	title, command := decodeRecord(string(i))
	dimCmd := lipgloss.NewStyle().Foreground(lipgloss.Color("#888780")).Render("$ " + command)
	brightCmd := lipgloss.NewStyle().Foreground(lipgloss.Color("#C8C8C3")).Render("$ " + command)

	if index == m.Index() {
		fmt.Fprint(w, d.styles.selectedItem.Render(title+"  "+brightCmd))
	} else {
		title = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A8A8A")).Render(title)
		fmt.Fprint(w, d.styles.item.Render(title+"  "+dimCmd))
	}
}
