package upload

import (
	"io"
	"os"
	"time"
)

type Uploader interface {
	// Stat remove Stat to `Exists` and return only bool and error?
	Stat(path string) (os.FileInfo, error)
	Upload(path string, body io.Reader, atime, mtime time.Time) error
	//Remove(path string) error
}
