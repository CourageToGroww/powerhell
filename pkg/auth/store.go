package auth

import (
	"fmt"
	"sync"
)

// Store manages account operations with SQLite backend
type Store struct {
	db *Database
	mu sync.RWMutex
}

// NewStore creates a new account store with SQLite backend
func NewStore() (*Store, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Store{
		db: db,
	}, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// CreateAccount creates a new account
func (s *Store) CreateAccount(name, email, accountNumber string) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if account already exists
	exists, err := s.db.AccountExists(accountNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateAccount
	}

	account := &Account{
		AccountNumber: accountNumber,
		Name:          name,
		Email:         email,
		IsActive:      true,
	}

	if err := s.db.CreateAccount(account); err != nil {
		return nil, err
	}

	return account, nil
}

// FindAccount finds an account by account number
func (s *Store) FindAccount(accountNumber string) (*Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.db.GetAccountByNumber(accountNumber)
}

// SignIn signs in a user and updates last login
func (s *Store) SignIn(accountNumber string) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, err := s.db.GetAccountByNumber(accountNumber)
	if err != nil {
		return nil, err
	}

	// Update last login
	if err := s.db.UpdateLastLogin(account.ID); err != nil {
		// Log error but don't fail sign in
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	return account, nil
}

// GenerateUniqueAccountNumber generates a unique account number
func (s *Store) GenerateUniqueAccountNumber(generator func() string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for attempts := 0; attempts < 100; attempts++ {
		accountNumber := generator()
		
		exists, err := s.db.AccountExists(accountNumber)
		if err != nil {
			return "", err
		}
		
		if !exists {
			return accountNumber, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique account number after 100 attempts")
}

// GetAccountCount returns the total number of accounts
func (s *Store) GetAccountCount() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.db.GetAccountCount()
}

// SaveProgress saves learning progress
func (s *Store) SaveProgress(accountID int, moduleID, lessonID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.SaveProgress(accountID, moduleID, lessonID)
}

// GetProgress retrieves learning progress
func (s *Store) GetProgress(accountID int) ([]Progress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.db.GetProgress(accountID)
}

// StartSession starts a new learning session
func (s *Store) StartSession(accountID int) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.StartSession(accountID)
}

// EndSession ends a learning session
func (s *Store) EndSession(sessionID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.EndSession(sessionID)
}

// GetStats retrieves account statistics
func (s *Store) GetStats(accountID int) (*AccountStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.db.GetStats(accountID)
}