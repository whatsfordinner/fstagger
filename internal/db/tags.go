package db

import (
	"context"
	"fmt"

	"github.com/whatsfordinner/fstagger/internal/tags"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// AddTags takes a slice of tags and tries adding them to the database. If a tag
// with the same name already exists it will be skipped. IDs are ignored. To
// change the details of an already existing tag use UpdateTags. The slice returned
// is the tags that were actually added along with their correct IDs.
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
		}
	}

	span.SetStatus(codes.Ok, "")
	return tagsAdded, nil
}

// DeleteTags takes a slice of tags and removes them from the database. If a
// tag doesn't match on ID then it will delete what it can and return an error.
func (tagDB *TagDB) DeleteTags(ctx context.Context, tags []tags.Tag) error {
	return nil
}

// UpdateTags takes a slice of tags and updates the tags with matching IDs. If a
// tag with a provided ID doesn't exist then it will update what it can and
// return an error.
func (tagDB *TagDB) UpdateTags(ctx context.Context, tags []tags.Tag) error {
	return nil
}

// GetTags returns a slice of all tags being tracked. There's no pagination on
// this right now because it's not expected to get way out of control for
// someone's local collection.
func (tagDB *TagDB) GetTags(ctx context.Context) ([]tags.Tag, error) {
	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	queryString := "SELECT * FROM tags"

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

// GetTagsByName returns a list of all tags which match a provided search string. Right now
// it only supports a % wildcard but that could be upgraded to use something more
// substantial later.
func (tagDB *TagDB) GetTagsByName(ctx context.Context, search string) ([]tags.Tag, error) {
	ctx, span := tracer.Start(ctx, "GetTagsByName")
	defer span.End()

	queryString := "SELECT * FROM tags WHERE name LIKE ?"

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
