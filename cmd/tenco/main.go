package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mix3/tenco"
)

var (
	defaultOffset    = -9
	defaultGenerator = &tenco.JsonGenerator{}
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage of tenco:
  tenco config1.yaml [config2.yaml ...]`)
	flag.PrintDefaults()
}

func run() error {
	var (
		outPath = flag.String("o", "-", "Write to FILE")
		offset  = flag.Int("offset", defaultOffset, "offset")
	)
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		return nil
	}

	var w io.WriteCloser
	switch *outPath {
	case "", "-":
		w = os.Stdout
	default:
		var err error
		w, err = os.Create(*outPath)
		if err != nil {
			return fmt.Errorf(`open file failed. %w`, err)
		}
		defer w.Close()
	}

	t, err := tenco.LoadWithEnv(args...)
	if err != nil {
		return err
	}

	return t.Write(w, *offset, defaultGenerator)
}
