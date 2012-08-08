package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"os"
	"strings"
)

var (
	polynomial = flag.String("polynomial", "ieee", "polynomial: ieee, castagnoli, koopman")
	output     = flag.String("output", "dec", "output format: hex, dec, oct")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 || flag.Arg(0) == "" {
		fmt.Printf("usage: crc32 <file>\n")
		os.Exit(1)
	}
	filename := flag.Arg(0)

	var poly uint32
	switch strings.ToLower(*polynomial) {
	case "ieee":
		poly = crc32.IEEE
	case "castagnoli":
		poly = crc32.Castagnoli
	case "koopman":
		poly = crc32.Koopman
	default:
		fmt.Printf("unknown -polynomial %s\n", *polynomial)
		os.Exit(1)
	}

	var format string
	switch strings.ToLower(*output) {
	case "hex":
		format = "%x\n"
	case "dec":
		format = "%d\n"
	case "oct":
		format = "%o\n"
	default:
		fmt.Printf("unknown -output %s\n", *output)
		os.Exit(1)
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s: %s\n", filename, err)
		os.Exit(1)
	}
	defer f.Close()

	// http://blog.vzv.ca/2012/06/crc64-file-hash-in-gogolang.html
	h := crc32.New(crc32.MakeTable(poly))
	buf := make([]byte, 8192)
	read, err := f.Read(buf)
	for read > -1 && err == nil {
		h.Write(buf)
		read, err = f.Read(buf)
	}

	s := h.Sum32()
	fmt.Printf(format, s)
}
