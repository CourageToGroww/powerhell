package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/couragetogroww/powerhell/pkg/auth"
	"github.com/couragetogroww/powerhell/pkg/menus/types"
	"github.com/couragetogroww/powerhell/pkg/views"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.TerminalHeight = msg.Height
		// Update view sizes
		if m.Dashboard != nil {
			userName := ""
			if m.CurrentAccount != nil {
				userName = m.CurrentAccount.Name
			} else if m.Dashboard.UserName != "" {
				userName = m.Dashboard.UserName
			}
			m.Dashboard = views.NewDashboardViewWithUser(m.TerminalWidth, m.TerminalHeight, userName)
		}
		if m.LessonView != nil && m.CurrentModule != nil {
			m.LessonView = views.NewLessonView(m.CurrentModule, m.TerminalWidth, m.TerminalHeight)
		}

	case tickMsg: // Handle animation tick
		if m.AppState == StateIntro {
			m.Flames = animateFlames(m.Flames)
			return m, tickCmd() // Continue ticking for intro animation
		} else if m.ShowAccountWarning {
			m.AnimationFrame++
			return m, tickCmd() // Continue ticking for warning animation
		}
		return m, nil // Ignore tick if not in other states

	case tea.KeyMsg:
		switch m.AppState {
		case StateIntro:
			if msg.String() == "enter" {
				m.AppState = StateAuthMenu // Transition to AuthMenu first
				// Initialize the auth menu
				m.MenuManager.SetCurrentMenu(StateAuthMenu)
				m.MenuCursor = 0
				return m, nil // Stop the tick for flame animation
			} else if msg.String() == "ctrl+c" || msg.String() == "q" {
				m.Quit = true
				return m, tea.Quit
			} else if msg.String() == "h" {
				m.ShowHelp = !m.ShowHelp
			}
			
		case StateAccountCreation:
			switch msg.String() {
			case "ctrl+c", "q":
				m.Quit = true
				return m, tea.Quit
			case "enter":
				if m.FocusedField == FocusDisplayInfo {
					// Transition to the new Dashboard after account creation
					m.AppState = StateDashboard
					if m.CurrentAccount != nil {
						m.Dashboard = views.NewDashboardViewWithUser(m.TerminalWidth, m.TerminalHeight, m.CurrentAccount.Name)
					} else {
						m.Dashboard = views.NewDashboardViewWithUser(m.TerminalWidth, m.TerminalHeight, m.NameInput.Value())
					}
					m.ShowAccountWarning = false // Stop the warning animation
					return m, nil
				} else if m.FocusedField == FocusEmail && m.NameInput.Value() != "" && m.EmailInput.Value() != "" {
					// All fields filled, generate unique account number and save
					if m.AccountStore != nil {
						accountNumber, err := m.AccountStore.GenerateUniqueAccountNumber(generateAccountNumber)
						if err != nil {
							// Handle error - maybe show to user
							return m, nil
						}
						m.GeneratedAccountNumber = accountNumber
						
						// Create the account
						account, err := m.AccountStore.CreateAccount(
							m.NameInput.Value(),
							m.EmailInput.Value(),
							m.GeneratedAccountNumber,
						)
						if err != nil {
							// Handle error - maybe show to user
							return m, nil
						}
						
						m.CurrentAccount = account
					} else {
						// Fallback if no database
						m.GeneratedAccountNumber = generateAccountNumber()
					}
					
					m.FocusedField = FocusDisplayInfo
					m.NameInput.Blur()
					m.EmailInput.Blur()
					return m, nil
				} else {
					// Move to next field or handle incomplete form
					if m.FocusedField == FocusName {
						m.FocusedField = FocusEmail
						m.NameInput.Blur()
						return m, m.EmailInput.Focus()
					}
				}
			case "tab", "shift+tab", "up", "down":
				if m.FocusedField == FocusDisplayInfo {
					return m, nil
				}
				s := msg.String()
				if s == "up" || s == "shift+tab" {
					m.FocusedField--
				} else {
					m.FocusedField++
				}
				if m.FocusedField > FocusEmail {
					m.FocusedField = FocusName
				} else if m.FocusedField < FocusName {
					m.FocusedField = FocusEmail
				}
				var focusCmd tea.Cmd
				if m.FocusedField == FocusName {
					m.EmailInput.Blur()
					focusCmd = m.NameInput.Focus()
				} else {
					m.NameInput.Blur()
					focusCmd = m.EmailInput.Focus()
				}
				return m, focusCmd
			}
			var cmd tea.Cmd
			if m.FocusedField == FocusName {
				m.NameInput, cmd = m.NameInput.Update(msg)
			} else if m.FocusedField == FocusEmail {
				m.EmailInput, cmd = m.EmailInput.Update(msg)
			}
			return m, cmd

		case StateAuthMenu:
			// Initialize auth menu if not already set
			if m.MenuManager.GetCurrentState() != StateAuthMenu {
				m.MenuManager.SetCurrentMenu(StateAuthMenu)
				m.MenuCursor = 0
			}
			
			switch msg.String() {
			case "ctrl+c", "q":
				m.Quit = true
				return m, tea.Quit
			case "h":
				m.ShowHelp = !m.ShowHelp
			case "up", "k":
				if m.MenuCursor > 0 {
					m.MenuCursor--
				}
			case "down", "j":
				menuOptions := m.MenuManager.GetMenuOptionsAsStrings()
				if m.MenuCursor < len(menuOptions)-1 {
					m.MenuCursor++
				}
			case "enter":
				result := m.MenuManager.HandleSelection(m.MenuCursor)
				switch result.Action {
				case types.ActionNavigate:
					m.AppState = result.NextState
					m.MenuCursor = 0 // Reset cursor for new menu
					// Special handling for certain states
					if result.NextState == StateAccountCreation {
						m.FocusedField = FocusName
						m.NameInput.SetValue("")
						m.EmailInput.SetValue("")
						m.GeneratedAccountNumber = ""
						return m, m.NameInput.Focus()
					} else if result.NextState == StateSignInPlaceholder {
						// Transition to Sign In view
						m.AppState = StateSignIn
						m.SignInView = views.NewSignInView(m.TerminalWidth, m.TerminalHeight)
						return m, nil
					}
				case types.ActionExit:
					m.Quit = true
					return m, tea.Quit
				}
			}

		case StateModuleExplorer:
			switch msg.String() {
			case "ctrl+c", "q":
				m.Quit = true
				return m, tea.Quit
			case "up", "k":
				if m.ModuleExplorerSidebarCursor > 0 {
					m.ModuleExplorerSidebarCursor--
				}
			case "down", "j":
				if m.ModuleExplorerSidebarCursor < len(m.ModuleExplorerSidebarOptions)-1 {
					m.ModuleExplorerSidebarCursor++
				}
			case "enter":
				selectedModule := m.ModuleExplorerSidebarOptions[m.ModuleExplorerSidebarCursor]
				switch selectedModule {
				case "Learn":
					m.ModuleExplorerSidebarOptions = []string{"On-Premise Learning", "MSGraph Module", "Exchange Module", "HTTP Client (Invoke-WebRequest, Graph API)", "Custom & Popular Modules (EntraExporter)", "PowerShell SDK Learning", "Exit"}
					m.ModuleExplorerSidebarCursor = 0
					m.ModuleExplorerContent = "Select a learning module."
				case "Exit":
					m.Quit = true
					return m, tea.Quit
				default:
					m.ModuleExplorerContent = fmt.Sprintf("Selected: %s\n(Press Q to quit)", selectedModule)
				}
			}
		
		case StateDashboard:
			if m.Dashboard == nil {
				userName := ""
				if m.CurrentAccount != nil {
					userName = m.CurrentAccount.Name
				}
				m.Dashboard = views.NewDashboardViewWithUser(m.TerminalWidth, m.TerminalHeight, userName)
			}
			switch msg.String() {
			case "ctrl+c", "q":
				m.Quit = true
				return m, tea.Quit
			case "h":
				m.ShowHelp = !m.ShowHelp
			case "enter":
				if selectedModule := m.Dashboard.GetSelectedModule(); selectedModule != nil {
					m.CurrentModule = selectedModule
					m.LessonView = views.NewLessonView(selectedModule, m.TerminalWidth, m.TerminalHeight)
					m.AppState = StateLesson
				}
			default:
				if !m.ShowHelp {
					m.Dashboard.Update(msg.String())
				}
			}
		
		case StateLesson:
			if m.LessonView == nil {
				return m, nil
			}
			switch msg.String() {
			case "ctrl+c":
				m.Quit = true
				return m, tea.Quit
			case "h":
				m.ShowHelp = !m.ShowHelp
			case "q":
				// Go back to dashboard
				m.AppState = StateDashboard
			default:
				if !m.ShowHelp {
					m.LessonView.Update(msg.String())
				}
			}
			
		case StateSignIn:
			if m.SignInView == nil {
				m.SignInView = views.NewSignInView(m.TerminalWidth, m.TerminalHeight)
			}
			switch msg.String() {
			case "ctrl+c":
				m.Quit = true
				return m, tea.Quit
			case "esc":
				// Go back to auth menu
				m.AppState = StateAuthMenu
				m.SignInView = nil
			case "enter":
				// Validate and sign in
				accountNumber := m.SignInView.GetAccountNumber()
				if accountNumber != "" && len(accountNumber) == 16 {
					if m.AccountStore != nil {
						// Sign in with database
						account, err := m.AccountStore.SignIn(accountNumber)
						if err == auth.ErrAccountNotFound {
							m.SignInView.SetError("Account not found. Please check your account number.")
						} else if err != nil {
							m.SignInView.SetError("Sign in failed. Please try again.")
						} else {
							// Successful sign in
							m.CurrentAccount = account
							
							// Start a new session
							if sessionID, err := m.AccountStore.StartSession(account.ID); err == nil {
								m.SessionID = sessionID
							}
							
							m.AppState = StateDashboard
							m.Dashboard = views.NewDashboardViewWithUser(m.TerminalWidth, m.TerminalHeight, account.Name)
						}
					} else {
						// No database - just proceed
						m.AppState = StateDashboard
						m.Dashboard = views.NewDashboardView(m.TerminalWidth, m.TerminalHeight)
					}
				} else {
					m.SignInView.SetError("Please enter a valid 16-digit account number")
				}
			default:
				// Pass the entire key message to the text input
				var cmd tea.Cmd
				m.SignInView.AccountInput, cmd = m.SignInView.AccountInput.Update(msg)
				return m, cmd
			}
		}
		
		// Common key handling for all states
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			if m.AppState != StateAccountCreation || m.FocusedField == FocusDisplayInfo {
				m.Quit = true
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}