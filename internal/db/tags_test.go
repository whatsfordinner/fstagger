package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/whatsfordinner/fstagger/internal/tags"
)

func TestTagDBAddTags(t *testing.T) {
	testMap := map[string]struct {
		shouldErr bool
		input     []tags.Tag
		expect    []tags.Tag
	}{
		"adding no tags": {
			false,
			[]tags.Tag{},
			[]tags.Tag{},
		},
		"adding a tag": {
			false,
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"adding multiple tags": {
			false,
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
				{
					Id:          3,
					Name:        "baz",
					Description: "a baz",
				},
			},
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
				{
					Id:          3,
					Name:        "baz",
					Description: "a baz",
				},
			},
		},
		"adding multiple tags where some already exist": {
			false,
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"adding the same tag multiple times": {
			false,
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"returning the correct id": {
			false,
			[]tags.Tag{
				{
					Id:          1,
					Name:        "bar",
					Description: "a bar",
				},
			},
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"returning correct id when no id provided": {
			false,
			[]tags.Tag{
				{
					Name: "bar",
				},
			},
			[]tags.Tag{
				{
					Id:   2,
					Name: "bar",
				},
			},
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, []string{"fixtures/add_tags.yml"})
			defer teardown()

			res, err := testDB.AddTags(context.Background(), testData.input)

			if err == nil && testData.shouldErr {
				t.Fatal("Expected error but got no error")
			}

			if err != nil && !testData.shouldErr {
				t.Fatalf("Expected no error but got: %s", err.Error())
			}

			if !reflect.DeepEqual(res, testData.expect) {
				t.Fatalf(
					"Result did not match expectation\nResult: %+v\nExpected: %+v",
					res,
					testData.expect,
				)
			}
		})
	}
}

func TestTagDBDeleteTags(t *testing.T) {
	testMap := map[string]struct {
		shouldErr bool
		input     []tags.Tag
		expect    []tags.Tag
	}{
		"delete nothing": {
			false,
			[]tags.Tag{},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
				{
					Id:          3,
					Name:        "baz",
					Description: "a baz",
				},
			},
		},
		"delete something that doesn't exist": {
			false,
			[]tags.Tag{
				{
					Id:          4,
					Name:        "qux",
					Description: "a qux",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
				{
					Id:          3,
					Name:        "baz",
					Description: "a baz",
				},
			},
		},
		"delete one thing": {
			false,
			[]tags.Tag{
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          3,
					Name:        "baz",
					Description: "a baz",
				},
			},
		},
		"delete many things": {
			false,
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
				{
					Id:          3,
					Name:        "baz",
					Description: "a baz",
				},
			},
			[]tags.Tag{},
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, []string{"fixtures/delete_tags.yml"})
			defer teardown()

			err := testDB.DeleteTags(context.Background(), testData.input)

			if err == nil && testData.shouldErr {
				t.Fatal("Expected error but got no error")
			}

			if err != nil && !testData.shouldErr {
				t.Fatalf("Expected no error but got: %s", err.Error())
			}

			res, err := testDB.GetTags(context.Background())
			if err != nil {
				t.Fatalf("Error retrieving remaining tags: %s", err.Error())
			}

			if !reflect.DeepEqual(res, testData.expect) {
				t.Fatalf(
					"Result did not match expectation\nResult: %+v\nExpected: %+v",
					res,
					testData.expect,
				)
			}
		})
	}
}

func TestTagDBUpdateTags(t *testing.T) {
	testMap := map[string]struct {
		shouldErr bool
		input     []tags.Tag
		expect    []tags.Tag
		expectDB  []tags.Tag
	}{
		"no updates": {
			false,
			[]tags.Tag{},
			[]tags.Tag{},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"one valid update": {
			false,
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"multiple valid updates": {
			false,
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
				{
					Id:          2,
					Name:        "qux",
					Description: "a qux",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
				{
					Id:          2,
					Name:        "qux",
					Description: "a qux",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
				{
					Id:          2,
					Name:        "qux",
					Description: "a qux",
				},
			},
		},
		"tag ID doesn't exist": {
			true,
			[]tags.Tag{
				{
					Id:          3,
					Name:        "baz",
					Description: "a baz",
				},
			},
			[]tags.Tag{},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"tag with same name already exists": {
			true,
			[]tags.Tag{
				{
					Id:          2,
					Name:        "foo",
					Description: "a foo",
				},
			},
			[]tags.Tag{},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
		"one valid, one invalid update": {
			true,
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
				{
					Id:          3,
					Name:        "qux",
					Description: "a qux",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
			},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "baz",
					Description: "a baz",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, []string{"fixtures/update_tags.yml"})
			defer teardown()

			res, err := testDB.UpdateTags(context.Background(), testData.input)

			if err == nil && testData.shouldErr {
				t.Fatal("Expected error but got no error")
			}

			if err != nil && !testData.shouldErr {
				t.Fatalf("Expected no error but got: %s", err.Error())
			}

			if !reflect.DeepEqual(res, testData.expect) {
				t.Fatalf(
					"Result did not match expectation\nResult: %+v\nExpected: %+v",
					res,
					testData.expect,
				)
			}

			dbRes, err := testDB.GetTags(context.Background())
			if err != nil {
				t.Fatalf("Error retrieving tags in DB: %s", err.Error())
			}

			if !reflect.DeepEqual(dbRes, testData.expectDB) {
				t.Fatalf(
					"DB contents did not match expectation\nResult: %+v\nExpected: %+v",
					res,
					testData.expectDB,
				)
			}
		})
	}
}

func TestTagDBGetTags(t *testing.T) {
	testMap := map[string]struct {
		shouldErr    bool
		fixtureFiles []string
		expect       []tags.Tag
	}{
		"no tags": {
			false,
			[]string{"fixtures/get_tags_no_tags.yml"},
			[]tags.Tag{},
		},
		"one tag": {
			false,
			[]string{"fixtures/get_tags_single_tag.yml"},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
			},
		},
		"some tags": {
			false,
			[]string{"fixtures/get_tags_many_tags.yml"},
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          2,
					Name:        "bar",
					Description: "a bar",
				},
			},
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, testData.fixtureFiles)
			defer teardown()

			res, err := testDB.GetTags(context.Background())

			if err == nil && testData.shouldErr {
				t.Fatal("Expected error but got no error")
			}

			if err != nil && !testData.shouldErr {
				t.Fatalf("Expected no error but got: %s", err.Error())
			}

			if !reflect.DeepEqual(res, testData.expect) {
				t.Fatalf(
					"Result did not match expectation\nResult: %+v\nExpected: %+v",
					res,
					testData.expect,
				)
			}
		})
	}
}

func TestTagDBGetTagById(t *testing.T) {
	testMap := map[string]struct {
		shouldErr bool
		input     int
		expect    tags.Tag
		expectErr string
	}{
		"tag exists": {
			false,
			1,
			tags.Tag{
				Id:          1,
				Name:        "foo",
				Description: "a foo",
			},
			"",
		},
		"tag doesn't exist": {
			true,
			2,
			tags.Tag{},
			"tag does not exist with id: 2",
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, []string{"fixtures/get_tag_by_id.yml"})
			defer teardown()

			res, err := testDB.GetTagById(context.Background(), testData.input)

			if err == nil && testData.shouldErr {
				t.Fatal("Expected error but got no error")
			}

			if err != nil && !testData.shouldErr {
				t.Fatalf("Expected no error but got: %s", err.Error())
			}

			if err != nil && err.Error() != testData.expectErr {
				t.Fatalf(
					"Expected error: %s but got: %s",
					testData.expectErr,
					err.Error(),
				)
			}

			if !reflect.DeepEqual(res, testData.expect) {
				t.Fatalf(
					"Result did not match expectation\nResult: %+v\nExpected: %+v",
					res,
					testData.expect,
				)
			}
		})
	}
}

func TestTagDBGetTagsByName(t *testing.T) {
	testMap := map[string]struct {
		shouldErr    bool
		searchString string
		expect       []tags.Tag
	}{
		"exact no match": {
			false,
			"qux",
			[]tags.Tag{},
		},
		"exact one match": {
			false,
			"foo",
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
			},
		},
		"wildcard no match": {
			false,
			"%qux%",
			[]tags.Tag{},
		},
		"wildcard many matches": {
			false,
			"foo%",
			[]tags.Tag{
				{
					Id:          1,
					Name:        "foo",
					Description: "a foo",
				},
				{
					Id:          3,
					Name:        "foobar",
					Description: "a foobar",
				},
				{
					Id:          5,
					Name:        "foobarbaz",
					Description: "a foobarbaz",
				},
			},
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, []string{"fixtures/get_tags_by_name.yml"})
			defer teardown()

			res, err := testDB.GetTagsByName(context.Background(), testData.searchString)

			if err == nil && testData.shouldErr {
				t.Fatal("Expected error but got no error")
			}

			if err != nil && !testData.shouldErr {
				t.Fatalf("Expected no error but got: %s", err.Error())
			}

			if !reflect.DeepEqual(res, testData.expect) {
				t.Fatalf(
					"Result did not match expectation\nResult: %+v\nExpected: %+v",
					res,
					testData.expect,
				)
			}
		})
	}
}

func setupDB(t *testing.T, fixtureFiles []string) (*TagDB, func()) {
	testDB := New()
	if err := testDB.Init(context.Background()); err != nil {
		t.Fatalf("Unable to init test DB: %s", err.Error())
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(testDB.client),
		testfixtures.Dialect("sqlite3"),
		testfixtures.FilesMultiTables(fixtureFiles...),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)

	if err != nil {
		t.Fatalf("Unable to create fixture manager: %s", err.Error())
	}

	if err := fixtures.Load(); err != nil {
		t.Fatalf("Unable to prep database: %s", err.Error())
	}

	return testDB, func() {
		testDB.Close(context.Background())
	}
}
