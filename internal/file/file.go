package file

import (
	"os"
)

type FileOps struct{}

func (r FileOps) Open(name string) (*os.File, error) {
	return os.Open(name)
}
