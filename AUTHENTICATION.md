# PowerHell Authentication System 🔐

## Overview

PowerHell now features a passwordless authentication system using unique 16-digit account numbers. This provides a simple yet secure way for users to access their learning progress.

## Features

### 1. **Account Creation (Sign Up)**
- Users provide their name and email
- System generates a unique 16-digit account number
- **FLASHING WARNING**: The account number is displayed with an animated warning in red/orange/yellow colors
- Warning message: "⚠️ KEEP YOUR ACCOUNT NUMBER SAFE ⚠️ IF YOU LOSE IT YOU WILL NOT BE ABLE TO SIGN BACK IN"
- Account is automatically saved to local storage

### 2. **Sign In**
- Users enter their 16-digit account number
- Format: `1234 5678 9012 3456` (spaces are automatically handled)
- Real-time validation ensures correct format
- Clear error messages for invalid or non-existent accounts

### 3. **Account Storage**
- Accounts are stored locally in `~/.powerhell/accounts.json`
- Each account includes:
  - Unique account number
  - User name
  - Email address
  - Creation timestamp
- Automatic persistence and loading

## User Flow

```
Start App → Intro Screen → Press Enter → Authentication Menu
                                             ↓
                              ┌─────────────┴─────────────┐
                              ↓                           ↓
                           Sign Up                     Sign In
                              ↓                           ↓
                         Enter Name                 Enter Account #
                              ↓                           ↓
                        Enter Email                   Validate
                              ↓                           ↓
                    Generate Account #          Success → Dashboard
                              ↓                     ↓
                    ⚠️ FLASHING WARNING ⚠️      Error → Try Again
                              ↓
                         Dashboard
```

## Security Features

1. **No Passwords**: Eliminates password-related vulnerabilities
2. **Unique Account Numbers**: Cryptographically random generation ensures uniqueness
3. **Local Storage**: Accounts stored securely on user's machine
4. **Format Validation**: Ensures account numbers are exactly 16 digits

## Technical Implementation

### Account Number Format
- 16 digits total
- Displayed as: `XXXX XXXX XXXX XXXX`
- First digit is never 0 (ensures consistent length)
- Generated using secure random number generator

### Flashing Warning Animation
- Cycles through colors: Red → Orange → Yellow
- Blink effect for maximum visibility
- Large, bold text with warning icons
- Continues until user presses Enter

### Storage Location
- **Local**: `~/.powerhell/accounts.json`
- **SSH Server**: Same location on server machine
- JSON format for easy backup/restore

## SSH Server Considerations

When running as an SSH server:
1. Each SSH session gets its own authentication flow
2. Accounts are stored on the server machine
3. Multiple users can have accounts on the same server
4. Consider implementing:
   - Account limits per server
   - Admin tools for account management
   - Backup/restore functionality

## Future Enhancements

1. **Account Recovery**
   - Security questions
   - Email-based recovery
   - Admin override capability

2. **Progress Tracking**
   - Link learning progress to accounts
   - Sync across devices
   - Export/import progress

3. **Multi-factor Authentication**
   - Optional email verification
   - Time-based codes
   - Biometric support (for local mode)

4. **Account Management**
   - Change name/email
   - View account creation date
   - Delete account option

## Demo Commands

```bash
# Run locally
./powerhell

# Run as SSH server
./powerhell -ssh

# Connect via SSH
ssh -p 2222 localhost

# View stored accounts (on server/local machine)
cat ~/.powerhell/accounts.json
```

## Important Notes

⚠️ **For Production Use**:
- Implement proper database storage (PostgreSQL, MySQL, etc.)
- Add account recovery mechanisms
- Consider rate limiting for account creation
- Implement audit logging
- Add admin tools for account management
- Consider GDPR compliance for email storage