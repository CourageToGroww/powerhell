package auth

import (
	"github.com/couragetogroww/powerhell/pkg/menus/types"
)

// AuthMenu represents the authentication menu
type AuthMenu struct {
	*types.BaseMenu
}

// NewAuthMenu creates a new authentication menu instance
func NewAuthMenu() *AuthMenu {
	base := types.NewBaseMenu("PowerHell Authentication", "Sign up for a new account or sign in to continue")
	
	menu := &AuthMenu{
		BaseMenu: base,
	}
	
	// Add authentication options
	menu.setupOptions()
	
	return menu
}

// setupOptions configures all authentication options
func (m *AuthMenu) setupOptions() {
	// Sign Up option - create new account
	m.AddExecuteOption(
		"Sign Up",
		"Create a new PowerHell learning account",
		m.handleSignUp,
	)
	
	// Login option - existing account login
	m.AddExecuteOption(
		"Login",
		"Login to your existing PowerHell account",
		m.handleSignIn,
	)
	
	// Exit option - quit application
	m.AddBackOption("Exit", StateExit)
}

// Authentication handler functions
func (m *AuthMenu) handleSignUp() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateAccountCreation,
		Message:   "Starting account creation process...",
	}
}

func (m *AuthMenu) handleSignIn() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateSignInPlaceholder,
		Message:   "Opening sign-in form...",
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
	StateExit = -1 // Special state for exiting
)
