package views

import (
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/couragetogroww/powerhell/pkg/ui"
)

// IntroView represents the animated intro screen
type IntroView struct {
	width       int
	height      int
	frame       int
	particles   []Particle
	initialized bool
}

// Particle represents a flame particle
type Particle struct {
	X, Y   float64
	VX, VY float64
	Life   float64
	Color  lipgloss.Color
}

// NewIntroView creates a new intro view
func NewIntroView(width, height int) *IntroView {
	return &IntroView{
		width:  width,
		height: height,
		frame:  0,
	}
}

// Update updates the intro animation
func (i *IntroView) Update() {
	i.frame++
	
	if !i.initialized {
		i.initParticles()
		i.initialized = true
	}
	
	// Update particles
	for j := range i.particles {
		p := &i.particles[j]
		p.Y -= p.VY
		p.X += p.VX
		p.Life -= 0.02
		
		// Reset particle if it dies
		if p.Life <= 0 {
			i.resetParticle(p)
		}
	}
}

func (i *IntroView) initParticles() {
	i.particles = make([]Particle, 100)
	for j := range i.particles {
		i.resetParticle(&i.particles[j])
	}
}

func (i *IntroView) resetParticle(p *Particle) {
	p.X = float64(i.width/2) + (rand.Float64()-0.5)*40
	p.Y = float64(i.height - 10)
	p.VX = (rand.Float64() - 0.5) * 0.5
	p.VY = rand.Float64()*1.5 + 0.5
	p.Life = rand.Float64()
	
	// Choose flame color based on height
	colors := []lipgloss.Color{ui.Accent, ui.Primary, ui.Secondary}
	p.Color = colors[rand.Intn(len(colors))]
}

// Render returns the intro view
func (i *IntroView) Render() string {
	// Create a canvas
	canvas := make([][]rune, i.height)
	for y := range canvas {
		canvas[y] = make([]rune, i.width)
		for x := range canvas[y] {
			canvas[y][x] = ' '
		}
	}
	
	// Draw particles
	for _, p := range i.particles {
		if p.Life > 0 && p.X >= 0 && p.X < float64(i.width) && p.Y >= 0 && p.Y < float64(i.height) {
			x, y := int(p.X), int(p.Y)
			intensity := p.Life
			char := '█'
			if intensity < 0.3 {
				char = '░'
			} else if intensity < 0.6 {
				char = '▒'
			} else if intensity < 0.8 {
				char = '▓'
			}
			canvas[y][x] = char
		}
	}
	
	// Convert canvas to string with colors
	var result strings.Builder
	for y, row := range canvas {
		for x, char := range row {
			if char != ' ' {
				// Find the nearest particle to color this character
				var nearestParticle *Particle
				minDist := math.MaxFloat64
				for j := range i.particles {
					p := &i.particles[j]
					dist := math.Sqrt(math.Pow(float64(x)-p.X, 2) + math.Pow(float64(y)-p.Y, 2))
					if dist < minDist && p.Life > 0 {
						minDist = dist
						nearestParticle = p
					}
				}
				if nearestParticle != nil {
					style := lipgloss.NewStyle().Foreground(nearestParticle.Color)
					result.WriteString(style.Render(string(char)))
				} else {
					result.WriteRune(char)
				}
			} else {
				result.WriteRune(' ')
			}
		}
		if y < len(canvas)-1 {
			result.WriteRune('\n')
		}
	}
	
	// Overlay title
	title := i.renderTitle()
	
	// Calculate position for centered title
	titleLines := strings.Split(title, "\n")
	titleHeight := len(titleLines)
	titleY := (i.height - titleHeight) / 3
	
	// Overlay title on canvas
	lines := strings.Split(result.String(), "\n")
	for y, titleLine := range titleLines {
		if titleY+y < len(lines) {
			titleWidth := lipgloss.Width(titleLine)
			titleX := (i.width - titleWidth) / 2
			if titleX >= 0 {
				// Replace part of the line with the title
				line := []rune(lines[titleY+y])
				titleRunes := []rune(titleLine)
				for x := 0; x < len(titleRunes) && titleX+x < len(line); x++ {
					if titleRunes[x] != ' ' {
						line[titleX+x] = titleRunes[x]
					}
				}
				lines[titleY+y] = string(line)
			}
		}
	}
	
	// Add subtitle and instructions
	subtitle := ui.SubtitleStyle.Copy().
		Foreground(ui.TextSecondary).
		MarginTop(2).
		Render("🔥 Master PowerShell Through Interactive Learning 🔥")
	
	instruction := ui.HelpStyle.Copy().
		Foreground(ui.Primary).
		Bold(true).
		Blink(i.frame%20 < 10). // Blinking effect
		MarginTop(4).
		Render("Press ENTER to begin your journey")
	
	// Combine everything
	flameCanvas := strings.Join(lines, "\n")
	
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		flameCanvas,
		"",
		"",
		subtitle,
		instruction,
	)
	
	return lipgloss.NewStyle().
		Width(i.width).
		Height(i.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

func (i *IntroView) renderTitle() string {
	// Your original POWERHELL ASCII art!
	titleArt := `
██████╗  ██████╗ ██╗    ██╗███████╗██████╗ ██╗  ██╗███████╗██╗     ██╗     
██╔══██╗██╔═══██╗██║    ██║██╔════╝██╔══██╗██║  ██║██╔════╝██║     ██║     
██████╔╝██║   ██║██║ █╗ ██║█████╗  ██████╔╝███████║█████╗  ██║     ██║     
██╔═══╝ ██║   ██║██║███╗██║██╔══╝  ██╔══██╗██╔══██║██╔══╝  ██║     ██║     
██║     ╚██████╔╝╚███╔███╔╝███████╗██║  ██║██║  ██║███████╗███████╗███████╗
╚═╝      ╚═════╝  ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝`
	
	// Apply colors line by line with gradient effect
	lines := strings.Split(titleArt, "\n")
	coloredLines := make([]string, len(lines))
	
	for i, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			coloredLines[i] = line
			continue
		}
		
		// Apply gradient colors based on position
		var color lipgloss.Color
		switch i {
		case 1, 2:
			color = ui.Secondary // Yellow for top
		case 3, 4:
			color = ui.Primary   // Orange for middle
		case 5, 6:
			color = ui.Accent    // Red for bottom
		default:
			color = ui.Primary
		}
		
		style := lipgloss.NewStyle().
			Foreground(color).
			Bold(true)
		
		coloredLines[i] = style.Render(line)
	}
	
	return strings.Join(coloredLines, "\n")
}

func (i *IntroView) getGradientColor() lipgloss.Color {
	// Cycle through flame colors
	colors := []lipgloss.Color{ui.Accent, ui.Primary, ui.Secondary}
	index := (i.frame / 10) % len(colors)
	return colors[index]
}