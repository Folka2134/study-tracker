package tui

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/folka2134/study-tracker/study-tracker/internal/notification"
	"github.com/folka2134/study-tracker/study-tracker/internal/storage"
)

type model struct {
	duration  time.Duration
	task      string
	remaining time.Duration
	startTime time.Time
	progress  progress.Model
	quitting  bool
}

func NewModel(duration time.Duration, task string) model {
	return model{
		duration:  duration,
		task:      task,
		remaining: duration,
		startTime: time.Now(),
		progress:  progress.New(progress.WithDefaultGradient()),
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			elapsed := time.Since(m.startTime)
			storage.SaveSession(m.task, elapsed)
			notification.Send("Timer Done!", fmt.Sprintf("Your timer for %s finished with %.f minutes left", m.task, m.duration.Minutes()-elapsed.Abs().Minutes()))
			notification.PlaySound()
			return m, tea.Quit
		}
	case tickMsg:
		if m.remaining <= 0 {
			elapsed := time.Since(m.startTime)
			storage.SaveSession(m.task, elapsed)
			notification.Send("Timer Done!", fmt.Sprintf("Your %.f minutes for %s are up", m.duration.Minutes(), m.task))
			notification.PlaySound()
			return m, tea.Quit
		}
		m.remaining -= time.Second
		cmd := m.progress.SetPercent(float64(m.duration-m.remaining) / float64(m.duration))
		// cmd := m.duration.Minutes()
		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n  Timer for '%s' | Remaining: %s\n\n%s\n\n  (q to quit)\n", m.task, m.remaining.String(), m.progress.View())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func Start(duration time.Duration, task string, w io.Writer) {
	p := tea.NewProgram(NewModel(duration, task), tea.WithOutput(w))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
