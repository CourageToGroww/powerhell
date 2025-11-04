package app

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/couragetogroww/powerhell/pkg/auth"
	"github.com/couragetogroww/powerhell/pkg/menus"
	mainmenu "github.com/couragetogroww/powerhell/pkg/menus/mainmenu"
	"github.com/couragetogroww/powerhell/pkg/modules"
	"github.com/couragetogroww/powerhell/pkg/views"
)

// Model represents the application state
type Model struct {
	Flames                   []string
	AppState                 int
	PreviousState            int
	MenuCursor               int
	TerminalWidth            int
	TerminalHeight           int
	NameInput                textinput.Model
	EmailInput               textinput.Model
	GeneratedAccountNumber   string
	FocusedField             int
	Quit                     bool
	ShowHelp                 bool

	// Menu manager for modular menu system
	MenuManager              *menus.MenuManager

	// Fields for stateModuleExplorer
	ModuleExplorerSidebarCursor  int
	ModuleExplorerSidebarOptions []string
	ModuleExplorerContent        string
	
	// New UI components
	Dashboard    *views.DashboardView
	LessonView   *views.LessonView
	SignInView   *views.SignInView
	CurrentModule *modules.Module
	
	// Animation states
	AnimationFrame int
	ShowAccountWarning bool
	
	// Account store
	AccountStore *auth.Store
	CurrentAccount *auth.Account
	SessionID int64
}

// App states
const (
	StateIntro = mainmenu.StateIntro
	StateAccountCreation = mainmenu.StateAccountCreation
	StateMainMenu = mainmenu.StateMainMenu
	StateLearnMenu = mainmenu.StateLearnMenu
	StateStudio = mainmenu.StateStudio
	StateSettings = mainmenu.StateSettings
	StateAuthMenu = mainmenu.StateAuthMenu
	StateSignInPlaceholder = mainmenu.StateSignInPlaceholder
	StateModuleExplorer = mainmenu.StateModuleExplorer
	StateDashboard = 100
	StateLesson = 101
	StateSignIn = 102
)

const (
	FocusName = iota
	FocusEmail
	FocusDisplayInfo
)

var asciiArt = `
                                                                               
                                                                               
                                                                               
                                                                               
                                        ..+-                                    
                                       .+#%-                                    
                                      .+@@%=        :                           
                                       *@%%%=                                   
                                    . :#%%%%%-       :*:.                       
                             .-      =%%%%%%%+   .   -%#-.                      
                             *#. :  :%%%#%%%%=  .-. :*%%#:                      
                             *%-*- :#%%%###%%+   +*-+#%%%#--:                   
                          +=:#%#%-+#%%%%###%%%==*#*#*#%%%%%#.                   
                         =%##%##%#%%%%%####%%%%%%%#%%%%%%#%#.                   
                      .--%%%%####%##################%%#%%%%#                    
                      =#%%*+**+***++++++++++*+*+++*+++*+#%%#.                   
                     .#%%+::.:.:......::..........:...::=#%#.                   
                      *%*:.:.:-::......................:-#%=                    
                     -*%=:...=*++-:......  ...  ........+#%=                    
                    .*##-:..:-+**++=:::.....    .. ....:*%%+.                   
                     *%*::....:=****=-:..... .  .......-*%#.                    
                     *%+-:... ..:=****=::...... ......:+#%+                     
                    .*#-:.........-+***+-.::.........::*#%+.                    
                    -%*:...........:=****=-:.. .. ....-#%#                      
                    *%+::.. ........:+****=:..     ...=#%+ :                    
                   :%#=:..  ......:=****+-:.... .....:+#%=:.                    
                 ..*%#-:.. ....::=**#*+-:............:*#%+.                     
                  =%%*:......:-=****+-.::............-*%%#-:                     
                  =%#+......:=++**+-...::..:::......:+#%#-                      
                  +%*-.....-+*#*+-:...:+++==+++=:..:-*##-                       
                 .#%*::....-**+-::...::++++++++=:..:=*%#:.                      
                .+%#+-:...:.....:..:......:--:.::.:-+*%%+                       
                 =%#*=----::::::::::::::.::--:::-:-+*#%#-                       
                 .+#%###*##*###**###*##***********####=  `

// NewModel creates a new application model
func NewModel() Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "Your Name"
	nameInput.Focus()
	nameInput.CharLimit = 50
	nameInput.Width = 30
	nameInput.Prompt = ""

	emailInput := textinput.New()
	emailInput.Placeholder = "Your Email"
	emailInput.CharLimit = 100
	emailInput.Width = 30
	emailInput.Prompt = ""

	initialFlames := strings.Split(asciiArt, "\n")

	// Initialize account store
	accountStore, err := auth.NewStore()
	if err != nil {
		// Log error but continue - app can work without persistence
		fmt.Printf("Warning: Failed to initialize account store: %v\n", err)
	}

	// Colorize the initial flames
	colorizedFlames := make([]string, len(initialFlames))
	for i, line := range initialFlames {
		colorizedFlames[i] = colorizeFlames(line)
	}
	
	m := Model{
		Flames:                       colorizedFlames,
		AppState:                     StateIntro,
		MenuCursor:                   0,
		NameInput:                    nameInput,
		EmailInput:                   emailInput,
		FocusedField:                 FocusName,
		MenuManager:                  menus.NewMenuManager(),
		ModuleExplorerSidebarOptions: []string{"Learn", "Studio", "Settings", "Exit"},
		ModuleExplorerSidebarCursor:  0,
		ModuleExplorerContent:        "Welcome to PowerHell! Select a module from the sidebar.",
		AccountStore:                 accountStore,
	}
	
	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Start ticking for animation and focus name input
	return tea.Batch(tickCmd(), m.NameInput.Focus())
}

// Cleanup performs cleanup operations
func (m *Model) Cleanup() {
	// End session if one is active
	if m.SessionID > 0 && m.AccountStore != nil {
		m.AccountStore.EndSession(m.SessionID)
	}
	
	// Close database connection
	if m.AccountStore != nil {
		m.AccountStore.Close()
	}
}

// TickMsg is used to advance the flame animation
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func generateAccountNumber() string {
	var b strings.Builder
	// Ensure the first digit is not 0
	b.WriteString(fmt.Sprintf("%d", rand.Intn(9)+1))
	for i := 1; i < 16; i++ {
		if i%4 == 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%d", rand.Intn(10)))
	}
	return b.String()
}

// Color definitions
var (
	darkGray = lipgloss.Color("#1e1b18")
	red      = lipgloss.Color("#ff4e00")
	orange   = lipgloss.Color("#fb923c")
	yellow   = lipgloss.Color("#fcd34d")
)

// Helper function to apply random fire colors to a line of the ASCII art
func colorizeFlames(line string) string {
	var coloredLine strings.Builder
	ansiEscape := regexp.MustCompile(`\x1b\[[0-9;]*m`) // Regex to match ANSI escape codes

	// Strip existing ANSI escape codes
	line = ansiEscape.ReplaceAllString(line, "")

	for _, char := range line {
		if char != ' ' { 
			coloredLine.WriteString(lipgloss.NewStyle().Foreground(flameColor()).Render(string(char)))
		} else {
			coloredLine.WriteString(string(char)) 
		}
	}

	return coloredLine.String()
}

const flameChars = " .:-=+*#%@"

// Helper function to strip ANSI escape codes from a string.
func stripANSI(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// Helper function to create a shimmering/pulsing animation effect
func animateFlames(currentFlames []string) []string {
	newFlames := make([]string, len(currentFlames))

	for i, line := range currentFlames {
		rawLine := stripANSI(line) // Strip existing colors

		var animatedLine strings.Builder
		for _, char := range rawLine {
			if char != ' ' {
				animatedLine.WriteString(lipgloss.NewStyle().Foreground(flameColor()).Render(string(char)))
			} else {
				animatedLine.WriteString(" ")
			}
		}

		newFlames[i] = animatedLine.String()
	}

	return newFlames
}

// flameColor returns a random fire color
func flameColor() lipgloss.Color {
	fireColors := []lipgloss.Color{
		lipgloss.Color("#ff4e00"), // Original Red/Orange
		lipgloss.Color("#fb923c"), // Original Orange
		lipgloss.Color("#fcd34d"), // Original Yellow
	}
	return fireColors[rand.Intn(len(fireColors))]
}