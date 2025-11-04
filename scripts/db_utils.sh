#!/bin/bash

# PowerHell Database Utilities

DB_PATH="$HOME/.powerhell/powerhell.db"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if database exists
check_db() {
    if [ ! -f "$DB_PATH" ]; then
        echo -e "${RED}Database not found at: $DB_PATH${NC}"
        echo "Run PowerHell at least once to create the database."
        exit 1
    fi
}

# Show database info
show_info() {
    check_db
    echo -e "${BLUE}PowerHell Database Information${NC}"
    echo "Location: $DB_PATH"
    echo "Size: $(du -h "$DB_PATH" | cut -f1)"
    echo ""
    
    # Show account count
    COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM accounts WHERE is_active = 1;" 2>/dev/null)
    echo -e "${GREEN}Total Active Accounts: $COUNT${NC}"
    
    # Show recent accounts
    echo -e "\n${YELLOW}Recent Accounts:${NC}"
    sqlite3 "$DB_PATH" -header -column "SELECT account_number, name, created_at FROM accounts ORDER BY created_at DESC LIMIT 5;"
}

# List all accounts
list_accounts() {
    check_db
    echo -e "${BLUE}All PowerHell Accounts${NC}\n"
    sqlite3 "$DB_PATH" -header -column "SELECT id, account_number, name, email, created_at, last_login FROM accounts WHERE is_active = 1 ORDER BY created_at DESC;"
}

# Search for an account
search_account() {
    check_db
    if [ -z "$1" ]; then
        echo "Usage: $0 search <account_number_or_name>"
        exit 1
    fi
    
    echo -e "${BLUE}Searching for: $1${NC}\n"
    sqlite3 "$DB_PATH" -header -column "SELECT * FROM accounts WHERE account_number LIKE '%$1%' OR name LIKE '%$1%' OR email LIKE '%$1%';"
}

# Show account stats
show_stats() {
    check_db
    if [ -z "$1" ]; then
        echo "Usage: $0 stats <account_number>"
        exit 1
    fi
    
    echo -e "${BLUE}Statistics for Account: $1${NC}\n"
    
    # Get account ID
    ACCOUNT_ID=$(sqlite3 "$DB_PATH" "SELECT id FROM accounts WHERE account_number = '$1';" 2>/dev/null)
    
    if [ -z "$ACCOUNT_ID" ]; then
        echo -e "${RED}Account not found!${NC}"
        exit 1
    fi
    
    # Show progress
    echo -e "${GREEN}Learning Progress:${NC}"
    sqlite3 "$DB_PATH" -header -column "SELECT module_id, lesson_id, completed_at FROM account_progress WHERE account_id = $ACCOUNT_ID ORDER BY completed_at DESC;"
    
    # Show sessions
    echo -e "\n${GREEN}Recent Sessions:${NC}"
    sqlite3 "$DB_PATH" -header -column "SELECT session_start, session_end, duration_seconds FROM account_sessions WHERE account_id = $ACCOUNT_ID ORDER BY session_start DESC LIMIT 10;"
}

# Backup database
backup_db() {
    check_db
    BACKUP_FILE="$HOME/.powerhell/backup_$(date +%Y%m%d_%H%M%S).db"
    cp "$DB_PATH" "$BACKUP_FILE"
    echo -e "${GREEN}Database backed up to: $BACKUP_FILE${NC}"
}

# Export accounts to CSV
export_accounts() {
    check_db
    OUTPUT_FILE="$HOME/.powerhell/accounts_export_$(date +%Y%m%d_%H%M%S).csv"
    sqlite3 "$DB_PATH" -header -csv "SELECT account_number, name, email, created_at FROM accounts WHERE is_active = 1;" > "$OUTPUT_FILE"
    echo -e "${GREEN}Accounts exported to: $OUTPUT_FILE${NC}"
}

# Main menu
case "$1" in
    info)
        show_info
        ;;
    list)
        list_accounts
        ;;
    search)
        search_account "$2"
        ;;
    stats)
        show_stats "$2"
        ;;
    backup)
        backup_db
        ;;
    export)
        export_accounts
        ;;
    *)
        echo "PowerHell Database Utilities"
        echo ""
        echo "Usage: $0 {info|list|search|stats|backup|export}"
        echo ""
        echo "Commands:"
        echo "  info              Show database information"
        echo "  list              List all accounts"
        echo "  search <term>     Search for an account"
        echo "  stats <account>   Show account statistics"
        echo "  backup            Backup the database"
        echo "  export            Export accounts to CSV"
        echo ""
        echo "Example:"
        echo "  $0 search john"
        echo "  $0 stats 1234567890123456"
        ;;
esac