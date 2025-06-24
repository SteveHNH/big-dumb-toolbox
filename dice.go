package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateDiceRoller(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "up", "k":
			if m.diceCursor > 0 {
				m.diceCursor--
			}
		case "down", "j":
			if m.diceCursor < len(m.diceTypes)-1 {
				m.diceCursor++
			}
		case "enter", " ":
			selectedDice := m.diceTypes[m.diceCursor]
			m.diceType = selectedDice
			m.diceRolling = true
			m.diceRollTime = time.Now()
			
			// Roll the dice based on type
			switch selectedDice {
			case "d4":
				m.diceResult = rand.Intn(4) + 1
			case "d6":
				m.diceResult = rand.Intn(6) + 1
			case "d8":
				m.diceResult = rand.Intn(8) + 1
			case "d10":
				m.diceResult = rand.Intn(10) + 1
			case "d12":
				m.diceResult = rand.Intn(12) + 1
			case "d20":
				m.diceResult = rand.Intn(20) + 1
			}
			
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return t
			})
		}
	case time.Time:
		if m.diceRolling && time.Since(m.diceRollTime) > time.Second*2 {
			m.diceRolling = false
		}
		if m.diceRolling {
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return t
			})
		}
	}
	return m, nil
}

func (m model) viewDiceRoller() string {
	// Define styles
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#FF6B6B")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6B6B")).
		Width(50).
		AlignHorizontal(lipgloss.Center)

	diceMenuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6B6B")).
		Padding(1, 2).
		MarginBottom(2).
		Width(50)

	selectedDiceStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#4ECDC4")).
		Padding(0, 1)

	normalDiceStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Padding(0, 1)

	resultStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#45B7D1")).
		Padding(2, 4).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#45B7D1")).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(2)

	rollingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA726")).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(2)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(50)

	// Build content
	title := titleStyle.Render("🎲 Dice Roller")
	
	// Dice selection menu
	var diceOptions []string
	for i, dice := range m.diceTypes {
		var style lipgloss.Style
		cursor := "  "
		if m.diceCursor == i {
			cursor = "🎯 "
			style = selectedDiceStyle
		} else {
			style = normalDiceStyle
		}
		diceOptions = append(diceOptions, style.Render(cursor+dice))
	}
	
	diceMenu := diceMenuStyle.Render("Choose your dice:\n\n" + strings.Join(diceOptions, "\n"))
	
	// Result display with visual flair
	var resultDisplay string
	if m.diceRolling {
		// Rolling animation
		rollingFrames := []string{"⚀", "⚁", "⚂", "⚃", "⚄", "⚅"}
		frame := rollingFrames[int(time.Since(m.diceRollTime)/time.Millisecond/100)%len(rollingFrames)]
		resultDisplay = rollingStyle.Render(fmt.Sprintf("🎲 Rolling %s... %s", m.diceType, frame))
	} else if m.diceResult > 0 {
		// Show result with visual dice
		diceVisual := getDiceVisual(m.diceResult)
		resultDisplay = resultStyle.Render(fmt.Sprintf("🎲 %s Result: %d\n\n%s", m.diceType, m.diceResult, diceVisual))
	}
	
	help := helpStyle.Render("Use ↑/↓ or j/k to navigate • Enter to roll • ESC to go back • Ctrl+C to quit")
	
	var content string
	if resultDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, diceMenu, resultDisplay, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, diceMenu, help)
	}
	
	return containerStyle.Render(content)
}

func getDiceVisual(result int) string {
	switch result {
	case 1:
		return "┌─────────┐\n│         │\n│    ●    │\n│         │\n└─────────┘"
	case 2:
		return "┌─────────┐\n│  ●      │\n│         │\n│      ●  │\n└─────────┘"
	case 3:
		return "┌─────────┐\n│  ●      │\n│    ●    │\n│      ●  │\n└─────────┘"
	case 4:
		return "┌─────────┐\n│  ●   ●  │\n│         │\n│  ●   ●  │\n└─────────┘"
	case 5:
		return "┌─────────┐\n│  ●   ●  │\n│    ●    │\n│  ●   ●  │\n└─────────┘"
	case 6:
		return "┌─────────┐\n│  ●   ●  │\n│  ●   ●  │\n│  ●   ●  │\n└─────────┘"
	default:
		// For dice with more than 6 sides, show a stylized number
		return fmt.Sprintf("┌─────────┐\n│         │\n│   %2d    │\n│         │\n└─────────┘", result)
	}
}