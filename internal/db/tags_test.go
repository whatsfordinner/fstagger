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

func TestTagDBDeleteTags(t *testing.T) { t.Fatal("test not implemented") }
func TestTagDBUpdateTags(t *testing.T) { t.Fatal("test not implemented") }

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
