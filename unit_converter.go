package main

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Conversion factors to base units
var conversions = map[string]map[string]float64{
	"Length": {
		"millimeter": 0.001,
		"centimeter": 0.01,
		"meter":      1.0,
		"kilometer":  1000.0,
		"inch":       0.0254,
		"foot":       0.3048,
		"yard":       0.9144,
		"mile":       1609.344,
	},
	"Weight": {
		"milligram": 0.001,
		"gram":      1.0,
		"kilogram":  1000.0,
		"ounce":     28.3495,
		"pound":     453.592,
		"stone":     6350.29,
		"ton":       1000000.0,
	},
	"Temperature": {
		"celsius":    1.0, // Special handling needed
		"fahrenheit": 1.0, // Special handling needed
		"kelvin":     1.0, // Special handling needed
	},
	"Volume": {
		"milliliter": 0.001,
		"liter":      1.0,
		"gallon":     3.78541,
		"quart":      0.946353,
		"pint":       0.473176,
		"cup":        0.236588,
		"fluid_ounce": 0.0295735,
		"tablespoon": 0.0147868,
		"teaspoon":   0.00492892,
	},
	"Area": {
		"square_millimeter": 0.000001,
		"square_centimeter": 0.0001,
		"square_meter":      1.0,
		"square_kilometer":  1000000.0,
		"square_inch":       0.00064516,
		"square_foot":       0.092903,
		"square_yard":       0.836127,
		"acre":              4046.86,
		"hectare":           10000.0,
	},
	"Speed": {
		"meters_per_second":     1.0,
		"kilometers_per_hour":   0.277778,
		"miles_per_hour":        0.44704,
		"feet_per_second":       0.3048,
		"knots":                 0.514444,
	},
}

func initUnitConverter(m model) model {
	m.unitConverterCategories = []string{"Length", "Weight", "Temperature", "Volume", "Area", "Speed"}
	m.unitConverterUnits = map[string][]string{
		"Length":      {"millimeter", "centimeter", "meter", "kilometer", "inch", "foot", "yard", "mile"},
		"Weight":      {"milligram", "gram", "kilogram", "ounce", "pound", "stone", "ton"},
		"Temperature": {"celsius", "fahrenheit", "kelvin"},
		"Volume":      {"milliliter", "liter", "gallon", "quart", "pint", "cup", "fluid_ounce", "tablespoon", "teaspoon"},
		"Area":        {"square_millimeter", "square_centimeter", "square_meter", "square_kilometer", "square_inch", "square_foot", "square_yard", "acre", "hectare"},
		"Speed":       {"meters_per_second", "kilometers_per_hour", "miles_per_hour", "feet_per_second", "knots"},
	}
	m.unitConverterCategory = "Length"
	m.unitConverterFromUnit = "meter"
	m.unitConverterToUnit = "foot"
	m.unitConverterInputMode = "value"
	m.unitConverterValue = ""
	m.unitConverterResult = ""
	m.unitConverterMessage = ""
	m.unitConverterCursor = 0
	return m
}

func convertUnits(value float64, fromUnit, toUnit, category string) (float64, error) {
	if category == "Temperature" {
		return convertTemperature(value, fromUnit, toUnit)
	}

	categoryConversions, exists := conversions[category]
	if !exists {
		return 0, fmt.Errorf("category not found")
	}

	fromFactor, fromExists := categoryConversions[fromUnit]
	toFactor, toExists := categoryConversions[toUnit]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unit not found")
	}

	// Convert to base unit, then to target unit
	baseValue := value * fromFactor
	result := baseValue / toFactor

	return result, nil
}

func convertTemperature(value float64, fromUnit, toUnit string) (float64, error) {
	var celsius float64

	// Convert from source to celsius
	switch fromUnit {
	case "celsius":
		celsius = value
	case "fahrenheit":
		celsius = (value - 32) * 5 / 9
	case "kelvin":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("unknown temperature unit")
	}

	// Convert from celsius to target
	switch toUnit {
	case "celsius":
		return celsius, nil
	case "fahrenheit":
		return celsius*9/5 + 32, nil
	case "kelvin":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("unknown temperature unit")
	}
}

func (m model) updateUnitConverter(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "tab":
			// Cycle through input modes
			switch m.unitConverterInputMode {
			case "value":
				m.unitConverterInputMode = "category"
			case "category":
				m.unitConverterInputMode = "from"
			case "from":
				m.unitConverterInputMode = "to"
			case "to":
				m.unitConverterInputMode = "value"
			}
			m.unitConverterCursor = 0
		case "enter":
			if m.unitConverterInputMode == "value" && m.unitConverterValue != "" {
				value, err := strconv.ParseFloat(m.unitConverterValue, 64)
				if err != nil {
					m.unitConverterMessage = "Invalid number format"
					return m, nil
				}

				result, err := convertUnits(value, m.unitConverterFromUnit, m.unitConverterToUnit, m.unitConverterCategory)
				if err != nil {
					m.unitConverterMessage = "Conversion error: " + err.Error()
					return m, nil
				}

				m.unitConverterResult = fmt.Sprintf("%.6f", result)
				m.unitConverterMessage = "Conversion completed"
			}
		case "up":
			if m.unitConverterInputMode == "category" {
				if m.unitConverterCursor > 0 {
					m.unitConverterCursor--
				} else {
					m.unitConverterCursor = len(m.unitConverterCategories) - 1
				}
				m.unitConverterCategory = m.unitConverterCategories[m.unitConverterCursor]
				// Reset units when category changes
				units := m.unitConverterUnits[m.unitConverterCategory]
				if len(units) > 1 {
					m.unitConverterFromUnit = units[0]
					m.unitConverterToUnit = units[1]
				}
			} else if m.unitConverterInputMode == "from" {
				units := m.unitConverterUnits[m.unitConverterCategory]
				if m.unitConverterCursor > 0 {
					m.unitConverterCursor--
				} else {
					m.unitConverterCursor = len(units) - 1
				}
				m.unitConverterFromUnit = units[m.unitConverterCursor]
			} else if m.unitConverterInputMode == "to" {
				units := m.unitConverterUnits[m.unitConverterCategory]
				if m.unitConverterCursor > 0 {
					m.unitConverterCursor--
				} else {
					m.unitConverterCursor = len(units) - 1
				}
				m.unitConverterToUnit = units[m.unitConverterCursor]
			}
		case "down":
			if m.unitConverterInputMode == "category" {
				if m.unitConverterCursor < len(m.unitConverterCategories)-1 {
					m.unitConverterCursor++
				} else {
					m.unitConverterCursor = 0
				}
				m.unitConverterCategory = m.unitConverterCategories[m.unitConverterCursor]
				// Reset units when category changes
				units := m.unitConverterUnits[m.unitConverterCategory]
				if len(units) > 1 {
					m.unitConverterFromUnit = units[0]
					m.unitConverterToUnit = units[1]
				}
			} else if m.unitConverterInputMode == "from" {
				units := m.unitConverterUnits[m.unitConverterCategory]
				if m.unitConverterCursor < len(units)-1 {
					m.unitConverterCursor++
				} else {
					m.unitConverterCursor = 0
				}
				m.unitConverterFromUnit = units[m.unitConverterCursor]
			} else if m.unitConverterInputMode == "to" {
				units := m.unitConverterUnits[m.unitConverterCategory]
				if m.unitConverterCursor < len(units)-1 {
					m.unitConverterCursor++
				} else {
					m.unitConverterCursor = 0
				}
				m.unitConverterToUnit = units[m.unitConverterCursor]
			}
		case "backspace":
			if m.unitConverterInputMode == "value" && len(m.unitConverterValue) > 0 {
				m.unitConverterValue = m.unitConverterValue[:len(m.unitConverterValue)-1]
				m.unitConverterResult = ""
				m.unitConverterMessage = ""
			}
		case "ctrl+r":
			m.unitConverterValue = ""
			m.unitConverterResult = ""
			m.unitConverterMessage = ""
		default:
			// Handle number input for value
			if m.unitConverterInputMode == "value" {
				char := msg.String()
				if len(char) == 1 && (char >= "0" && char <= "9" || char == "." || char == "-") {
					m.unitConverterValue += char
					m.unitConverterResult = ""
					m.unitConverterMessage = ""
				}
			}
		}
	}
	return m, nil
}

func (m model) viewUnitConverter() string {
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
		Width(70).
		AlignHorizontal(lipgloss.Center)

	activeStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9B59B6")).
		Background(lipgloss.Color("#9B59B6")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(1, 2).
		MarginBottom(1).
		Width(32)

	inactiveStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#626262")).
		Padding(1, 2).
		MarginBottom(1).
		Width(32)

	resultStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#27AE60")).
		Background(lipgloss.Color("#D5F5E3")).
		Padding(1, 2).
		MarginBottom(2).
		Width(70).
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

	title := titleStyle.Render("🔄 Unit Converter")

	// Value input
	var valueContent string
	if m.unitConverterInputMode == "value" {
		valueContent = fmt.Sprintf("Value: %s█", m.unitConverterValue)
	} else {
		valueContent = fmt.Sprintf("Value: %s", m.unitConverterValue)
	}
	
	var valueBox string
	if m.unitConverterInputMode == "value" {
		valueBox = activeStyle.Render(valueContent)
	} else {
		valueBox = inactiveStyle.Render(valueContent)
	}

	// Category selection
	var categoryContent string
	if m.unitConverterInputMode == "category" {
		categoryContent = fmt.Sprintf("Category: ▶ %s", m.unitConverterCategory)
	} else {
		categoryContent = fmt.Sprintf("Category: %s", m.unitConverterCategory)
	}
	
	var categoryBox string
	if m.unitConverterInputMode == "category" {
		categoryBox = activeStyle.Render(categoryContent)
	} else {
		categoryBox = inactiveStyle.Render(categoryContent)
	}

	// From unit selection
	fromDisplayName := strings.ReplaceAll(m.unitConverterFromUnit, "_", " ")
	var fromContent string
	if m.unitConverterInputMode == "from" {
		fromContent = fmt.Sprintf("From: ▶ %s", fromDisplayName)
	} else {
		fromContent = fmt.Sprintf("From: %s", fromDisplayName)
	}
	
	var fromBox string
	if m.unitConverterInputMode == "from" {
		fromBox = activeStyle.Render(fromContent)
	} else {
		fromBox = inactiveStyle.Render(fromContent)
	}

	// To unit selection
	toDisplayName := strings.ReplaceAll(m.unitConverterToUnit, "_", " ")
	var toContent string
	if m.unitConverterInputMode == "to" {
		toContent = fmt.Sprintf("To: ▶ %s", toDisplayName)
	} else {
		toContent = fmt.Sprintf("To: %s", toDisplayName)
	}
	
	var toBox string
	if m.unitConverterInputMode == "to" {
		toBox = activeStyle.Render(toContent)
	} else {
		toBox = inactiveStyle.Render(toContent)
	}

	// Layout inputs in 2x2 grid
	topRow := lipgloss.JoinHorizontal(lipgloss.Left, valueBox, "  ", categoryBox)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Left, fromBox, "  ", toBox)
	inputGrid := lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)

	// Result display
	var resultBox string
	if m.unitConverterResult != "" {
		resultContent := fmt.Sprintf("Result: %s %s", m.unitConverterResult, strings.ReplaceAll(m.unitConverterToUnit, "_", " "))
		resultBox = resultStyle.Render(resultContent)
	}

	// Help text
	helpText := "Tab to switch fields • ↑/↓ to change selection • Enter to convert • Ctrl+R to clear • ESC to go back"
	help := helpStyle.Render(helpText)

	// Status message
	var messageDisplay string
	if m.unitConverterMessage != "" {
		messageDisplay = messageStyle.Render(m.unitConverterMessage)
	}

	// Combine all elements
	var content string
	if resultBox != "" && messageDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputGrid, resultBox, messageDisplay, help)
	} else if resultBox != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputGrid, resultBox, help)
	} else if messageDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputGrid, messageDisplay, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, inputGrid, help)
	}

	return containerStyle.Render(content)
}