package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/couragetogroww/powerhell/pkg/ui"
)

// SignInView handles the sign-in process
type SignInView struct {
	AccountInput textinput.Model
	width        int
	height       int
	error        string
}

// NewSignInView creates a new sign-in view
func NewSignInView(width, height int) *SignInView {
	accountInput := textinput.New()
	accountInput.Placeholder = "1234 5678 9012 3456"
	accountInput.Focus()
	accountInput.CharLimit = 19 // 16 digits + 3 spaces
	accountInput.Width = 30
	accountInput.Prompt = ""

	return &SignInView{
		AccountInput: accountInput,
		width:        width,
		height:       height,
	}
}

// Update handles input
func (s *SignInView) Update(msg string) {
	switch msg {
	case "enter":
		// Validate account number format
		s.validateAccount()
	default:
		// Let the text input handle it
		// This would be handled in the main update loop
	}
}

// UpdateTextInput updates the text input model
func (s *SignInView) UpdateTextInput(msg interface{}) {
	var cmd interface{}
	s.AccountInput, cmd = s.AccountInput.Update(msg)
	_ = cmd
}

// Render returns the sign-in view
func (s *SignInView) Render() string {
	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.Primary).
		Align(lipgloss.Center).
		Width(32).
		Render("🔐 Login to PowerHell")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(ui.TextSecondary).
		Align(lipgloss.Center).
		Width(32).
		Render("Enter your 16-digit account number")

	// Account input field
	inputField := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Primary).
		Padding(0, 1).
		Width(30).
		Render(s.AccountInput.View())

	// Error message if any
	var errorMsg string
	if s.error != "" {
		errorMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff4444")).
			Align(lipgloss.Center).
			Width(32).
			Render(s.error)
	}

	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(ui.TextSecondary).
		Italic(true).
		Align(lipgloss.Center).
		Width(32).
		Render("Format: 1234 5678 9012 3456")

	// Combine all elements
	var elements []string
	elements = append(elements, title, "", instructions, "", inputField)
	if errorMsg != "" {
		elements = append(elements, "", errorMsg)
	}
	elements = append(elements, "", helpText)

	content := lipgloss.JoinVertical(lipgloss.Center, elements...)

	// Container
	container := lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Primary).
		Width(40).
		Render(content)

	// Navigation hints
	navHints := lipgloss.NewStyle().
		Foreground(ui.TextSecondary).
		Render("Enter: Login • Esc: Back • Ctrl+C: Quit")

	// Center everything
	fullContent := lipgloss.JoinVertical(
		lipgloss.Center,
		container,
		"",
		navHints,
	)

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		fullContent,
	)
}

// validateAccount validates the account number format
func (s *SignInView) validateAccount() bool {
	account := s.AccountInput.Value()
	// Remove spaces for validation
	cleaned := strings.ReplaceAll(account, " ", "")
	
	if len(cleaned) != 16 {
		s.error = "Account number must be exactly 16 digits"
		return false
	}
	
	// Check if all characters are digits
	for _, c := range cleaned {
		if c < '0' || c > '9' {
			s.error = "Account number must contain only digits"
			return false
		}
	}
	
	s.error = ""
	return true
}

// GetAccountNumber returns the entered account number
func (s *SignInView) GetAccountNumber() string {
	return strings.ReplaceAll(s.AccountInput.Value(), " ", "")
}

// SetError sets an error message
func (s *SignInView) SetError(err string) {
	s.error = err
}