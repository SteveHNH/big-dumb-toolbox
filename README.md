# 🎯 Big Dumb Toolbox

A collection of useful utilities packed into a single, beautiful terminal user interface (TUI) application built with Go and Bubble Tea.

## 🚀 Quick Start

```bash
# Build the application (using Makefile)
make build

# Or build directly with Go
go build -o bdt .

# Run the toolbox
./bdt

# Run tests (for todo persistence)
./bdt test
```

## 🏗️ Build Commands

```bash
# Build the application
make build          # Creates 'bdt' executable

# Clean build artifacts
make clean          # Removes 'bdt' executable

# Run tests
make test           # Runs go test ./...

# Install to ~/bin/
make install        # Builds and installs to ~/bin/

# Build everything
make all            # Same as make build
```

## 🛠️ Available Tools (9 Total)

### 1. 📱 QR Code Generator
Generate QR codes from any text input with visual ASCII art display.

**Features:**
- Real-time QR code generation as you type
- ASCII art display in terminal
- PNG file generation for clipboard copying
- Cross-platform clipboard support (macOS, Linux with xclip/wl-copy)

**Controls:**
- Type text to generate QR code
- `Enter` to generate
- `Ctrl+D` to copy QR image to clipboard
- `ESC` to go back

### 2. 🎲 Dice Roller
Virtual dice roller with visual dice faces for tabletop gaming.

**Features:**
- Multiple dice types: d4, d6, d8, d10, d12, d20
- Visual dice face representation for results 1-6
- Rolling animation with random frames
- Clean, game-themed interface

**Controls:**
- `↑/↓` or `j/k` to select dice type
- `Enter` or `Space` to roll
- `ESC` to go back

### 3. 🎡 Wheel Spinner
Customizable decision wheel for random selections.

**Features:**
- Add custom items to the wheel
- Spinning animation with variable speed
- Remove items with backspace
- Persistent wheel state during session

**Controls:**
- `Tab` to add new items
- `Enter` to spin (when items exist)
- `Backspace` to remove last item
- `ESC` to go back or cancel input

### 4. ⚔️ RPG Character Creator
D&D 5E character generator with full equipment and stats.

**Features:**
- 8 character classes: Barbarian, Rogue, Wizard, Paladin, Warlock, Cleric, Monk, Ranger
- Smart stat allocation (highest rolls to primary/secondary stats)
- Class-specific starting equipment and weapons
- Gold generation with realistic distributions
- Export to text and HTML formats
- 4d6 drop lowest stat rolling with reroll 1s

**Controls:**
- Select class, then generate character
- `Enter/R` to reroll stats
- `B` to change class
- `S` to save as text file
- `P` to save as HTML file
- `ESC` to go back

### 5. 📝 Todo List
Persistent task management with filtering and local storage.

**Features:**
- Persistent storage in `~/.big-dumb-toolbox-todos.json`
- Add, complete, and delete todos
- Filter by all/active/completed
- JSON-based data persistence
- Real-time status updates

**Controls:**
- `Tab` to add new todo
- `Enter` to toggle completion
- `D` to delete selected todo
- `F` to cycle filters (all → active → completed → all)
- `↑/↓` to navigate
- `ESC` to go back

### 6. 🍅 Pomodoro Timer
Classic productivity timer following the Pomodoro Technique.

**Features:**
- 25-minute work sessions
- 5-minute short breaks
- 15-minute long breaks (every 4th session)
- Real-time progress bar
- Session counter
- Motivational messages

**Controls:**
- `Enter/Space` to start/pause timer
- `R` to reset current timer
- `S` to skip to next phase
- `ESC` to go back

### 7. 🔐 Base64 Encoder/Decoder
Real-time Base64 encoding and decoding tool.

**Features:**
- Dual mode: Encode text to Base64 or decode Base64 to text
- Real-time processing as you type
- Input validation for decode mode
- Error handling with clear messages
- Support for long text (with display truncation)

**Controls:**
- `Tab` to switch between encode/decode modes
- Type normally for real-time processing
- `Enter` for manual processing
- `Ctrl+R` to clear all
- `Backspace` to edit with live updates
- `ESC` to go back

### 8. 💻 System Information
Display comprehensive system and environment information.

**Features:**
- Operating system and architecture details
- CPU core count and Go version
- System hostname and current user
- Directory paths (home, working, temp)
- Real-time refresh capability
- Timestamp of last update

**System Details Shown:**
- Operating System (Linux, macOS, Windows)
- Architecture (amd64, arm64, etc.)
- CPU Cores available
- Go runtime version
- Hostname and username
- Home, working, and temp directories

**Controls:**
- `R` to refresh system information
- `ESC` to go back
- Auto-loads on entry

### 9. 🌐 Network Information
Display detailed network interface information.

**Features:**
- All network interfaces with status indicators
- IP addresses (IPv4 and IPv6)
- MAC addresses for physical interfaces
- Interface status (UP/DOWN) with visual indicators
- Loopback interface identification
- Real-time refresh capability

**Interface Details Shown:**
- Interface name and status (🟢 UP / 🔴 DOWN)
- MAC/Hardware addresses
- All assigned IP addresses with CIDR notation
- Loopback interface marking
- Interface type indicators

**Controls:**
- `R` to refresh network information
- `ESC` to go back
- Auto-loads on entry

## 🎨 Design Philosophy

**Big Dumb Toolbox** follows these principles:
- **Simple**: Each tool does one thing well
- **Beautiful**: Clean TUI design with consistent styling
- **Useful**: Practical utilities for developers and users
- **Fast**: Instant feedback and responsive interactions
- **Persistent**: Data survives between sessions where appropriate

## 🎯 Navigation

**Main Menu:**
- `↑/↓` or `j/k` - Navigate menu items
- `Enter` or `Space` - Select tool
- `/` - Start filtering tools
- `q` - Quit application

**Filter Mode:**
- Type to search tool names (case-insensitive)
- `ESC` - Clear filter and return to full menu
- `Enter` - Select filtered tool
- `Backspace` - Edit filter text

**Global Controls:**
- `↑/↓` or `j/k` - Navigate menus
- `Enter` or `Space` - Select/activate
- `ESC` - Go back to previous screen
- `Ctrl+C` - Quit application

## 🔍 Quick Filter Feature

The main menu includes a powerful filter system to quickly find tools:

**How to use:**
1. Press `/` from the main menu to enter filter mode
2. Type part of any tool name (e.g., "base", "dice", "todo")
3. See real-time filtered results with match highlighting
4. Use `↑/↓` to navigate filtered results
5. Press `Enter` to select, or `ESC` to clear filter

**Examples:**
- Type `"qr"` → Shows "QR Code Generator"
- Type `"timer"` → Shows "Pomodoro Timer"  
- Type `"info"` → Shows "System Info" and "Network Info"
- Type `"roll"` → Shows "Dice Roller"

**Visual indicators:**
- 🔍 Orange filter box shows current search
- Match count display (e.g., "3 matches")
- **Bold yellow** highlighting of matched text
- "No matches found" when filter has no results

## 📁 File Structure

```
big-dumb-toolbox/
├── main.go              # Main application entry point and core tools
├── types.go             # Data structures and model definitions
├── menu.go              # Main menu and filter functionality
├── todo.go              # Todo list tool implementation
├── system_info.go       # System and network info tools
├── utils.go             # Shared utilities and helper functions
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── README.md           # This documentation
└── ~/.big-dumb-toolbox-todos.json  # Persistent todo storage
```

### Code Organization

The codebase has been refactored into focused, manageable files:

- **`main.go`** - Application entry point, QR Code and Base64 tools
- **`types.go`** - All data structures, constants, and the main model
- **`menu.go`** - Main menu navigation and filtering system
- **`todo.go`** - Complete todo list functionality with persistence
- **`system_info.go`** - System and network information tools
- **`utils.go`** - Shared utilities like clipboard functions and test helpers

This modular structure makes the code easier to:
- **Navigate** - Find specific functionality quickly
- **Maintain** - Update individual tools without affecting others
- **Extend** - Add new tools by creating focused files
- **Test** - Unit test individual components

## 🔧 Technical Details

**Built with:**
- [Go](https://golang.org/) - Core language
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [go-qrcode](https://github.com/skip2/go-qrcode) - QR code generation

**Dependencies:**
```bash
go mod tidy  # Install all dependencies
```

## 🚀 Future Tool Ideas

Potential additions to consider:
- Password Generator
- Hash Generator (MD5, SHA256)
- JSON Formatter/Validator
- URL Encoder/Decoder
- Color Picker/Converter
- Unit Converter
- Time Zone Converter
- ASCII Art Generator
- Morse Code Translator
- Port Scanner

## 🤝 Contributing

Feel free to add new tools! Each tool should:
1. Add a new session state constant
2. Add menu choice and navigation
3. Implement `update{ToolName}` and `view{ToolName}` methods
4. Follow the existing UI patterns and styling
5. Update this README with documentation

## 📄 License

This project is a personal utility toolbox. Use and modify as needed!