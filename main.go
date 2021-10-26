package main

import (
	"flag"
	"fmt"
	"github.com/vldrus/golang/image/ico"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <flags> <input_file>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
	}

	var output string
	flag.StringVar(&output, "o", "test.ico", "output file name")

	flag.Parse()

	tail := flag.Args()

	if len(tail) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(tail[0])
	if err != nil {
		fmt.Printf("Cannot open specified file '%s': %v\n", tail[0], err)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Cannot read image from specified file '%s': %v\n", tail[0], err)
		return
	}

	o, err := os.Create(output)
	if err != nil {
		fmt.Printf("Cannot create output file '%s.ico': %v\n", output, err)
		return
	}
	defer o.Close()

	if err := ico.Encode(o, img); err != nil {
		fmt.Printf("Cannot write output file '%s.ico': %v\n", output, err)
		return
	}

	fmt.Printf("Icon created: %s\n", o.Name())
}
