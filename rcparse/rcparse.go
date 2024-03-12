/*
 * Copyright (c) 2024. Christopher Stillson <stillson@gmail.com>
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 * Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
