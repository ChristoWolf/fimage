// Simple Go package for encoding and decoding images directly from and to file paths.
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

// Supported image file types.
const (
	gifExt  = "gif"
	jpegExt = "jpeg"
	jpgExt  = "jpg"
	pngExt  = "png"
)

// Instance of an io.LogSink implementation which will be used for logging.
// By default, this will be initialized as io.Discard, but can be set
// to any desired io.LogSink implementation.
// This is designed to NOT be a logger instance, because the package user
// might be using a logging library which does not offer drop-in replacement
// of the log standard package.
var LogSink io.Writer

// Initialized the package upon use.
// For now, this is limited to initialization of LogSink to io.Discard.
// See LogSink for details.
func init() {
	LogSink = io.Discard
}

// Parses a file's file type from its name or path.
// Returns an empty string if filePath does not contain a file type extension.
func ParseFileType(filePath string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
}

// Reads the image file from the provided file path if possible.
// If successful, the first return value represents the image
// as image.Image.
// Supported file types: gif, jpeg, jpg, png.
func ReadImage(filePath string) (image.Image, error) {
	filePath = filepath.Clean(filePath)
	fmt.Fprintf(LogSink, "provided file path: %s\n", filePath)
	fmt.Fprintln(LogSink, "checking if file type is supported")
	if fileType := ParseFileType(filePath); fileType != gifExt &&
		fileType != jpegExt &&
		fileType != jpgExt &&
		fileType != pngExt {
		return nil, fmt.Errorf("unsupported file type '%s'", fileType)
	}
	fmt.Fprintln(LogSink, "opening file")
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()
	defer fmt.Fprintln(LogSink, "closing file")
	fmt.Fprintln(LogSink, "decoding file")
	image, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file as image: %v", err)
	}
	return image, err
}

// Writes an image.Image to the provided file path if possible.
// Supported file types: gif, jpeg, jpg, png.
func WriteImage(i image.Image, filePath string) error {
	filePath = filepath.Clean(filePath)
	fmt.Fprintf(LogSink, "provided target file path: %s\n", filePath)
	fmt.Fprintln(LogSink, "detecting desired target encoding")
	fileType := ParseFileType(filePath)
	switch fileType {
	case gifExt, jpegExt, jpgExt, pngExt:
	default:
		return fmt.Errorf("file type %s is not supported", fileType)
	}
	fmt.Fprintln(LogSink, "creating file")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()
	defer fmt.Fprintln(LogSink, "closing file")
	fmt.Fprintf(LogSink, "encoding file as %s\n", fileType)
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
