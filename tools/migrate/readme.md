# Migrate

A simple CLI tool based on https://github.com/golang-migrate to allow migrating a Postgres DB based on the migration files located in repository.

## Usage

`go run tools/migrate/main.go -dsn="postgres://user:pass@host:5432/schema" -source ./database`

Pretty straightforward, the tool has just a minimal functionality of migrating up and down.

The following flags are supported:
- `dsn` (required, a DSN to a postgres DB in the format of `postgres://user:pass@host:5432/schema`)
- `source` (required, a path to a directory containing the DB migrations)
- `step` (optional, number of migrations to execute - 0 is default and means all available)
