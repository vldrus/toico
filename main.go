package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"io"

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
	// ...
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
