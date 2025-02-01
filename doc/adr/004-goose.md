# Title

This decision describes using the [`goose` package](https://github.com/pressly/goose) for managing DB schema and migrations

# Status

Active

# Date

2025-02-01

# Context

I want to be able to manage the DB schema for `fstagger` without having to rely on hand writing SQL commands in Go to run at runtime. A migration tool (think Flyway, Liquibase, etc.) that can be used like a package will give me the convenience of having migrations run without user intervention or upgrade scripts. It will also keep track of the current schema and which migrations have been applied.

# Decision

I will use `goose` to manage migrations for `fstagger` programmatically.
