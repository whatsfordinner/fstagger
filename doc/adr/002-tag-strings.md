# Title

This decision defines the acceptable strings for tags used by `fstagger`

# Status

Active

# Date

2025-01-26

# Context

Tags are meant to be a list of free-form metadata to let a user arbitrarily find files with the same tag (or combination of tags). That is, tags only have to be meaningful to the user and must be storeable in a JSON list in the SQLite datastore.

# Decision

Because tags only have to be meaningful to the user then it's the path of least resistance to let the user provide whatever value they want - even emoji if that's their thing
