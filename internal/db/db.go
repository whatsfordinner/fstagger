package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/XSAM/otelsql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	name                    = "github.com/whatsfordinner/fstagger/internal/db"
	version                 = "0.0.1"
	defaultMigrationsDir    = "migrations"
	defaultConnectionString = ":memory:"
)

var (
	//go:embed migrations/*.sql
	defaultMigrationsFS embed.FS
	tracer              trace.Tracer
)

type TagDB struct {
	client           *sql.DB
	connectionString string
	migrationsFS     fs.FS
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

	tagDB.connectionString = tagDB.connectionString + "?_foreign_keys=true"

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
func WithMigrationsFS(migrationsFS fs.FS) func(*TagDB) {
	return func(t *TagDB) {
		t.migrationsFS = migrationsFS
	}
}

// WithMigrationsDir is used exclusively for testing this package. The default
// directory in the embedded fielsystem is never going to be overrideen.
func WithMigrationsDir(migrationsDir string) func(*TagDB) {
	return func(t *TagDB) {
		t.migrationsDir = migrationsDir
	}
}

// Init prepares fstagger's database for use. It will create a new database if
// one doesn't already exist and run migrations if required. It will error
// it is unable to create or open the configured file OR if it's unable to
// successfully run migrations to the latest version.
func (tagDB *TagDB) Init(ctx context.Context) error {
	tracer = otel.GetTracerProvider().Tracer(name)
	ctx, span := tracer.Start(ctx, "Init")
	defer span.End()

	db, err := otelsql.Open("sqlite3", tagDB.connectionString)
	if err != nil {
		span.SetStatus(
			codes.Error,
			fmt.Sprintf("unable to open sqlite3 file at %s", tagDB.connectionString),
		)
		return err
	}
	span.AddEvent(
		fmt.Sprintf("unable to open sqlite3 file at %s", tagDB.connectionString),
	)

	tagDB.client = db

	goose.SetBaseFS(tagDB.migrationsFS)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		tagDB.Close(ctx)
		span.SetStatus(
			codes.Error,
			err.Error(),
		)
		return err
	}

	if err := goose.Up(tagDB.client, tagDB.migrationsDir); err != nil {
		tagDB.Close(ctx)
		span.SetStatus(
			codes.Error,
			err.Error(),
		)
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// Close will clean up the connection to the DB if one's been established
func (tagDB *TagDB) Close(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "Close")
	defer span.End()

	if tagDB.client != nil {
		tagDB.client.Close()
		span.AddEvent(
			fmt.Sprintf("closed sqlite3 file at %s", tagDB.connectionString),
		)
	}

	span.SetStatus(codes.Ok, "")
}
