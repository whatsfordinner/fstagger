# Name

Add tag(s) to a file

# Status

To Do

# Considerations

* Need to make sure the file is in the DB -> Search by file hash and add it if it isn't
* Need to make sure the tag is in the DB -> Search by tag name and add it if it isn't
* File and tag relationship should be unique in DB -> Table constraint

# Required functionality

* Adding a file to the DB
    * Finding the absolute path to the file
    * Getting a hash of the file
    * Writing the file to the `files` table in the DB
* Finding if a tag already exists in the DB
    * Search DB for a tag with the same name as provided
* Adding a tag to the DB if it doesn't already exist
    * Should eventually be able to add a description for a new tag
* Linking a tag to a file

# Examples

```shell
fstagger tag add [FILENAME] [TAGS...]
```

Could be one tag:

```shell
fstagger tag add pie.jpg food
```

Or many tags:

```shell
fstagger tag add pie.jpg food dessert pies
```

Could be absolute paths too:

```shell
fstagger tag add /home/whatsfordinner/pictures/pie.jpg food
```

