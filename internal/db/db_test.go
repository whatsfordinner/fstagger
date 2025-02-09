package db

import (
	"context"
	"testing"
)

// TestTagDBInit does NOT test production functionality of migrations. It only
// tests that the Init function does what it says on the tin:
// 1. Open a given DB file
// 2. Run migrations
func TestTagDBInit(t *testing.T) {
	testMap := map[string]struct {
		migrationsDir string
		shouldErr     bool
		errString     string
	}{
		"valid configuration": {
			migrationsDir: "fixtures/migrations",
			shouldErr:     false,
			errString:     "",
		},
		"invalid migrations directory": {
			migrationsDir: "does/not/exist",
			shouldErr:     true,
			errString:     "does/not/exist directory does not exist",
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB := New(
				WithMigrationsDir(testData.migrationsDir),
				WithMigrationsFS(nil), // use the os filesystem
			)
			err := testDB.Init(context.Background())
			defer testDB.Close(context.Background())

			if err == nil && testData.shouldErr {
				t.Fatalf("Expected error but got no error")
			}

			if err != nil && !testData.shouldErr {
				t.Fatalf("Expected no error but got: %s", err.Error())
			}

			if err != nil && err.Error() != testData.errString {
				t.Fatalf(
					"Expected error: %s but got: %s",
					testData.errString,
					err.Error(),
				)
			}

			if err == nil && !testData.shouldErr {
				_, err = testDB.client.Query("SELECT * FROM test")
				if err != nil {
					t.Fatalf("Error while running test query: %s", err.Error())
				}
			}
		})
	}
}

func TestTagDBClose(t *testing.T) {
	testMap := map[string]struct {
		initDB bool
	}{
		"database initialised": {
			initDB: true,
		},
		"database uninitialised": {
			initDB: false,
		},
	}

	for testName, testData := range testMap {
		t.Run(testName, func(t *testing.T) {
			testDB := New(
				WithMigrationsDir("migrations_test"),
				WithMigrationsFS(nil),
			)

			if testData.initDB {
				testDB.Init(context.Background())
			}

			testDB.Close(context.Background())
		})
	}
}
