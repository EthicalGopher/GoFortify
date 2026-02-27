package main_menu

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

type Model struct {
	List      list.Model
	StartLogs bool
	StartCite bool
	StartInit bool
	Quit      bool
}

func New() Model {
	items := []list.Item{
		Item{"All Logs", "Views all the logs"},
		Item{"Reset Logs", "Remove previous vulnerability and rate-limit logs"},
		Item{"Configure Protected Paths", "Set sensitive or hidden application routes"},
		Item{"Cite GoFortify", "Display how to cite this project"},
		Item{"Exit", "Quit GoFortify"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "GoFortify"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)

	return Model{
		List: l,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			selected := m.List.SelectedItem().(Item)
			if selected.title == "All Logs" {
				m.StartLogs = true
			}
			if selected.title == "Cite GoFortify" {
				m.StartCite = true
			}
			if selected.title == "Exit" {
				return m, tea.Quit
			}

		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.List.View())
}
