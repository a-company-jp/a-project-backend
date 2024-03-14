DB_NAME = $(shell cat .env | grep DB_NAME | cut -d = -f 2)
DB_USER = $(shell cat .env | grep DB_USER | cut -d = -f 2)
DB_PASSWORD = $(shell cat .env | grep DB_PASS | cut -d = -f 2)
DB_PORT = $(shell cat .env | grep DB_PORT | cut -d = -f 2)

DSN = mysql://$(DB_USER):$(DB_PASSWORD)@tcp(127.0.0.1:$(DB_PORT))/$(DB_NAME)?parseTime=true

.PHONY: create_migration

check_env:
	echo $(DB_NAME), $(DB_USER), $(DB_PASSWORD)

check_dsn:
	echo $(DSN)

create_migration:
ifndef NAME
	@echo "Usage: make NAME=migration_name create_migration"
else
	migrate create -ext sql -dir mysql/migrations -seq $(NAME)
endif

.PHONY: migrate_up
migrate_up:
	migrate -path mysql/migrations -database "$(DSN)" up

.PHONY: migrate_down
migrate_down:
	migrate -path mysql/migrations -database "$(DSN)" down 1

# setup golang-migrate
setup:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
