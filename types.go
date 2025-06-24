package main

import "time"

type sessionState int

const (
	menuView sessionState = iota
	qrCodeView
	diceRollerView
	wheelSpinnerView
	rpgCharacterView
	rpgClassSelectionView
	todoListView
	pomodoroView
	base64View
	systemInfoView
	networkInfoView
	unitConverterView
)

type ClassStats struct {
	Primary   string
	Secondary string
}

type StartingGear struct {
	Weapons []string
	Armor   []string
	Items   []string
}

type TodoItem struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type SystemInfo struct {
	OS           string
	Arch         string
	NumCPU       int
	GoVersion    string
	Hostname     string
	Username     string
	HomeDir      string
	WorkingDir   string
	TempDir      string
}

type NetworkInterface struct {
	Name         string
	Addresses    []string
	HardwareAddr string
	IsUp         bool
	IsLoopback   bool
}

type model struct {
	state    sessionState
	cursor   int
	choices  []string
	selected map[int]struct{}
	
	// Main menu filter
	filterMode      bool
	filterInput     string
	filteredChoices []int // indices of choices that match filter
	
	qrInput    string
	qrCode     string
	qrCopied   bool
	qrImagePath string
	
	diceCursor   int
	diceTypes    []string
	diceResult   int
	diceType     string
	diceRolling  bool
	diceRollTime time.Time
	
	wheelItems     []string
	wheelInput     string
	wheelSpinning  bool
	wheelSpinTime  time.Time
	wheelResult    string
	wheelSpinIndex int
	wheelInputMode bool
	
	rpgCharacter     map[string]int
	rpgRolling       bool
	rpgRollTime      time.Time
	rpgClasses       []string
	rpgClassCursor   int
	rpgSelectedClass string
	rpgGear          StartingGear
	rpgGold          int
	rpgExportStatus  string
	
	todoItems     []TodoItem
	todoInput     string
	todoInputMode bool
	todoCursor    int
	todoMessage   string
	todoFilter    string
	
	pomodoroRunning   bool
	pomodoroStartTime time.Time
	pomodoroDuration  time.Duration
	pomodoroIsBreak   bool
	pomodoroSession   int
	pomodoroMessage   string
	pomodoroCompleted bool
	
	base64Input     string
	base64Output    string
	base64Mode      string // "encode" or "decode"
	base64Message   string
	base64InputMode bool
	
	systemInfo            SystemInfo
	systemInfoMessage     string
	systemInfoLastUpdate  time.Time
	
	networkInterfaces     []NetworkInterface
	networkInfoMessage    string
	networkInfoLastUpdate time.Time
	
	unitConverterValue      string
	unitConverterFromUnit   string
	unitConverterToUnit     string
	unitConverterResult     string
	unitConverterCategory   string
	unitConverterCategories []string
	unitConverterUnits      map[string][]string
	unitConverterCursor     int
	unitConverterInputMode  string // "value", "from", "to", "category"
	unitConverterMessage    string
	
	width  int
	height int
}