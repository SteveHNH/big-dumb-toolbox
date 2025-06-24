package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/skip2/go-qrcode"
)

func initialModel() model {
	rand.Seed(time.Now().UnixNano())
	choices := []string{"QR Code Generator", "Dice Roller", "Wheel Spinner", "RPG Character Creator", "Todo List", "Pomodoro Timer", "Base64 Encoder/Decoder", "Unit Converter", "System Info", "Network Info", "Quit"}
	m := model{
		state:           menuView,
		choices:         choices,
		selected:        make(map[int]struct{}),
		filteredChoices: make([]int, len(choices)), // Initialize with all choices
		diceTypes:       []string{"d4", "d6", "d8", "d10", "d12", "d20"},
		wheelItems:      []string{}, // Start empty
		rpgCharacter:    make(map[string]int),
		rpgClasses:      []string{"Barbarian", "Rogue", "Wizard", "Paladin", "Warlock", "Cleric", "Monk", "Ranger"},
		todoItems:       loadTodos(),
		todoFilter:      "all",
		pomodoroDuration: 25 * time.Minute, // Default 25-minute work session
		pomodoroSession:  1,
		base64Mode:       "encode",
		base64InputMode:  true,
	}
	
	// Initialize unit converter
	m = initUnitConverter(m)
	
	// Initialize filtered choices with all indices
	for i := range choices {
		m.filteredChoices[i] = i
	}
	
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	
	switch m.state {
	case menuView:
		return m.updateMenu(msg)
	case qrCodeView:
		return m.updateQRCode(msg)
	case diceRollerView:
		return m.updateDiceRoller(msg)
	case wheelSpinnerView:
		return m.updateWheelSpinner(msg)
	case rpgCharacterView:
		return m.updateRPGCharacter(msg)
	case rpgClassSelectionView:
		return m.updateRPGClassSelection(msg)
	case todoListView:
		return m.updateTodoList(msg)
	case pomodoroView:
		return m.updatePomodoro(msg)
	case base64View:
		return m.updateBase64(msg)
	case systemInfoView:
		return m.updateSystemInfo(msg)
	case networkInfoView:
		return m.updateNetworkInfo(msg)
	case unitConverterView:
		return m.updateUnitConverter(msg)
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case menuView:
		return m.viewMenu()
	case qrCodeView:
		return m.viewQRCode()
	case diceRollerView:
		return m.viewDiceRoller()
	case wheelSpinnerView:
		return m.viewWheelSpinner()
	case rpgCharacterView:
		return m.viewRPGCharacter()
	case rpgClassSelectionView:
		return m.viewRPGClassSelection()
	case todoListView:
		return m.viewTodoList()
	case pomodoroView:
		return m.viewPomodoro()
	case base64View:
		return m.viewBase64()
	case systemInfoView:
		return m.viewSystemInfo()
	case networkInfoView:
		return m.viewNetworkInfo()
	case unitConverterView:
		return m.viewUnitConverter()
	}
	return ""
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		testTodoPersistence()
		return
	}
	
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

// QR Code functions
func (m model) updateQRCode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "enter":
			if m.qrInput != "" {
				if qr, err := qrcode.New(m.qrInput, qrcode.Medium); err == nil {
					m.qrCode = qr.ToSmallString(false)
					// Also generate PNG image for clipboard
					tempDir := os.TempDir()
					m.qrImagePath = filepath.Join(tempDir, "qrcode.png")
					qrcode.WriteFile(m.qrInput, qrcode.Medium, 256, m.qrImagePath)
				}
			}
		case "ctrl+d", "ctrl+shift+c":
			if m.qrImagePath != "" {
				if err := copyImageToClipboard(m.qrImagePath); err == nil {
					m.qrCopied = true
				}
			}
		case "backspace":
			if len(m.qrInput) > 0 {
				m.qrInput = m.qrInput[:len(m.qrInput)-1]
				m.qrCode = ""
				m.qrCopied = false
				m.qrImagePath = ""
			}
		default:
			if len(msg.String()) == 1 {
				m.qrInput += msg.String()
				m.qrCopied = false
				m.qrImagePath = ""
			}
		}
	}
	return m, nil
}

func (m model) viewQRCode() string {
	// Define styles
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#00D4AA")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00D4AA")).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	inputBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00D4AA")).
		Padding(1, 2).
		MarginBottom(1).
		Width(60)

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D4AA")).
		Bold(true)

	qrStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF")).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00D4AA")).
		MarginBottom(1).
		AlignHorizontal(lipgloss.Center)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(60)

	copiedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D4AA")).
		Bold(true).
		AlignHorizontal(lipgloss.Center)

	// Build content
	title := titleStyle.Render("📱 QR Code Generator")
	
	inputPrompt := "Enter text to generate QR code:"
	inputDisplay := inputStyle.Render(fmt.Sprintf("▶ %s", m.qrInput))
	inputCursor := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D4AA")).
		Render("█")
	
	inputBox := inputBoxStyle.Render(inputPrompt + "\n" + inputDisplay + inputCursor)
	
	var qrDisplay string
	if m.qrCode != "" {
		qrDisplay = qrStyle.Render(m.qrCode)
	}
	
	var copiedMsg string
	if m.qrCopied {
		copiedMsg = copiedStyle.Render("✓ QR code image copied to clipboard!")
	}
	
	help := helpStyle.Render("Press Enter to generate QR code • Ctrl+D to copy QR image • ESC to go back • Ctrl+C to quit")
	
	var content string
	if qrDisplay != "" && copiedMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputBox, qrDisplay, copiedMsg, help)
	} else if qrDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputBox, qrDisplay, help)
	} else if copiedMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputBox, copiedMsg, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputBox, help)
	}
	
	return containerStyle.Render(content)
}

// Base64 functions
func (m model) updateBase64(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "tab":
			// Toggle between encode and decode modes
			if m.base64Mode == "encode" {
				m.base64Mode = "decode"
			} else {
				m.base64Mode = "encode"
			}
			m.base64Input = ""
			m.base64Output = ""
			m.base64Message = fmt.Sprintf("Switched to %s mode", m.base64Mode)
		case "enter":
			if strings.TrimSpace(m.base64Input) == "" {
				m.base64Message = "Please enter some text to process"
				return m, nil
			}
			
			if m.base64Mode == "encode" {
				// Encode to base64
				m.base64Output = base64.StdEncoding.EncodeToString([]byte(m.base64Input))
				m.base64Message = "✅ Text encoded to Base64"
			} else {
				// Decode from base64
				decoded, err := base64.StdEncoding.DecodeString(m.base64Input)
				if err != nil {
					m.base64Message = "❌ Invalid Base64 input: " + err.Error()
					m.base64Output = ""
				} else {
					m.base64Output = string(decoded)
					m.base64Message = "✅ Base64 decoded to text"
				}
			}
		case "ctrl+shift+c":
			// Copy output to clipboard (placeholder - would need platform-specific implementation)
			if m.base64Output != "" {
				m.base64Message = "📋 Output copied to clipboard (feature not implemented)"
			}
		case "ctrl+r":
			// Reset/clear all
			m.base64Input = ""
			m.base64Output = ""
			m.base64Message = "Cleared"
		case "backspace":
			if len(m.base64Input) > 0 {
				m.base64Input = m.base64Input[:len(m.base64Input)-1]
				// Auto-process on backspace if there's still content
				if len(m.base64Input) > 0 {
					if m.base64Mode == "encode" {
						m.base64Output = base64.StdEncoding.EncodeToString([]byte(m.base64Input))
					} else {
						decoded, err := base64.StdEncoding.DecodeString(m.base64Input)
						if err != nil {
							m.base64Output = ""
						} else {
							m.base64Output = string(decoded)
						}
					}
				} else {
					m.base64Output = ""
				}
			}
		default:
			// Add character to input
			if len(msg.String()) == 1 {
				m.base64Input += msg.String()
				// Auto-process as user types for immediate feedback
				if m.base64Mode == "encode" {
					m.base64Output = base64.StdEncoding.EncodeToString([]byte(m.base64Input))
					m.base64Message = ""
				} else {
					decoded, err := base64.StdEncoding.DecodeString(m.base64Input)
					if err != nil {
						m.base64Output = ""
						m.base64Message = "Invalid Base64..."
					} else {
						m.base64Output = string(decoded)
						m.base64Message = ""
					}
				}
			}
		}
	}
	return m, nil
}

func (m model) viewBase64() string {
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#3498DB")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3498DB")).
		Width(70).
		AlignHorizontal(lipgloss.Center)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3498DB")).
		Padding(1, 2).
		MarginBottom(1).
		Width(70).
		Height(6)

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#27AE60")).
		Padding(1, 2).
		MarginBottom(2).
		Width(70).
		Height(6)

	modeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#E67E22")).
		Padding(0, 2).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E67E22")).
		AlignHorizontal(lipgloss.Center)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(70)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#27AE60")).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("🔐 Base64 Encoder/Decoder")
	
	// Mode indicator
	modeText := strings.ToUpper(m.base64Mode) + " MODE"
	mode := modeStyle.Render(modeText)
	
	// Input area
	inputLabel := "Input (type here):"
	if m.base64Mode == "decode" {
		inputLabel = "Base64 Input (type here):"
	}
	
	inputContent := m.base64Input + "█" // cursor
	if len(m.base64Input) > 200 {
		// Truncate display if too long, but keep full input
		displayInput := m.base64Input[:200] + "..."
		inputContent = displayInput + "█"
	}
	
	inputDisplay := inputStyle.Render(inputLabel + "\n\n" + inputContent)
	
	// Output area
	outputLabel := "Base64 Output:"
	if m.base64Mode == "decode" {
		outputLabel = "Decoded Text Output:"
	}
	
	outputContent := m.base64Output
	if len(outputContent) == 0 {
		outputContent = "(output will appear here)"
	} else if len(outputContent) > 200 {
		// Truncate display if too long
		outputContent = outputContent[:200] + "..."
	}
	
	outputDisplay := outputStyle.Render(outputLabel + "\n\n" + outputContent)
	
	// Help text
	helpText := "Tab to switch modes • Enter to process • Ctrl+R to clear • ESC to go back"
	help := helpStyle.Render(helpText)
	
	// Status message
	var messageDisplay string
	if m.base64Message != "" {
		messageDisplay = messageStyle.Render(m.base64Message)
	}
	
	// Combine all elements
	var content string
	if messageDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, mode, inputDisplay, outputDisplay, messageDisplay, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, mode, inputDisplay, outputDisplay, help)
	}
	
	return containerStyle.Render(content)
}

// All other tool implementations are now in separate files:
// - dice.go: Dice roller functionality
// - wheel.go: Wheel spinner functionality  
// - rpg.go: RPG character creator functionality
// - pomodoro.go: Pomodoro timer functionality
// - todo.go: Todo list functionality
// - system_info.go: System and network info functionality