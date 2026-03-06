package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/whatsfordinner/fstagger/internal/files"
	"github.com/whatsfordinner/fstagger/internal/links"
	"github.com/whatsfordinner/fstagger/internal/tags"

	"github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/otel/codes"
)

func (tagDB *TagDB) AddLinks(ctx context.Context, newLinks []links.Link) ([]links.Link, error) {
	const (
		insertString = "INSERT INTO filetags(fileid, tagid) VALUES(?, ?) RETURNING fileid"
	)
	ctx, span := tracer.Start(ctx, "AddLinks")
	defer span.End()

	addedLinks := []links.Link{}

	tx, err := tagDB.client.BeginTx(ctx, nil)
	if err != nil {
		return []links.Link{}, err
	}

	txErrors := []error{}

	for _, newLink := range newLinks {
		span.AddEvent(fmt.Sprintf("adding tag ID %d to file ID %d", newLink.Tag, newLink.File))
		row := tx.QueryRowContext(ctx, insertString, newLink.File, newLink.Tag)
		var fileId int64
		err := row.Scan(&fileId)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					collisionTag, err := tagDB.GetTagById(ctx, newLink.Tag)
					if err != nil {
					}
					collisionFile, err := tagDB.GetFileById(ctx, newLink.File)
					if err != nil {
					}
					collisionErr := fmt.Errorf(
						"file at path %s already tagged with %s",
						collisionFile.Path,
						collisionTag.Name,
					)
					txErrors = append(txErrors, collisionErr)
					continue
				}
			}
			txErrors = append(txErrors, err)
			continue
		}
		addedLinks = append(addedLinks, newLink)
	}
	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		return []links.Link{}, err
	}

	span.SetStatus(codes.Ok, "")
	return addedLinks, errors.Join(txErrors...)
}

func (tagDB *TagDB) GetLinksForFile(ctx context.Context, targetFile files.File) ([]links.Link, error) {
	return nil, nil
}

func (tagDB *TagDB) GetLinksForTag(ctx context.Context, targetTag tags.Tag) ([]links.Link, error) {
	return nil, nil
}
