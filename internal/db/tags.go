package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/whatsfordinner/fstagger/internal/tags"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// AddTags takes a slice of tags, tries adding them to the datastore and returns a slice of
// tags. It ignores any IDs in the input slice. The output slice is the same tags with
// the IDs assigned to them in the datastore. If a tag with the same name is supplied
// multiple times it will only be added once.
func (tagDB *TagDB) AddTags(ctx context.Context, newTags []tags.Tag) ([]tags.Tag, error) {
	ctx, span := tracer.Start(ctx, "AddTags")
	defer span.End()

	queryString := "INSERT INTO tags(name, description) VALUES(?, ?) RETURNING id"
	tagsAdded := []tags.Tag{}
	tagNamesAdded := []string{}
	tagNamesToAdd := []string{}

	span.SetAttributes(
		attribute.String("db.file", tagDB.connectionString),
		attribute.String("db.query", queryString),
	)

	for _, tag := range newTags {
		tagNamesToAdd = append(tagNamesToAdd, tag.Name)
		res, err := tagDB.GetTagsByName(ctx, tag.Name)
		if err != nil {
			// TODO: this should add an error to a slice and continue rather than bail
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		if len(res) == 0 {
			span.AddEvent(fmt.Sprintf("adding new tag: %s", tag.Name))
			row := tagDB.client.QueryRowContext(ctx, queryString, tag.Name, tag.Description)

			var tagId int64
			err := row.Scan(&tagId)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}

			tag.Id = int(tagId)
			tagsAdded = append(tagsAdded, tag)
			tagNamesAdded = append(tagNamesAdded, tag.Name)
		} else {
			span.AddEvent(fmt.Sprintf("tag already exists: %s", tag.Name))
			found := false
			for _, test := range tagsAdded {
				if test.Name == res[0].Name {
					found = true
					break
				}
			}

			if !found {
				tagsAdded = append(tagsAdded, res[0])
			}
		}
	}

	span.SetAttributes(
		attribute.StringSlice("db.tags_to_add", tagNamesToAdd),
		attribute.StringSlice("db.tags_added", tagNamesAdded),
	)
	span.SetStatus(codes.Ok, "")
	return tagsAdded, nil
}

// DeleteTags takes a slice of tags and removes them from the database.
func (tagDB *TagDB) DeleteTags(ctx context.Context, deleteTags []tags.Tag) error {
	ctx, span := tracer.Start(ctx, "DeleteTags")
	defer span.End()

	idsToDelete := []string{}
	tagNamesToDelete := []string{}
	for _, tag := range deleteTags {
		idsToDelete = append(idsToDelete, strconv.Itoa(tag.Id))
		tagNamesToDelete = append(tagNamesToDelete, tag.Name)
	}

	queryString := fmt.Sprintf(
		"DELETE FROM tags WHERE id IN(%s)",
		strings.Join(idsToDelete, ","),
	)

	span.SetAttributes(
		attribute.String("db.file", tagDB.connectionString),
		attribute.String("db.query", queryString),
		attribute.StringSlice("db.tags_to_delete", tagNamesToDelete),
	)

	_, err := tagDB.client.ExecContext(ctx, queryString)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// UpdateTags takes a slice of tags and updates the tags with matching IDs. If a
// tag with a provided ID doesn't exist then it will update what it can and
// return an error.
func (tagDB *TagDB) UpdateTags(ctx context.Context, updateTags []tags.Tag) ([]tags.Tag, error) {
	ctx, span := tracer.Start(ctx, "UpdateTags")
	defer span.End()

	queryString := "UPDATE tags SET name = ?, description = ? WHERE id = ?"
	queryErrors := []error{}
	updatedTags := []tags.Tag{}

	span.SetAttributes(
		attribute.String("db.file", tagDB.connectionString),
		attribute.String("db.query", queryString),
	)

	for _, tag := range updateTags {
		_, err := tagDB.GetTagById(ctx, tag.Id)
		if err != nil {
			queryErrors = append(queryErrors, err)
			continue
		}

		if _, err = tagDB.client.ExecContext(
			ctx,
			queryString,
			tag.Name,
			tag.Description,
			tag.Id,
		); err != nil {
			queryErrors = append(queryErrors, err)
			continue
		}

		updatedTags = append(updatedTags, tag)
	}

	if len(queryErrors) > 0 {
		returnErr := errors.Join(queryErrors...)
		span.SetStatus(codes.Error, returnErr.Error())
		return updatedTags, returnErr
	}

	span.SetStatus(codes.Ok, "")
	return updatedTags, nil
}

// GetTags returns a slice of all tags being tracked. There's no pagination on
// this right now because it's not expected to get way out of control for
// someone's local collection.
func (tagDB *TagDB) GetTags(ctx context.Context) ([]tags.Tag, error) {
	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	queryString := "SELECT id, name, description FROM tags"

	span.SetAttributes(
		attribute.String("db.file", tagDB.connectionString),
		attribute.String("db.query", queryString),
	)

	rows, err := tagDB.client.QueryContext(
		ctx,
		queryString,
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

	span.SetAttributes(attribute.Int("db.rows_returned", len(ret)))
	span.SetStatus(codes.Ok, "")
	return ret, nil
}

// GetTagById returns a single tag whose ID matches the input or an error if no tag has
// that ID
func (tagDB *TagDB) GetTagById(ctx context.Context, search int) (tags.Tag, error) {
	ctx, span := tracer.Start(ctx, "GetTagById")
	defer span.End()

	queryString := "SELECT id, name, description FROM tags WHERE id = ?"

	span.SetAttributes(
		attribute.String("db.file", tagDB.connectionString),
		attribute.String("db.query", queryString),
		attribute.Int("db.tag_id", search),
	)

	row := tagDB.client.QueryRowContext(ctx, queryString, search)
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
	ctx, span := tracer.Start(ctx, "GetTagsByName")
	defer span.End()

	queryString := "SELECT id, name, description FROM tags WHERE name LIKE ?"

	span.SetAttributes(
		attribute.String("db.file", tagDB.connectionString),
		attribute.String("db.query", queryString),
		attribute.String("db.search_value", search),
	)

	rows, err := tagDB.client.QueryContext(
		ctx,
		queryString,
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

	span.SetAttributes(
		attribute.Int("db.rows_returned", len(ret)),
	)

	span.SetStatus(codes.Ok, "")
	return ret, nil
}
