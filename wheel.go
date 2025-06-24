package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateWheelSpinner(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.wheelInputMode {
				m.wheelInputMode = false
				m.wheelInput = ""
			} else {
				m.state = menuView
			}
		case "tab":
			m.wheelInputMode = !m.wheelInputMode
			m.wheelInput = ""
		case "enter":
			if m.wheelInputMode {
				// Add new item
				if m.wheelInput != "" {
					m.wheelItems = append(m.wheelItems, m.wheelInput)
					m.wheelInput = ""
					m.wheelInputMode = false
				}
			} else if len(m.wheelItems) > 0 && !m.wheelSpinning {
				// Start spinning
				m.wheelSpinning = true
				m.wheelSpinTime = time.Now()
				m.wheelSpinIndex = 0
				
				// Choose random result
				m.wheelResult = m.wheelItems[rand.Intn(len(m.wheelItems))]
				
				return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
					return t
				})
			}
		case "backspace":
			if m.wheelInputMode {
				if len(m.wheelInput) > 0 {
					m.wheelInput = m.wheelInput[:len(m.wheelInput)-1]
				}
			} else if len(m.wheelItems) > 0 {
				// Remove last item
				m.wheelItems = m.wheelItems[:len(m.wheelItems)-1]
				// Clear result if list becomes empty
				if len(m.wheelItems) == 0 {
					m.wheelResult = ""
				}
			}
		default:
			if m.wheelInputMode && len(msg.String()) == 1 {
				m.wheelInput += msg.String()
			}
		}
	case time.Time:
		if m.wheelSpinning {
			elapsed := time.Since(m.wheelSpinTime)
			if elapsed > time.Second*3 { // Spin for 3 seconds
				m.wheelSpinning = false
			} else {
				// Cycle through items faster early on, slower later
				speed := time.Millisecond * time.Duration(50 + int64(elapsed/time.Millisecond)/20)
				m.wheelSpinIndex = (m.wheelSpinIndex + 1) % len(m.wheelItems)
				return m, tea.Tick(speed, func(t time.Time) tea.Msg {
					return t
				})
			}
		}
	}
	return m, nil
}

func (m model) viewWheelSpinner() string {
	// Define styles
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#9B59B6")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9B59B6")).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	wheelStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#9B59B6")).
		Padding(2, 4).
		MarginBottom(2).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	spinningItemStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#E74C3C")).
		Background(lipgloss.Color("#FFF3CD")).
		Padding(1, 2).
		AlignHorizontal(lipgloss.Center)

	resultStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#27AE60")).
		Padding(2, 4).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#27AE60")).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(2)

	itemsListStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9B59B6")).
		Padding(1, 2).
		MarginBottom(2).
		Width(60)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3498DB")).
		Padding(1, 2).
		MarginBottom(1).
		Width(60)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(60)

	// Build content
	title := titleStyle.Render("🎡 Wheel Spinner")
	
	// Current items list
	var itemsDisplay string
	if len(m.wheelItems) == 0 {
		itemsDisplay = "No items yet!\n\nPress Tab to add your first item"
	} else {
		itemsDisplay = "Current Items:\n"
		for i, item := range m.wheelItems {
			itemsDisplay += fmt.Sprintf("  %d. %s\n", i+1, item)
		}
	}
	itemsList := itemsListStyle.Render(itemsDisplay)
	
	// Wheel display
	var wheelDisplay string
	if m.wheelSpinning {
		// Show spinning animation
		currentItem := m.wheelItems[m.wheelSpinIndex]
		spinSymbols := []string{"🔄", "⭮", "⭯", "🔃"}
		spinSymbol := spinSymbols[int(time.Since(m.wheelSpinTime)/time.Millisecond/125)%len(spinSymbols)]
		wheelDisplay = wheelStyle.Render(fmt.Sprintf("🎡 SPINNING %s\n\n%s", spinSymbol, spinningItemStyle.Render(currentItem)))
	} else if m.wheelResult != "" {
		// Show result
		wheelDisplay = resultStyle.Render(fmt.Sprintf("🎉 WINNER! 🎉\n\n%s", m.wheelResult))
	} else if len(m.wheelItems) == 0 {
		// Show empty state
		wheelDisplay = wheelStyle.Render("🎡 Wheel is Empty\n\nAdd some items first!")
	} else {
		// Show ready to spin
		wheelDisplay = wheelStyle.Render("🎡 Ready to Spin!\n\nPress Enter to start")
	}
	
	// Input area
	var inputDisplay string
	if m.wheelInputMode {
		inputPrompt := "Add new item:"
		inputText := fmt.Sprintf("▶ %s█", m.wheelInput)
		inputDisplay = inputStyle.Render(inputPrompt + "\n" + inputText)
	}
	
	// Help text
	var helpText string
	if m.wheelInputMode {
		helpText = "Type item name • Enter to add • ESC to cancel"
	} else if len(m.wheelItems) == 0 {
		helpText = "Tab to add items • ESC to go back"
	} else {
		helpText = "Enter to spin • Tab to add item • Backspace to remove last • ESC to go back"
	}
	help := helpStyle.Render(helpText)
	
	// Combine all elements
	var content string
	if m.wheelInputMode {
		content = lipgloss.JoinVertical(lipgloss.Center, title, itemsList, inputDisplay, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, itemsList, wheelDisplay, help)
	}
	
	return containerStyle.Render(content)
}