package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	Primary       = lipgloss.Color("#fb923c") // Orange
	Secondary     = lipgloss.Color("#fcd34d") // Yellow
	Accent        = lipgloss.Color("#ff4e00") // Red
	Background    = lipgloss.Color("#1a1a1a") // Dark background
	Surface       = lipgloss.Color("#2d2d2d") // Card/surface background
	TextPrimary   = lipgloss.Color("#f3f4f6") // Light text
	TextSecondary = lipgloss.Color("#9ca3af") // Muted text
	Success       = lipgloss.Color("#10b981") // Green
	Error         = lipgloss.Color("#ef4444") // Red
	Info          = lipgloss.Color("#3b82f6") // Blue
	Border        = lipgloss.Color("#4b5563") // Border color
)

// Title styles
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Primary).
	PaddingTop(1).
	PaddingBottom(1).
	Align(lipgloss.Center)

var SubtitleStyle = lipgloss.NewStyle().
	Foreground(TextSecondary).
	Align(lipgloss.Center)

// Card/Container styles
var CardStyle = lipgloss.NewStyle().
	Background(Surface).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Border).
	Padding(1, 2)

var FocusedCardStyle = CardStyle.Copy().
	BorderForeground(Primary)

// Button styles
var ButtonStyle = lipgloss.NewStyle().
	Foreground(TextPrimary).
	Background(Surface).
	Padding(0, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Border)

var ActiveButtonStyle = ButtonStyle.Copy().
	Foreground(Background).
	Background(Primary).
	BorderForeground(Primary).
	Bold(true)

// Progress bar styles
var ProgressBarStyle = lipgloss.NewStyle().
	Foreground(Primary).
	Background(Surface)

var ProgressBarContainerStyle = lipgloss.NewStyle().
	Background(Surface).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Border).
	Padding(0, 1)

// Code block styles
var CodeBlockStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#0d0d0d")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Border).
	Padding(1, 2)

var CodeHighlightStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#3a3a3a"))

// Status indicator styles
var SuccessIndicatorStyle = lipgloss.NewStyle().
	Foreground(Success).
	Bold(true)

var ErrorIndicatorStyle = lipgloss.NewStyle().
	Foreground(Error).
	Bold(true)

var InfoIndicatorStyle = lipgloss.NewStyle().
	Foreground(Info).
	Bold(true)

// Tab styles
var TabStyle = lipgloss.NewStyle().
	Foreground(TextSecondary).
	Background(Surface).
	Padding(0, 2).
	Border(lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  " ",
		BottomRight: " ",
	}).
	BorderForeground(Border)

var ActiveTabStyle = TabStyle.Copy().
	Foreground(TextPrimary).
	Background(Primary).
	BorderForeground(Primary).
	Bold(true)

// Help text styles
var HelpStyle = lipgloss.NewStyle().
	Foreground(TextSecondary).
	Italic(true)

var KeybindStyle = lipgloss.NewStyle().
	Foreground(Primary).
	Bold(true)

// Module card styles
var ModuleCardStyle = lipgloss.NewStyle().
	Width(40).
	Height(10).
	Background(Surface).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Border).
	Padding(1, 2)

var ModuleCardHoverStyle = ModuleCardStyle.Copy().
	BorderForeground(Primary).
	Background(lipgloss.Color("#3a3a3a"))

// Badge styles
var BadgeStyle = lipgloss.NewStyle().
	Foreground(Primary).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Primary).
	Padding(0, 1).
	Bold(true)

var CompletedBadgeStyle = lipgloss.NewStyle().
	Foreground(Success).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Success).
	Padding(0, 1).
	Bold(true)

var InProgressBadgeStyle = lipgloss.NewStyle().
	Foreground(Info).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Info).
	Padding(0, 1).
	Bold(true)

// Layout helpers
func CenterHorizontal(width int, content string) string {
	contentWidth := lipgloss.Width(content)
	if contentWidth >= width {
		return content
	}
	padding := (width - contentWidth) / 2
	return lipgloss.NewStyle().PaddingLeft(padding).Render(content)
}

func CenterVertical(height int, content string) string {
	contentHeight := lipgloss.Height(content)
	if contentHeight >= height {
		return content
	}
	padding := (height - contentHeight) / 2
	return lipgloss.NewStyle().PaddingTop(padding).Render(content)
}