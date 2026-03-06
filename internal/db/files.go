package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/whatsfordinner/fstagger/internal/files"

	"github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/otel/codes"
)

// AddFiles takes a slice of files, tries adding them to the datastore and return a slice of
// files. It ignores any IDs in the input slice. The output slice is the same files with the
// IDs assigned to them in the datastore. If a file with the same path OR the same hash already
// exists in the datastore then AddFiles will return an error instead of writing that file. It
// will attempt to write every file and return any errors it encounters as a joined error. The
// reason this behaves differently to AddTags is that a new file could collide on path OR on hash
// and it's impossible to know which was the intended one to keep.
func (tagDB *TagDB) AddFiles(ctx context.Context, newFiles []files.File) ([]files.File, error) {
	const (
		insertString     = "INSERT INTO files(path, hash) VALUES (?, ?) RETURNING id"
		searchPathString = "SELECT id, path, hash FROM files WHERE path = ?"
		searchHashString = "SELECT id, path, hash FROM files WHERE hash = ?"
	)

	ctx, span := tracer.Start(ctx, "AddFile")
	defer span.End()

	addedFiles := []files.File{}

	tx, err := tagDB.client.BeginTx(ctx, nil)
	if err != nil {
		return []files.File{}, err
	}
	txErrors := []error{}

	for _, newFile := range newFiles {
		row := tx.QueryRowContext(ctx, insertString, newFile.Path, newFile.Hash)
		var fileId int64
		err := row.Scan(&fileId)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					var collidingFile files.File
					if strings.Contains(sqliteErr.Error(), "path") {
						row := tx.QueryRowContext(ctx, searchPathString, newFile.Path)
						if err := row.Scan(&collidingFile); err != nil {
							txErrors = append(txErrors, err)
							continue
						}
						collisionErr := fmt.Errorf(
							"file at path %s already being tracked: %w",
							collidingFile.Path,
							sqliteErr,
						)
						txErrors = append(txErrors, collisionErr)
					} else if strings.Contains(sqliteErr.Error(), "hash") {
						row := tx.QueryRowContext(ctx, searchHashString, newFile.Hash)
						if err := row.Scan(&collidingFile); err != nil {
							txErrors = append(txErrors, err)
							continue
						}
						collisionErr := fmt.Errorf(
							"file with hash %s already being tracked: %w",
							collidingFile.Hash,
							sqliteErr,
						)
						txErrors = append(txErrors, collisionErr)
					}
				}
			}
			txErrors = append(txErrors, err)
			continue
		}
		newFile.Id = int(fileId)
		addedFiles = append(addedFiles, newFile)
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		return []files.File{}, err
	}

	span.SetStatus(codes.Ok, "")
	return addedFiles, errors.Join(txErrors...)
}

func (tagDB *TagDB) GetFileById(ctx context.Context, search int) (files.File, error) {
	ctx, span := tracer.Start(ctx, "GetFileById")
	defer span.End()

	ret := files.File{}
	queryString := "SELECT id, path, hash FROM files WHERE id = ?"
	row := tagDB.client.QueryRowContext(ctx, queryString, search)
	if err := row.Scan(&ret.Id, &ret.Path, &ret.Hash); err != nil {
	}

	return ret, nil
}
