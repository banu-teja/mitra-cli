DB_NAME=history.db
MIGRATION_DIR=db/migration

.PHONY: createdb dropdb migrateup migratedown new_migration sqlc test run build

createdb:
	touch $(DB_NAME)

dropdb:
	rm -f $(DB_NAME)

migrateup:
	@echo "Applying migrations..."
	@for file in $(MIGRATION_DIR)/*up.sql; do \
		echo "Applying $$file"; \
		sqlite3 $(DB_NAME) < $$file; \
	done

migratedown:
	@echo "Reverting migrations..."
	@for file in $$(ls -r $(MIGRATION_DIR)/*down.sql 2>/dev/null); do \
		echo "Applying $$file"; \
		sqlite3 $(DB_NAME) < $$file; \
	done


new_migration:
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	up_file="$(MIGRATION_DIR)/$${timestamp}_$${name}_up.sql"; \
	down_file="$(MIGRATION_DIR)/$${timestamp}_$${name}_down.sql"; \
	touch $$up_file $$down_file; \
	echo "Created migration files: $$up_file and $$down_file"

sqlc:
	sqlc generate

test:
	go test -v -count=1 -cover ./...

run:
	go run main.go

build:
	go build .