# This line includes the .envrc file, which typically contains environment variables. 
# Variables like DB_ADDR (database address) would be defined there.
include .envrc

MIGRATIONS_PATH = ./cmd/migrate/migrations

# Each target follows this pattern:
# .PHONY: target-name
# target-name:
# 	@command
#
# Usage
# You'd run these with commands like make test, make migrate-up, or make migration create_users_table.
#  


# .PHONY: - Declares that the target name isn't a file.
# This prevents Make from looking for actual files with those names 
# and ensures the commands always run.
.PHONY: test
# test: - Runs Go tests with verbose output (-v) for all packages (./...)
test:
# Always USE TABS not spaces
# @ prefix - Suppresses echoing the command to the terminal. Without @, Make would print each command before executing it.
	@go test -v ./...

# migration: - Creates new database migration files using the migrate tool:
# .PHONY: migrate-create
.PHONY: migration 
migration:
# -seq creates sequential numbering
# -ext sql sets file extension to .sql
# -dir $(MIGRATIONS_PATH) specifies the directory
# $(filter-out $@,$(MAKECMDGOALS)) captures additional arguments passed to make
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

# migrate-up/migrate-down: - Runs database migrations up or down using the stored DB_ADDR environment variable
.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

# seed: - Executes a Go program that seeds the database with initial data
.PHONY: seed
seed: 
	@go run cmd/migrate/seed/main.go

# gen-docs: - Generates API documentation using Swagger:
.PHONY: gen-docs
gen-docs:
# swag init initializes documentation
# -g ./api/main.go specifies the main API file
# -d cmd,internal sets directories to scan
# swag fmt formats the generated files
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
