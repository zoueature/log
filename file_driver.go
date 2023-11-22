package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
)

const fileDriverName = "file"

type FileDriver struct {
	Path    string `json:"path"`
	Loggers map[logrus.Level]logrus.Logger
}

func (d FileDriver) GetWriter() (io.Writer, error) {
	fd, err := os.OpenFile(d.Path, os.O_WRONLY, fs.ModePerm)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func (d FileDriver) String() string {
	return fileDriverName
}
