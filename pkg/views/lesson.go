package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/couragetogroww/powerhell/pkg/modules"
	"github.com/couragetogroww/powerhell/pkg/ui"
)

// LessonView represents an interactive lesson view
type LessonView struct {
	module        *modules.Module
	currentLesson int
	lesson        *modules.Lesson
	width         int
	height        int
	showHints     bool
	currentHint   int
	userCode      string
	outputBuffer  string
	isRunning     bool
	activeTab     int // 0: lesson, 1: code editor, 2: output
}

// NewLessonView creates a new lesson view
func NewLessonView(module *modules.Module, width, height int) *LessonView {
	return &LessonView{
		module:        module,
		currentLesson: 0,
		lesson:        &module.Lessons[0],
		width:         width,
		height:        height,
		activeTab:     0,
	}
}

// Update handles input for the lesson view
func (l *LessonView) Update(key string) {
	switch key {
	case "tab":
		l.activeTab = (l.activeTab + 1) % 3
	case "shift+tab":
		l.activeTab--
		if l.activeTab < 0 {
			l.activeTab = 2
		}
	case "?":
		l.showHints = !l.showHints
	case "n":
		if l.currentLesson < len(l.module.Lessons)-1 {
			l.currentLesson++
			l.lesson = &l.module.Lessons[l.currentLesson]
			l.showHints = false
			l.currentHint = 0
		}
	case "p":
		if l.currentLesson > 0 {
			l.currentLesson--
			l.lesson = &l.module.Lessons[l.currentLesson]
			l.showHints = false
			l.currentHint = 0
		}
	case "r":
		// Run code (simulated)
		l.isRunning = true
		l.outputBuffer = "PS C:\\> " + l.userCode + "\n\n" + 
			"# Output will appear here when connected to PowerShell runtime\n" +
			"# This is a simulation for now"
		l.isRunning = false
	}
}

// Render returns the lesson view
func (l *LessonView) Render() string {
	// Header with module info
	header := l.renderHeader()

	// Progress bar
	progress := l.renderProgress()

	// Tab bar
	tabs := ui.TabBar([]string{"📖 Lesson", "💻 Code Editor", "📊 Output"}, l.activeTab)

	// Content based on active tab
	var content string
	switch l.activeTab {
	case 0:
		content = l.renderLessonContent()
	case 1:
		content = l.renderCodeEditor()
	case 2:
		content = l.renderOutput()
	}

	// Help bar
	helpBar := ui.HelpBar([][2]string{
		{"Tab", "Switch Tabs"},
		{"n/p", "Next/Prev Lesson"},
		{"?", "Show Hints"},
		{"r", "Run Code"},
		{"q", "Back to Dashboard"},
	})

	// Combine all elements
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		progress,
		"",
		tabs,
		content,
	)

	// Apply main container styling
	mainStyled := lipgloss.NewStyle().
		Width(l.width).
		Height(l.height - 3).
		Padding(1, 2).
		Render(mainContent)

	helpBarStyled := lipgloss.NewStyle().
		Width(l.width).
		Align(lipgloss.Center).
		Background(ui.Surface).
		Padding(0, 2).
		Render(helpBar)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainStyled,
		helpBarStyled,
	)
}

func (l *LessonView) renderHeader() string {
	title := fmt.Sprintf("%s %s", l.module.Icon, l.module.Title)
	subtitle := fmt.Sprintf("Lesson %d of %d: %s", 
		l.currentLesson+1, 
		len(l.module.Lessons), 
		l.lesson.Title,
	)
	
	return ui.Header(title, subtitle)
}

func (l *LessonView) renderProgress() string {
	lessonProgress := float64(l.currentLesson) / float64(len(l.module.Lessons))
	return ui.ProgressBar(l.width-4, lessonProgress)
}

func (l *LessonView) renderLessonContent() string {
	// Lesson description
	descStyle := lipgloss.NewStyle().
		Foreground(ui.TextPrimary).
		MarginBottom(2)
	
	description := descStyle.Render(l.lesson.Description)

	// Lesson content (would be loaded from lesson data)
	content := l.renderLessonText()

	// Code example
	codeExample := ""
	if l.lesson.CodeExample != "" {
		codeExample = ui.CodeBlock(l.lesson.CodeExample, "PowerShell")
	}

	// Exercise section
	exercise := l.renderExercise()

	// Hints section (if enabled)
	hints := ""
	if l.showHints && l.lesson.Exercise.Hints != nil {
		hints = l.renderHints()
	}

	sections := []string{description, content}
	if codeExample != "" {
		sections = append(sections, codeExample)
	}
	if exercise != "" {
		sections = append(sections, exercise)
	}
	if hints != "" {
		sections = append(sections, hints)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)
}

func (l *LessonView) renderLessonText() string {
	// Get dynamic content from the lesson
	lessonText := l.lesson.GetContent(l.module.ID)

	return lipgloss.NewStyle().
		Foreground(ui.TextPrimary).
		Render(lessonText)
}

func (l *LessonView) renderExercise() string {
	exercise := l.lesson.GetExercise(l.module.ID)
	if exercise.Instructions == "" {
		return ""
	}

	exerciseBox := ui.CardStyle.Copy().
		BorderForeground(ui.Info).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				ui.InfoIndicatorStyle.Render("📝 Exercise"),
				"",
				exercise.Instructions,
			),
		)

	return exerciseBox
}

func (l *LessonView) renderHints() string {
	exercise := l.lesson.GetExercise(l.module.ID)
	if len(exercise.Hints) == 0 {
		return ""
	}

	hint := exercise.Hints[l.currentHint%len(exercise.Hints)]
	
	hintBox := ui.CardStyle.Copy().
		BorderForeground(ui.Secondary).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				ui.KeybindStyle.Render("💡 Hint"),
				"",
				hint,
			),
		)

	return hintBox
}

func (l *LessonView) renderCodeEditor() string {
	editorStyle := lipgloss.NewStyle().
		Width(l.width - 4).
		Height(l.height - 15).
		Background(lipgloss.Color("#0d0d0d")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Primary).
		Padding(1)

	// Initialize with starter code from exercise
	if l.userCode == "" {
		exercise := l.lesson.GetExercise(l.module.ID)
		if exercise.StarterCode != "" {
			l.userCode = exercise.StarterCode
		} else {
			l.userCode = "# Type your PowerShell code here\n" +
				"$greeting = \"Hello, PowerHell!\"\n" +
				"Write-Host $greeting"
		}
	}

	// Line numbers
	lines := strings.Split(l.userCode, "\n")
	var numberedLines []string
	for i, line := range lines {
		lineNum := lipgloss.NewStyle().
			Foreground(ui.TextSecondary).
			Width(4).
			Align(lipgloss.Right).
			Render(fmt.Sprintf("%d ", i+1))
		
		numberedLines = append(numberedLines, lineNum + line)
	}

	content := strings.Join(numberedLines, "\n")

	// Status bar
	status := ""
	if l.isRunning {
		status = ui.LoadingSpinner(0) + " Running..."
	} else {
		status = ui.SuccessIndicatorStyle.Render("Ready")
	}

	statusBar := lipgloss.NewStyle().
		Width(l.width - 4).
		Background(ui.Surface).
		Padding(0, 1).
		Render(status + " | Press 'r' to run")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		editorStyle.Render(content),
		statusBar,
	)
}

func (l *LessonView) renderOutput() string {
	outputStyle := lipgloss.NewStyle().
		Width(l.width - 4).
		Height(l.height - 15).
		Background(lipgloss.Color("#0d0d0d")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Border).
		Padding(1)

	if l.outputBuffer == "" {
		l.outputBuffer = "PS C:\\> # Output will appear here when you run your code"
	}

	return outputStyle.Render(l.outputBuffer)
}