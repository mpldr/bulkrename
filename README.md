# bulkrename

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/poldi1405/bulkrename)](https://goreportcard.com/report/gitlab.com/poldi1405/bulkrename)

[Report a bug](https://todo.sr.ht/~poldi1405/issues/?title=bulkrename:%20&description=**Type%3A**+BUG%0A%0A**Version%3A**+%0A%0A%3C!--+insert+description+here+--%3E)
or
[Request a feature](https://todo.sr.ht/~poldi1405/issues/?title=bulkrename:%20&description=**Type%3A**+Feature+Request%0A%0A**Version%3A**+%0A%0A%3C!--+insert+description+here+--%3E)

## Installation

```
go install mpldr.codes/br@latest
```

## Usage

1. To move and rename all Episodes from the series Psycho-Pass:

```
br -r ~/NAS/Multimedia/Series/Anime/Psycho-Pass/
```

2. To automatically rename all .xml files to \*.drawio

```
br --editor sed --arg=-i --arg 's/\.xml$/\.drawio/' --arg '{}' ~/DrawIO/*
```

This executes `sed -i 's/\.xml$/\.drawio/' [Tempfile]`

## Moving forward

- GUI (for those desktop peasants)
- check for:
	- loops through symlinks
	- find redundant scans (`~/abc/` when `~/` is also scanned)
- speed up scans through concurrency
