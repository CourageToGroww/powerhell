package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBar creates a visual progress bar
func ProgressBar(width int, percent float64) string {
	if percent > 1.0 {
		percent = 1.0
	}
	if percent < 0 {
		percent = 0
	}

	filled := int(float64(width) * percent)
	if filled > width {
		filled = width
	}

	filledBar := lipgloss.NewStyle().Foreground(Primary).Render(strings.Repeat("█", filled))
	emptyBar := lipgloss.NewStyle().Foreground(Border).Render(strings.Repeat("░", width-filled))
	percentText := lipgloss.NewStyle().Foreground(TextSecondary).Render(fmt.Sprintf(" %3.0f%%", percent*100))
	
	return filledBar + emptyBar + percentText
}

// ModuleCard creates a card for displaying module information
func ModuleCard(title, description string, progress float64, isActive bool) string {
	// Use different border color for active
	borderColor := Border
	if isActive {
		borderColor = Primary
	}

	// Title
	titleRendered := lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true).
		MarginBottom(1).
		Render(title)

	// Description
	descRendered := lipgloss.NewStyle().
		Foreground(TextPrimary).
		MarginBottom(1).
		Render(description)

	// Progress
	progressText := fmt.Sprintf("Progress: %.0f%%", progress*100)
	progressRendered := ProgressBar(30, progress)

	// Status badge
	var badge string
	if progress >= 1.0 {
		badge = CompletedBadgeStyle.Render("✓ Completed")
	} else if progress > 0 {
		badge = InProgressBadgeStyle.Render("In Progress")
	} else {
		badge = BadgeStyle.Render("Not Started")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleRendered,
		descRendered,
		"",
		progressText,
		progressRendered,
		"",
		badge,
	)

	// Create a custom style without dark background
	cardStyle := lipgloss.NewStyle().
		Width(40).
		Height(10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2)
	
	if isActive {
		// Add a subtle indicator for active card
		cardStyle = cardStyle.BorderForeground(Primary).Bold(true)
	}

	return cardStyle.Render(content)
}

// Tab creates a tab component
func Tab(label string, isActive bool) string {
	if isActive {
		return ActiveTabStyle.Render(label)
	}
	return TabStyle.Render(label)
}

// TabBar creates a horizontal tab bar
func TabBar(tabs []string, activeIndex int) string {
	var renderedTabs []string
	for i, tab := range tabs {
		renderedTabs = append(renderedTabs, Tab(tab, i == activeIndex))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// KeyBinding displays a key binding help text
func KeyBinding(key, description string) string {
	k := KeybindStyle.Render(key)
	d := HelpStyle.Render(description)
	return fmt.Sprintf("%s %s", k, d)
}

// HelpBar creates a help bar with key bindings
func HelpBar(bindings [][2]string) string {
	var helps []string
	for _, binding := range bindings {
		helps = append(helps, KeyBinding(binding[0], binding[1]))
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		strings.Join(helps, "  "),
	)
}

// CodeBlock renders a code block with optional syntax highlighting
func CodeBlock(code string, language string) string {
	// For now, simple rendering. Can be enhanced with actual syntax highlighting
	header := lipgloss.NewStyle().
		Foreground(TextSecondary).
		Render(fmt.Sprintf("// %s", language))
	
	return CodeBlockStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, header, "", code),
	)
}

// StatusMessage renders a status message with appropriate styling
func StatusMessage(message string, status string) string {
	var style lipgloss.Style
	var icon string
	
	switch status {
	case "success":
		style = SuccessIndicatorStyle
		icon = "✓"
	case "error":
		style = ErrorIndicatorStyle
		icon = "✗"
	case "info":
		style = InfoIndicatorStyle
		icon = "ℹ"
	default:
		style = lipgloss.NewStyle().Foreground(TextPrimary)
		icon = "•"
	}
	
	return style.Render(fmt.Sprintf("%s %s", icon, message))
}

// SplitView creates a split view layout
func SplitView(left, right string, leftWidth int) string {
	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(lipgloss.Height(right)).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(Border).
		PaddingRight(1)
	
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		right,
	)
}

// Header creates a styled header with title and subtitle
func Header(title, subtitle string) string {
	t := TitleStyle.Render(title)
	s := SubtitleStyle.Render(subtitle)
	
	return lipgloss.JoinVertical(lipgloss.Center, t, s)
}

// HeaderWithUser creates a styled header with title, subtitle, and user info
func HeaderWithUser(title, subtitle, userName string) string {
	// Create the main title and subtitle
	t := TitleStyle.Render(title)
	s := SubtitleStyle.Render(subtitle)
	
	// Create the user info that will be positioned on the right
	userInfo := ""
	if userName != "" {
		userInfo = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Render("👤 " + userName)
	}
	
	// Get the terminal width (we'll use a reasonable default if not available)
	termWidth := 80 // This should be passed from the view
	
	// Create the header content
	headerContent := lipgloss.JoinVertical(lipgloss.Center, t, s)
	
	// If we have user info, create a layout with it on the right
	if userInfo != "" {
		// Calculate the width for proper spacing
		headerWidth := lipgloss.Width(headerContent)
		userWidth := lipgloss.Width(userInfo)
		spacerWidth := termWidth - headerWidth - userWidth - 4 // 4 for padding
		
		if spacerWidth < 0 {
			spacerWidth = 1
		}
		
		// Create the full header with user info on the right
		topLine := lipgloss.JoinHorizontal(
			lipgloss.Top,
			strings.Repeat(" ", spacerWidth/2),
			headerContent,
			strings.Repeat(" ", spacerWidth/2),
			userInfo,
		)
		
		return topLine
	}
	
	return headerContent
}

// LoadingSpinner creates an animated loading indicator
func LoadingSpinner(frame int) string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return lipgloss.NewStyle().
		Foreground(Primary).
		Render(spinners[frame%len(spinners)])
}