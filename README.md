## Development

### Launch MySQL Dev Instance

```bash
# Fill .env before running MySQL Instance
cd mysql
docker-compose up -d
```

### Setup `golang-migrate`

```bash
make setup
```
### Run the service

```bash
# Make sure that you have set DSN of active MySQL instance in .env
go run svc/cmd/dev/main.go
```

## Migrations

### Run Migrations

```bash
# Migrate all the way up
make migrate_up

# Migrate one step down
make migrate_down
```

### Create New Migration

```bash
# Replace MIGRATION_NAME with the name of the migration
make create_migration NAME="MIGRATION_NAME"
```
