# Title

Decision to try and ensure files are unique 

# Status

Active

# Date

2025-04-04

# Context

`fstagger` tracks files in two ways: their path on the filesystem and the MD5 hash of the file's contents. In this way it's possible to check if a file might have moved or been updated. A moved file will have the same hash but different contents. An updated file with have the same path but a different hash. 

As much as possible I want files to be unique so that they only get tagged once. It'd be annoying to have two files with the same contents tagged separately and the goal of `fstagger` is to try and make it easier to find and manage things.

# Decision

Files that are being tagged need to be unique both in terms of their path _and_ hash. This will help make the user aware of duplicate files and help consolidate tags.

