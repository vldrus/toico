package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/ccitt"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/nfnt/resize"
)

const MaxWidth = 256
const MaxHeight = 256

type FileHeader struct {
	Reserved  uint16
	ImageType uint16
	NumImages uint16
}

type IconDirEntry struct {
	ImageWidth   uint8
	ImageHeight  uint8
	NumColors    uint8
	Reserved     uint8
	ColorPlanes  uint16
	BitsPerPixel uint16
	SizeInBytes  uint32
	Offset       uint32
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s image.png\n", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[1]

	f, err := os.Open(name)
	if err != nil {
		fmt.Printf("Cannot read file %s\n", name)
		os.Exit(1)
	}
	defer f.Close()

	im, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Cannot decode file %s\n", name)
		os.Exit(1)
	}

	o, err := os.Create(fmt.Sprintf("%s.ico", name))
	if err != nil {
		fmt.Printf("Cannot create output file for %s\n", name)
		os.Exit(1)
	}
	defer o.Close()

	err = Encode(o, im)
	if err != nil {
		fmt.Printf("Cannot write output file for %s\n", name)
		os.Exit(1)
	}
}

func Encode(w io.Writer, im image.Image) error {
	if im.Bounds().Dx() > MaxWidth || im.Bounds().Dy() > MaxHeight {
		im = resize.Thumbnail(MaxWidth, MaxHeight, im, resize.Lanczos3)
	}

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	var buf bytes.Buffer

	err := encoder.Encode(&buf, im)
	if err != nil {
		return err
	}

	fh := FileHeader{
		Reserved:  0,
		ImageType: 1,
		NumImages: 1,
	}
	ide := IconDirEntry{
		ImageWidth:   uint8(im.Bounds().Dx()),
		ImageHeight:  uint8(im.Bounds().Dy()),
		NumColors:    0,
		Reserved:     0,
		ColorPlanes:  0,
		BitsPerPixel: 0,
		SizeInBytes:  uint32(len(buf.Bytes())),
		Offset:       22,
	}

	err = binary.Write(w, binary.LittleEndian, fh)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, ide)
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
