package fimage

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	"image/gif"
	"image/jpeg"
	"image/png"
)

const (
	gifExt  = "gif"
	jpegExt = "jpeg"
	jpgExt  = "jpg"
	pngExt  = "png"
)

var Writer io.Writer

func init() {
	Writer = io.Discard
}

func ReadImage(filePath string) (image.Image, error) {
	fmt.Fprintf(Writer, "provided file path: %s\n", filePath)
	fmt.Fprintln(Writer, "checking if file type is supported")
	if fileType := ParseFileType(filePath); fileType != gifExt &&
		fileType != jpegExt &&
		fileType != jpgExt &&
		fileType != pngExt {
		return nil, fmt.Errorf("unsupported file type '%s'", fileType)
	}
	fmt.Fprintln(Writer, "opening file")
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()
	defer fmt.Fprintln(Writer, "closing file")
	fmt.Fprintln(Writer, "decoding file")
	image, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file as image: %v", err)
	}
	return image, err
}

func WriteImage(i image.Image, filePath string) error {
	fmt.Fprintf(Writer, "provided target file path: %s\n", filePath)
	fmt.Fprintln(Writer, "detecting desired target encoding")
	fileType := ParseFileType(filePath)
	switch fileType {
	case gifExt, jpegExt, jpgExt, pngExt:
	default:
		return fmt.Errorf("file type %s is not supported", fileType)
	}
	fmt.Fprintln(Writer, "creating file")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()
	defer fmt.Fprintln(Writer, "closing file")
	fmt.Fprintf(Writer, "encoding file as %s\n", fileType)
	switch fileType {
	case gifExt:
		err = gif.Encode(f, i, nil)
	case jpegExt, jpgExt:
		err = jpeg.Encode(f, i, nil)
	case pngExt:
		err = png.Encode(f, i)
	}
	if err != nil {
		return fmt.Errorf("failed to encode file: %v", err)
	}
	return nil
}

func ParseFileType(filePath string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
}
