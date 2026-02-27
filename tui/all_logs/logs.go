package all_logs

import (
	"fmt"
	"strings"

	"github.com/EthicalGopher/GoFortify/shared"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)
)

type logMsg string

func waitForLog() tea.Cmd {
	return func() tea.Msg {
		msg := <-shared.LogChan
		return logMsg(msg)
	}
}

type Model struct {
	viewport      viewport.Model
	lines         []string
	Done          bool
	ready         bool
	width, height int
}

func New() Model {
	return Model{
		lines: []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return waitForLog()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Done = true
			return m, nil

		case "q", "ctrl+c":
			return m, tea.Quit

		case "c":
			m.lines = []string{}
			m.viewport.SetContent("")
			return m, nil
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(strings.Join(m.lines, "\n"))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		m.width = msg.Width
		m.height = msg.Height

	case logMsg:
		cleanMsg := strings.TrimSpace(string(msg))
		if cleanMsg != "" {
			m.lines = append(m.lines, cleanMsg)
			if m.ready {
				m.viewport.SetContent(strings.Join(m.lines, "\n"))
				m.viewport.GotoBottom()
			}
		}
		cmds = append(cmds, waitForLog())
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing logs..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m Model) headerView() string {
	title := headerStyle.Render("🛡️  GoFortify | Firewall Logs")
	line := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Repeat("─", max(0, m.width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := footerStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	help := helpStyle.Render(" ↑/↓: scroll • c: clear • esc: back • q: quit")
	line := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Repeat("─", max(0, m.width-lipgloss.Width(info)-lipgloss.Width(help)-4)))
	return lipgloss.JoinHorizontal(lipgloss.Center, " ", help, " ", line, " ", info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
