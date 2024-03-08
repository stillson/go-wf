// Package rcparse holds the interface for the rcfile
// and the implementations for handling various rcfile types
// i.e. simple, yaml, klingon, etc.
package rcparse

import "io"

// PlainRCFile is based around a reader
// instead of a filename; this is much more testable.
type PlainRCFile struct {
	Contents io.Reader
	Commands map[string]string
}

type RCFile interface {
	Parse() error
	GetCommand(rubric string) string
}

func (rc *PlainRCFile) Parse() error {
	return nil
}

func (rc *PlainRCFile) GetCommand(rubric string) string {
	return rc.Commands[rubric]
}
