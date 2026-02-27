package tui

import (
	"github.com/EthicalGopher/GoFortify/tui/all_logs"
	"github.com/EthicalGopher/GoFortify/tui/cite"
	"github.com/EthicalGopher/GoFortify/tui/main_menu"
	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	Menu View = iota
	Logs
	Cite
)

type RootModel struct {
	view   View
	menu   main_menu.Model
	logs   all_logs.Model
	cite   cite.Model
	width  int
	height int
}

func NewRoot() tea.Model {
	return RootModel{
		view: Menu,
		menu: main_menu.New(),
		logs: all_logs.New(),
		cite: cite.New(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate to current view and update other models' sizes if needed
		// For now, let's just update the models so they have the latest size
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.menu, cmd = m.menu.Update(msg)
		cmds = append(cmds, cmd)

		updatedLogs, cmd := m.logs.Update(msg)
		m.logs = updatedLogs.(all_logs.Model)
		cmds = append(cmds, cmd)

		updatedCite, cmd := m.cite.Update(msg)
		m.cite = updatedCite.(cite.Model)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	switch m.view {

	case Menu:
		{
			var cmd tea.Cmd
			m.menu, cmd = m.menu.Update(msg)

			if m.menu.StartLogs {
				m.menu.StartLogs = false
				m.view = Logs
				// Ensure logs model has correct size
				updatedLogs, _ := m.logs.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				m.logs = updatedLogs.(all_logs.Model)
				// Explicitly initialize the logs model when switching to its view.
				return m, m.logs.Init()
			}
			if m.menu.StartCite {
				m.menu.StartCite = false
				m.view = Cite
				updatedCite, _ := m.cite.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				m.cite = updatedCite.(cite.Model)
				return m, m.cite.Init()
			}
			return m, cmd
		}
	case Logs:
		{
			var cmd tea.Cmd
			var updatedModel tea.Model
			updatedModel, cmd = m.logs.Update(msg)
			m.logs = updatedModel.(all_logs.Model)
			if m.logs.Done {
				m.logs.Done = false
				m.view = Menu
			}
			return m, cmd
		}
	case Cite:
		{
			var cmd tea.Cmd
			var updatedModel tea.Model
			updatedModel, cmd = m.cite.Update(msg)
			m.cite = updatedModel.(cite.Model)
			if m.cite.Done {
				m.cite.Done = false
				m.view = Menu
			}
			return m, cmd
		}
	default:
		return m, nil
	}
}
func (m RootModel) View() string {
	if m.view == Logs {
		return m.logs.View()
	}
	if m.view == Cite {
		return m.cite.View()
	}
	return m.menu.View()

}