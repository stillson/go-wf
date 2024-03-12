// Package rcparse holds the interface for the rcfile
// and the implementations for handling various rcfile types
// i.e. simple, yaml, klingon, etc.
package rcparse

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PlainRCFile is based around a reader
// instead of a filename; this is much more testable.
type PlainRCFile struct {
	Commands map[string]string
}

type RCFile interface {
	Parse(r io.Reader) error
	GetCommand(rubric string) string
}

// Only works with full path.
func NewPlainRcFile(filename string) (*PlainRCFile, error) {
	filename = filepath.Join("/", filepath.Clean(filename))

	var fp, err = os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fp.Close()
	}()

	rv := PlainRCFile{
		Commands: make(map[string]string),
	}

	err = rv.Parse(fp)

	return &rv, err
}

func (rc *PlainRCFile) Parse(r io.Reader) error {

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()

		// Put in # comments
		if line[0] == '#' {
			continue
		}

		rubric, cmd, found := strings.Cut(line, ",")
		if !found {
			return fmt.Errorf("could not parse line in rcfile %v", line)
		}

		rubric = strings.Trim(rubric, " \n\t")
		cmd = strings.Trim(cmd, " \n\t")
		rc.Commands[rubric] = cmd
	}

	return nil
}

func (rc *PlainRCFile) GetCommand(rubric string) (string, bool) {
	val, exists := rc.Commands[rubric]
	return val, exists
}
