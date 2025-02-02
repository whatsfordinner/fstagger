# Title

Decision to use closures for DAO dependency injection

# Status

Active

# Date

2025-02-02

# Context

[ADR-003](003-dao.md) described using a global DAO in the `db` package but while it was quick to implement initially it was difficult to test. Code should be easily testable so it wasn't a great fit.

# Decision

`fstagger` will use closures around the `cobra` commands to inject a DAO into the commands. It will be more complicated up front but it will make it easier to test database operations and reduce the amount of repeated code for initialising the DAO
