# go-xunit-reader [![CircleCI](https://circleci.com/gh/josegonzalez/go-xunit-reader.svg?style=svg)](https://circleci.com/gh/josegonzalez/go-xunit-reader)

A tool for reading xunit xml output files.

## Installation

Install it using the "go get" command:

    go get github.com/josegonzalez/go-xunit-reader

## Usage

```
# xunit-reader

usage: xunit-reader [-h|--help] [-v|--version] [-s|--show-skipped] -p|--pattern
                    "<value>"

                    Prints provided string to stdout

Arguments:

  -h  --help          Print help information
  -v  --version       show version
  -s  --show-skipped  show skipped tests
  -p  --pattern       pattern referencing files to process
```
