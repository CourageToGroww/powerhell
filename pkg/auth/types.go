package auth

import (
	"errors"
)

// Common errors
var (
	ErrAccountNotFound = errors.New("account not found")
	ErrDuplicateAccount = errors.New("account already exists")
	ErrInvalidAccountNumber = errors.New("invalid account number")
)

// Account represents a user account with database fields
type Account struct {
	ID            int    `json:"id"`
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	CreatedAt     string `json:"created_at"`
	LastLogin     string `json:"last_login,omitempty"`
	IsActive      bool   `json:"is_active"`
}

// Progress represents learning progress
type Progress struct {
	ModuleID    string `json:"module_id"`
	LessonID    string `json:"lesson_id"`
	CompletedAt string `json:"completed_at"`
}

// AccountStats represents account statistics
type AccountStats struct {
	TotalLessonsCompleted int `json:"total_lessons_completed"`
	TotalTimeSeconds      int `json:"total_time_seconds"`
	AchievementCount      int `json:"achievement_count"`
	CurrentStreak         int `json:"current_streak"`
}

// Session represents a learning session
type Session struct {
	ID               int64  `json:"id"`
	AccountID        int    `json:"account_id"`
	SessionStart     string `json:"session_start"`
	SessionEnd       string `json:"session_end,omitempty"`
	DurationSeconds  int    `json:"duration_seconds,omitempty"`
}