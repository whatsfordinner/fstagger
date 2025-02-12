# Name

Add tag(s) to a file

# Status

To Do

# Considerations

* Need to make sure the file is in the DB -> Search by file hash and add it if it isn't
* Need to make sure the tag is in the DB -> Search by tag name and add it if it isn't
* File and tag relationship should be unique in DB -> Table constraint

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

