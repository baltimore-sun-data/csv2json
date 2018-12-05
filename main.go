package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/carlmjohnson/errors"
	"github.com/carlmjohnson/flagext"
)

func main() {
	enc, err := FromArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initial error: %v\n", err)
		os.Exit(3)
	}
	if err = enc.Encode(); err != nil {
		fmt.Fprintf(os.Stderr, "Run time error: %v\n", err)
		os.Exit(1)
	}
}

type Encoder struct {
	src         io.ReadCloser
	dest        io.WriteCloser
	transformer func() error
}

func FromArgs(args []string) (*Encoder, error) {
	fl := flag.NewFlagSet("csv2json", flag.ExitOnError)

	src := flagext.FileOrURL(flagext.StdIO, nil)
	fl.Var(src, "src", "Source `path` for file or URL")
	dest := fl.String("dest", "", "Destination file (default: stdout)")
	array := fl.Bool("no-headers", false, "Return each row as an array")
	_ = fl.Parse(args)

	e := Encoder{
		src:  src,
		dest: os.Stdout,
	}
	if *dest != "" && *dest != "-" {
		f, err := os.Create(*dest)
		if err != nil {
			return nil, err
		}
		e.dest = f
	}

	if *array {
		e.transformer = e.NoHeaders
	} else {
		e.transformer = e.WithHeaders
	}
	return &e, nil
}

func (e *Encoder) Encode() error {
	return e.transformer()
}

func (e *Encoder) WithHeaders() (err error) {
	defer errors.Defer(&err, e.dest.Close)

	data, err := makeWithHeader(e.src)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(e.dest)
	return enc.Encode(&data)
}

func makeWithHeader(src io.ReadCloser) (data []map[string]string, err error) {
	defer errors.Defer(&err, src.Close)

	cr := csv.NewReader(src)
	cr.Comment = '#'
	cr.FieldsPerRecord = -1
	cr.ReuseRecord = true

	fields, err := cr.Read()

	// Save headers for each row of dict
	dataHeader := make(map[int]string, len(fields))
	for i, field := range fields {
		dataHeader[i] = field
	}

	for {
		fields, err = cr.Read()
		if err == io.EOF {
			return data, nil
		}

		if err != nil {
			return nil, err
		}

		datum := make(map[string]string, len(fields))
		for i, val := range fields {
			datum[dataHeader[i]] = val
		}
		data = append(data, datum)
	}
}

func (e *Encoder) NoHeaders() (err error) {
	defer errors.Defer(&err, e.dest.Close)

	data, err := makeWithoutHeader(e.src)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(e.dest)
	return enc.Encode(&data)
}

func makeWithoutHeader(src io.ReadCloser) (data [][]string, err error) {
	defer errors.Defer(&err, src.Close)

	cr := csv.NewReader(src)
	cr.Comment = '#'
	cr.FieldsPerRecord = -1

	return cr.ReadAll()
}
