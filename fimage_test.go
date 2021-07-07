package fimage_test

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"os"
	"testing"

	"github.com/christowolf/fimage"
	"github.com/google/go-cmp/cmp"
)

const testFile = "./testdata/testimage.png"
const tempFile = "./testdata/tempimage.png"

var maxUint8 = ^uint8(0)
var red = color.RGBA{maxUint8, 0, 0, maxUint8}
var green = color.RGBA{0, maxUint8, 0, maxUint8}
var blue = color.RGBA{0, 0, maxUint8, maxUint8}
var gray = color.Gray{maxUint8 / 2}

type ModelSpy struct {
}

func (ModelSpy) Convert(c color.Color) color.Color { return c }

type ImageSpy struct {
	colors [][]color.Color
}

func (i *ImageSpy) ColorModel() color.Model { return ModelSpy{} }

func (i *ImageSpy) Bounds() image.Rectangle { return image.Rect(0, 0, len(i.colors), len(i.colors[0])) }

func (i *ImageSpy) At(x, y int) color.Color { return i.colors[x][y] }

var testImage image.Image = &ImageSpy{[][]color.Color{
	{color.White, gray, color.Black, red, green, blue},
	{blue, green, red, color.Black, gray, color.White}}}

// Tests if fimage.LogSink is initialized correctly as io.Discard.
func TestInit(t *testing.T) {
	t.Parallel()
	want := io.Discard
	got := fimage.LogSink
	if got == nil {
		t.Error("expected LogSink to not be nil")
	}
	if !cmp.Equal(got, want) {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}

func TestParseFileType(t *testing.T) {
	t.Parallel()
	var data = []struct {
		path string
		want string
	}{
		{`C:\Test Folder\testfile.go`, "go"},
		{"C:/Test Folder/testfile.png", "png"},
		{"./../testfile.jpg", "jpg"},
		{"//Test Share/testfile.tar.gz", "gz"},
		{"testfile", ""},
	}
	for _, row := range data {
		row := row
		t.Run(row.path, func(t *testing.T) {
			t.Parallel()
			if got := fimage.ParseFileType(row.path); got != row.want {
				t.Errorf("got : %s, want: %s", got, row.want)
			}
		})
	}
}

// Verifies that the result of reading a physical test image file
// via fimage.ReadImage matches some image.Image mock data.
// Reference image file is a png instead of jpg,
// as compression makes testing quite a bit harder.
func TestReadImage(t *testing.T) {
	t.Parallel()
	got, err := fimage.ReadImage(testFile)
	if err != nil {
		t.Errorf("got: %#v, want: nil", err)
	}
	if !cmp.Equal(&got, &testImage, imCmp) {
		t.Errorf("got: %#v, want: %#v", got, testImage)
	}
}

// Verifies that the result of writing some image.Image mock data
// via fimage.WriteImage matches a physical physical test image file
// w.r.t. their hashes.
// Reference image file is a png instead of jpg,
// as compression makes testing quite a bit harder.
func TestWriteImage(t *testing.T) {
	t.Parallel()
	err := fimage.WriteImage(testImage, tempFile)
	if err != nil {
		t.Errorf("got: %#v, want: nil", err)
	}
	want, _ := os.ReadFile(testFile)
	got, _ := os.ReadFile(tempFile)
	if !bytes.Equal(got, want) {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}

// Comparer implementation for comparing 2 concrete image.Images.
var imCmp = cmp.Comparer(func(x, y *image.Image) bool {
	// ColorModel is ignored, as we only need to check the image's visual properties.
	if !(cmp.Equal((*x).Bounds(), (*y).Bounds())) {
		return false
	}
	b := (*x).Bounds()
	for x1 := b.Min.X; x1 < b.Max.X; x1++ {
		for x2 := b.Min.Y; x2 < b.Max.Y; x2++ {
			xC := (*x).At(x1, x2)
			yC := (*y).At(x1, x2)
			if !cmp.Equal(&xC, &yC, colorCmp) {
				return false
			}
		}
	}
	return true
})

// Comparer implementation for comparing 2 concrete color.Colors.
var colorCmp = cmp.Comparer(func(x, y *color.Color) bool {
	xR, xG, xB, xA := (*x).RGBA()
	yR, yG, yB, yA := (*y).RGBA()
	return xR == yR &&
		xG == yG &&
		xB == yB &&
		xA == yA
})
