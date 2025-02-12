# Name

See all files with provided tag(s)

# Status

To Do

# Considerations



# Examples

## Input

```shell
fstagger search [TAGS...]
```

Could be one exact tag:

```shell
fstagger search food
```

Could be multiple exact tags:

```shell
fstagger search food pies
```

## Output

One tag should have all files with that tag:

```shell
$ fstagger search food
pie.jpg
burger.jpg
cookie.jpg
```

Multiple tags are logical `AND`:

```shell
$ fstagger search food dessert
pie.jpg
cookie.jpg
```

A search with no results is empty but a non-zero return code:

```shell
$ fstagger search does_not_exist

$ echo $?
1
```
