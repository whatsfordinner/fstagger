package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/whatsfordinner/fstagger/internal/tags"

	"github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/otel/codes"
)

// AddTags takes a slice of tags, tries adding them to the datastore and returns a slice of
// tags. It ignores any IDs in the input slice. The output slice is the same tags with
// the IDs assigned to them in the datastore. If a tag with the same name is supplied
// multiple times it will only be added once.
func (tagDB *TagDB) AddTags(ctx context.Context, newTags []tags.Tag) ([]tags.Tag, error) {
	const (
		insertString = "INSERT INTO tags(name, description) VALUES(?, ?) RETURNING id"
		searchString = "SELECT id, name, description FROM tags WHERE name = ?"
	)

	ctx, span := tracer.Start(ctx, "AddTags")
	defer span.End()

	returnTags := []tags.Tag{}

	tx, err := tagDB.client.BeginTx(ctx, nil)
	if err != nil {
		return []tags.Tag{}, err
	}
	txErrors := []error{}

	for _, tag := range newTags {
		tagAlreadyProcessed := false
		for _, addedTag := range returnTags {
			if tag.Name == addedTag.Name {
				tagAlreadyProcessed = true
				continue
			}
		}

		if tagAlreadyProcessed {
			continue
		}

		span.AddEvent(fmt.Sprintf("adding new tag: %s", tag.Name))
		row := tx.QueryRowContext(ctx, insertString, tag.Name, tag.Description)
		var tagId int64
		err := row.Scan(&tagId)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				switch sqliteErr.Code {
				case sqlite3.ErrConstraint:
					span.AddEvent(fmt.Sprintf("tag already exists: %s", tag.Name))
					searchRow := tx.QueryRowContext(
						ctx,
						searchString,
						tag.Name,
					)
					err := searchRow.Scan(
						&tag.Id,
						&tag.Name,
						&tag.Description,
					)
					if err != nil {
						txErrors = append(txErrors, err)
						continue
					}
				default:
					txErrors = append(txErrors, sqliteErr)
				}
			} else {
				txErrors = append(txErrors, err)
				continue
			}
		} else {
			tag.Id = int(tagId)
		}

		returnTags = append(returnTags, tag)
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}

		return []tags.Tag{}, err
	}

	if len(txErrors) > 0 {
		span.SetStatus(codes.Error, "encountered DB errors during transaction")
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return returnTags, errors.Join(txErrors...)
}

// DeleteTags takes a slice of tags and removes them from the database.
func (tagDB *TagDB) DeleteTags(ctx context.Context, deleteTags []tags.Tag) error {
	const (
		deleteString = "DELETE FROM tags WHERE id = ?"
	)

	ctx, span := tracer.Start(ctx, "DeleteTags")
	defer span.End()

	tx, err := tagDB.client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txErrors := []error{}

	for _, tag := range deleteTags {
		_, err := tx.ExecContext(ctx, deleteString, tag.Id)
		if err != nil {
			txErrors = append(txErrors, err)
		}
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		span.SetStatus(codes.Error, "encountered error finalising transaction")
		return err
	}

	if len(txErrors) > 0 {
		span.SetStatus(codes.Error, "encountered DB errors during transaction")
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return errors.Join(txErrors...)
}

// UpdateTags takes a slice of tags and updates the tags with matching IDs. If a
// tag with a provided ID doesn't exist then it will update what it can and
// return an error.
func (tagDB *TagDB) UpdateTags(ctx context.Context, updateTags []tags.Tag) ([]tags.Tag, error) {
	const (
		updateString = "UPDATE tags SET name = ?, description = ? WHERE id = ?"
		searchString = "SELECT id FROM tags WHERE id = ?"
	)

	ctx, span := tracer.Start(ctx, "UpdateTags")
	defer span.End()

	updatedTags := []tags.Tag{}

	tx, err := tagDB.client.BeginTx(ctx, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return []tags.Tag{}, err
	}
	txErrors := []error{}

	for _, tag := range updateTags {
		row := tx.QueryRowContext(ctx, searchString, tag.Id)
		var tagId int64
		err := row.Scan(&tagId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				span.AddEvent(fmt.Sprintf("tag doesn't exist in db: %s", tag.Name))
			}
			txErrors = append(txErrors, err)
			continue
		}

		if _, err = tx.ExecContext(
			ctx,
			updateString,
			tag.Name,
			tag.Description,
			tag.Id,
		); err != nil {
			txErrors = append(txErrors, err)
			continue
		}

		updatedTags = append(updatedTags, tag)
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		span.SetStatus(codes.Error, err.Error())
		return []tags.Tag{}, err
	}

	if len(txErrors) > 0 {
		span.SetStatus(codes.Error, "")
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return updatedTags, errors.Join(txErrors...)
}

// GetTags returns a slice of all tags being tracked. There's no pagination on
// this right now because it's not expected to get way out of control for
// someone's local collection.
func (tagDB *TagDB) GetTags(ctx context.Context) ([]tags.Tag, error) {
	const (
		searchString = "SELECT id, name, description FROM tags"
	)

	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	rows, err := tagDB.client.QueryContext(
		ctx,
		searchString,
	)
	if err != nil {
		span.SetStatus(
			codes.Error,
			err.Error(),
		)
		return nil, err
	}

	ret := []tags.Tag{}

	for rows.Next() {
		tag := tags.Tag{}
		if err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.Description,
		); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		ret = append(ret, tag)
	}

	span.SetStatus(codes.Ok, "")
	return ret, nil
}

// GetTagById returns a single tag whose ID matches the input or an error if no tag has
// that ID
func (tagDB *TagDB) GetTagById(ctx context.Context, search int) (tags.Tag, error) {
	const (
		searchString = "SELECT id, name, description FROM tags WHERE id = ?"
	)

	ctx, span := tracer.Start(ctx, "GetTagById")
	defer span.End()

	row := tagDB.client.QueryRowContext(ctx, searchString, search)
	ret := tags.Tag{}

	if err := row.Scan(&ret.Id, &ret.Name, &ret.Description); err != nil {
		if err == sql.ErrNoRows {
			span.SetStatus(codes.Error, "tag not found")
			return ret, errors.New(fmt.Sprintf("tag does not exist with id: %d", search))
		}

		span.SetStatus(codes.Error, err.Error())
		return ret, err
	}

	span.SetStatus(codes.Ok, "")
	return ret, nil
}

// GetTagsByName returns a list of all tags which match a provided search string. Right now
// it only supports a % wildcard but that could be upgraded to use something more
// substantial later.
func (tagDB *TagDB) GetTagsByName(ctx context.Context, search string) ([]tags.Tag, error) {
	const (
		searchString = "SELECT id, name, description FROM tags WHERE name LIKE ?"
	)

	ctx, span := tracer.Start(ctx, "GetTagsByName")
	defer span.End()

	rows, err := tagDB.client.QueryContext(
		ctx,
		searchString,
		search,
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	ret := []tags.Tag{}

	for rows.Next() {
		tag := tags.Tag{}
		if err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.Description,
		); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		ret = append(ret, tag)
	}

	span.SetStatus(codes.Ok, "")
	return ret, nil
}
