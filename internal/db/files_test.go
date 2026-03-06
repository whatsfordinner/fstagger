package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/whatsfordinner/fstagger/internal/files"
)

func TestTagDBAddFiles(t *testing.T) {
	testMap := map[string]struct {
		shouldErr bool
		input     []files.File
		expect    []files.File
	}{
		"adding no files": {
			false,
			[]files.File{},
			[]files.File{},
		},
		"adding a file": {
			false,
			[]files.File{
				{
					Path: "/path/to/bar",
					Hash: "barhash",
				},
			},
			[]files.File{
				{
					Id:   2,
					Path: "/path/to/bar",
					Hash: "barhash",
				},
			},
		},
		"adding multiple files": {
			false,
			[]files.File{
				{
					Path: "/path/to/bar",
					Hash: "barhash",
				},
				{
					Path: "/path/to/baz",
					Hash: "bazhash",
				},
			},
			[]files.File{
				{
					Id:   2,
					Path: "/path/to/bar",
					Hash: "barhash",
				},
				{
					Id:   3,
					Path: "/path/to/baz",
					Hash: "bazhash",
				},
			},
		},
		"adding a file that already exists": {
			true,
			[]files.File{
				{
					Path: "/path/to/foo",
					Hash: "foohash",
				},
			},
			[]files.File{},
		},
		"adding multiple files and one exists": {
			true,
			[]files.File{
				{
					Path: "/path/to/foo",
					Hash: "foohash",
				},
				{
					Path: "/path/to/bar",
					Hash: "barhash",
				},
			},
			[]files.File{
				{
					Id:   2,
					Path: "/path/to/bar",
					Hash: "barhash",
				},
			},
		},
		"returning the correct id": {
			false,
			[]files.File{
				{
					Id:   8,
					Path: "/path/to/bar",
					Hash: "barhash",
				},
			},
			[]files.File{
				{
					Id:   2,
					Path: "/path/to/bar",
					Hash: "barhash",
				},
			},
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, []string{"fixtures/add_files.yml"})
			defer teardown()

			res, err := testDB.AddFiles(context.Background(), testData.input)

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
