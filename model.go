package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	viewDashboard viewState = iota
	viewDetail
	viewHelp
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Padding(0, 1)
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	criticalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	detailLabel   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	detailValue   = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	cursorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
)

type tickMsg time.Time

type model struct {
	systems      []UrbanSystem
	cursor       int
	view         viewState
	rng          *rand.Rand
	tickInterval time.Duration
}

func newModel(tickInterval time.Duration) model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return model{
		systems:      GenerateSystems(rng),
		rng:          rng,
		tickInterval: tickInterval,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Tick(m.tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.systems = GenerateSystems(m.rng)
		return m, tea.Tick(m.tickInterval, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.view != viewDashboard {
				m.view = viewDashboard
				return m, nil
			}
			return m, tea.Quit
		case "j", "down":
			if m.view == viewDashboard && m.cursor < len(m.systems)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.view == viewDashboard && m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.view == viewDashboard {
				m.view = viewDetail
			}
		case "esc":
			m.view = viewDashboard
		case "?":
			if m.view == viewHelp {
				m.view = viewDashboard
			} else {
				m.view = viewHelp
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.view {
	case viewDetail:
		return m.renderDetail()
	case viewHelp:
		return m.renderHelp()
	default:
		return m.renderDashboard()
	}
}

func severityStyle(s Severity) lipgloss.Style {
	switch s {
	case SeverityWarning:
		return warningStyle
	case SeverityCritical:
		return criticalStyle
	default:
		return normalStyle
	}
}

func (m model) renderDashboard() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("SmartCity Orchestrator Dashboard"))
	b.WriteString("\n\n")

	header := fmt.Sprintf("  %-14s %-14s %s", "SYSTEM", "STATUS", "KEY METRICS")
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(headerStyle.Render(strings.Repeat("─", 72)))
	b.WriteString("\n")

	for i, sys := range m.systems {
		prefix := "  "
		if i == m.cursor {
			prefix = cursorStyle.Render("▸ ")
		}

		var metricParts []string
		for _, metric := range sys.Metrics {
			metricParts = append(metricParts, metric.Name+": "+metric.Value)
		}

		style := severityStyle(sys.Severity)
		statusText := fmt.Sprintf("● %-10s", sys.Severity.String())
		styledStatus := style.Render(statusText)

		line := fmt.Sprintf("%s%-14s%s  %s", prefix, sys.Name, styledStatus, strings.Join(metricParts, "  "))
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("j/k: navigate  enter: details  ?: help  q: quit"))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderDetail() string {
	if m.cursor >= len(m.systems) {
		return ""
	}
	sys := m.systems[m.cursor]
	style := severityStyle(sys.Severity)

	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("System Detail: %s", sys.Name)))
	b.WriteString("\n\n")

	b.WriteString(detailLabel.Render("Status: "))
	b.WriteString(style.Render("● " + sys.Severity.String()))
	b.WriteString("\n\n")

	b.WriteString(detailLabel.Render("Metrics:"))
	b.WriteString("\n")
	for _, metric := range sys.Metrics {
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			detailLabel.Render(fmt.Sprintf("%-20s", metric.Name+":")),
			detailValue.Render(metric.Value),
		))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("esc/q: back  ?: help"))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderHelp() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Help"))
	b.WriteString("\n\n")

	keys := []struct{ key, desc string }{
		{"j / ↓", "Move cursor down"},
		{"k / ↑", "Move cursor up"},
		{"enter", "View system details"},
		{"esc / q", "Back / Quit"},
		{"?", "Toggle help"},
		{"ctrl+c", "Force quit"},
	}

	for _, k := range keys {
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			headerStyle.Render(fmt.Sprintf("%-10s", k.key)),
			k.desc,
		))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press ? or esc to return"))
	b.WriteString("\n")

	return b.String()
}
