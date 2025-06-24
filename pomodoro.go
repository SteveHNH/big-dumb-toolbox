package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updatePomodoro(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "enter", " ":
			if !m.pomodoroRunning {
				// Start timer
				m.pomodoroRunning = true
				m.pomodoroStartTime = time.Now()
				m.pomodoroCompleted = false
				if m.pomodoroIsBreak {
					m.pomodoroMessage = "Break time! Relax and recharge 😌"
				} else {
					m.pomodoroMessage = "Focus time! Stay productive 🎯"
				}
				return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
					return t
				})
			} else {
				// Stop timer
				m.pomodoroRunning = false
				m.pomodoroMessage = "Timer stopped"
			}
		case "r":
			// Reset timer
			m.pomodoroRunning = false
			m.pomodoroCompleted = false
			m.pomodoroMessage = "Timer reset"
		case "s":
			// Skip to next phase
			if m.pomodoroRunning {
				m.pomodoroRunning = false
				m.pomodoroCompleted = true
				if m.pomodoroIsBreak {
					m.pomodoroIsBreak = false
					m.pomodoroDuration = 25 * time.Minute
					m.pomodoroSession++
					m.pomodoroMessage = "Break skipped! Ready for next work session"
				} else {
					m.pomodoroIsBreak = true
					if m.pomodoroSession%4 == 0 {
						m.pomodoroDuration = 15 * time.Minute // Long break
						m.pomodoroMessage = "Work session complete! Time for a long break"
					} else {
						m.pomodoroDuration = 5 * time.Minute // Short break
						m.pomodoroMessage = "Work session complete! Time for a short break"
					}
				}
			}
		}
	case time.Time:
		if m.pomodoroRunning {
			elapsed := time.Since(m.pomodoroStartTime)
			if elapsed >= m.pomodoroDuration {
				// Timer completed
				m.pomodoroRunning = false
				m.pomodoroCompleted = true
				if m.pomodoroIsBreak {
					m.pomodoroIsBreak = false
					m.pomodoroDuration = 25 * time.Minute
					m.pomodoroSession++
					m.pomodoroMessage = "Break complete! Ready for next work session 💪"
				} else {
					m.pomodoroIsBreak = true
					if m.pomodoroSession%4 == 0 {
						m.pomodoroDuration = 15 * time.Minute // Long break every 4 sessions
						m.pomodoroMessage = "Work session complete! Time for a long break ☕"
					} else {
						m.pomodoroDuration = 5 * time.Minute // Short break
						m.pomodoroMessage = "Work session complete! Time for a short break 🌱"
					}
				}
			} else {
				return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
					return t
				})
			}
		}
	}
	return m, nil
}

func (m model) viewPomodoro() string {
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#E74C3C")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E74C3C")).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	timerStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#E74C3C")).
		Padding(3, 6).
		MarginBottom(2).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	progressStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E74C3C")).
		Padding(1, 2).
		MarginBottom(2).
		Width(60)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(60)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#27AE60")).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("🍅 Pomodoro Timer")
	
	// Timer display
	var timeDisplay string
	var remaining time.Duration
	
	if m.pomodoroRunning {
		elapsed := time.Since(m.pomodoroStartTime)
		remaining = m.pomodoroDuration - elapsed
		if remaining < 0 {
			remaining = 0
		}
	} else {
		remaining = m.pomodoroDuration
	}
	
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60
	
	timerText := fmt.Sprintf("%02d:%02d", minutes, seconds)
	phaseText := "Work Session"
	if m.pomodoroIsBreak {
		if m.pomodoroSession%4 == 0 && m.pomodoroSession > 0 {
			phaseText = "Long Break"
		} else {
			phaseText = "Short Break"
		}
	}
	
	sessionText := fmt.Sprintf("Session %d", m.pomodoroSession)
	
	var statusEmoji string
	if m.pomodoroRunning {
		statusEmoji = "⏰"
	} else if m.pomodoroCompleted {
		statusEmoji = "✅"
	} else {
		statusEmoji = "⏸️"
	}
	
	timeDisplay = timerStyle.Render(fmt.Sprintf("%s\n\n%s\n%s\n\n%s", statusEmoji, timerText, phaseText, sessionText))
	
	// Progress bar
	var progressDisplay string
	if m.pomodoroRunning {
		elapsed := time.Since(m.pomodoroStartTime)
		progress := float64(elapsed) / float64(m.pomodoroDuration)
		if progress > 1 {
			progress = 1
		}
		
		barWidth := 50
		filled := int(progress * float64(barWidth))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		
		progressDisplay = progressStyle.Render(fmt.Sprintf("Progress:\n[%s] %.1f%%", bar, progress*100))
	}
	
	// Help text
	var helpText string
	if m.pomodoroRunning {
		helpText = "Enter to pause • S to skip • R to reset • ESC to go back"
	} else if m.pomodoroCompleted {
		helpText = "Enter to start next phase • R to reset • ESC to go back"
	} else {
		helpText = "Enter to start • R to reset • ESC to go back"
	}
	help := helpStyle.Render(helpText)
	
	// Status message
	var messageDisplay string
	if m.pomodoroMessage != "" {
		messageDisplay = messageStyle.Render(m.pomodoroMessage)
	}
	
	// Combine all elements
	var content string
	if progressDisplay != "" && messageDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, timeDisplay, progressDisplay, messageDisplay, help)
	} else if progressDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, timeDisplay, progressDisplay, help)
	} else if messageDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, timeDisplay, messageDisplay, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, timeDisplay, help)
	}
	
	return containerStyle.Render(content)
}