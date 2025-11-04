package studio

import (
	"github.com/couragetogroww/powerhell/pkg/menus/types"
)

// StudioMenu represents the PowerShell studio/practice environment menu
type StudioMenu struct {
	*types.BaseMenu
}

// NewStudioMenu creates a new studio menu instance
func NewStudioMenu() *StudioMenu {
	base := types.NewBaseMenu("PowerHell Studio", "Practice PowerShell in interactive environments")
	
	menu := &StudioMenu{
		BaseMenu: base,
	}
	
	// Add studio options
	menu.setupOptions()
	
	return menu
}

// setupOptions configures all studio options
func (m *StudioMenu) setupOptions() {
	// Interactive PowerShell Console
	m.AddExecuteOption(
		"Interactive Console",
		"Launch an interactive PowerShell console for practice",
		m.handleInteractiveConsole,
	)
	
	// Script Editor
	m.AddExecuteOption(
		"Script Editor",
		"Create and edit PowerShell scripts with syntax highlighting",
		m.handleScriptEditor,
	)
	
	// Sandbox Environment
	m.AddExecuteOption(
		"Sandbox Environment",
		"Practice in a safe, isolated PowerShell environment",
		m.handleSandboxEnvironment,
	)
	
	// Code Challenges
	m.AddExecuteOption(
		"Code Challenges",
		"Solve PowerShell coding challenges to test your skills",
		m.handleCodeChallenges,
	)
	
	// Snippet Library
	m.AddExecuteOption(
		"Snippet Library",
		"Browse and save useful PowerShell code snippets",
		m.handleSnippetLibrary,
	)
	
	// Project Templates
	m.AddExecuteOption(
		"Project Templates",
		"Start new PowerShell projects from pre-built templates",
		m.handleProjectTemplates,
	)
	
	// Back to main menu
	m.AddBackOption("Back to Main Menu", StateMainMenu)
}

// Studio handler functions
func (m *StudioMenu) handleInteractiveConsole() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Launching interactive PowerShell console...",
		Data:    "interactive_console",
	}
}

func (m *StudioMenu) handleScriptEditor() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Opening PowerShell script editor...",
		Data:    "script_editor",
	}
}

func (m *StudioMenu) handleSandboxEnvironment() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Initializing sandbox environment...",
		Data:    "sandbox_environment",
	}
}

func (m *StudioMenu) handleCodeChallenges() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Loading code challenges...",
		Data:    "code_challenges",
	}
}

func (m *StudioMenu) handleSnippetLibrary() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Opening snippet library...",
		Data:    "snippet_library",
	}
}

func (m *StudioMenu) handleProjectTemplates() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Loading project templates...",
		Data:    "project_templates",
	}
}

// State constants (should match main menu states)
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
)
