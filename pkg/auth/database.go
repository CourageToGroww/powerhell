package auth

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database handles SQLite database operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase() (*Database, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(homeDir, ".powerhell")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "powerhell.db")
	
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	d := &Database{db: db}
	
	// Create tables if they don't exist
	if err := d.createTables(); err != nil {
		db.Close()
		return nil, err
	}

	return d, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// createTables creates the necessary database tables
func (d *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_number TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login DATETIME,
		is_active BOOLEAN DEFAULT 1
	);

	CREATE INDEX IF NOT EXISTS idx_account_number ON accounts(account_number);
	CREATE INDEX IF NOT EXISTS idx_email ON accounts(email);

	CREATE TABLE IF NOT EXISTS account_progress (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_id INTEGER NOT NULL,
		module_id TEXT NOT NULL,
		lesson_id TEXT NOT NULL,
		completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (account_id) REFERENCES accounts(id),
		UNIQUE(account_id, module_id, lesson_id)
	);

	CREATE TABLE IF NOT EXISTS account_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_id INTEGER NOT NULL,
		session_start DATETIME DEFAULT CURRENT_TIMESTAMP,
		session_end DATETIME,
		duration_seconds INTEGER,
		FOREIGN KEY (account_id) REFERENCES accounts(id)
	);

	CREATE TABLE IF NOT EXISTS account_achievements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_id INTEGER NOT NULL,
		achievement_id TEXT NOT NULL,
		earned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (account_id) REFERENCES accounts(id),
		UNIQUE(account_id, achievement_id)
	);
	`

	_, err := d.db.Exec(query)
	return err
}

// CreateAccount creates a new account in the database
func (d *Database) CreateAccount(account *Account) error {
	query := `
		INSERT INTO accounts (account_number, name, email) 
		VALUES (?, ?, ?)
	`
	
	result, err := d.db.Exec(query, account.AccountNumber, account.Name, account.Email)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	account.ID = int(id)
	return nil
}

// GetAccountByNumber retrieves an account by account number
func (d *Database) GetAccountByNumber(accountNumber string) (*Account, error) {
	// Try both with and without spaces to handle legacy accounts
	query := `
		SELECT id, account_number, name, email, created_at, last_login, is_active
		FROM accounts 
		WHERE (account_number = ? OR account_number = ?) AND is_active = 1
	`

	// Format account number with spaces for legacy compatibility
	formattedNumber := ""
	if len(accountNumber) == 16 {
		formattedNumber = fmt.Sprintf("%s %s %s %s", 
			accountNumber[0:4], 
			accountNumber[4:8], 
			accountNumber[8:12], 
			accountNumber[12:16])
	}

	var account Account
	var createdAt, lastLogin sql.NullTime

	err := d.db.QueryRow(query, accountNumber, formattedNumber).Scan(
		&account.ID,
		&account.AccountNumber,
		&account.Name,
		&account.Email,
		&createdAt,
		&lastLogin,
		&account.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	if createdAt.Valid {
		account.CreatedAt = createdAt.Time.Format(time.RFC3339)
	}
	if lastLogin.Valid {
		account.LastLogin = lastLogin.Time.Format(time.RFC3339)
	}

	return &account, nil
}

// UpdateLastLogin updates the last login time for an account
func (d *Database) UpdateLastLogin(accountID int) error {
	query := `UPDATE accounts SET last_login = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := d.db.Exec(query, accountID)
	return err
}

// AccountExists checks if an account number already exists
func (d *Database) AccountExists(accountNumber string) (bool, error) {
	// Check both with and without spaces
	query := `SELECT COUNT(*) FROM accounts WHERE account_number = ? OR account_number = ?`
	
	// Format account number with spaces for legacy compatibility
	formattedNumber := ""
	if len(accountNumber) == 16 {
		formattedNumber = fmt.Sprintf("%s %s %s %s", 
			accountNumber[0:4], 
			accountNumber[4:8], 
			accountNumber[8:12], 
			accountNumber[12:16])
	}
	
	var count int
	err := d.db.QueryRow(query, accountNumber, formattedNumber).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetAccountCount returns the total number of accounts
func (d *Database) GetAccountCount() (int, error) {
	query := `SELECT COUNT(*) FROM accounts WHERE is_active = 1`
	
	var count int
	err := d.db.QueryRow(query).Scan(&count)
	return count, err
}

// SaveProgress saves learning progress for an account
func (d *Database) SaveProgress(accountID int, moduleID, lessonID string) error {
	query := `
		INSERT OR REPLACE INTO account_progress (account_id, module_id, lesson_id)
		VALUES (?, ?, ?)
	`
	
	_, err := d.db.Exec(query, accountID, moduleID, lessonID)
	return err
}

// GetProgress retrieves progress for an account
func (d *Database) GetProgress(accountID int) ([]Progress, error) {
	query := `
		SELECT module_id, lesson_id, completed_at
		FROM account_progress
		WHERE account_id = ?
		ORDER BY completed_at DESC
	`

	rows, err := d.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []Progress
	for rows.Next() {
		var p Progress
		var completedAt time.Time
		
		err := rows.Scan(&p.ModuleID, &p.LessonID, &completedAt)
		if err != nil {
			return nil, err
		}
		
		p.CompletedAt = completedAt.Format(time.RFC3339)
		progress = append(progress, p)
	}

	return progress, nil
}

// StartSession starts a new learning session
func (d *Database) StartSession(accountID int) (int64, error) {
	query := `INSERT INTO account_sessions (account_id) VALUES (?)`
	
	result, err := d.db.Exec(query, accountID)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// EndSession ends a learning session
func (d *Database) EndSession(sessionID int64) error {
	query := `
		UPDATE account_sessions 
		SET session_end = CURRENT_TIMESTAMP,
		    duration_seconds = CAST((julianday(CURRENT_TIMESTAMP) - julianday(session_start)) * 86400 AS INTEGER)
		WHERE id = ?
	`
	
	_, err := d.db.Exec(query, sessionID)
	return err
}

// GetStats retrieves statistics for an account
func (d *Database) GetStats(accountID int) (*AccountStats, error) {
	stats := &AccountStats{}

	// Get total completed lessons
	query := `SELECT COUNT(*) FROM account_progress WHERE account_id = ?`
	err := d.db.QueryRow(query, accountID).Scan(&stats.TotalLessonsCompleted)
	if err != nil {
		return nil, err
	}

	// Get total time spent
	query = `SELECT COALESCE(SUM(duration_seconds), 0) FROM account_sessions WHERE account_id = ?`
	err = d.db.QueryRow(query, accountID).Scan(&stats.TotalTimeSeconds)
	if err != nil {
		return nil, err
	}

	// Get achievement count
	query = `SELECT COUNT(*) FROM account_achievements WHERE account_id = ?`
	err = d.db.QueryRow(query, accountID).Scan(&stats.AchievementCount)
	if err != nil {
		return nil, err
	}

	// Get current streak (simplified - counts consecutive days with sessions)
	// This is a basic implementation - you might want something more sophisticated
	query = `
		SELECT COUNT(DISTINCT DATE(session_start)) 
		FROM account_sessions 
		WHERE account_id = ? 
		AND session_start >= date('now', '-7 days')
	`
	err = d.db.QueryRow(query, accountID).Scan(&stats.CurrentStreak)
	if err != nil {
		return nil, err
	}

	return stats, nil
}