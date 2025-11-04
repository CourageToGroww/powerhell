package settings

import (
	"github.com/couragetogroww/powerhell/pkg/menus/types"
)

// SettingsMenu represents the application settings menu
type SettingsMenu struct {
	*types.BaseMenu
}

// NewSettingsMenu creates a new settings menu instance
func NewSettingsMenu() *SettingsMenu {
	base := types.NewBaseMenu("PowerHell Settings", "Configure your PowerHell learning experience")
	
	menu := &SettingsMenu{
		BaseMenu: base,
	}
	
	// Add settings options
	menu.setupOptions()
	
	return menu
}

// setupOptions configures all settings options
func (m *SettingsMenu) setupOptions() {
	// Theme Settings
	m.AddExecuteOption(
		"Theme Settings",
		"Configure color themes and visual preferences",
		m.handleThemeSettings,
	)
	
	// Learning Preferences
	m.AddExecuteOption(
		"Learning Preferences",
		"Set your skill level and learning pace preferences",
		m.handleLearningPreferences,
	)
	
	// Notification Settings
	m.AddExecuteOption(
		"Notification Settings",
		"Configure learning reminders and progress notifications",
		m.handleNotificationSettings,
	)
	
	// Account Settings
	m.AddExecuteOption(
		"Account Settings",
		"Manage your account information and privacy settings",
		m.handleAccountSettings,
	)
	
	// Export/Import Progress
	m.AddExecuteOption(
		"Export/Import Progress",
		"Backup or restore your learning progress",
		m.handleProgressManagement,
	)
	
	// Reset Settings
	m.AddExecuteOption(
		"Reset to Defaults",
		"Reset all settings to their default values",
		m.handleResetSettings,
	)
	
	// Back to main menu
	m.AddBackOption("Back to Main Menu", StateMainMenu)
}

// Settings handler functions
func (m *SettingsMenu) handleThemeSettings() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Theme settings functionality coming soon...",
		Data:    "theme_settings",
	}
}

func (m *SettingsMenu) handleLearningPreferences() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Learning preferences functionality coming soon...",
		Data:    "learning_preferences",
	}
}

func (m *SettingsMenu) handleNotificationSettings() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Notification settings functionality coming soon...",
		Data:    "notification_settings",
	}
}

func (m *SettingsMenu) handleAccountSettings() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Account settings functionality coming soon...",
		Data:    "account_settings",
	}
}

func (m *SettingsMenu) handleProgressManagement() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Progress management functionality coming soon...",
		Data:    "progress_management",
	}
}

func (m *SettingsMenu) handleResetSettings() types.MenuResult {
	return types.MenuResult{
		Action:  types.ActionExecute,
		Message: "Settings reset functionality coming soon...",
		Data:    "reset_settings",
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
