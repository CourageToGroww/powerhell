package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/couragetogroww/powerhell/pkg/menus"
	"github.com/couragetogroww/powerhell/pkg/menus/types"
	mainmenu "github.com/couragetogroww/powerhell/pkg/menus/mainmenu"
	"github.com/couragetogroww/powerhell/pkg/modules"
	"github.com/couragetogroww/powerhell/pkg/ui"
	"github.com/couragetogroww/powerhell/pkg/views"
)

const flameChars = " .:-=+*#%@"

var (
	darkGray = lipgloss.Color("#1e1b18")
	red      = lipgloss.Color("#ff4e00")
	orange   = lipgloss.Color("#fb923c")
	yellow   = lipgloss.Color("#fcd34d")
)

// AppState defines the current view/state of the application.
const (
	stateIntro = mainmenu.StateIntro
	stateAccountCreation = mainmenu.StateAccountCreation
	stateMainMenu = mainmenu.StateMainMenu
	stateLearnMenu = mainmenu.StateLearnMenu
	stateStudio = mainmenu.StateStudio
	stateSettings = mainmenu.StateSettings
	stateAuthMenu = mainmenu.StateAuthMenu
	stateSignInPlaceholder = mainmenu.StateSignInPlaceholder
	stateModuleExplorer = mainmenu.StateModuleExplorer
	stateDashboard = 100 // New dashboard state
	stateLesson = 101    // New lesson state
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
                 .+#%###*##*###**###*##***********####=  
`
// Styles
var fieryPowerhellArtStyle = lipgloss.NewStyle().Bold(true).PaddingTop(1) // For the main ASCII art, colors applied dynamically

// Stealthy UI Styles for Menus and other UI elements
var menuTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(orange).PaddingTop(1).PaddingBottom(1) // Muted Blue (default) -> CHANGED TO ORANGE
var orangeTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(orange).PaddingTop(1).PaddingBottom(1) // Orange for specific titles
var promptStyle = lipgloss.NewStyle().Foreground(orange).PaddingTop(1) // Orange for prompts
var menuOptionStyle = lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(lipgloss.Color("#A0A0A0")) // Light Gray
var selectedMenuOptionStyle = lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(lipgloss.Color("#fb923c")).Bold(true) // Orange FG, no background change
var menuContainerStyle = lipgloss.NewStyle().Margin((24 / 10), (80 / 4)).Background(lipgloss.Color("#2D2D2D")) // Dark Gray BG

// Styles for stateModuleExplorer layout
var powerhellTabStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(orange).
	Padding(0, 1)

var sidebarStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#505050")) // Dim border for sidebar

var contentAreaStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#505050")) // Dim border for content

var footerStyle = lipgloss.NewStyle().
	MarginTop(1).
	Padding(0,1).
	Border(lipgloss.ThickBorder(), true, false, false, false).
	BorderTop(true).
	BorderForeground(lipgloss.Color("#505050"))

// inputFieldStyle for text inputs
var inputFieldStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("#fb923c")).Padding(0,1).MarginRight(2)
var focusedInputFieldStyle = inputFieldStyle.Copy().BorderForeground(lipgloss.Color("#FFB347")) // Brighter orange for focus

// Placeholders for dynamic calculation if needed, or remove if Margin handles it well enough
const terminalHeightPlaceholder = 24 
const terminalWidthPlaceholder = 80

var fireColors = []lipgloss.Color{
	lipgloss.Color("#ff4e00"), // Original Red/Orange
	lipgloss.Color("#fb923c"), // Original Orange
	lipgloss.Color("#fcd34d"), // Original Yellow
}

func flameColor() lipgloss.Color {
	return fireColors[rand.Intn(len(fireColors))]
}

// Helper function to strip ANSI escape codes from a string.
// This is a simplified version; a robust one would handle more cases.
func stripANSI(s string) string {
	ansiPhone := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiPhone.ReplaceAllString(s, "")
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

func colorizePowerShellTitle(titleArt string) string {
	var result strings.Builder
	lines := strings.Split(titleArt, "\n")
	for i, line := range lines {
		for _, char := range line {
			if char == ' ' {
				result.WriteString(" ")
			} else {
				color := lipgloss.NewStyle().Foreground(flameColor()).Render(string(char))
				result.WriteString(color)
			}
		}
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}

type model struct {
	flames                   []string
	appState                 int
	previousState            int
	menuCursor               int    // Used for menu navigation
	terminalWidth            int
	terminalHeight           int
	nameInput                textinput.Model
	emailInput               textinput.Model
	generatedAccountNumber   string
	focusedField             int
	quit                     bool
	showHelp                 bool   // Show help overlay

	// Menu manager for modular menu system
	menuManager              *menus.MenuManager

	// Fields for stateModuleExplorer
	moduleExplorerSidebarCursor  int
	moduleExplorerSidebarOptions []string
	moduleExplorerContent        string
	
	// New UI components
	dashboard    *views.DashboardView
	lessonView   *views.LessonView
	currentModule *modules.Module
}

const (
	focusName = iota
	focusEmail
	focusDisplayInfo
)

// tickMsg is used to advance the flame animation.
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func generateAccountNumber() string {
	var b strings.Builder
	// Ensure the first digit is not 0 to be within "1000..." to "9999..." range concept
	b.WriteString(fmt.Sprintf("%d", rand.Intn(9)+1)) // First digit 1-9
	for i := 1; i < 16; i++ {
		if i%4 == 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%d", rand.Intn(10))) // Subsequent digits 0-9
	}
	return b.String()
}

func initialModel() model {
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
	for i, line := range initialFlames {
		initialFlames[i] = colorizeFlames(line)
	}

	m := model{
		flames:                       initialFlames,
		appState:                     stateIntro,
		menuCursor:                   0,
		nameInput:                    nameInput,
		emailInput:                   emailInput,
		focusedField:                 focusName,
		menuManager:                  menus.NewMenuManager(),
		moduleExplorerSidebarOptions: []string{"Learn", "Studio", "Settings", "Exit"},
		moduleExplorerSidebarCursor:  0,
		moduleExplorerContent:        "Welcome to PowerHell! Select a module from the sidebar.",
	}
	// Get terminal size will be updated later
	return m
}

func (m model) Init() tea.Cmd {
	// Initialize flames for the intro
	for i, line := range m.flames {
		m.flames[i] = colorizeFlames(line)
	}
	// Start ticking for animation and focus name input
	return tea.Batch(tickCmd(), m.nameInput.Focus())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		// Update view sizes
		if m.dashboard != nil {
			m.dashboard = views.NewDashboardView(m.terminalWidth, m.terminalHeight)
		}
		if m.lessonView != nil && m.currentModule != nil {
			m.lessonView = views.NewLessonView(m.currentModule, m.terminalWidth, m.terminalHeight)
		}

	case tickMsg: // Handle animation tick
		if m.appState == stateIntro {
			m.flames = animateFlames(m.flames)
			return m, tickCmd() // Continue ticking for intro animation
		}
		return m, nil // Ignore tick if not in intro state

	case tea.KeyMsg:
		switch m.appState {
		case stateIntro:
			if msg.String() == "enter" {
				m.appState = stateAuthMenu // Transition to AuthMenu first
				return m, nil // Stop the tick for flame animation
			} else if msg.String() == "ctrl+c" || msg.String() == "q" {
				m.quit = true
				return m, tea.Quit
			} else if msg.String() == "h" {
				m.showHelp = !m.showHelp
			}
		case stateAccountCreation:
			switch msg.String() {
			case "ctrl+c", "q":
				m.quit = true
				return m, tea.Quit
			case "enter":
				if m.focusedField == focusDisplayInfo {
					// Transition to the new Dashboard after account creation
					m.appState = stateDashboard
					m.dashboard = views.NewDashboardView(m.terminalWidth, m.terminalHeight)
					return m, nil
				} else if m.focusedField == focusEmail && m.nameInput.Value() != "" && m.emailInput.Value() != "" {
					// All fields filled, generate account number and move to display info
					m.generatedAccountNumber = generateAccountNumber()
					m.focusedField = focusDisplayInfo
					m.nameInput.Blur()
					m.emailInput.Blur()
					return m, nil
				} else {
					// Move to next field or handle incomplete form
					if m.focusedField == focusName {
						m.focusedField = focusEmail
						m.nameInput.Blur()
						return m, m.emailInput.Focus()
					} else if m.focusedField == focusEmail {
						// Potentially do nothing or show a message if name is empty
					}
				}
			case "tab", "shift+tab", "up", "down": // Simplified navigation
				// ... (account creation field navigation logic remains mostly the same)
				if m.focusedField == focusDisplayInfo { // Don't tab out of display info
					return m, nil
				}
				s := msg.String()
				if s == "up" || s == "shift+tab" {
					m.focusedField--
				} else {
					m.focusedField++
				}
				if m.focusedField > focusEmail {
					m.focusedField = focusName
				} else if m.focusedField < focusName {
					m.focusedField = focusEmail
				}
				var focusCmd tea.Cmd
				if m.focusedField == focusName {
					m.emailInput.Blur() // Call Blur directly, it doesn't return a command
					focusCmd = m.nameInput.Focus()
					m.nameInput.PromptStyle = focusedInputFieldStyle
					m.emailInput.PromptStyle = inputFieldStyle
				} else { // focusedField == focusEmail
					m.nameInput.Blur() // Call Blur directly
					focusCmd = m.emailInput.Focus()
					m.nameInput.PromptStyle = inputFieldStyle
					m.emailInput.PromptStyle = focusedInputFieldStyle
				}
				return m, focusCmd // Return the single focus command
			}
			var cmd tea.Cmd
			if m.focusedField == focusName {
				m.nameInput, cmd = m.nameInput.Update(msg)
			} else if m.focusedField == focusEmail {
				m.emailInput, cmd = m.emailInput.Update(msg)
			}
			return m, cmd

		case stateAuthMenu:
			// Initialize auth menu if not already set
			if m.menuManager.GetCurrentState() != stateAuthMenu {
				m.menuManager.SetCurrentMenu(stateAuthMenu)
				m.menuCursor = 0
			}
			
			switch msg.String() {
			case "ctrl+c", "q":
				m.quit = true
				return m, tea.Quit
			case "h":
				m.showHelp = !m.showHelp
			case "up", "k":
				if m.menuCursor > 0 {
					m.menuCursor--
				}
			case "down", "j":
				menuOptions := m.menuManager.GetMenuOptionsAsStrings()
				if m.menuCursor < len(menuOptions)-1 {
					m.menuCursor++
				}
			case "enter":
				result := m.menuManager.HandleSelection(m.menuCursor)
				switch result.Action {
				case types.ActionNavigate:
					m.appState = result.NextState
					m.menuCursor = 0
					if result.NextState == stateAccountCreation {
						m.focusedField = focusName
						m.nameInput.SetValue("")
						m.emailInput.SetValue("")
						m.generatedAccountNumber = ""
						return m, m.nameInput.Focus()
					} else if result.NextState == stateSignInPlaceholder {
						// Transition to Module Explorer on Sign In
						m.appState = stateModuleExplorer
						m.moduleExplorerSidebarCursor = 0
						m.moduleExplorerSidebarOptions = []string{"Learn", "Studio", "Settings", "Log Out", "Exit"}
						m.moduleExplorerContent = "Signed in! Welcome to PowerHell."
						return m, nil
					}
				case types.ActionExit:
					m.quit = true
					return m, tea.Quit
				}
			}

		// Deprecated states - redirect to stateModuleExplorer or handle as error/legacy
		case stateMainMenu, stateLearnMenu, stateSignInPlaceholder, stateStudio, stateSettings:
			m.appState = stateModuleExplorer 
			m.moduleExplorerSidebarCursor = 0 // Reset cursor
			m.moduleExplorerContent = fmt.Sprintf("Redirected from a legacy state (%v).", m.appState) // Inform user
			return m, nil

		case stateModuleExplorer:
			switch msg.String() {
			case "ctrl+c", "q":
				m.quit = true
				return m, tea.Quit
			case "up", "k":
				if m.moduleExplorerSidebarCursor > 0 {
					m.moduleExplorerSidebarCursor--
				}
			case "down", "j":
				if m.moduleExplorerSidebarCursor < len(m.moduleExplorerSidebarOptions)-1 {
					m.moduleExplorerSidebarCursor++
				}
			case "enter":
				selectedModule := m.moduleExplorerSidebarOptions[m.moduleExplorerSidebarCursor]
				switch selectedModule {
				case "Learn":
					// Switch sidebar to learn menu options. For future submenus, add similar cases here.
					m.moduleExplorerSidebarOptions = []string{"On-Premise Learning", "MSGraph Module", "Exchange Module", "HTTP Client (Invoke-WebRequest, Graph API)", "Custom & Popular Modules (EntraExporter)", "PowerShell SDK Learning", "Exit"}
					m.moduleExplorerSidebarCursor = 0
					m.moduleExplorerContent = "Select a learning module."
				case "Exit":
					m.quit = true
					return m, tea.Quit
				default:
					// In the future, this will launch the actual module UI or a sub-menu
					m.moduleExplorerContent = fmt.Sprintf("Selected: %s\n(Press Q to quit)", selectedModule)
				}
			}
		
		case stateDashboard:
			if m.dashboard == nil {
				m.dashboard = views.NewDashboardView(m.terminalWidth, m.terminalHeight)
			}
			switch msg.String() {
			case "ctrl+c", "q":
				m.quit = true
				return m, tea.Quit
			case "h":
				m.showHelp = !m.showHelp
			case "enter":
				if selectedModule := m.dashboard.GetSelectedModule(); selectedModule != nil {
					m.currentModule = selectedModule
					m.lessonView = views.NewLessonView(selectedModule, m.terminalWidth, m.terminalHeight)
					m.appState = stateLesson
				}
			default:
				if !m.showHelp {
					m.dashboard.Update(msg.String())
				}
			}
		
		case stateLesson:
			if m.lessonView == nil {
				return m, nil
			}
			switch msg.String() {
			case "ctrl+c":
				m.quit = true
				return m, tea.Quit
			case "h":
				m.showHelp = !m.showHelp
			case "q":
				// Go back to dashboard
				m.appState = stateDashboard
			default:
				if !m.showHelp {
					m.lessonView.Update(msg.String())
				}
			}
		}
		// Common key handling for all states (like quit)
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			if m.appState != stateAccountCreation || m.focusedField == focusDisplayInfo {
				m.quit = true
				return m, tea.Quit
			}
		}
	} // End of tea.KeyMsg case

	// Add other top-level message types (e.g., tea.MouseMsg) if needed

	return m, cmd
}

func (m model) View() string {
	if m.quit {
		return "Exiting PowerHell... Goodbye!\n"
	}

	var mainView string
	switch m.appState {
	case stateIntro:
		// Use the original intro with animated flames
		flameBlockLines := make([]string, len(m.flames))
		for i, line := range m.flames {
			flameBlockLines[i] = colorizeFlames(line)
		}

		// Construct the flames art block
		flameBlock := lipgloss.JoinVertical(lipgloss.Left, flameBlockLines...)

		// Center the flame block
		centeredFlameBlock := lipgloss.Place(
			m.terminalWidth,
			len(flameBlockLines),
			lipgloss.Center,
			lipgloss.Top,
			flameBlock,
		)

		// Center the "POWERHELL" title below the flames
		centeredPowerhellTitle := lipgloss.Place(
			m.terminalWidth,
			1,
			lipgloss.Center,
			lipgloss.Top,
			colorizePowerShellTitle(`
██████╗  ██████╗ ██╗    ██╗███████╗██████╗ ██╗  ██╗███████╗██╗     ██╗     
██╔══██╗██╔═══██╗██║    ██║██╔════╝██╔══██╗██║  ██║██╔════╝██║     ██║     
██████╔╝██║   ██║██║ █╗ ██║█████╗  ██████╔╝███████║█████╗  ██║     ██║     
██╔═══╝ ██║   ██║██║███╗██║██╔══╝  ██╔══██╗██╔══██║██╔══╝  ██║     ██║     
██║     ╚██████╔╝╚███╔███╔╝███████╗██║  ██║██║  ██║███████╗███████╗███████╗
╚═╝      ╚═════╝  ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝
`),
		)

		// Center the "Press Enter to Continue" text under the "POWERHELL" block
		promptText := "Press Enter to continue..."
		renderedPromptText := promptStyle.Render(promptText)
		promptContainer := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Render(renderedPromptText)

		finalCenteredPrompt := lipgloss.Place(
			m.terminalWidth,
			1, // Single line for the prompt
			lipgloss.Center,
			lipgloss.Top,
			promptContainer,
		)

		// Combine all components with vertical spacing
		content := lipgloss.JoinVertical(
			lipgloss.Top,
			centeredFlameBlock,
			centeredPowerhellTitle,
			finalCenteredPrompt,
		)

		mainView = lipgloss.Place(
			m.terminalWidth,
			m.terminalHeight,
			lipgloss.Center,
			lipgloss.Top,
			content,
		)
	case stateSignInPlaceholder:
		// This state should now be redirected by Update(), but as a fallback:
		mainView = renderModuleExplorerView(m)
	case stateAuthMenu:
		mainView = renderMenuView(m, "Main Menu", "Use ↑/↓ to navigate and Enter to select.")
	case stateModuleExplorer:
		mainView = renderModuleExplorerView(m)
	case stateDashboard:
		if m.dashboard != nil {
			mainView = m.dashboard.Render()
		} else {
			mainView = "Loading dashboard..."
		}
	case stateLesson:
		if m.lessonView != nil {
			mainView = m.lessonView.Render()
		} else {
			mainView = "Loading lesson..."
		}
	case stateMainMenu, stateLearnMenu, stateStudio, stateSettings: // Catch-all for other deprecated states
		// These states should now be redirected by Update(), but as a fallback:
		mainView = renderModuleExplorerView(m)
	default:
		mainView = "Unknown state."
	}

	// Overlay help if needed
	if m.showHelp {
		var context string
		switch m.appState {
		case stateIntro:
			context = "intro"
		case stateDashboard:
			context = "dashboard"
		case stateLesson:
			context = "lesson"
		case stateAuthMenu, stateMainMenu, stateLearnMenu:
			context = "menu"
		default:
			context = "general"
		}
		return ui.HelpOverlay(m.terminalWidth, m.terminalHeight, context)
	}

	return mainView
}

func renderSharedExplorerLayout(m model, sidebarOptions []string, sidebarCursor int, mainContent string, footerHint string) string {
	// Header
	powerhellTab := powerhellTabStyle.Render("PowerHell")
	// The headerView is essentially the container for our tab(s), providing margin below.
	headerView := lipgloss.NewStyle().MarginBottom(1).Render(powerhellTab)

	// Sidebar content
	var sidebarContent strings.Builder
	for i, choice := range sidebarOptions {
		style := menuOptionStyle
		if sidebarCursor == i {
			style = selectedMenuOptionStyle
		}
		sidebarContent.WriteString(style.Render(choice) + "\n")
	}

	// Define heights and widths
	headerHeight := lipgloss.Height(headerView)
	footerHeight := 3 // Approximate height for footer (1 for border, 1 for text, 1 for margin/padding)
	availableHeight := m.terminalHeight - headerHeight - footerHeight
	if availableHeight < 0 {
		availableHeight = 0
	}

	sidebarWidth := int(float64(m.terminalWidth) * 0.33) // Sidebar takes 1/3 of the width
	if sidebarWidth < 20 { // Minimum sidebar width
		sidebarWidth = 20
	}
	contentWidth := m.terminalWidth - sidebarWidth - 1 // Account for the single space character used in JoinHorizontal
	if contentWidth < 0 {
		contentWidth = 0
	}

	sidebarFormatted := sidebarStyle.Width(sidebarWidth).Height(availableHeight).Render(sidebarContent.String())

	// Content Area
	contentFormatted := contentAreaStyle.Width(contentWidth).Height(availableHeight).Render(mainContent)

	// Main body (Sidebar + Content)
	mainBodyView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarFormatted, " ", contentFormatted)

	// Footer
	styledFooterHint := promptStyle.Render(footerHint)
	footerView := footerStyle.Width(m.terminalWidth).Render(lipgloss.PlaceHorizontal(m.terminalWidth, lipgloss.Center, styledFooterHint))

	// Final layout
	finalView := lipgloss.JoinVertical(lipgloss.Left,
		headerView,
		mainBodyView,
		footerView,
	)

	return finalView
}

func renderMenuView(m model, title string, promptText string) string {
	var menuContent strings.Builder
	var currentTitleStyle lipgloss.Style

	if title == "Learn PowerShell" || title == "Main Menu" || title == "Authentication Menu" {
		currentTitleStyle = orangeTitleStyle
	} else {
		currentTitleStyle = menuTitleStyle
	}

	menuContent.WriteString(currentTitleStyle.Render(title) + "\n\n")
	for i, option := range m.moduleExplorerSidebarOptions {
		if i == m.menuCursor {
			menuContent.WriteString(selectedMenuOptionStyle.Render("> " + option))
		} else {
			menuContent.WriteString(menuOptionStyle.Render("  " + option))
		}
		menuContent.WriteString("\n")
	}
	menuContent.WriteString("\n" + promptStyle.Render(promptText))

	if m.terminalWidth > 0 {
		menuBlockWidth := m.terminalWidth * 2 / 3
		if menuBlockWidth < 40 {
			menuBlockWidth = 40
		}
		if menuBlockWidth > 100 {
			menuBlockWidth = 100
		}

		styledMenu := lipgloss.NewStyle().
			Width(menuBlockWidth).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#fb923c")).
			Render(menuContent.String())

		return lipgloss.Place(m.terminalWidth, m.terminalHeight, lipgloss.Center, lipgloss.Center, styledMenu)
	}
	return lipgloss.NewStyle().Align(lipgloss.Center).Render(menuContent.String())
}

func renderAccountCreationView(m model) string {
	var s strings.Builder
	title := orangeTitleStyle.Render("Create Your PowerHell Account")

	if m.focusedField == focusDisplayInfo {
		title = orangeTitleStyle.Render("Account Created Successfully!")
		s.WriteString(title + "\n\n")
		s.WriteString(fmt.Sprintf("Name: %s\n", m.nameInput.Value()))
		s.WriteString(fmt.Sprintf("Email: %s\n", m.emailInput.Value()))
		s.WriteString(fmt.Sprintf("Account Number: %s\n\n", m.generatedAccountNumber))
		s.WriteString(promptStyle.Render("IMPORTANT: Save your Account Number! It's your key to PowerHell.\nPress Enter to continue to Main Menu."))
	} else {
		s.WriteString(title + "\n\n")
		s.WriteString("Name:  " + m.nameInput.View() + "\n")
		s.WriteString("Email: " + m.emailInput.View() + "\n\n")
		if m.focusedField == focusName {
			s.WriteString(promptStyle.Render("Enter your name, then press Enter."))
		} else if m.focusedField == focusEmail {
			s.WriteString(promptStyle.Render("Enter your email, then press Enter."))
		}
		s.WriteString(promptStyle.Render("\n\nPress 'M' for Main Menu at any time."))
	}

	if m.terminalWidth > 0 {
		formWidth := m.terminalWidth * 2 / 3
		if formWidth < 50 {
			formWidth = 50
		}
		if formWidth > 100 {
			formWidth = 100
		}

		styledForm := lipgloss.NewStyle().
			Width(formWidth).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#fb923c")).
			Render(s.String())

		return lipgloss.Place(m.terminalWidth, m.terminalHeight, lipgloss.Center, lipgloss.Center, styledForm)
	}
	return lipgloss.NewStyle().Align(lipgloss.Center).Render(s.String())
}

func renderPlaceholderView(m model, titleText string, bodyText string) string {
	var s strings.Builder

	// Ensure title and body are not empty to avoid panic with lipgloss operations if they are
	if titleText == "" {
		titleText = "Notification"
	}
	if bodyText == "" {
		bodyText = "(No details provided)"
	}

	title := menuTitleStyle.Render(titleText)
	body := menuOptionStyle.Render(bodyText)

	// Centering the content
	s.WriteString(lipgloss.PlaceHorizontal(m.terminalWidth, lipgloss.Center, title) + "\n\n")
	s.WriteString(lipgloss.PlaceHorizontal(m.terminalWidth, lipgloss.Center, body))

	return s.String()
}

func renderModuleExplorerView(m model) string {
	return renderSharedExplorerLayout(
		m,
		m.moduleExplorerSidebarOptions,
		m.moduleExplorerSidebarCursor,
		m.moduleExplorerContent,
		"↑/↓ Navigate | Enter Select | Q Quit",
	)
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	// Initialize the model
	m := initialModel()

	// Start the Bubble Tea program with AltScreen
	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
