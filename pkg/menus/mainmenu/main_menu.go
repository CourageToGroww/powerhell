package mainmenu

import (
	"github.com/couragetogroww/powerhell/pkg/menus/types"
)

// MainMenu represents the primary application menu
type MainMenu struct {
	*types.BaseMenu
}

// NewMainMenu creates a new main menu instance
func NewMainMenu() *MainMenu {
	base := types.NewBaseMenu("PowerHell Main Menu", "Choose your learning path or manage your session")
	
	menu := &MainMenu{
		BaseMenu: base,
	}
	
	// Add main menu options
	menu.setupOptions()
	
	return menu
}

// setupOptions configures all main menu options
func (m *MainMenu) setupOptions() {
	// Learn option - navigate to learning modules
	m.AddSimpleOption("Learn", StateLearnMenu)
	
	// Studio option - navigate to PowerShell studio
	m.AddSimpleOption("Studio", StateStudio)
	
	// Settings option - navigate to settings
	m.AddSimpleOption("Settings", StateSettings)
	
	// Log Out option - return to auth menu
	m.AddSimpleOption("Log Out", StateAuthMenu)
	
	// Exit option - quit application
	m.AddBackOption("Exit", StateExit)
}

// State constants for navigation
const (
	StateIntro = iota
	StateAccountCreation
	StateMainMenu
	StateLearnMenu
	StateStudio
	StateSettings
	StateAuthMenu
	StateSignInPlaceholder
	StateModuleExplorer
	StateExit = -1 // Special state for exiting
)

// GetStateConstants returns all state constants for use in other packages
func GetStateConstants() map[string]int {
	return map[string]int{
		"Intro":              StateIntro,
		"AccountCreation":    StateAccountCreation,
		"MainMenu":           StateMainMenu,
		"LearnMenu":          StateLearnMenu,
		"Studio":             StateStudio,
		"Settings":           StateSettings,
		"AuthMenu":           StateAuthMenu,
		"SignInPlaceholder":  StateSignInPlaceholder,
		"ModuleExplorer":     StateModuleExplorer,
		"Exit":               StateExit,
	}
}
