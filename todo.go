package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getTodoFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "todos.json"
	}
	return filepath.Join(homeDir, ".big-dumb-toolbox-todos.json")
}

func loadTodos() []TodoItem {
	filePath := getTodoFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []TodoItem{}
	}
	
	var todos []TodoItem
	if err := json.Unmarshal(data, &todos); err != nil {
		return []TodoItem{}
	}
	
	return todos
}

func saveTodos(todos []TodoItem) error {
	filePath := getTodoFilePath()
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filePath, data, 0644)
}

func generateTodoID() string {
	return fmt.Sprintf("todo_%d", time.Now().UnixNano())
}

func (m model) updateTodoList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.todoInputMode {
				m.todoInputMode = false
				m.todoInput = ""
				m.todoMessage = ""
			} else {
				m.state = menuView
			}
		case "tab":
			m.todoInputMode = !m.todoInputMode
			m.todoInput = ""
			m.todoMessage = ""
		case "enter":
			if m.todoInputMode {
				if strings.TrimSpace(m.todoInput) != "" {
					newTodo := TodoItem{
						ID:        generateTodoID(),
						Text:      strings.TrimSpace(m.todoInput),
						Completed: false,
						CreatedAt: time.Now(),
					}
					m.todoItems = append(m.todoItems, newTodo)
					if err := saveTodos(m.todoItems); err != nil {
						m.todoMessage = "❌ Failed to save todo"
					} else {
						m.todoMessage = "✅ Todo added successfully"
					}
					m.todoInput = ""
					m.todoInputMode = false
				}
			} else if len(m.getFilteredTodos()) > 0 {
				filtered := m.getFilteredTodos()
				if m.todoCursor < len(filtered) {
					targetID := filtered[m.todoCursor].ID
					for i := range m.todoItems {
						if m.todoItems[i].ID == targetID {
							m.todoItems[i].Completed = !m.todoItems[i].Completed
							break
						}
					}
					if err := saveTodos(m.todoItems); err != nil {
						m.todoMessage = "❌ Failed to save changes"
					} else {
						m.todoMessage = "✅ Todo updated"
					}
				}
			}
		case "backspace":
			if m.todoInputMode && len(m.todoInput) > 0 {
				m.todoInput = m.todoInput[:len(m.todoInput)-1]
			}
		case "up":
			if !m.todoInputMode && m.todoCursor > 0 {
				m.todoCursor--
			}
		case "down":
			if !m.todoInputMode {
				filtered := m.getFilteredTodos()
				if m.todoCursor < len(filtered)-1 {
					m.todoCursor++
				}
			}
		case "k":
			if m.todoInputMode {
				m.todoInput += "k"
			} else if m.todoCursor > 0 {
				m.todoCursor--
			}
		case "j":
			if m.todoInputMode {
				m.todoInput += "j"
			} else {
				filtered := m.getFilteredTodos()
				if m.todoCursor < len(filtered)-1 {
					m.todoCursor++
				}
			}
		case "d":
			if m.todoInputMode {
				m.todoInput += "d"
			} else if len(m.getFilteredTodos()) > 0 {
				filtered := m.getFilteredTodos()
				if m.todoCursor < len(filtered) {
					targetID := filtered[m.todoCursor].ID
					for i := len(m.todoItems) - 1; i >= 0; i-- {
						if m.todoItems[i].ID == targetID {
							m.todoItems = append(m.todoItems[:i], m.todoItems[i+1:]...)
							break
						}
					}
					if err := saveTodos(m.todoItems); err != nil {
						m.todoMessage = "❌ Failed to delete todo"
					} else {
						m.todoMessage = "✅ Todo deleted"
					}
					if m.todoCursor >= len(m.getFilteredTodos()) && m.todoCursor > 0 {
						m.todoCursor--
					}
				}
			}
		case "f":
			if m.todoInputMode {
				m.todoInput += "f"
			} else {
				switch m.todoFilter {
				case "all":
					m.todoFilter = "active"
				case "active":
					m.todoFilter = "completed"
				case "completed":
					m.todoFilter = "all"
				}
				m.todoCursor = 0
				m.todoMessage = fmt.Sprintf("Filter: %s", m.todoFilter)
			}
		default:
			if m.todoInputMode && len(msg.String()) == 1 {
				m.todoInput += msg.String()
			}
		}
	}
	return m, nil
}

func (m model) getFilteredTodos() []TodoItem {
	var filtered []TodoItem
	for _, todo := range m.todoItems {
		switch m.todoFilter {
		case "active":
			if !todo.Completed {
				filtered = append(filtered, todo)
			}
		case "completed":
			if todo.Completed {
				filtered = append(filtered, todo)
			}
		default: // "all"
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

func (m model) viewTodoList() string {
	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#FF9500")).
		Padding(1, 2).
		MarginBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF9500")).
		Width(60).
		AlignHorizontal(lipgloss.Center)

	todoListStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF9500")).
		Padding(1, 2).
		MarginBottom(2).
		Width(60).
		Height(15)

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

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(1)

	completedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Strikethrough(true)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#FF9500")).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF9500")).
		Padding(0, 1)

	title := titleStyle.Render("📝 Todo List")
	
	var todoDisplay strings.Builder
	todoDisplay.WriteString(fmt.Sprintf("Filter: %s\n\n", strings.ToUpper(m.todoFilter)))
	
	filtered := m.getFilteredTodos()
	if len(filtered) == 0 {
		todoDisplay.WriteString("No todos found.\n\nPress Tab to add your first todo!")
	} else {
		for i, todo := range filtered {
			cursor := "  "
			style := normalStyle
			if i == m.todoCursor && !m.todoInputMode {
				cursor = "▶ "
				style = selectedStyle
			}
			
			status := "☐"
			text := todo.Text
			if todo.Completed {
				status = "✅"
				text = completedStyle.Render(text)
			}
			
			todoDisplay.WriteString(style.Render(fmt.Sprintf("%s%s %s", cursor, status, text)) + "\n")
		}
	}
	
	todoList := todoListStyle.Render(todoDisplay.String())
	
	var inputDisplay string
	if m.todoInputMode {
		inputPrompt := "Add new todo:"
		inputText := fmt.Sprintf("▶ %s█", m.todoInput)
		inputDisplay = inputStyle.Render(inputPrompt + "\n" + inputText)
	}
	
	var helpText string
	if m.todoInputMode {
		helpText = "Type todo text • Enter to add • ESC to cancel"
	} else {
		helpText = "Enter to toggle • D to delete • F to filter • Tab to add • ↑/↓ to navigate • ESC to go back"
	}
	help := helpStyle.Render(helpText)
	
	var messageDisplay string
	if m.todoMessage != "" {
		messageDisplay = messageStyle.Render(m.todoMessage)
	}
	
	var content string
	if m.todoInputMode {
		if messageDisplay != "" {
			content = lipgloss.JoinVertical(lipgloss.Center, title, todoList, inputDisplay, messageDisplay, help)
		} else {
			content = lipgloss.JoinVertical(lipgloss.Center, title, todoList, inputDisplay, help)
		}
	} else {
		if messageDisplay != "" {
			content = lipgloss.JoinVertical(lipgloss.Center, title, todoList, messageDisplay, help)
		} else {
			content = lipgloss.JoinVertical(lipgloss.Center, title, todoList, help)
		}
	}
	
	return containerStyle.Render(content)
}