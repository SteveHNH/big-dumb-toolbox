package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func copyImageToClipboard(imagePath string) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd := exec.Command("osascript", "-e", fmt.Sprintf(`set the clipboard to (read (POSIX file "%s") as JPEG picture)`, imagePath))
		return cmd.Run()
	case "linux":
		// Try xclip first, then wl-clipboard for Wayland
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png", "-i", imagePath)
			return cmd.Run()
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd := exec.Command("wl-copy", "--type", "image/png")
			file, err := os.Open(imagePath)
			if err != nil {
				return err
			}
			defer file.Close()
			cmd.Stdin = file
			return cmd.Run()
		}
		return fmt.Errorf("no suitable clipboard tool found (xclip or wl-copy required)")
	case "windows":
		return fmt.Errorf("windows image clipboard not implemented yet")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func testTodoPersistence() {
	fmt.Println("Testing todo persistence...")
	
	testTodos := []TodoItem{
		{
			ID:        "test1",
			Text:      "Test todo 1",
			Completed: false,
			CreatedAt: time.Now(),
		},
		{
			ID:        "test2", 
			Text:      "Test todo 2",
			Completed: true,
			CreatedAt: time.Now(),
		},
	}
	
	fmt.Printf("Saving %d todos...\n", len(testTodos))
	err := saveTodos(testTodos)
	if err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
		return
	}
	fmt.Println("✅ Todos saved successfully")
	
	fmt.Println("Loading todos...")
	loadedTodos := loadTodos()
	fmt.Printf("✅ Loaded %d todos\n", len(loadedTodos))
	
	for _, todo := range loadedTodos {
		status := "☐"
		if todo.Completed {
			status = "✅"
		}
		fmt.Printf("  %s %s (ID: %s)\n", status, todo.Text, todo.ID)
	}
	
	fmt.Printf("Todo file location: %s\n", getTodoFilePath())
	
	fmt.Println("Cleaning up test file...")
	os.Remove(getTodoFilePath())
	fmt.Println("✅ Test completed")
}