DB_FILE := ./storage/storage.db
MIGRATIONS_DIR := ./internal/adapters/repository/sqlite/migrations
GOOSE := goose

DB_DSN := sqlite3 $(DB_FILE)

migrate-up:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DSN) up

migrate-down:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DSN) down

clean-db:
	@rm -f $(DB_FILE)
	@echo "Databse $(DB_FILE) deleted"

reset-db: clean-db
	@$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DSN) up
	@echo "New DB initialized"

status:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DSN) status
