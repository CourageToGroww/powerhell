package learn

import (
	"github.com/couragetogroww/powerhell/pkg/menus/types"
)

// LearnMenu represents the learning modules menu
type LearnMenu struct {
	*types.BaseMenu
}

// NewLearnMenu creates a new learn menu instance
func NewLearnMenu() *LearnMenu {
	base := types.NewBaseMenu("PowerHell Learning Modules", "Select a PowerShell learning module to explore")
	
	menu := &LearnMenu{
		BaseMenu: base,
	}
	
	// Add learning module options
	menu.setupOptions()
	
	return menu
}

// setupOptions configures all learning module options
func (m *LearnMenu) setupOptions() {
	// On-Premise Learning - Active Directory, file servers, etc.
	m.AddExecuteOption(
		"On-Premise Learning",
		"Learn PowerShell for on-premise environments (AD, file servers, local management)",
		m.handleOnPremiseModule,
	)
	
	// MSGraph Module - Microsoft Graph PowerShell SDK
	m.AddExecuteOption(
		"MSGraph Module",
		"Learn Microsoft Graph PowerShell SDK for cloud management",
		m.handleMSGraphModule,
	)
	
	// Exchange Module - Exchange PowerShell management
	m.AddExecuteOption(
		"Exchange Module",
		"Learn Exchange PowerShell for email system management",
		m.handleExchangeModule,
	)
	
	// HTTP Client - Invoke-WebRequest and direct API calls
	m.AddExecuteOption(
		"HTTP Client (Invoke-WebRequest, Graph API)",
		"Learn web requests and API interactions with PowerShell",
		m.handleHTTPClientModule,
	)
	
	// Custom & Popular Modules - EntraExporter, etc.
	m.AddExecuteOption(
		"Custom & Popular Modules (EntraExporter)",
		"Learn popular community modules and custom PowerShell tools",
		m.handleCustomModulesModule,
	)
	
	// PowerShell SDK Learning - C# development
	m.AddExecuteOption(
		"PowerShell SDK Learning",
		"Learn PowerShell SDK development with C#",
		m.handleSDKLearningModule,
	)
	
	// Back to main menu
	m.AddBackOption("Back to Main Menu", StateMainMenu)
}

// Module handler functions - these will launch specific learning modules
func (m *LearnMenu) handleOnPremiseModule() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateModuleExplorer,
		Data:      "onprem",
		Message:   "Loading On-Premise Learning Module...",
	}
}

func (m *LearnMenu) handleMSGraphModule() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateModuleExplorer,
		Data:      "mggraph",
		Message:   "Loading MSGraph Learning Module...",
	}
}

func (m *LearnMenu) handleExchangeModule() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateModuleExplorer,
		Data:      "exchange",
		Message:   "Loading Exchange Learning Module...",
	}
}

func (m *LearnMenu) handleHTTPClientModule() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateModuleExplorer,
		Data:      "httpclient",
		Message:   "Loading HTTP Client Learning Module...",
	}
}

func (m *LearnMenu) handleCustomModulesModule() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateModuleExplorer,
		Data:      "custommodules",
		Message:   "Loading Custom Modules Learning...",
	}
}

func (m *LearnMenu) handleSDKLearningModule() types.MenuResult {
	return types.MenuResult{
		Action:    types.ActionNavigate,
		NextState: StateModuleExplorer,
		Data:      "sdklearning",
		Message:   "Loading PowerShell SDK Learning...",
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
