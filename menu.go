package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateFilter() {
	m.filteredChoices = m.filteredChoices[:0] // Clear slice
	
	if m.filterInput == "" {
		// Show all choices when no filter
		for i := range m.choices {
			m.filteredChoices = append(m.filteredChoices, i)
		}
	} else {
		// Filter choices based on input
		filterLower := strings.ToLower(m.filterInput)
		for i, choice := range m.choices {
			choiceLower := strings.ToLower(choice)
			if strings.Contains(choiceLower, filterLower) {
				m.filteredChoices = append(m.filteredChoices, i)
			}
		}
	}
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "/":
			// Start filter mode
			m.filterMode = true
			m.filterInput = ""
			m.updateFilter()
		case "esc":
			if m.filterMode {
				// Exit filter mode
				m.filterMode = false
				m.filterInput = ""
				m.updateFilter()
				m.cursor = 0
			}
		case "backspace":
			if m.filterMode && len(m.filterInput) > 0 {
				m.filterInput = m.filterInput[:len(m.filterInput)-1]
				m.updateFilter()
				m.cursor = 0
			}
		case "up", "k":
			if !m.filterMode {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if !m.filterMode {
				maxIndex := len(m.filteredChoices) - 1
				if m.cursor < maxIndex {
					m.cursor++
				}
			}
		case "enter", " ":
			if len(m.filteredChoices) > 0 && m.cursor < len(m.filteredChoices) {
				// Get the actual choice index from filtered results
				actualChoice := m.filteredChoices[m.cursor]
				
				// Reset filter mode when selecting
				m.filterMode = false
				m.filterInput = ""
				m.updateFilter()
				
				switch actualChoice {
				case 0: // QR Code Generator
					m.state = qrCodeView
					m.qrInput = ""
					m.qrCode = ""
					m.qrCopied = false
					m.qrImagePath = ""
				case 1: // Dice Roller
					m.state = diceRollerView
					m.diceCursor = 0
					m.diceResult = 0
					m.diceType = ""
					m.diceRolling = false
				case 2: // Wheel Spinner
					m.state = wheelSpinnerView
					m.wheelSpinning = false
					m.wheelResult = ""
					m.wheelInputMode = false
					m.wheelInput = ""
				case 3: // RPG Character Creator
					m.state = rpgClassSelectionView
					m.rpgClassCursor = 0
					m.rpgSelectedClass = ""
					m.rpgCharacter = make(map[string]int)
					m.rpgRolling = false
				case 4: // Todo List
					m.state = todoListView
					m.todoInputMode = false
					m.todoInput = ""
					m.todoCursor = 0
					m.todoMessage = ""
				case 5: // Pomodoro Timer
					m.state = pomodoroView
					m.pomodoroMessage = ""
				case 6: // Base64 Encoder/Decoder
					m.state = base64View
					m.base64Input = ""
					m.base64Output = ""
					m.base64Message = ""
					m.base64InputMode = true
				case 7: // Unit Converter
					m.state = unitConverterView
					m = initUnitConverter(m)
				case 8: // System Info
					m.state = systemInfoView
					m.systemInfo = getSystemInfo()
					m.systemInfoMessage = "System information loaded"
					m.systemInfoLastUpdate = time.Now()
				case 9: // Network Info
					m.state = networkInfoView
					m.networkInterfaces = getNetworkInfo()
					m.networkInfoMessage = "Network information loaded"
					m.networkInfoLastUpdate = time.Now()
				case 10: // Quit
					return m, tea.Quit
				}
			}
		default:
			// Handle text input for filter
			if m.filterMode && len(msg.String()) == 1 {
				m.filterInput += msg.String()
				m.updateFilter()
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m model) viewMenu() string {
	// Define styles
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(50).
		AlignHorizontal(lipgloss.Center)

	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1).
		Width(50)

	filterStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E67E22")).
		Padding(1, 2).
		MarginBottom(1).
		Width(50)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#F25D94")).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(50)

	// Build content
	title := titleStyle.Render("🎯 Big Dumb Toolbox")
	
	// Filter input
	var filterDisplay string
	if m.filterMode {
		filterText := m.filterInput + "█" // cursor
		filterLabel := "🔍 Filter: " + filterText
		if len(m.filteredChoices) == 0 {
			filterLabel += " (no matches)"
		} else {
			filterLabel += fmt.Sprintf(" (%d matches)", len(m.filteredChoices))
		}
		filterDisplay = filterStyle.Render(filterLabel)
	}
	
	// Menu items (show filtered results)
	var menuItems []string
	displayChoices := m.filteredChoices
	
	for i, choiceIdx := range displayChoices {
		choice := m.choices[choiceIdx]
		var style lipgloss.Style
		cursor := "  "
		if m.cursor == i {
			cursor = "▶ "
			style = selectedStyle
		} else {
			style = normalStyle
		}
		
		// Highlight matching text in filter mode
		if m.filterMode && m.filterInput != "" {
			choice = m.highlightMatch(choice, m.filterInput)
		}
		
		menuItems = append(menuItems, style.Render(cursor+choice))
	}
	
	if len(menuItems) == 0 {
		menuItems = append(menuItems, normalStyle.Render("  No matches found"))
	}
	
	menu := menuStyle.Render(strings.Join(menuItems, "\n"))
	
	// Help text
	var helpText string
	if m.filterMode {
		helpText = "Type to filter • ESC to clear filter • Enter to select • Ctrl+C to quit"
	} else {
		helpText = "↑/↓ or j/k to navigate • Enter to select • / to filter • q to quit"
	}
	help := helpStyle.Render(helpText)
	
	// Combine content
	var content string
	if filterDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, filterDisplay, menu, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, menu, help)
	}
	
	return containerStyle.Render(content)
}

func (m model) highlightMatch(text, filter string) string {
	if filter == "" {
		return text
	}
	
	// Simple highlighting - make matching text bold
	filterLower := strings.ToLower(filter)
	textLower := strings.ToLower(text)
	
	if strings.Contains(textLower, filterLower) {
		// Find the position of the match
		index := strings.Index(textLower, filterLower)
		if index >= 0 {
			before := text[:index]
			match := text[index : index+len(filter)]
			after := text[index+len(filter):]
			
			highlightStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
			return before + highlightStyle.Render(match) + after
		}
	}
	
	return text
}