package logfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	dirName  = "tcpeek"
	fileName = "tcpeek.log"
	maxSize  = 30 * 1024
	keep     = 3
	dirPerm  = 0o755
	filePerm = 0o644
)

// Path returns the active log file path.
func Path() (string, error) {
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolving home directory: %w", err)
		}
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, dirName, fileName), nil
}

// Writer appends to the log file and rotates it once it exceeds maxSize,
// keeping `keep` files in total.
type Writer struct {
	mu   sync.Mutex
	path string
	file *os.File
	size int64
}

func Open() (*Writer, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return nil, err
	}
	w := &Writer{path: path}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) > maxSize && w.size > 0 {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}

func (w *Writer) openFile() error {
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, filePerm)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	w.file = f
	w.size = info.Size()
	return nil
}

func (w *Writer) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}
	for i := keep - 1; i >= 1; i-- {
		src := w.path
		if i > 1 {
			src = fmt.Sprintf("%s.%d", w.path, i-1)
		}
		dst := fmt.Sprintf("%s.%d", w.path, i)
		if err := os.Rename(src, dst); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return w.openFile()
}
