package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateRPGClassSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "up", "k":
			if m.rpgClassCursor > 0 {
				m.rpgClassCursor--
			}
		case "down", "j":
			if m.rpgClassCursor < len(m.rpgClasses)-1 {
				m.rpgClassCursor++
			}
		case "enter", " ":
			m.rpgSelectedClass = m.rpgClasses[m.rpgClassCursor]
			m.state = rpgCharacterView
			// Generate new character with the selected class
			m.rpgCharacter = generateCharacter(m.rpgSelectedClass)
			m.rpgGear = getStartingGear(m.rpgSelectedClass)
			m.rpgGold = generateGold()
			m.rpgExportStatus = ""
		}
	}
	return m, nil
}

func getClassStats(className string) ClassStats {
	// For now, using placeholder values - you can customize these later
	classMap := map[string]ClassStats{
		"Barbarian": {Primary: "Strength", Secondary: "Constitution"},
		"Rogue":     {Primary: "Dexterity", Secondary: "Intelligence"},
		"Wizard":    {Primary: "Intelligence", Secondary: "Wisdom"},
		"Paladin":   {Primary: "Strength", Secondary: "Charisma"},
		"Warlock":   {Primary: "Charisma", Secondary: "Constitution"},
		"Cleric":    {Primary: "Wisdom", Secondary: "Constitution"},
		"Monk":      {Primary: "Dexterity", Secondary: "Wisdom"},
		"Ranger":    {Primary: "Dexterity", Secondary: "Wisdom"},
	}
	return classMap[className]
}

func getStartingGear(className string) StartingGear {
	gearMap := map[string]StartingGear{
		"Barbarian": {
			Weapons: []string{"Greataxe", "Handaxe (2)", "Javelin (4)"},
			Armor:   []string{"Leather armor", "Shield"},
			Items:   []string{"Explorer's pack", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
		"Rogue": {
			Weapons: []string{"Rapier", "Shortbow", "Arrows (20)", "Dagger (2)"},
			Armor:   []string{"Leather armor"},
			Items:   []string{"Burglar's pack", "Thieves' tools", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
		"Wizard": {
			Weapons: []string{"Quarterstaff", "Dagger (2)"},
			Armor:   []string{},
			Items:   []string{"Scholar's pack", "Spellbook", "Component pouch", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
		"Paladin": {
			Weapons: []string{"Longsword", "Shield", "Javelin (5)"},
			Armor:   []string{"Chain mail"},
			Items:   []string{"Explorer's pack", "Holy symbol", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
		"Warlock": {
			Weapons: []string{"Light crossbow", "Crossbow bolts (20)", "Dagger (2)"},
			Armor:   []string{"Leather armor"},
			Items:   []string{"Scholar's pack", "Arcane focus", "Simple weapon", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
		"Cleric": {
			Weapons: []string{"Mace", "Light crossbow", "Crossbow bolts (20)", "Shield"},
			Armor:   []string{"Scale mail"},
			Items:   []string{"Explorer's pack", "Holy symbol", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
		"Monk": {
			Weapons: []string{"Shortsword", "Dart (10)"},
			Armor:   []string{},
			Items:   []string{"Dungeoneer's pack", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
		"Ranger": {
			Weapons: []string{"Longbow", "Arrows (20)", "Shortsword (2)", "Handaxe (2)"},
			Armor:   []string{"Leather armor"},
			Items:   []string{"Explorer's pack", "Bedroll", "Mess kit", "Tinderbox", "Torches (10)", "Rations (10 days)", "Waterskin", "Hemp rope (50 feet)"},
		},
	}
	return gearMap[className]
}

func generateGold() int {
	roll := rand.Intn(100)
	if roll < 10 { // 10% chance
		return 60
	} else if roll < 40 { // 30% chance (10 + 30)
		return 40
	} else { // 60% chance (remaining)
		return 20
	}
}

func exportCharacterText(className string, character map[string]int, gear StartingGear, gold int) (string, error) {
	var content strings.Builder
	
	content.WriteString("===============================\n")
	content.WriteString("       D&D 5E CHARACTER SHEET\n")
	content.WriteString("===============================\n\n")
	
	if className != "" {
		content.WriteString(fmt.Sprintf("Class: %s\n\n", className))
	}
	
	content.WriteString("ABILITY SCORES:\n")
	content.WriteString("---------------\n")
	
	stats := []string{"Strength", "Constitution", "Intelligence", "Wisdom", "Charisma", "Dexterity"}
	classStats := getClassStats(className)
	
	for _, stat := range stats {
		if value, exists := character[stat]; exists {
			marker := ""
			if stat == classStats.Primary {
				marker = " (Primary)"
			} else if stat == classStats.Secondary {
				marker = " (Secondary)"
			}
			content.WriteString(fmt.Sprintf("%-13s: %2d%s\n", stat, value, marker))
		}
	}
	
	content.WriteString(fmt.Sprintf("\nGOLD: %d gp\n\n", gold))
	
	if len(gear.Weapons) > 0 {
		content.WriteString("WEAPONS:\n")
		content.WriteString("--------\n")
		for _, weapon := range gear.Weapons {
			content.WriteString(fmt.Sprintf("• %s\n", weapon))
		}
		content.WriteString("\n")
	}
	
	if len(gear.Armor) > 0 {
		content.WriteString("ARMOR:\n")
		content.WriteString("------\n")
		for _, armor := range gear.Armor {
			content.WriteString(fmt.Sprintf("• %s\n", armor))
		}
		content.WriteString("\n")
	}
	
	if len(gear.Items) > 0 {
		content.WriteString("EQUIPMENT:\n")
		content.WriteString("----------\n")
		for _, item := range gear.Items {
			content.WriteString(fmt.Sprintf("• %s\n", item))
		}
		content.WriteString("\n")
	}
	
	content.WriteString("Generated by Big Dumb Toolbox RPG Character Creator\n")
	
	// Generate unique filename with class and timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_Character_%s.txt", className, timestamp)
	
	err := os.WriteFile(filename, []byte(content.String()), 0644)
	return filename, err
}

func exportCharacterPDF(className string, character map[string]int, gear StartingGear, gold int) (string, error) {
	// For PDF export, we'll create an HTML file and suggest using a browser to print to PDF
	// This is a simple approach that works across all platforms
	var content strings.Builder
	
	content.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>D&D 5E Character Sheet</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; line-height: 1.6; }
        .header { text-align: center; border-bottom: 2px solid #333; padding-bottom: 10px; margin-bottom: 20px; }
        .section { margin-bottom: 20px; }
        .section h3 { color: #333; border-bottom: 1px solid #ccc; padding-bottom: 5px; }
        .stats { display: grid; grid-template-columns: repeat(2, 1fr); gap: 10px; }
        .stat { padding: 8px; background: #f5f5f5; border-radius: 4px; }
        .primary { background: #ffd700; font-weight: bold; }
        .secondary { background: #c0c0c0; font-weight: bold; }
        ul { padding-left: 20px; }
        .footer { margin-top: 30px; text-align: center; color: #666; font-size: 12px; }
        @media print { body { margin: 0; } }
    </style>
</head>
<body>
    <div class="header">
        <h1>D&D 5E CHARACTER SHEET</h1>`)
	
	if className != "" {
		content.WriteString(fmt.Sprintf(`        <h2>%s</h2>`, className))
	}
	
	content.WriteString(`    </div>
    
    <div class="section">
        <h3>Ability Scores</h3>
        <div class="stats">`)
	
	stats := []string{"Strength", "Constitution", "Intelligence", "Wisdom", "Charisma", "Dexterity"}
	classStats := getClassStats(className)
	
	for _, stat := range stats {
		if value, exists := character[stat]; exists {
			class := "stat"
			if stat == classStats.Primary {
				class = "stat primary"
			} else if stat == classStats.Secondary {
				class = "stat secondary"
			}
			content.WriteString(fmt.Sprintf(`            <div class="%s">%s: %d</div>`, class, stat, value))
		}
	}
	
	content.WriteString(`        </div>
    </div>`)
	
	content.WriteString(fmt.Sprintf(`    
    <div class="section">
        <h3>Gold</h3>
        <p><strong>%d gp</strong></p>
    </div>`, gold))
	
	if len(gear.Weapons) > 0 {
		content.WriteString(`    
    <div class="section">
        <h3>Weapons</h3>
        <ul>`)
		for _, weapon := range gear.Weapons {
			content.WriteString(fmt.Sprintf(`            <li>%s</li>`, weapon))
		}
		content.WriteString(`        </ul>
    </div>`)
	}
	
	if len(gear.Armor) > 0 {
		content.WriteString(`    
    <div class="section">
        <h3>Armor</h3>
        <ul>`)
		for _, armor := range gear.Armor {
			content.WriteString(fmt.Sprintf(`            <li>%s</li>`, armor))
		}
		content.WriteString(`        </ul>
    </div>`)
	}
	
	if len(gear.Items) > 0 {
		content.WriteString(`    
    <div class="section">
        <h3>Equipment</h3>
        <ul>`)
		for _, item := range gear.Items {
			content.WriteString(fmt.Sprintf(`            <li>%s</li>`, item))
		}
		content.WriteString(`        </ul>
    </div>`)
	}
	
	content.WriteString(`    
    <div class="footer">
        <p>Generated by Big Dumb Toolbox RPG Character Creator</p>
        <p><em>To convert to PDF: Open this file in your browser and use Print → Save as PDF</em></p>
    </div>
</body>
</html>`)
	
	// Generate unique filename with class and timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_Character_%s.html", className, timestamp)
	
	err := os.WriteFile(filename, []byte(content.String()), 0644)
	return filename, err
}

func rollStat() int {
	var rolls []int
	for i := 0; i < 4; i++ {
		roll := rand.Intn(6) + 1
		if roll == 1 {
			roll = rand.Intn(6) + 1 // Reroll ones
		}
		rolls = append(rolls, roll)
	}
	
	// Sort in descending order and take highest 3
	for i := 0; i < len(rolls); i++ {
		for j := i + 1; j < len(rolls); j++ {
			if rolls[i] < rolls[j] {
				rolls[i], rolls[j] = rolls[j], rolls[i]
			}
		}
	}
	
	return rolls[0] + rolls[1] + rolls[2]
}

func generateCharacter(className string) map[string]int {
	stats := []string{"Strength", "Constitution", "Intelligence", "Wisdom", "Charisma", "Dexterity"}
	
	// Roll all stats
	var rolls []int
	for range stats {
		rolls = append(rolls, rollStat())
	}
	
	// Sort rolls in descending order to get highest values first
	for i := 0; i < len(rolls); i++ {
		for j := i + 1; j < len(rolls); j++ {
			if rolls[i] < rolls[j] {
				rolls[i], rolls[j] = rolls[j], rolls[i]
			}
		}
	}
	
	// Initialize character map
	character := make(map[string]int)
	
	// If class is selected, prioritize primary and secondary stats
	if className != "" {
		classStats := getClassStats(className)
		
		// Assign highest roll to primary stat
		character[classStats.Primary] = rolls[0]
		
		// Assign second highest to secondary stat
		character[classStats.Secondary] = rolls[1]
		
		// Assign remaining rolls to other stats
		rollIndex := 2
		for _, stat := range stats {
			if stat != classStats.Primary && stat != classStats.Secondary {
				character[stat] = rolls[rollIndex]
				rollIndex++
			}
		}
	} else {
		// Random assignment if no class selected
		for i, stat := range stats {
			character[stat] = rolls[i]
		}
	}
	
	return character
}

func (m model) updateRPGCharacter(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "enter", " ":
			if !m.rpgRolling {
				m.rpgRolling = true
				m.rpgRollTime = time.Now()
				m.rpgCharacter = generateCharacter(m.rpgSelectedClass)
				m.rpgGear = getStartingGear(m.rpgSelectedClass)
				m.rpgGold = generateGold()
				m.rpgExportStatus = ""
				
				return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
					return t
				})
			}
		case "r":
			if !m.rpgRolling {
				m.rpgCharacter = generateCharacter(m.rpgSelectedClass)
				m.rpgGear = getStartingGear(m.rpgSelectedClass)
				m.rpgGold = generateGold()
				m.rpgExportStatus = ""
			}
		case "b":
			m.state = rpgClassSelectionView
		case "s":
			if len(m.rpgCharacter) > 0 {
				filename, err := exportCharacterText(m.rpgSelectedClass, m.rpgCharacter, m.rpgGear, m.rpgGold)
				if err != nil {
					m.rpgExportStatus = "❌ Export failed: " + err.Error()
				} else {
					m.rpgExportStatus = "✅ Character saved to " + filename
				}
			}
		case "p":
			if len(m.rpgCharacter) > 0 {
				filename, err := exportCharacterPDF(m.rpgSelectedClass, m.rpgCharacter, m.rpgGear, m.rpgGold)
				if err != nil {
					m.rpgExportStatus = "❌ PDF export failed: " + err.Error()
				} else {
					m.rpgExportStatus = "✅ Character saved to " + filename + " (open in browser to print as PDF)"
				}
			}
		}
	case time.Time:
		if m.rpgRolling && time.Since(m.rpgRollTime) > time.Second*2 {
			m.rpgRolling = false
		}
		if m.rpgRolling {
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return t
			})
		}
	}
	return m, nil
}

func (m model) viewRPGCharacter() string {
	// Define styles
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#8B5CF6")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	characterStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Padding(2, 4).
		MarginBottom(2).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	statStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#10B981")).
		MarginBottom(1)

	rollingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B")).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(2)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(60)

	// Build content
	title := titleStyle.Render("⚔️  RPG Character Creator")
	
	// Character display
	var characterDisplay string
	if m.rpgRolling {
		// Rolling animation
		rollingFrames := []string{"🎲", "🎯", "⚡", "🔥", "✨", "🌟"}
		frame := rollingFrames[int(time.Since(m.rpgRollTime)/time.Millisecond/200)%len(rollingFrames)]
		characterDisplay = rollingStyle.Render(fmt.Sprintf("Rolling character stats... %s", frame))
	} else if len(m.rpgCharacter) > 0 {
		// Show character stats
		stats := []string{
			"Strength", "Constitution", "Intelligence", 
			"Wisdom", "Charisma", "Dexterity",
		}
		
		var statLines []string
		for _, stat := range stats {
			if value, exists := m.rpgCharacter[stat]; exists {
				// Highlight primary and secondary stats
				if m.rpgSelectedClass != "" {
					classStats := getClassStats(m.rpgSelectedClass)
					if stat == classStats.Primary {
						statLines = append(statLines, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700")).Render(fmt.Sprintf("%-13s: %2d", stat, value)))
					} else if stat == classStats.Secondary {
						statLines = append(statLines, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C0C0C0")).Render(fmt.Sprintf("%-13s: %2d", stat, value)))
					} else {
						statLines = append(statLines, statStyle.Render(fmt.Sprintf("%-13s: %2d", stat, value)))
					}
				} else {
					statLines = append(statLines, statStyle.Render(fmt.Sprintf("%-13s: %2d", stat, value)))
				}
			}
		}
		
		classTitle := "🧙 Your Character Stats:"
		if m.rpgSelectedClass != "" {
			classTitle = fmt.Sprintf("🧙 %s Character Stats:", m.rpgSelectedClass)
		}
		
		// Add gear and gold information
		var gearDisplay strings.Builder
		gearDisplay.WriteString(classTitle + "\n\n")
		gearDisplay.WriteString(strings.Join(statLines, "\n"))
		
		if m.rpgGold > 0 {
			gearDisplay.WriteString(fmt.Sprintf("\n\n💰 Gold: %d gp", m.rpgGold))
		}
		
		if len(m.rpgGear.Weapons) > 0 {
			gearDisplay.WriteString("\n\n⚔️  Weapons:")
			for _, weapon := range m.rpgGear.Weapons {
				gearDisplay.WriteString(fmt.Sprintf("\n  • %s", weapon))
			}
		}
		
		if len(m.rpgGear.Armor) > 0 {
			gearDisplay.WriteString("\n\n🛡️  Armor:")
			for _, armor := range m.rpgGear.Armor {
				gearDisplay.WriteString(fmt.Sprintf("\n  • %s", armor))
			}
		}
		
		if len(m.rpgGear.Items) > 0 {
			gearDisplay.WriteString("\n\n🎒 Equipment:")
			for _, item := range m.rpgGear.Items {
				gearDisplay.WriteString(fmt.Sprintf("\n  • %s", item))
			}
		}
		
		characterDisplay = characterStyle.Render(gearDisplay.String())
	} else {
		// Show initial state
		characterDisplay = characterStyle.Render("🧙 Ready to create your character!\n\nPress Enter to roll stats")
	}
	
	// Help text
	var helpText string
	if m.rpgRolling {
		helpText = "Rolling stats using 4d6, reroll 1s, take highest 3..."
	} else if len(m.rpgCharacter) > 0 {
		helpText = "Enter/R to reroll • B to change class • S to save as text • P to save as HTML • ESC to go back"
	} else {
		helpText = "Enter to roll character • B to change class • ESC to go back"
	}
	help := helpStyle.Render(helpText)
	
	// Export status message
	var content string
	if m.rpgExportStatus != "" {
		exportStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			MarginBottom(1)
		exportMsg := exportStyle.Render(m.rpgExportStatus)
		content = lipgloss.JoinVertical(lipgloss.Center, title, characterDisplay, exportMsg, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, characterDisplay, help)
	}
	
	return containerStyle.Render(content)
}

func (m model) viewRPGClassSelection() string {
	// Define styles
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#8B5CF6")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Width(50).
		AlignHorizontal(lipgloss.Center)

	classMenuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Padding(1, 2).
		MarginBottom(2).
		Width(50)

	selectedClassStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#10B981")).
		Padding(0, 1)

	normalClassStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B5CF6")).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(50)

	// Build content
	title := titleStyle.Render("⚔️  Choose Your Class")
	
	// Class selection menu
	var classOptions []string
	for i, class := range m.rpgClasses {
		var style lipgloss.Style
		cursor := "  "
		if m.rpgClassCursor == i {
			cursor = "▶ "
			style = selectedClassStyle
		} else {
			style = normalClassStyle
		}
		
		classOptions = append(classOptions, style.Render(cursor+class))
	}
	
	classMenu := classMenuStyle.Render("Select your character class:\n\n" + strings.Join(classOptions, "\n"))
	
	help := helpStyle.Render("Use ↑/↓ or j/k to navigate • Enter to select • ESC to go back • Ctrl+C to quit")
	
	content := lipgloss.JoinVertical(lipgloss.Center, title, classMenu, help)
	
	return containerStyle.Render(content)
}