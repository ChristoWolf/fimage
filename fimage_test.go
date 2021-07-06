package fimage_test

import (
	"io"
	"testing"

	"github.com/christowolf/fimage"
	"github.com/google/go-cmp/cmp"
)

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
