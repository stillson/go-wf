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
	"sort"
	"strings"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
)

type RCFile interface {
	Parse(r io.Reader) error
	GetCommand(rule string) ([]string, bool)
	GetCommandEnv(rule string) ([]string, map[string]string, bool)
	ListRules() ([]string, error)
}

type CmdEnv struct {
	Cmd  []string
	Envs map[string]string
}

type YRCfile struct {
	G        map[string]string
	Commands map[string]CmdEnv
}

func NewYRCFile(filename string) (*YRCfile, error) {
	filename = filepath.Join("/", filepath.Clean(filename))
	var fp, err = os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fp.Close()
	}()

	rv := YRCfile{
		Commands: make(map[string]CmdEnv),
		G:        make(map[string]string),
	}

	err = rv.Parse(fp)

	return &rv, err
}

type YRCFileEntry struct {
	Rule     string            `yaml:"rule"`
	Commands []string          `yaml:"c"`
	Env      map[string]string `yaml:"env,omitempty"`
}

type YRCFormat struct {
	Items   []YRCFileEntry    `yaml:"wf_file"`
	Globals map[string]string `yaml:"globals,omitempty"`
}

func (rc *YRCfile) Parse(r io.Reader) error {
	entries := YRCFormat{
		Items:   make([]YRCFileEntry, 10, 11),
		Globals: make(map[string]string, 10),
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
		newRule := CmdEnv{Cmd: []string{}, Envs: map[string]string{}}

		newRule.Cmd = append(newRule.Cmd, entry.Commands...)

		for k, v := range entry.Env {
			newRule.Envs[k] = v
		}

		rc.Commands[entry.Rule] = newRule
	}

	for k, v := range entries.Globals {
		rc.G[k] = v
	}

	return nil
}

func (rc *YRCfile) GetCommand(rule string) ([]string, bool) {
	cmd, _, exists := rc.GetCommandEnv(rule)
	return cmd, exists
}

func (rc *YRCfile) GetCommandEnv(rule string) ([]string, map[string]string, bool) {
	val, exists := rc.Commands[rule]
	if !exists {
		return []string{}, nil, exists
	}

	rv := []string{}

	for _, c := range val.Cmd {
		t := template.New("Cmd").Funcs(sprig.FuncMap())
		tmlp, err := t.Parse(c)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error in template %v", err)
			return []string{}, nil, false
		}

		var b strings.Builder
		err = tmlp.Execute(&b, rc)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error executing template: %v", err)
		}

		rv = append(rv, b.String())
		b.Reset()
	}
	return rv, val.Envs, exists
}

func (rc *YRCfile) ListRules() ([]string, error) {
	rv := []string{}

	for k := range rc.Commands {
		rv = append(rv, k)
	}

	sort.Strings(rv)
	return rv, nil
}
