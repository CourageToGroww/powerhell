package types

import tea "github.com/charmbracelet/bubbletea"

// MenuAction represents different actions a menu can perform
type MenuAction int

const (
	ActionNone MenuAction = iota
	ActionNavigate
	ActionExecute
	ActionBack
	ActionExit
)

// MenuResult contains the result of a menu selection
type MenuResult struct {
	Action      MenuAction
	NextState   int
	Data        interface{}
	Message     string
	TeaCmd      tea.Cmd
}

// MenuOption represents a single menu option
type MenuOption struct {
	Label       string
	Description string
	Action      MenuAction
	Target      int    // Target state for navigation
	Handler     func() MenuResult // Custom handler function
}

// Menu interface that all menus must implement
type Menu interface {
	// GetOptions returns all menu options
	GetOptions() []MenuOption
	
	// GetTitle returns the menu title
	GetTitle() string
	
	// GetDescription returns optional menu description
	GetDescription() string
	
	// HandleSelection processes a menu selection by index
	HandleSelection(index int) MenuResult
	
	// GetBackOption returns the index of the back/exit option (-1 if none)
	GetBackOption() int
	
	// Initialize performs any setup needed for the menu
	Initialize() error
	
	// Cleanup performs any cleanup when leaving the menu
	Cleanup() error
}

// BaseMenu provides common functionality for all menus
type BaseMenu struct {
	title       string
	description string
	options     []MenuOption
	backIndex   int
}

// NewBaseMenu creates a new base menu
func NewBaseMenu(title, description string) *BaseMenu {
	return &BaseMenu{
		title:       title,
		description: description,
		options:     make([]MenuOption, 0),
		backIndex:   -1,
	}
}

// GetTitle returns the menu title
func (m *BaseMenu) GetTitle() string {
	return m.title
}

// GetDescription returns the menu description
func (m *BaseMenu) GetDescription() string {
	return m.description
}

// GetOptions returns all menu options
func (m *BaseMenu) GetOptions() []MenuOption {
	return m.options
}

// GetBackOption returns the back option index
func (m *BaseMenu) GetBackOption() int {
	return m.backIndex
}

// AddOption adds a new option to the menu
func (m *BaseMenu) AddOption(option MenuOption) {
	m.options = append(m.options, option)
}

// AddSimpleOption adds a simple navigation option
func (m *BaseMenu) AddSimpleOption(label string, targetState int) {
	m.AddOption(MenuOption{
		Label:  label,
		Action: ActionNavigate,
		Target: targetState,
	})
}

// AddExecuteOption adds an option that executes a function
func (m *BaseMenu) AddExecuteOption(label, description string, handler func() MenuResult) {
	m.AddOption(MenuOption{
		Label:       label,
		Description: description,
		Action:      ActionExecute,
		Handler:     handler,
	})
}

// AddBackOption adds a back/exit option
func (m *BaseMenu) AddBackOption(label string, targetState int) {
	m.backIndex = len(m.options)
	m.AddOption(MenuOption{
		Label:  label,
		Action: ActionBack,
		Target: targetState,
	})
}

// HandleSelection provides default selection handling
func (m *BaseMenu) HandleSelection(index int) MenuResult {
	if index < 0 || index >= len(m.options) {
		return MenuResult{Action: ActionNone}
	}
	
	option := m.options[index]
	
	// If option has a custom handler, use it
	if option.Handler != nil {
		return option.Handler()
	}
	
	// Default handling based on action type
	return MenuResult{
		Action:    option.Action,
		NextState: option.Target,
		Message:   "Selected: " + option.Label,
	}
}

// Initialize provides default initialization (no-op)
func (m *BaseMenu) Initialize() error {
	return nil
}

// Cleanup provides default cleanup (no-op)
func (m *BaseMenu) Cleanup() error {
	return nil
}
