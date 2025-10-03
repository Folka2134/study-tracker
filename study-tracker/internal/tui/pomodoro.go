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
	duration       time.Duration
	task           string
	remaining      time.Duration
	startTime      time.Time
	progress       progress.Model
	isBreak        bool
	quitting       bool
	paused         bool
	pausedDuration time.Duration
	pauseStartTime time.Time
}

func NewModel(duration time.Duration, task string, color1 string, color2 string, isBreak bool) model {
	return model{
		duration:       duration,
		task:           task,
		remaining:      duration,
		startTime:      time.Now(),
		progress:       progress.New(progress.WithGradient(color1, color2)),
		isBreak:        isBreak,
		paused:         false,
		pausedDuration: 0,
	}
}

// progress:  progress.New(progress.WithGradient("#3005F2", "#7256F2")),

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
			if m.paused {
				m.pausedDuration += time.Since(m.pauseStartTime)
			}
			if !m.isBreak {
				elapsed := time.Since(m.startTime) - m.pausedDuration
				storage.SaveSession(m.task, elapsed)
				notification.Send("Timer Done!", fmt.Sprintf("Your timer for %s finished with %.f minutes left", m.task, m.duration.Minutes()-elapsed.Abs().Minutes()))
			} else {
				notification.Send("Break ended!", "")
			}
			notification.PlaySound()
			return m, tea.Quit
		}
		if msg.String() == "p" {
			m.paused = !m.paused
			if m.paused {
				m.pauseStartTime = time.Now()
			} else {
				m.pausedDuration += time.Since(m.pauseStartTime)
			}
		}
	case tickMsg:
		if m.paused {
			return m, tickCmd()
		}
		if m.remaining <= 0 {
			if !m.isBreak {
				elapsed := time.Since(m.startTime)
				storage.SaveSession(m.task, elapsed)
			}
			notification.Send("Timer Done!", fmt.Sprintf("Your %.f minutes for %s are up", m.duration.Minutes(), m.task))
			notification.PlaySound()
			return m, tea.Quit
		}
		m.remaining -= time.Second
		cmd := m.progress.SetPercent(float64(m.duration-m.remaining) / float64(m.duration))
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
	paused := ""
	if m.paused {
		paused = "| PAUSED"
	}
	return fmt.Sprintf("\n  Timer for '%s' | Remaining: %s %s\n\n%s\n\n  (q to quit, p to pause)\n", m.task, m.remaining.String(), paused, m.progress.View())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func Start(duration time.Duration, task string, w io.Writer, color1 string, color2 string, isBreak bool) {
	p := tea.NewProgram(NewModel(duration, task, color1, color2, isBreak), tea.WithOutput(w))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
