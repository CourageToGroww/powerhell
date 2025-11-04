package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/couragetogroww/powerhell/pkg/ui"
)

// EnhancedIntroView combines the original flame art with the POWERHELL logo
type EnhancedIntroView struct {
	width      int
	height     int
	frame      int
	flameArt   string
	introView  *IntroView
}

// NewEnhancedIntroView creates a new enhanced intro view
func NewEnhancedIntroView(width, height int) *EnhancedIntroView {
	flameArt := `
                                                                               
                                                                               
                                                                               
                                                                               
                                        ..+-                                    
                                       .+#%-                                    
                                      .+@@%=        :                           
                                       *@%%%=                                   
                                    . :#%%%%%-       :*:.                       
                             .-      =%%%%%%%+   .   -%#-.                      
                             *#. :  :%%%#%%%%=  .-. :*%%#:                      
                             *%-*- :#%%%###%%+   +*-+#%%%#--:                   
                          +=:#%#%-+#%%%%###%%%==*#*#*#%%%%%#.                   
                         =%##%##%#%%%%%####%%%%%%%#%%%%%%#%#.                   
                      .--%%%%####%##################%%#%%%%#                    
                      =#%%*+**+***++++++++++*+*+++*+++*+#%%#.                   
                     .#%%+::.:.:......::..........:...::=#%#.                   
                      *%*:.:.:-::......................:-#%=                    
                     -*%=:...=*++-:......  ...  ........+#%=                    
                    .*##-:..:-+**++=:::.....    .. ....:*%%+.                   
                     *%*::....:=****=-:..... .  .......-*%#.                    
                     *%+-:... ..:=****=::...... ......:+#%+                     
                    .*#-:.........-+***+-.::.........::*#%+.                    
                    -%*:...........:=****=-:.. .. ....-#%#                      
                    *%+::.. ........:+****=:..     ...=#%+ :                    
                   :%#=:..  ......:=****+-:.... .....:+#%=:.                    
                 ..*%#-:.. ....::=**#*+-:............:*#%+.                     
                  =%%*:......:-=****+-.::............-*%%#-:                     
                  =%#+......:=++**+-...::..:::......:+#%#-                      
                  +%*-.....-+*#*+-:...:+++==+++=:..:-*##-                       
                 .#%*::....-**+-::...::++++++++=:..:=*%#:.                      
                .+%#+-:...:.....:..:......:--:.::.:-+*%%+                       
                 =%#*=----::::::::::::::.::--:::-:-+*#%#-                       
                 .+#%###*##*###**###*##***********####=  `
	
	return &EnhancedIntroView{
		width:     width,
		height:    height,
		flameArt:  flameArt,
		introView: NewIntroView(width, height),
	}
}

// Update updates the animation
func (e *EnhancedIntroView) Update() {
	e.frame++
	e.introView.Update()
}

// Render returns the enhanced intro view with original flame art and POWERHELL logo
func (e *EnhancedIntroView) Render() string {
	// Colorize the flame art
	flameLines := strings.Split(e.flameArt, "\n")
	coloredFlames := make([]string, len(flameLines))
	
	for i, line := range flameLines {
		coloredFlames[i] = e.colorizeFlames(line)
	}
	
	flameBlock := lipgloss.JoinVertical(lipgloss.Center, coloredFlames...)
	
	// Get the POWERHELL title
	title := e.introView.renderTitle()
	
	// Subtitle
	subtitle := ui.SubtitleStyle.Copy().
		Foreground(ui.TextSecondary).
		MarginTop(1).
		Render("🔥 Master PowerShell Through Interactive Learning 🔥")
	
	// Instructions with blinking effect
	instruction := ui.HelpStyle.Copy().
		Foreground(ui.Primary).
		Bold(true).
		Blink(e.frame%20 < 10).
		MarginTop(2).
		Render("Press ENTER to begin your journey")
	
	// Combine everything
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		flameBlock,
		"",
		title,
		"",
		subtitle,
		instruction,
	)
	
	return lipgloss.NewStyle().
		Width(e.width).
		Height(e.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

func (e *EnhancedIntroView) colorizeFlames(line string) string {
	if strings.TrimSpace(line) == "" {
		return line
	}
	
	// Apply fire colors based on character density
	var result strings.Builder
	for _, char := range line {
		var color lipgloss.Color
		switch char {
		case '#', '@', '%':
			color = ui.Accent // Red for dense parts
		case '*', '+', '=':
			color = ui.Primary // Orange for medium density
		case '-', ':', '.':
			color = ui.Secondary // Yellow for light parts
		default:
			result.WriteRune(char)
			continue
		}
		
		style := lipgloss.NewStyle().Foreground(color)
		result.WriteString(style.Render(string(char)))
	}
	
	return result.String()
}