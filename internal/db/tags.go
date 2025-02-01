package db

import (
	"context"

	"github.com/whatsfordinner/fstagger/internal/tags"
)

func AddTags(ctx context.Context, tag []*tags.Tag) error {
	return nil
}

func DeleteTags(ctx context.Context, tagNames []string) error {
	return nil
}

func UpdateTags(ctx context.Context, tags []*tags.Tag) error {
	return nil
}

func GetTags(ctx context.Context) ([]*tags.Tag, error) {
	return nil, nil
}
