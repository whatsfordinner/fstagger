-- +goose Up
CREATE TABLE IF NOT EXISTS test(
	foo TEXT NOT NULL
);

-- +goose Down
DROP TABLE test;
