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
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
)

type RCFile interface {
	Parse(r io.Reader) error
	GetCommand(rubric string) (string, bool)
}

type RCFileEnv interface {
	RCFile
	GetCommandEnv(rubric string) (string, map[string]string, bool)
}

// PlainRCFile
// (<rubric>,<cmd>\n)* .
type PlainRCFile struct {
	Commands map[string]string
}

// Yaml RCFile.
type YamlRCFile struct {
	Commands map[string]string
}

type cmdEnv struct {
	cmd  string
	envs map[string]string
}

// Yaml and Template RCFile.
type YTRCFile struct {
	G        map[string]string
	Commands map[string]cmdEnv
}

// NewPlainRcFile Only works with full path.
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

// NewYamlRcFile Only works with full path.
func _(filename string) (*YamlRCFile, error) {
	filename = filepath.Join("/", filepath.Clean(filename))
	var fp, err = os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fp.Close()
	}()

	rv := YamlRCFile{
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

// YamlRCFile
// An rcfile formatted in yaml.

type YamlFileEntry struct {
	Rubric string `yaml:"rubric"`
	Cmd    string `yaml:"cmd"`
}

type YamlFileFormat struct {
	Items []YamlFileEntry `yaml:"wf_file"`
}

func (rc *YamlRCFile) Parse(r io.Reader) error {
	entries := YamlFileFormat{
		Items: make([]YamlFileEntry, 10, 11),
	}
	br := bufio.NewReader(r)
	fileBuf := make([]byte, 4095, 4096)

	size, err := br.Read(fileBuf)

	if err != nil {
		return err
	}
	if size == 0 {
		return fmt.Errorf("empty rc file")
	}

	fileBuf = fileBuf[0:size]

	err = yaml.Unmarshal(fileBuf, &entries)
	if err != nil {
		return err
	}

	for _, entry := range entries.Items {
		rc.Commands[entry.Rubric] = entry.Cmd
	}

	return nil
}

func (rc *YamlRCFile) GetCommand(rubric string) (string, bool) {
	val, exists := rc.Commands[rubric]
	return val, exists
}

// YamlTmplRCFile

func NewYTFile(filename string) (*YTRCFile, error) {
	filename = filepath.Join("/", filepath.Clean(filename))
	var fp, err = os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fp.Close()
	}()

	rv := YTRCFile{
		Commands: make(map[string]cmdEnv),
		G:        make(map[string]string),
	}

	err = rv.Parse(fp)

	return &rv, err
}

type YTFileEntry struct {
	Rubric   string            `yaml:"rubric"`
	Commands string            `yaml:"c"`
	Env      map[string]string `yaml:"env,omitempty"`
}

type YTFormat struct {
	Items   []YTFileEntry `yaml:"wf_file"`
	Globals []string      `yaml:"globals,omitempty"`
}

func (rc *YTRCFile) Parse(r io.Reader) error {
	entries := YTFormat{
		Items:   make([]YTFileEntry, 10, 11),
		Globals: make([]string, 10, 11),
	}

	br := bufio.NewReader(r)
	fileBuf := make([]byte, 4095, 4096)

	size, err := br.Read(fileBuf)
	if err != nil {
		return err
	}
	if size == 0 {
		return fmt.Errorf("empty rc file")
	}

	fileBuf = fileBuf[0:size]

	err = yaml.Unmarshal(fileBuf, &entries)
	if err != nil {
		return err
	}

	for _, entry := range entries.Items {
		rc.Commands[entry.Rubric] = cmdEnv{cmd: entry.Commands, envs: map[string]string{}}
		for k, v := range entry.Env {
			rc.Commands[entry.Rubric].envs[k] = v
		}
	}

	for _, global := range entries.Globals {

		name := strings.SplitN(global, "=", 2)
		key := strings.Trim(name[0], " \n\t")
		val := strings.Trim(name[1], " \n\t")

		// make sure that there are only two items in split?
		if len(name) == 2 {
			rc.G[key] = val
		}
	}

	return nil
}

func (rc *YTRCFile) GetCommand(rubric string) (string, bool) {
	cmd, _, exists := rc.GetCommandEnv(rubric)
	return cmd, exists
}

func (rc *YTRCFile) GetCommandEnv(rubric string) (string, map[string]string, bool) {
	val, exists := rc.Commands[rubric]
	if !exists {
		return "", nil, exists
	}

	t := template.New("cmd").Funcs(sprig.FuncMap())
	tmlp, err := t.Parse(val.cmd)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error in template %v", err)
		return "", nil, false
	}

	var b strings.Builder
	err = tmlp.Execute(&b, rc)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error executing template: %v", err)
	}

	return b.String(), val.envs, exists
}
