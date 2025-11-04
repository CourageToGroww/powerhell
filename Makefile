.PHONY: build run ssh-server clean generate-key

# Build the application
build:
	go build -o powerhell cmd/powerhell/main.go

# Run locally
run: build
	./powerhell

# Run as SSH server
ssh-server: build
	./powerhell -ssh

# Run SSH server with custom port
ssh-custom: build
	./powerhell -ssh -port 2222

# Generate SSH host key
generate-key:
	@cd scripts && ./generate_host_key.sh

# Clean build artifacts
clean:
	rm -f powerhell
	rm -f ssh_host_key ssh_host_key.pub

# Database utilities
db-info:
	@cd scripts && ./db_utils.sh info

db-list:
	@cd scripts && ./db_utils.sh list

db-backup:
	@cd scripts && ./db_utils.sh backup

db-export:
	@cd scripts && ./db_utils.sh export

# Clean database (WARNING: removes all accounts!)
db-clean:
	@echo "WARNING: This will delete all accounts!"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read confirm
	rm -f ~/.powerhell/powerhell.db
	@echo "Database deleted."

# Install dependencies
deps:
	go mod download

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o powerhell-linux-amd64 cmd/powerhell/main.go
	GOOS=darwin GOARCH=amd64 go build -o powerhell-darwin-amd64 cmd/powerhell/main.go
	GOOS=windows GOARCH=amd64 go build -o powerhell-windows-amd64.exe cmd/powerhell/main.go