package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type phase int

const (
	focus phase = iota
	shortBreak
	longBreak
)

type timerConfig struct {
	focusSeconds      int
	shortBreakSeconds int
	longBreakSeconds  int
	cycles            int
}

type model struct {
	config    timerConfig
	phase     phase
	cycle     int
	remaining int
	running   bool
}

func initialModel(config timerConfig) model {
	return model{
		config:    config,
		phase:     focus,
		cycle:     1,
		remaining: config.focusSeconds,
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

func (m *model) reset() {
	m.phase = focus
	m.cycle = 1
	m.remaining = m.config.focusSeconds
	m.running = true
}

func (m *model) advancePhase() {
	switch m.phase {
	case focus:
		if m.cycle >= m.config.cycles {
			m.phase = longBreak
			m.remaining = m.config.longBreakSeconds
		} else {
			m.phase = shortBreak
			m.remaining = m.config.shortBreakSeconds
		}
	case shortBreak:
		m.phase = focus
		m.cycle++
		m.remaining = m.config.focusSeconds
	case longBreak:
		m.phase = focus
		m.cycle = 1
		m.remaining = m.config.focusSeconds
	}
	m.running = true
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
		case "r":
			m.reset()
		case "n":
			m.advancePhase()
		}

	case tickMsg:
		if m.running && m.remaining > 0 {
			m.remaining--
		}

		if m.remaining <= 0 {
			m.advancePhase()
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

	status := map[phase]string{
		focus:      "FOCUS",
		shortBreak: "SHORT BREAK",
		longBreak:  "LONG BREAK",
	}[m.phase]
	if !m.running {
		status += " (PAUSED)"
	}
	cycle := fmt.Sprintf("Cycle %d/%d", m.cycle, m.config.cycles)

	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A0A0A0")).
		Render(status + "  •  " + cycle + "\nSpace: pause/resume  R: reset  N: skip  Q: quit")

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
	focusMinutes := flag.Int("focus", 25, "focus duration in minutes")
	shortBreakMinutes := flag.Int("short-break", 5, "short break duration in minutes")
	longBreakMinutes := flag.Int("long-break", 15, "long break duration in minutes")
	cycles := flag.Int("cycles", 4, "focus sessions before a long break")
	flag.Parse()

	if *focusMinutes <= 0 || *shortBreakMinutes <= 0 || *longBreakMinutes <= 0 || *cycles <= 0 {
		fmt.Fprintln(os.Stderr, "durations and cycles must be greater than zero")
		os.Exit(1)
	}

	config := timerConfig{
		focusSeconds:      *focusMinutes * 60,
		shortBreakSeconds: *shortBreakMinutes * 60,
		longBreakSeconds:  *longBreakMinutes * 60,
		cycles:            *cycles,
	}
	p := tea.NewProgram(initialModel(config))

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
