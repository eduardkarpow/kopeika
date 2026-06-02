DB_DRIVER=postgres
DB_DSN="postgres://user:admin@pg:5432/kopeika?sslmode=disable"
MIGRATIONS_DIR=./migrations

.PRONY: migration-create migration-up migration-down migration-status

migration-create:
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

migration-up:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) up

migration-down:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) down

migration-status:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) status