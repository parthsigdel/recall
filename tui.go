package main

import (
	"bufio"
	"bytes"
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	textInput textinput.Model
	err       error
	quitting  bool

	list   list.Model
	choice string
	styles styles

	buf        *[]byte
	recallFile *os.File
}

func initialModel(buf []byte, recallFile *os.File) model {
	ti := textinput.New()
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.SetWidth(20)

	items := []list.Item{}
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)

	m := model{
		textInput:  ti,
		buf:        &buf,
		list:       l,
		recallFile: recallFile,
	}
	m.updateStyles(true) // default to dark styles.
	return m
}

func (m *model) updateStyles(isDark bool) {
	m.styles = newStyles(isDark)
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "ctrl+c": // quit
			m.quitting = true
			return m, tea.Quit

		case "enter": // print the selected command to the terminal.
			if len(m.list.Items()) == 0 {
				return m, nil
			}
			i, ok := m.list.SelectedItem().(item)
			if ok {
				_, command := decodeRecord(string(i))
				m.choice = command
			}
			return m, tea.Quit

		case "ctrl+d": // delete the selected record.
			i, ok := m.list.SelectedItem().(item)
			if ok {
				deleteRecord(m.recallFile, m.buf, string(i))
			}

			v := m.textInput.Value()
			if len(v) > 0 {
				items := filterList(m, v)
				m.list.SetItems(items)
				if len(items) > 0 {
					m.list.Select(0)
				}
			}

			return m, cmd

		case "ctrl+j":
			m.list.CursorDown()

		case "ctrl+k":
			m.list.CursorUp()

		case "up", "down", "left", "right":
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	v := m.textInput.Value()
	if len(v) > 0 {
		items := filterList(m, v)
		m.list.SetItems(items)
	}

	return m, cmd
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
	}

	if m.choice != "" || m.quitting {
		return tea.NewView("")
	}

	inputBox := inputBoxStyle.Render(m.textInput.View())
	str := lipgloss.JoinVertical(lipgloss.Top, inputBox, m.listView(), m.helpView())

	v := tea.NewView(str)
	if c != nil {
		c.X += 2 // border (1) + padding (1)
		c.Y += 1 // border top (1)
	}
	v.Cursor = c
	return v
}

func (m model) listView() string {
	if len(m.textInput.Value()) > 0 {
		return m.list.View()
	}
	// return empty space of same height to keep "help" section in a fixed position.
	return lipgloss.NewStyle().Height(listHeight).Render("")
}

func (m model) helpView() string {
	var help = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(
		`
        • ctrl+j/k, ↑/↓: navigate  
        • enter: select 
        • ctrl+d: delete 
        • ctrl+c: quit

`)
	return help
}

func filterList(m model, v string) []list.Item {
	items := []list.Item{}
	reader := bufio.NewReader(bytes.NewReader(*(m.buf)))
	match := searchInBuffer(reader, v)
	for _, line := range match {
		items = append(items, item(line))
	}
	return items
}
