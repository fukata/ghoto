# ghoto

Photo data transfer tool.

## Status

[![.github/workflows/build.yml](https://github.com/fukata/ghoto/actions/workflows/build.yml/badge.svg)](https://github.com/fukata/ghoto/actions/workflows/build.yml)

## Requirement

- exiftool

## Install

```bash
$ go get github.com/fukata/ghoto
```

## Usage

```bash
$ ghoto --from /path/to/src --to /path/to/dst --photo-dir photo_raw --video-dir video_raw --recursive --exclude lightroom
2016/03/19 13:10:11 /path/to/src/P3060621.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060621.ORF
2016/03/19 13:10:11 /path/to/src/P3060622.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060622.ORF
2016/03/19 13:10:11 /path/to/src/P3060623.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060623.ORF
2016/03/19 13:10:12 /path/to/src/P3060624.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060624.ORF
2016/03/19 13:10:12 /path/to/src/P3060625.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060625.ORF
2016/03/19 13:10:12 /path/to/src/P3060626.MOV -> /path/to/dst/video_raw/2016/03/06/P3060626.MOV
```

## Help

```bash
$ ghoto --help

NAME:
   ghoto - Transfer photo(video)

USAGE:
   ghoto [global options] command [command options] [arguments...]

VERSION:
   0.0.5

AUTHOR:
   fukata <tatsuya.fukata@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --from value                   Source directory (default: "/path/to/src")
   --to value                     Destination directory (default: "/path/to/dst")
   --photo-dir value, -P value    Destination photo directory (default: "photo")
   --video-dir value, -V value    Destination video directory (default: "video")
   --exclude value, -x value      Exclude dir/file separate comma.
   --concurrency value, -c value  Concurrency num. (default: 8)
   --recursive, -r                Resursive (default: false)
   --force                        Force (default: false)
   --skip-invalid-data            SkipInvalidData (default: false)
   --dry-run                      Dry Run (default: false)
   --verbose                      Verbose (default: false)
   --help, -h                     show help (default: false)
   --version, -v                  print the version (default: false)
```

## Build

```bash
$ go install
$ go build
```
