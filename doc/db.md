# Filesystem Tagger Database Documentation

## Entity Relationship

```mermaid
erDiagram
    FILE ||--o{ FILETAGS
    TAG ||--o{ FILETAGS

    FILE {
        INTEGER id PK
        TEXT path
        TEST hash
    }

    TAG {
        INTEGER id PK
        TEXT name
    }

    FILETAGS {
        INTEGER fildId FK,PK
        INTEGER tagId FK,PK
    }
```

## Notes

* the file and tag IDs are named ROWIDs that can be used as a composite primary key with the tags
* `file.hash` to be used for re-scanning a file if it's been moved
