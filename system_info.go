package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // Windows fallback
	}
	homeDir, _ := os.UserHomeDir()
	workingDir, _ := os.Getwd()
	tempDir := os.TempDir()

	return SystemInfo{
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		NumCPU:     runtime.NumCPU(),
		GoVersion:  runtime.Version(),
		Hostname:   hostname,
		Username:   username,
		HomeDir:    homeDir,
		WorkingDir: workingDir,
		TempDir:    tempDir,
	}
}

func getNetworkInfo() []NetworkInterface {
	var interfaces []NetworkInterface
	
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return interfaces
	}
	
	for _, iface := range netInterfaces {
		var addresses []string
		
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		
		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}
		
		interfaces = append(interfaces, NetworkInterface{
			Name:         iface.Name,
			Addresses:    addresses,
			HardwareAddr: iface.HardwareAddr.String(),
			IsUp:         iface.Flags&net.FlagUp != 0,
			IsLoopback:   iface.Flags&net.FlagLoopback != 0,
		})
	}
	
	return interfaces
}

func (m model) updateSystemInfo(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "r":
			// Refresh system info
			m.systemInfo = getSystemInfo()
			m.systemInfoMessage = "System information refreshed"
			m.systemInfoLastUpdate = time.Now()
		}
	}
	return m, nil
}

func (m model) updateNetworkInfo(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuView
		case "r":
			// Refresh network info
			m.networkInterfaces = getNetworkInfo()
			m.networkInfoMessage = "Network information refreshed"
			m.networkInfoLastUpdate = time.Now()
		}
	}
	return m, nil
}

func (m model) viewSystemInfo() string {
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#2ECC71")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2ECC71")).
		Width(70).
		AlignHorizontal(lipgloss.Center)

	infoStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2ECC71")).
		Padding(2, 3).
		MarginBottom(2).
		Width(70)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(70)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2ECC71")).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("💻 System Information")
	
	// System info display
	var infoContent strings.Builder
	infoContent.WriteString("🖥️  SYSTEM DETAILS\n")
	infoContent.WriteString("─────────────────────────────────────────────────────────\n")
	infoContent.WriteString(fmt.Sprintf("Operating System:     %s\n", strings.Title(m.systemInfo.OS)))
	infoContent.WriteString(fmt.Sprintf("Architecture:         %s\n", m.systemInfo.Arch))
	infoContent.WriteString(fmt.Sprintf("CPU Cores:            %d\n", m.systemInfo.NumCPU))
	infoContent.WriteString(fmt.Sprintf("Go Version:           %s\n", m.systemInfo.GoVersion))
	infoContent.WriteString("\n")
	
	infoContent.WriteString("🏠 ENVIRONMENT\n")
	infoContent.WriteString("─────────────────────────────────────────────────────────\n")
	infoContent.WriteString(fmt.Sprintf("Hostname:             %s\n", m.systemInfo.Hostname))
	infoContent.WriteString(fmt.Sprintf("Username:             %s\n", m.systemInfo.Username))
	infoContent.WriteString(fmt.Sprintf("Home Directory:       %s\n", m.systemInfo.HomeDir))
	infoContent.WriteString(fmt.Sprintf("Working Directory:    %s\n", m.systemInfo.WorkingDir))
	infoContent.WriteString(fmt.Sprintf("Temp Directory:       %s\n", m.systemInfo.TempDir))
	infoContent.WriteString("\n")
	
	lastUpdate := m.systemInfoLastUpdate.Format("15:04:05")
	infoContent.WriteString(fmt.Sprintf("Last Updated:         %s", lastUpdate))
	
	infoDisplay := infoStyle.Render(infoContent.String())
	
	// Help text
	helpText := "R to refresh • ESC to go back • Ctrl+C to quit"
	help := helpStyle.Render(helpText)
	
	// Status message
	var messageDisplay string
	if m.systemInfoMessage != "" {
		messageDisplay = messageStyle.Render(m.systemInfoMessage)
	}
	
	// Combine all elements
	var content string
	if messageDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, infoDisplay, messageDisplay, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, infoDisplay, help)
	}
	
	return containerStyle.Render(content)
}

func (m model) viewNetworkInfo() string {
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
		Width(80).
		AlignHorizontal(lipgloss.Center)

	interfaceStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3498DB")).
		Padding(2, 3).
		MarginBottom(2).
		Width(80).
		Height(20)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(80)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3498DB")).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("🌐 Network Information")
	
	// Network interfaces display
	var interfacesContent strings.Builder
	interfacesContent.WriteString("📡 NETWORK INTERFACES\n")
	interfacesContent.WriteString("─────────────────────────────────────────────────────────────────────────\n")
	
	if len(m.networkInterfaces) == 0 {
		interfacesContent.WriteString("No network interfaces found.\n")
	} else {
		for i, iface := range m.networkInterfaces {
			if i > 0 {
				interfacesContent.WriteString("\n")
			}
			
			// Interface name with status
			status := "DOWN"
			statusIcon := "🔴"
			if iface.IsUp {
				status = "UP"
				statusIcon = "🟢"
			}
			
			ifaceType := ""
			if iface.IsLoopback {
				ifaceType = " (Loopback)"
			}
			
			interfacesContent.WriteString(fmt.Sprintf("%s %s %s%s\n", statusIcon, iface.Name, status, ifaceType))
			
			// Hardware address
			if iface.HardwareAddr != "" {
				interfacesContent.WriteString(fmt.Sprintf("    MAC: %s\n", iface.HardwareAddr))
			}
			
			// IP addresses
			if len(iface.Addresses) > 0 {
				interfacesContent.WriteString("    Addresses:\n")
				for _, addr := range iface.Addresses {
					interfacesContent.WriteString(fmt.Sprintf("      • %s\n", addr))
				}
			}
		}
	}
	
	interfacesContent.WriteString("\n")
	lastUpdate := m.networkInfoLastUpdate.Format("15:04:05")
	interfacesContent.WriteString(fmt.Sprintf("Last Updated: %s", lastUpdate))
	
	interfacesDisplay := interfaceStyle.Render(interfacesContent.String())
	
	// Help text
	helpText := "R to refresh • ESC to go back • Ctrl+C to quit"
	help := helpStyle.Render(helpText)
	
	// Status message
	var messageDisplay string
	if m.networkInfoMessage != "" {
		messageDisplay = messageStyle.Render(m.networkInfoMessage)
	}
	
	// Combine all elements
	var content string
	if messageDisplay != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, interfacesDisplay, messageDisplay, help)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, interfacesDisplay, help)
	}
	
	return containerStyle.Render(content)
}