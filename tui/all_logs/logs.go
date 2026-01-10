package all_logs

import (
	"strings"

	"github.com/EthicalGopher/SentinelShield/tui/shared"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type logMsg string

func waitForLog() tea.Cmd {
	return func() tea.Msg {
		return logMsg(<-shared.LogChan)
	}
}

type Model struct {
	viewport viewport.Model
	lines    []string
	Done     bool
}

func New() Model {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return Model{
		viewport: vp,
	}
}

func (m Model) Init() tea.Cmd {
	// This is called when the view becomes active.
	return waitForLog()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 2
		m.viewport.Height = msg.Height - 5
		// Re-render content on resize
		m.viewport.SetContent(strings.Join(m.lines, "\n"))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Done = true
			return m, nil

		case "q":
			return m, tea.Quit
		}

	case logMsg:
		m.lines = append(m.lines, string(msg))
		m.viewport.SetContent(strings.Join(m.lines, "\n"))
		m.viewport.GotoBottom()
		cmds = append(cmds, waitForLog())
	}

	// This will handle all messages we don't explicitly handle,
	// including arrow keys, page up/down, etc. for scrolling.
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View() + m.helpView()
}

func (m Model) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • esc: Back . q:quit\n")
}
