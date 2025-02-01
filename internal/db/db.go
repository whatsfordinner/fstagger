package db

import (
	"context"
	"database/sql"
	"embed"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var defaultMigrationsFS embed.FS
var tagDB *TagDB

const defaultMigrationsDir = "migrations"
const defaultConnectionString = ":memory:"

type TagDB struct {
	client           *sql.DB
	connectionString string
	migrationsFS     embed.FS
	migrationsDir    string
}

func New(options ...func(*TagDB)) *TagDB {
	tagDB := &TagDB{
		connectionString: defaultConnectionString,
		migrationsFS:     defaultMigrationsFS,
		migrationsDir:    defaultMigrationsDir,
	}
	for _, o := range options {
		o(tagDB)
	}

	return tagDB
}

// WithConnectionString sets the sqlite3 connection string used to open a
// connection to the DB. Ordinarily configured via fstagger config but
// used in testing.
func WithConnectionString(connectionString string) func(*TagDB) {
	return func(t *TagDB) {
		t.connectionString = connectionString
	}
}

// WithMigrationsFS is used exclusively for testing this package. The default
// embedded migration filesystem is never going to be overridden.
func WithMigrationsFS(migrationsFS embed.FS) func(*TagDB) {
	return func(t *TagDB) {
		t.migrationsFS = migrationsFS
	}
}

func WithMigrationsDir(migrationsDir string) func(*TagDB) {
	return func(t *TagDB) {
		t.migrationsDir = migrationsDir
	}
}

// Init prepares fstagger's database for use. It will create a new database if
// one doesn't already exist and run migrations if required. It will error
// it is unable to create or open the configured file OR if it's unable to
// successfully run migrations to the latest version.
func Init(ctx context.Context, options ...func(*TagDB)) error {
	tagDB = New(options...)

	db, err := sql.Open("sqlite3", tagDB.connectionString)
	if err != nil {
		return err
	}

	tagDB.client = db

	goose.SetBaseFS(tagDB.migrationsFS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		Close(ctx)
		return err
	}

	if err := goose.Up(tagDB.client, "migrations"); err != nil {
		Close(ctx)
		return err
	}

	return nil
}

// Close will clean up the connection to the DB
func Close(ctx context.Context) {
	tagDB.client.Close()
}
