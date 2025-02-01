# Filesystem Tagger Database Documentation

## Entity Relationship

```mermaid
erDiagram
    FILES ||--o{ FILETAGS : tagged
    TAGS ||--o{ FILETAGS : tags

    FILES {
        INTEGER id PK
        TEXT path
        TEST hash
    }

    TAGS {
        INTEGER id PK
        TEXT name
        TEXT description
    }

    FILETAGS {
        INTEGER fileid FK
        INTEGER tagid FK
    }
```

## Notes

* the file and tag IDs are named ROWIDs that can be used as a composite primary key with the tags
* `files.hash` to be used for re-scanning a file if it's been moved
