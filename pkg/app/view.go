package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/couragetogroww/powerhell/pkg/ui"
)

// Styles
var (
	promptStyle = lipgloss.NewStyle().Foreground(orange).PaddingTop(1)
	menuTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(orange).PaddingTop(1).PaddingBottom(1)
	orangeTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(orange).PaddingTop(1).PaddingBottom(1)
	menuOptionStyle = lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(lipgloss.Color("#A0A0A0"))
	selectedMenuOptionStyle = lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(lipgloss.Color("#fb923c")).Bold(true)
	menuContainerStyle = lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(orange).
		Width(40) // Fixed width for compact look
	
	// Account creation styles
	inputFieldStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	focusedInputFieldStyle = inputFieldStyle.Copy().Foreground(orange)
	containerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4A4A4A")).Padding(1, 2).Width(40)
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).PaddingBottom(1)
	
	// Module explorer styles
	powerhellTabStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(orange).
		Padding(0, 1)
	
	sidebarStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#505050"))
	
	contentAreaStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#505050"))
	
	footerStyle = lipgloss.NewStyle().
		MarginTop(1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("#3A3A3A"))
)

// View renders the current state
func (m Model) View() string {
	if m.Quit {
		return "Exiting PowerHell... Goodbye!\n"
	}

	var mainView string
	switch m.AppState {
	case StateIntro:
		mainView = m.renderIntro()
	case StateAccountCreation:
		mainView = m.renderAccountCreation()
	case StateSignIn:
		if m.SignInView != nil {
			mainView = m.SignInView.Render()
		} else {
			mainView = "Loading sign-in..."
		}
	case StateSignInPlaceholder:
		mainView = m.renderModuleExplorer()
	case StateAuthMenu:
		mainView = m.renderMenu("Authentication", "Choose an option to continue")
	case StateModuleExplorer:
		mainView = m.renderModuleExplorer()
	case StateDashboard:
		if m.Dashboard != nil {
			mainView = m.Dashboard.Render()
		} else {
			mainView = "Loading dashboard..."
		}
	case StateLesson:
		if m.LessonView != nil {
			mainView = m.LessonView.Render()
		} else {
			mainView = "Loading lesson..."
		}
	case StateMainMenu, StateLearnMenu, StateStudio, StateSettings:
		mainView = m.renderModuleExplorer()
	default:
		mainView = "Unknown state."
	}

	// Overlay help if needed
	if m.ShowHelp {
		var context string
		switch m.AppState {
		case StateIntro:
			context = "intro"
		case StateDashboard:
			context = "dashboard"
		case StateLesson:
			context = "lesson"
		case StateAuthMenu, StateMainMenu, StateLearnMenu:
			context = "menu"
		default:
			context = "general"
		}
		return ui.HelpOverlay(m.TerminalWidth, m.TerminalHeight, context)
	}

	return mainView
}

func (m Model) renderIntro() string {
	// Use the original intro with animated flames
	flameBlockLines := make([]string, len(m.Flames))
	for i, line := range m.Flames {
		flameBlockLines[i] = colorizeFlames(line)
	}

	// Construct the flames art block
	flameBlock := lipgloss.JoinVertical(lipgloss.Left, flameBlockLines...)

	// Center the flame block
	centeredFlameBlock := lipgloss.Place(
		m.TerminalWidth,
		len(flameBlockLines),
		lipgloss.Center,
		lipgloss.Top,
		flameBlock,
	)

	// Center the "POWERHELL" title below the flames
	centeredPowerhellTitle := lipgloss.Place(
		m.TerminalWidth,
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
		m.TerminalWidth,
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

	mainView := lipgloss.Place(
		m.TerminalWidth,
		m.TerminalHeight,
		lipgloss.Center,
		lipgloss.Top,
		content,
	)

	return mainView
}

func (m Model) renderAccountCreation() string {
	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(orange).
		Align(lipgloss.Center).
		Width(32).
		Render("Create Your Account")

	if m.FocusedField == FocusDisplayInfo {
		// Success screen
		m.ShowAccountWarning = true
		
		// Account number display
		accountText := fmt.Sprintf("Account Number:\n%s", m.GeneratedAccountNumber)
		accountStyle := lipgloss.NewStyle().
			Foreground(orange).
			Bold(true).
			Align(lipgloss.Center).
			Width(32)
		
		// Animated warning
		warningColors := []lipgloss.Color{red, orange, yellow}
		warningColor := warningColors[(m.AnimationFrame/5)%len(warningColors)]
		
		warningText := lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Align(lipgloss.Center).
			Width(32).
			Render("⚠️  SAVE THIS NUMBER  ⚠️\nYou'll need it to login")
		
		continueText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			Align(lipgloss.Center).
			Width(32).
			Render("Press Enter to continue")
		
		content := lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			"",
			"✅ Account Created!",
			"",
			accountStyle.Render(accountText),
			"",
			warningText,
			"",
			continueText,
		)

		container := lipgloss.NewStyle().
			Padding(2, 4).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(orange).
			Width(40).
			Render(content)

		return lipgloss.Place(
			m.TerminalWidth,
			m.TerminalHeight,
			lipgloss.Center,
			lipgloss.Center,
			container,
		)
	}

	// Input fields
	nameField := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(func() lipgloss.Color {
			if m.FocusedField == FocusName {
				return orange
			}
			return lipgloss.Color("#4A4A4A")
		}()).
		Padding(0, 1).
		Width(30).
		Render(m.NameInput.View())

	emailField := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(func() lipgloss.Color {
			if m.FocusedField == FocusEmail {
				return orange
			}
			return lipgloss.Color("#4A4A4A")
		}()).
		Padding(0, 1).
		Width(30).
		Render(m.EmailInput.View())

	// Labels
	nameLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Render("Name")
	
	emailLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Render("Email")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true).
		Align(lipgloss.Center).
		Width(32).
		Render("Tab: Next field • Enter: Submit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		nameLabel,
		nameField,
		"",
		emailLabel,
		emailField,
		"",
		instructions,
	)

	container := lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(orange).
		Width(40).
		Render(content)

	return lipgloss.Place(
		m.TerminalWidth,
		m.TerminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		container,
	)
}

func (m Model) renderMenu(title string, promptText string) string {
	// Title with centered alignment
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(orange).
		Align(lipgloss.Center).
		Width(32) // Match content width
	
	titleRendered := titleStyle.Render(title)

	// Menu options
	var options []string
	menuOptions := m.MenuManager.GetMenuOptionsAsStrings()
	for i, option := range menuOptions {
		var optionText string
		if m.MenuCursor == i {
			// Selected option with arrow
			optionText = lipgloss.NewStyle().
				Foreground(orange).
				Bold(true).
				Render("▶ " + option)
		} else {
			// Unselected option with spacing
			optionText = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Render("  " + option)
		}
		options = append(options, optionText)
	}

	// Join options with minimal spacing
	optionsBlock := lipgloss.JoinVertical(lipgloss.Left, options...)

	// Prompt text
	promptRendered := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true).
		MarginTop(1).
		Align(lipgloss.Center).
		Width(32).
		Render(promptText)

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleRendered,
		"", // Single empty line
		optionsBlock,
		promptRendered,
	)

	// Apply compact container style
	styledContent := menuContainerStyle.Render(content)

	// Center the entire menu
	return lipgloss.Place(
		m.TerminalWidth,
		m.TerminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		styledContent,
	)
}

func (m Model) renderModuleExplorer() string {
	footerHint := "↑/↓: Navigate • Enter: Select • Q: Quit"
	return m.renderSharedExplorerLayout(
		m.ModuleExplorerSidebarOptions,
		m.ModuleExplorerSidebarCursor,
		m.ModuleExplorerContent,
		footerHint,
	)
}

func (m Model) renderSharedExplorerLayout(sidebarOptions []string, sidebarCursor int, mainContent string, footerHint string) string {
	// Header
	powerhellTab := powerhellTabStyle.Render("PowerHell")
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
	footerHeight := 3
	availableHeight := m.TerminalHeight - headerHeight - footerHeight
	if availableHeight < 0 {
		availableHeight = 0
	}

	sidebarWidth := int(float64(m.TerminalWidth) * 0.33)
	if sidebarWidth < 20 {
		sidebarWidth = 20
	}
	contentWidth := m.TerminalWidth - sidebarWidth - 1
	if contentWidth < 0 {
		contentWidth = 0
	}

	sidebarFormatted := sidebarStyle.Width(sidebarWidth).Height(availableHeight).Render(sidebarContent.String())
	contentFormatted := contentAreaStyle.Width(contentWidth).Height(availableHeight).Render(mainContent)

	// Main body
	mainBodyView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarFormatted, " ", contentFormatted)

	// Footer
	styledFooterHint := promptStyle.Render(footerHint)
	footerView := footerStyle.Width(m.TerminalWidth).Render(lipgloss.PlaceHorizontal(m.TerminalWidth, lipgloss.Center, styledFooterHint))

	// Final layout
	return lipgloss.JoinVertical(lipgloss.Left,
		headerView,
		mainBodyView,
		footerView,
	)
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


