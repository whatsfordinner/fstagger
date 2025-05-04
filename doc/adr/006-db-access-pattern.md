# Title

Decision to use batch operations for DAO

# Status

Active

# Date

2025-04-03

# Context

For every action with the DB it's possible that the user could be wanting to work with one-to-many entities. E.g. adding a single tag or adding multiple tags. Even for things that the existing user journeys don't necessarily support: i.e. bulk registering of files.

# Decision

The DAO will treat every operation as a bulk operation. Any mutation operation will accept a slice of inputs, attempt to process each element in a single transaction and return a slice of updated entities (where relevant) and a custom error that wraps a slice of errors which may have occurred for each entity.
