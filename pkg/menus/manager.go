package menus

import (
	"fmt"

	"github.com/couragetogroww/powerhell/pkg/menus/auth"
	"github.com/couragetogroww/powerhell/pkg/menus/learn"
	mainmenu "github.com/couragetogroww/powerhell/pkg/menus/mainmenu"
	"github.com/couragetogroww/powerhell/pkg/menus/settings"
	"github.com/couragetogroww/powerhell/pkg/menus/studio"
	"github.com/couragetogroww/powerhell/pkg/menus/types"
)

// MenuManager handles all menu instances and navigation
type MenuManager struct {
	menus       map[int]types.Menu
	currentMenu types.Menu
	currentState int
}

// NewMenuManager creates a new menu manager with all menus initialized
func NewMenuManager() *MenuManager {
	manager := &MenuManager{
		menus: make(map[int]types.Menu),
	}
	
	// Initialize all menus
	manager.initializeMenus()
	
	return manager
}

// initializeMenus creates and registers all menu instances
func (m *MenuManager) initializeMenus() {
	// Register all menus with their corresponding states
	m.menus[mainmenu.StateMainMenu] = mainmenu.NewMainMenu()
	m.menus[mainmenu.StateLearnMenu] = learn.NewLearnMenu()
	m.menus[mainmenu.StateAuthMenu] = auth.NewAuthMenu()
	m.menus[mainmenu.StateSettings] = settings.NewSettingsMenu()
	m.menus[mainmenu.StateStudio] = studio.NewStudioMenu()
	
	// Initialize all menus
	for _, menu := range m.menus {
		if err := menu.Initialize(); err != nil {
			fmt.Printf("Error initializing menu: %v\n", err)
		}
	}
}

// SetCurrentMenu sets the active menu by state
func (m *MenuManager) SetCurrentMenu(state int) error {
	menu, exists := m.menus[state]
	if !exists {
		return fmt.Errorf("menu for state %d not found", state)
	}
	
	// Cleanup previous menu if exists
	if m.currentMenu != nil {
		if err := m.currentMenu.Cleanup(); err != nil {
			fmt.Printf("Error cleaning up previous menu: %v\n", err)
		}
	}
	
	m.currentMenu = menu
	m.currentState = state
	
	return nil
}

// GetCurrentMenu returns the currently active menu
func (m *MenuManager) GetCurrentMenu() types.Menu {
	return m.currentMenu
}

// GetCurrentState returns the current menu state
func (m *MenuManager) GetCurrentState() int {
	return m.currentState
}

// HandleSelection processes a menu selection and returns the result
func (m *MenuManager) HandleSelection(index int) types.MenuResult {
	if m.currentMenu == nil {
		return types.MenuResult{Action: types.ActionNone}
	}
	
	return m.currentMenu.HandleSelection(index)
}

// GetMenuOptions returns the options for the current menu
func (m *MenuManager) GetMenuOptions() []types.MenuOption {
	if m.currentMenu == nil {
		return []types.MenuOption{}
	}
	
	return m.currentMenu.GetOptions()
}

// GetMenuTitle returns the title of the current menu
func (m *MenuManager) GetMenuTitle() string {
	if m.currentMenu == nil {
		return ""
	}
	
	return m.currentMenu.GetTitle()
}

// GetMenuDescription returns the description of the current menu
func (m *MenuManager) GetMenuDescription() string {
	if m.currentMenu == nil {
		return ""
	}
	
	return m.currentMenu.GetDescription()
}

// GetBackOptionIndex returns the back option index for the current menu
func (m *MenuManager) GetBackOptionIndex() int {
	if m.currentMenu == nil {
		return -1
	}
	
	return m.currentMenu.GetBackOption()
}

// GetMenuOptionsAsStrings returns menu options as string slice for compatibility
func (m *MenuManager) GetMenuOptionsAsStrings() []string {
	options := m.GetMenuOptions()
	strings := make([]string, len(options))
	
	for i, option := range options {
		strings[i] = option.Label
	}
	
	return strings
}

// Cleanup performs cleanup for all menus
func (m *MenuManager) Cleanup() {
	for _, menu := range m.menus {
		if err := menu.Cleanup(); err != nil {
			fmt.Printf("Error cleaning up menu: %v\n", err)
		}
	}
}
