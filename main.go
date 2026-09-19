package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type model struct {
	remaining int
	running   bool
}

func initialModel() model {
	return model{
		remaining: 25 * 60,
		running:   true,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ":
			if m.remaining > 0 {
				m.running = !m.running
			}
		}

	case tickMsg:
		if m.running && m.remaining > 0 {
			m.remaining--
		}

		if m.remaining <= 0 {
			m.remaining = 0
			m.running = false
			return m, nil
		}

		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7CFFB2")).
		Render("POMODORO")

	timer := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(fmt.Sprintf("%02d:%02d", m.remaining/60, m.remaining%60))

	status := "FOCUS"
	if m.remaining == 0 {
		status = "COMPLETE"
	} else if !m.running {
		status = "PAUSED"
	}

	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A0A0A0")).
		Render(status + "  •  Space: pause/resume  •  Q: quit")

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Padding(1, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7CFFB2")).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			"",
			timer,
			"",
			info,
		))
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}