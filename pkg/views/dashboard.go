package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/couragetogroww/powerhell/pkg/modules"
	"github.com/couragetogroww/powerhell/pkg/ui"
)

// DashboardView represents the main dashboard
type DashboardView struct {
	modules        []modules.Module
	selectedModule int
	width          int
	height         int
	UserName       string // Made public for access from update.go
}

// NewDashboardView creates a new dashboard view
func NewDashboardView(width, height int) *DashboardView {
	return &DashboardView{
		modules: modules.GetAvailableModules(),
		width:   width,
		height:  height,
	}
}

// NewDashboardViewWithUser creates a new dashboard view with user info
func NewDashboardViewWithUser(width, height int, userName string) *DashboardView {
	return &DashboardView{
		modules:  modules.GetAvailableModules(),
		width:    width,
		height:   height,
		UserName: userName,
	}
}

// Update handles input for the dashboard
func (d *DashboardView) Update(key string) {
	switch key {
	case "left", "h":
		if d.selectedModule > 0 {
			d.selectedModule--
		}
	case "right", "l":
		if d.selectedModule < len(d.modules)-1 {
			d.selectedModule++
		}
	case "up", "k":
		if d.selectedModule >= 3 {
			d.selectedModule -= 3
		}
	case "down", "j":
		if d.selectedModule < len(d.modules)-3 {
			d.selectedModule += 3
		}
	}
}

// Render returns the dashboard view
func (d *DashboardView) Render() string {
	// Header - use HeaderWithUser if we have a user name
	var header string
	if d.UserName != "" {
		// Create a simpler header layout with user info on the right
		title := ui.TitleStyle.Render("🔥 PowerHell Learning Platform")
		subtitle := ui.SubtitleStyle.Render("Master PowerShell through interactive lessons and real-world examples")
		userInfo := lipgloss.NewStyle().
			Foreground(ui.Primary).
			Bold(true).
			Render("👤 " + d.UserName)
		
		// Create header content
		headerContent := lipgloss.JoinVertical(lipgloss.Center, title, subtitle)
		
		// Calculate spacing
		headerWidth := lipgloss.Width(headerContent)
		userWidth := lipgloss.Width(userInfo)
		spacerWidth := d.width - headerWidth - userWidth - 6
		
		if spacerWidth < 0 {
			// If not enough space, put user on next line
			header = lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top, headerContent, strings.Repeat(" ", 4), userInfo),
			)
		} else {
			// Put user info on the right
			header = lipgloss.JoinHorizontal(
				lipgloss.Top,
				headerContent,
				strings.Repeat(" ", spacerWidth),
				userInfo,
			)
		}
	} else {
		header = ui.Header(
			"🔥 PowerHell Learning Platform",
			"Master PowerShell through interactive lessons and real-world examples",
		)
	}

	// Stats bar
	stats := d.renderStats()

	// Module grid
	moduleGrid := d.renderModuleGrid()

	// Help bar
	helpBar := ui.HelpBar([][2]string{
		{"↑↓←→", "Navigate"},
		{"Enter", "Select Module"},
		{"Tab", "Categories"},
		{"?", "Help"},
		{"q", "Quit"},
	})

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		stats,
		"",
		moduleGrid,
	)

	// Center content and add help bar at bottom
	mainContent := lipgloss.NewStyle().
		Width(d.width).
		Height(d.height - 3).
		Align(lipgloss.Center, lipgloss.Top).
		Render(content)

	helpBarStyled := lipgloss.NewStyle().
		Width(d.width).
		Align(lipgloss.Center).
		BorderTop(true).
		BorderForeground(ui.Border).
		Padding(0, 2).
		Render(helpBar)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		helpBarStyled,
	)
}

func (d *DashboardView) renderStats() string {
	// Calculate overall progress
	totalProgress := 0.0
	completedModules := 0
	for _, module := range d.modules {
		totalProgress += module.Progress
		if module.IsCompleted {
			completedModules++
		}
	}
	avgProgress := totalProgress / float64(len(d.modules))

	// Create stat cards
	stats := []string{
		d.createStatCard("📊 Overall Progress", fmt.Sprintf("%.0f%%", avgProgress*100), ui.Primary),
		d.createStatCard("✅ Completed", fmt.Sprintf("%d/%d modules", completedModules, len(d.modules)), ui.Success),
		d.createStatCard("🔥 Current Streak", "5 days", ui.Accent),
		d.createStatCard("⭐ Points", "1,250", ui.Secondary),
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		stats...,
	)
}

func (d *DashboardView) createStatCard(title, value string, color lipgloss.Color) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.TextSecondary).
		MarginBottom(1)

	valueStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(title),
		valueStyle.Render(value),
	)

	return lipgloss.NewStyle().
		Width(20).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Primary).
		Margin(0, 1).
		Render(content)
}

func (d *DashboardView) renderModuleGrid() string {
	// Create rows of 3 modules each
	var rows []string
	for i := 0; i < len(d.modules); i += 3 {
		var row []string
		for j := 0; j < 3 && i+j < len(d.modules); j++ {
			idx := i + j
			module := d.modules[idx]
			isSelected := idx == d.selectedModule
			
			card := ui.ModuleCard(
				fmt.Sprintf("%s %s", module.Icon, module.Title),
				module.Description,
				module.Progress,
				isSelected,
			)
			
			row = append(row, card)
		}
		
		// Fill empty slots if needed
		for len(row) < 3 {
			row = append(row, strings.Repeat(" ", 40))
		}
		
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// GetSelectedModule returns the currently selected module
func (d *DashboardView) GetSelectedModule() *modules.Module {
	if d.selectedModule >= 0 && d.selectedModule < len(d.modules) {
		return &d.modules[d.selectedModule]
	}
	return nil
}