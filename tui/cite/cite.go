package cite

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1)

	contentStyle = lipgloss.NewStyle().
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 2)
)

type Model struct {
	viewport      viewport.Model
	ready         bool
	Done          bool
	width, height int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Done = true
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		headerHeight := 3
		footerHeight := 2
		contentHeight := msg.Height - headerHeight - footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, contentHeight)
			m.viewport.SetContent(m.getContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = contentHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := titleStyle.Render("🛡️  GoFortify | Citation Information")
	line := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Repeat("─", max(0, m.width)))
	
	footer := helpStyle.Render(" esc: back • q: quit")

	return fmt.Sprintf("%s\n%s\n%s\n%s", header, line, m.viewport.View(), footer)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) getContent() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Plain Text Citation:") + "\n")
	b.WriteString("EthicalGopher. (2026). GoFortify: A High-Performance Security Reverse Proxy & Traffic Inspector. https://github.com/EthicalGopher/GoFortify\n\n")

	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("BibTeX Entry:") + "\n")
	
	bibtex := `@software{GoFortify_2026,
  author = {EthicalGopher},
  title = {{GoFortify: A High-Performance Security Reverse Proxy \& Traffic Inspector}},
  url = {https://github.com/EthicalGopher/GoFortify},
  year = {2026},
  version = {1.0.0}
}`
	b.WriteString(bibtex)

	return contentStyle.Render(b.String())
}
