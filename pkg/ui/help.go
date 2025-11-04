package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpOverlay creates a help overlay showing keyboard shortcuts
func HelpOverlay(width, height int, context string) string {
	// Define help content based on context
	var helpContent string
	
	switch context {
	case "intro":
		helpContent = renderIntroHelp()
	case "dashboard":
		helpContent = renderDashboardHelp()
	case "lesson":
		helpContent = renderLessonHelp()
	case "menu":
		helpContent = renderMenuHelp()
	default:
		helpContent = renderGeneralHelp()
	}
	
	// Create the help box
	helpBox := CardStyle.Copy().
		Width(60).
		BorderForeground(Primary).
		Background(Surface).
		Render(helpContent)
	
	// Center it on screen
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		helpBox,
	)
}

func renderIntroHelp() string {
	title := TitleStyle.Render("🔥 PowerHell Help")
	
	shortcuts := [][]string{
		{"Enter", "Start the application"},
		{"q, Ctrl+C", "Quit"},
		{"h", "Toggle this help menu"},
	}
	
	return formatHelp(title, shortcuts)
}

func renderDashboardHelp() string {
	title := TitleStyle.Render("📚 Dashboard Navigation")
	
	shortcuts := [][]string{
		{"↑/↓/←/→", "Navigate between modules"},
		{"j/k/h/l", "Vim-style navigation"},
		{"Enter", "Select module"},
		{"Tab", "Switch to categories view"},
		{"q", "Quit to main menu"},
		{"Ctrl+C", "Exit application"},
		{"h", "Toggle this help menu"},
	}
	
	return formatHelp(title, shortcuts)
}

func renderLessonHelp() string {
	title := TitleStyle.Render("📖 Lesson Navigation")
	
	shortcuts := [][]string{
		{"Tab", "Switch between Lesson/Code/Output tabs"},
		{"Shift+Tab", "Switch tabs backwards"},
		{"n", "Next lesson"},
		{"p", "Previous lesson"},
		{"?", "Show/hide exercise hints"},
		{"r", "Run code (when connected)"},
		{"q", "Back to dashboard"},
		{"h", "Toggle this help menu"},
	}
	
	return formatHelp(title, shortcuts)
}

func renderMenuHelp() string {
	title := TitleStyle.Render("📋 Menu Navigation")
	
	shortcuts := [][]string{
		{"↑/↓", "Navigate menu items"},
		{"j/k", "Vim-style navigation"},
		{"Enter", "Select item"},
		{"Esc", "Go back"},
		{"q", "Quit"},
		{"h", "Toggle this help menu"},
	}
	
	return formatHelp(title, shortcuts)
}

func renderGeneralHelp() string {
	title := TitleStyle.Render("⌨️  General Navigation")
	
	shortcuts := [][]string{
		{"↑/↓/←/→", "Navigate"},
		{"j/k/h/l", "Vim-style navigation"},
		{"Enter", "Select/Confirm"},
		{"Esc", "Go back"},
		{"Tab", "Next field/section"},
		{"q", "Quit current view"},
		{"Ctrl+C", "Exit application"},
		{"h", "Toggle this help menu"},
	}
	
	return formatHelp(title, shortcuts)
}

func formatHelp(title string, shortcuts [][]string) string {
	var content strings.Builder
	
	content.WriteString(title)
	content.WriteString("\n\n")
	
	// Find the longest key for alignment
	maxKeyLen := 0
	for _, shortcut := range shortcuts {
		if len(shortcut[0]) > maxKeyLen {
			maxKeyLen = len(shortcut[0])
		}
	}
	
	// Format each shortcut
	for _, shortcut := range shortcuts {
		key := KeybindStyle.Render(padRight(shortcut[0], maxKeyLen))
		desc := HelpStyle.Render(shortcut[1])
		content.WriteString(key)
		content.WriteString("  ")
		content.WriteString(desc)
		content.WriteString("\n")
	}
	
	content.WriteString("\n")
	content.WriteString(SubtitleStyle.Render("Press 'h' again to close"))
	
	return content.String()
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}