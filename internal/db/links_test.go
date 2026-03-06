package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/whatsfordinner/fstagger/internal/links"
)

func TestTagDBAddLinks(t *testing.T) {
	testMap := map[string]struct {
		shouldErr bool
		input     []links.Link
		expect    []links.Link
	}{
		"adding no links": {
			false,
			[]links.Link{},
			[]links.Link{},
		},
		"adding a link": {
			false,
			[]links.Link{
				{
					File: 2,
					Tag:  2,
				},
			},
			[]links.Link{
				{
					File: 2,
					Tag:  2,
				},
			},
		},
		"adding multiple links": {
			false,
			[]links.Link{
				{
					File: 1,
					Tag:  2,
				},
				{
					File: 2,
					Tag:  1,
				},
			},
			[]links.Link{
				{
					File: 1,
					Tag:  2,
				},
				{
					File: 2,
					Tag:  1,
				},
			},
		},
		"adding a link that already exists": {
			true,
			[]links.Link{
				{

					File: 1,
					Tag:  1,
				},
			},
			[]links.Link{},
		},
		"adding multiple links and one exists": {
			true,
			[]links.Link{
				{
					File: 2,
					Tag:  2,
				},
				{
					File: 1,
					Tag:  1,
				},
			},
			[]links.Link{
				{
					File: 2,
					Tag:  2,
				},
			},
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB, teardown := setupDB(t, []string{"fixtures/add_links.yml"})
			defer teardown()

			res, err := testDB.AddLinks(context.Background(), testData.input)

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
