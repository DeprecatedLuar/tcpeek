package logfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPathHonorsXDGStateHome(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", "/x")
	got, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("/x", dirName, fileName); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRotationKeepsThreeFiles(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	w, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	line := strings.Repeat("a", 1023) + "\n"
	for i := 0; i < 150; i++ {
		if _, err := w.Write([]byte(line)); err != nil {
			t.Fatal(err)
		}
	}

	dir := filepath.Dir(w.path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != keep {
		t.Fatalf("got %d files, want %d", len(entries), keep)
	}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			t.Fatal(err)
		}
		if info.Size() > maxSize {
			t.Errorf("%s is %d bytes, over cap", e.Name(), info.Size())
		}
	}
}
