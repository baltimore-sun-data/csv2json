# csv2json

Simple tool for converting CSVs to JSON

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```shell
GOBIN=$(pwd) GOPATH=$(mktemp -d) go get github.com/baltimore-sun-data/csv2json
```

## Usage

```shell
$ csv2json
Usage of csv2json:
  -dest string
        Destination file (default: stdout)
  -no-headers
        Return each row as an array
  -src string
        Source file (default: stdin)

$ more test.csv
a,b,c
1,2,3

$ csv2json | json-tidy
[
        {
                "a": "1",
                "b": "2",
                "c": "3"
        }
]

$ csv2json -src test.csv -no-headers | json-tidy
[
        [
                "a",
                "b",
                "c"
        ],
        [
                "1",
                "2",
                "3"
        ]
]
```
