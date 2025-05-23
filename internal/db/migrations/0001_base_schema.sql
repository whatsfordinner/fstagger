-- +goose Up
CREATE TABLE IF NOT EXISTS files(
	id INTEGER PRIMARY KEY,
	path TEXT UNIQUE NOT NULL,
	hash TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS tags(
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT
);

CREATE TABLE IF NOT EXISTS filetags(
	fileid INTEGER NOT NULL,
	tagid INTEGER NOT NULL,
	FOREIGN KEY(fileid) REFERENCES files(id) ON DELETE CASCADE,
	FOREIGN KEY(tagid) REFERENCES tags(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE files;
DROP TABLE tags;
DROP TABLE filetags;
