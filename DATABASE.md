# PowerHell SQLite Database 🗄️

## Overview

PowerHell now uses SQLite for persistent storage of user accounts, learning progress, and session tracking. The database provides a robust, serverless solution for managing user data.

## Database Location

- **Local Mode**: `~/.powerhell/powerhell.db`
- **SSH Server**: Same location on the server machine
- **Auto-created**: Database and tables are automatically created on first run

## Database Schema

### 1. **accounts** Table
Stores user account information:
```sql
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_number TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME,
    is_active BOOLEAN DEFAULT 1
);
```

### 2. **account_progress** Table
Tracks learning progress:
```sql
CREATE TABLE account_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    module_id TEXT NOT NULL,
    lesson_id TEXT NOT NULL,
    completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    UNIQUE(account_id, module_id, lesson_id)
);
```

### 3. **account_sessions** Table
Records learning sessions:
```sql
CREATE TABLE account_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    session_start DATETIME DEFAULT CURRENT_TIMESTAMP,
    session_end DATETIME,
    duration_seconds INTEGER,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);
```

### 4. **account_achievements** Table
Stores earned achievements:
```sql
CREATE TABLE account_achievements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    achievement_id TEXT NOT NULL,
    earned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    UNIQUE(account_id, achievement_id)
);
```

## Features

### Account Management
- ✅ Unique account number generation
- ✅ Duplicate prevention
- ✅ Last login tracking
- ✅ Account deactivation support

### Progress Tracking
- ✅ Module/lesson completion tracking
- ✅ No duplicate progress entries
- ✅ Timestamp for each completion

### Session Management
- ✅ Automatic session start on login
- ✅ Session duration calculation
- ✅ Historical session data

### Data Integrity
- ✅ Foreign key constraints
- ✅ Unique constraints where needed
- ✅ Indexes for performance
- ✅ Transaction support

## Database Utilities

Use the provided utilities via Make:

```bash
# Show database info and stats
make db-info

# List all accounts
make db-list

# Backup database
make db-backup

# Export accounts to CSV
make db-export

# Clean database (WARNING: deletes all data!)
make db-clean
```

Or use the script directly:

```bash
# Search for an account
./scripts/db_utils.sh search john

# Show account statistics
./scripts/db_utils.sh stats 1234567890123456
```

## Direct Database Access

You can also query the database directly using SQLite:

```bash
# Open database
sqlite3 ~/.powerhell/powerhell.db

# Show all tables
.tables

# Show table schema
.schema accounts

# Query examples
SELECT * FROM accounts;
SELECT COUNT(*) FROM accounts WHERE created_at > date('now', '-7 days');
SELECT * FROM account_progress WHERE account_id = 1;
```

## Backup and Recovery

### Automatic Backup
```bash
make db-backup
```
Creates timestamped backup in `~/.powerhell/backup_YYYYMMDD_HHMMSS.db`

### Manual Backup
```bash
cp ~/.powerhell/powerhell.db ~/powerhell_backup.db
```

### Restore from Backup
```bash
cp ~/powerhell_backup.db ~/.powerhell/powerhell.db
```

## Performance Considerations

1. **Indexes**: Created on frequently queried columns
   - `account_number` for fast login lookup
   - `email` for potential future features

2. **Connection Pooling**: Single connection per app instance

3. **Write-Ahead Logging**: SQLite uses WAL mode for better concurrency

## Security Notes

⚠️ **Important Security Considerations**:

1. **File Permissions**: Database file should be readable/writable only by owner
   ```bash
   chmod 600 ~/.powerhell/powerhell.db
   ```

2. **No Encryption**: SQLite database is not encrypted by default
   - Consider SQLCipher for encryption needs
   - Don't store sensitive data without encryption

3. **SQL Injection**: All queries use parameterized statements

4. **Backup Security**: Secure your backups as they contain all user data

## Migration from JSON

If you were using the previous JSON-based storage:

1. Old data location: `~/.powerhell/accounts.json`
2. New data location: `~/.powerhell/powerhell.db`
3. Migration is not automatic - accounts need to be recreated

## Future Enhancements

1. **Data Encryption**: Implement SQLCipher for encrypted storage
2. **Cloud Sync**: Sync progress across devices
3. **Analytics**: Learning analytics and insights
4. **Import/Export**: More data formats (JSON, XML)
5. **Audit Trail**: Track all account changes
6. **Data Retention**: Automatic cleanup of old sessions

## Troubleshooting

### Database Locked
If you get "database is locked" errors:
```bash
# Check for other processes
lsof ~/.powerhell/powerhell.db

# Force unlock (use carefully)
rm ~/.powerhell/powerhell.db-wal
rm ~/.powerhell/powerhell.db-shm
```

### Corruption
If database becomes corrupted:
```bash
# Try to recover
sqlite3 ~/.powerhell/powerhell.db ".recover" | sqlite3 ~/.powerhell/powerhell_recovered.db

# Or restore from backup
make db-backup
```

### Reset Everything
To start fresh:
```bash
make db-clean
```