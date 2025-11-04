package modules

import (
	"time"
)

// Module represents a learning module
type Module struct {
	ID          string
	Title       string
	Description string
	Category    string
	Duration    time.Duration
	Difficulty  string // "Beginner", "Intermediate", "Advanced"
	Lessons     []Lesson
	Progress    float64 // 0.0 to 1.0
	IsCompleted bool
	Icon        string // Unicode icon for visual representation
}

// Lesson represents a single lesson within a module
type Lesson struct {
	ID          string
	Title       string
	Description string
	Content     string
	CodeExample string
	Exercise    Exercise
	IsCompleted bool
	Duration    time.Duration
}

// GetContent returns the full content for this lesson
func (l *Lesson) GetContent(moduleID string) string {
	if l.Content == "" {
		l.Content = GetLessonContent(moduleID, l.ID)
	}
	return l.Content
}

// GetExercise returns the exercise for this lesson
func (l *Lesson) GetExercise(moduleID string) Exercise {
	if l.Exercise.ID == "" {
		l.Exercise = GetExerciseForLesson(moduleID, l.ID)
	}
	return l.Exercise
}

// Exercise represents an interactive exercise
type Exercise struct {
	ID           string
	Instructions string
	StarterCode  string
	Solution     string
	Hints        []string
	TestCases    []TestCase
}

// TestCase represents a test case for an exercise
type TestCase struct {
	Input    string
	Expected string
}

// UserProgress tracks user progress through modules
type UserProgress struct {
	UserID           string
	CompletedLessons map[string]bool
	ModuleProgress   map[string]float64
	LastAccessed     map[string]time.Time
	TotalTimeSpent   time.Duration
	CurrentStreak    int
	LongestStreak    int
	Points           int
	Achievements     []Achievement
}

// Achievement represents a user achievement
type Achievement struct {
	ID          string
	Title       string
	Description string
	Icon        string
	UnlockedAt  time.Time
	Points      int
}

// ModuleCategory represents a category of modules
type ModuleCategory struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Modules     []Module
}

// GetAvailableModules returns all available learning modules
func GetAvailableModules() []Module {
	return []Module{
		{
			ID:          "basics",
			Title:       "PowerShell Basics",
			Description: "Learn the fundamentals of PowerShell scripting",
			Category:    "Foundation",
			Duration:    2 * time.Hour,
			Difficulty:  "Beginner",
			Icon:        "📚",
			Lessons: []Lesson{
				{
					ID:          "basics-1",
					Title:       "Introduction to PowerShell",
					Description: "Understanding what PowerShell is and why it's powerful",
					Duration:    15 * time.Minute,
				},
				{
					ID:          "basics-2",
					Title:       "Variables and Data Types",
					Description: "Working with variables, strings, numbers, and arrays",
					Duration:    20 * time.Minute,
				},
				{
					ID:          "basics-3",
					Title:       "Basic Cmdlets",
					Description: "Essential PowerShell commands you need to know",
					Duration:    25 * time.Minute,
				},
			},
		},
		{
			ID:          "active-directory",
			Title:       "Active Directory Management",
			Description: "Master AD administration with PowerShell",
			Category:    "On-Premise",
			Duration:    3 * time.Hour,
			Difficulty:  "Intermediate",
			Icon:        "🏢",
			Lessons: []Lesson{
				{
					ID:          "ad-1",
					Title:       "AD Module Overview",
					Description: "Introduction to the Active Directory PowerShell module",
					Duration:    20 * time.Minute,
				},
				{
					ID:          "ad-2",
					Title:       "User Management",
					Description: "Creating, modifying, and managing AD users",
					Duration:    30 * time.Minute,
				},
				{
					ID:          "ad-3",
					Title:       "Group Management",
					Description: "Working with AD groups and memberships",
					Duration:    25 * time.Minute,
				},
			},
		},
		{
			ID:          "msgraph",
			Title:       "Microsoft Graph PowerShell",
			Description: "Modern cloud management with MS Graph",
			Category:    "Cloud",
			Duration:    4 * time.Hour,
			Difficulty:  "Advanced",
			Icon:        "☁️",
			Lessons: []Lesson{
				{
					ID:          "graph-1",
					Title:       "Graph API Fundamentals",
					Description: "Understanding Microsoft Graph and authentication",
					Duration:    30 * time.Minute,
				},
				{
					ID:          "graph-2",
					Title:       "User and Group Management",
					Description: "Managing Azure AD/Entra ID resources",
					Duration:    35 * time.Minute,
				},
			},
		},
		{
			ID:          "scripting",
			Title:       "Advanced Scripting",
			Description: "Build robust PowerShell scripts and tools",
			Category:    "Advanced",
			Duration:    5 * time.Hour,
			Difficulty:  "Advanced",
			Icon:        "🚀",
			Lessons: []Lesson{
				{
					ID:          "script-1",
					Title:       "Functions and Modules",
					Description: "Creating reusable PowerShell code",
					Duration:    40 * time.Minute,
				},
				{
					ID:          "script-2",
					Title:       "Error Handling",
					Description: "Implementing robust error handling",
					Duration:    30 * time.Minute,
				},
			},
		},
		{
			ID:          "automation",
			Title:       "Automation & DevOps",
			Description: "Automate everything with PowerShell",
			Category:    "DevOps",
			Duration:    4 * time.Hour,
			Difficulty:  "Intermediate",
			Icon:        "⚙️",
			Lessons: []Lesson{
				{
					ID:          "auto-1",
					Title:       "Scheduled Tasks",
					Description: "Automating scripts with Task Scheduler",
					Duration:    25 * time.Minute,
				},
				{
					ID:          "auto-2",
					Title:       "CI/CD Pipelines",
					Description: "PowerShell in modern DevOps workflows",
					Duration:    35 * time.Minute,
				},
			},
		},
	}
}

// GetCategories returns all module categories
func GetCategories() []ModuleCategory {
	modules := GetAvailableModules()
	categoryMap := make(map[string][]Module)
	
	for _, module := range modules {
		categoryMap[module.Category] = append(categoryMap[module.Category], module)
	}
	
	categories := []ModuleCategory{
		{ID: "Foundation", Name: "Foundation", Icon: "📚", Description: "Start here to build a solid PowerShell foundation"},
		{ID: "On-Premise", Name: "On-Premise", Icon: "🏢", Description: "Manage traditional infrastructure"},
		{ID: "Cloud", Name: "Cloud", Icon: "☁️", Description: "Modern cloud management and APIs"},
		{ID: "Advanced", Name: "Advanced", Icon: "🚀", Description: "Take your skills to the next level"},
		{ID: "DevOps", Name: "DevOps", Icon: "⚙️", Description: "Automation and CI/CD practices"},
	}
	
	for i := range categories {
		categories[i].Modules = categoryMap[categories[i].ID]
	}
	
	return categories
}