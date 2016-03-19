# ghoto 
Photo manage cli command.

## Requirement

- exiftool 
- mv 

## Install

```bash
$ go get github.com/fukata/ghoto
```

## Usage

```bash
$ ghoto -from /path/to/src -to /path/to/dst -photo-dir photo_raw -video-dir video_raw
2016/03/19 13:10:11 /path/to/src/P3060621.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060621.ORF
2016/03/19 13:10:11 /path/to/src/P3060622.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060622.ORF
2016/03/19 13:10:11 /path/to/src/P3060623.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060623.ORF
2016/03/19 13:10:12 /path/to/src/P3060624.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060624.ORF
2016/03/19 13:10:12 /path/to/src/P3060625.ORF -> /path/to/dst/photo_raw/2016/03/06/P3060625.ORF
2016/03/19 13:10:12 /path/to/src/P3060626.MOV -> /path/to/dst/video_raw/2016/03/06/P3060626.MOV
```

## Help

```bash
$ ghoto -h
NAME:
   ghoto - Transfer photo(video)

USAGE:
   ghoto [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR(S):
   fukata <tatsuya.fukata@gmail.com>

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --from "/path/to/src"        Source directory
   --to "/path/to/dst"          Destination directory
   --photo-dir, -P "photo"      Destination photo directory
   --video-dir, -V "video"      Destination video directory
   --exclude, -x                Exclude dir/file separate comma.
   --concurrency, -c "8"        Concurrency num.
   --recursive, -r              Resursive
   --dry-run                    Dry Run
   --verbose, --vvv             Verbose
   --help, -h                   show help
   --version, -v                print the version
```
