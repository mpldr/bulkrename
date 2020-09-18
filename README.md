# bulkrename

## Installation

```
go get -u -v gitlab.com/poldi1405/bulkrename
```

## Usage

1. To move and rename all Episodes from the series Psycho-Pass:

```
br -r ~/NAS/Multimedia/Series/Anime/Psycho-Pass/
```

2. To automatically rename all .xml files to \*.drawio

```
br --editor vim --arg '+%s/\.xml$/\.drawio/' --arg +x --arg '{}' ~/DrawIO/*
```

This executes `vim +%s/\.xml$/\.drawio/ +x [Tempfile]`

## Moving forward

- GUI (for those desktop peasants)
- check for:
	- loops through symlinks
	- find redundant scans (`~/abc/` when `~/` is also scanned)
- speed up scans through concurrency
