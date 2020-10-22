# bulkrename

![codecov](https://codecov.io/gl/poldi1405/bulkrename/branch/develop/graph/badge.svg?token=656MXKJG7U)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=poldi1405_bulkrename&metric=alert_status)](https://sonarcloud.io/dashboard?id=poldi1405_bulkrename)

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/poldi1405/bulkrename)](https://goreportcard.com/report/gitlab.com/poldi1405/bulkrename)

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
