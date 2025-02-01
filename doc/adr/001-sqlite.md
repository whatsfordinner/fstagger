# Title

This decision defines the storage used for `fstagger`

# Status

Active

# Date

2025-01-25

# Context

I need a local datastore to track files that have been tagged and what tags have been added to them. A tag is an arbitrary string. Once tagged, a file can have one-or-more arbitrary tags. Users need to be able to:

* Add a tag or tags to a file
* Remove a tag or tags from a file
* List all the tags on a file
* List all the files with a tag or combination of tags

The datstore must not need internet (or any network) access and ideally won't need any additional daemons or services.

# Decision

I'll use SQLite for the datastore. A relational datastore is the best fit for these sorts of clearly defined relationships. Additionally, SQLite only needs to store a single file in a directory the user can write to. I can provide a sensible default location or let the user configure it with a dotfile or environment variable.
