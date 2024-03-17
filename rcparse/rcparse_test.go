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

package rcparse

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestPlainRCFile_GetCommand(t *testing.T) {
	type fields struct {
		Commands map[string]string
	}
	type args struct {
		rubric string
	}

	gf := func() fields {
		commands := make(map[string]string)
		commands["a"] = "b"
		return fields{Commands: commands}
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name:   "test1",
			fields: gf(),
			args:   args{rubric: "a"},
			want:   "b",
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &PlainRCFile{
				Commands: tt.fields.Commands,
			}
			got, got1 := rc.GetCommand(tt.args.rubric)
			if got != tt.want {
				t.Errorf("GetCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCommand() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPlainRCFile_Parse(t *testing.T) {
	type fields struct {
		Commands map[string]string
	}
	type args struct {
		r io.Reader
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		rubric  string
		cmd     string
	}{
		{
			name:    "test1",
			fields:  fields{map[string]string{}},
			args:    args{bytes.NewBufferString("a,b")},
			wantErr: false,
			rubric:  "a",
			cmd:     "b",
		},
		{
			name:    "test2",
			fields:  fields{map[string]string{}},
			args:    args{bytes.NewBufferString("#foo\na,b")},
			wantErr: false,
			rubric:  "a",
			cmd:     "b",
		},
		{
			name:    "test3",
			fields:  fields{map[string]string{}},
			args:    args{bytes.NewBufferString("ab")},
			wantErr: true,
			rubric:  "a",
			cmd:     "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &PlainRCFile{
				Commands: tt.fields.Commands,
			}
			if err := rc.Parse(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if cmd, exists := rc.GetCommand(tt.rubric); cmd != tt.cmd || !exists {
				t.Errorf("Parse()-get \"%v\":%v == wanted \"%v\"", cmd, exists, tt.cmd)
			}
		})
	}
}

func TestNewPlainRcFile(t *testing.T) {
	type args struct {
		filename string
	}

	const WF_FILE = ".workflow.yaml"

	pwd, _ := os.Getwd()

	f, err := os.OpenFile(WF_FILE, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		t.Fatal("Could not open .workflow.yaml for writing")
	}
	_, _ = f.WriteString("# for testing                       \n")
	_ = f.Close()
	defer func() {
		_ = os.Remove(WF_FILE)
	}()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{filepath.Join(pwd, ".workflow.yaml")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPlainRcFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlainRcFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

const YamlFile1 = `
# This is a test comment
wf_file:
  -
    rubric: a
    cmd: b
  -
    rubric: c
    cmd: >
      This is a special
      test. tada!
`

func TestYamlRCFile_Parse(t *testing.T) {
	type fields struct {
		Commands map[string]string
	}
	type args struct {
		r io.Reader
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		rubric  string
		cmd     string
	}{
		{
			name:    "test1",
			fields:  fields{map[string]string{}},
			args:    args{bytes.NewBufferString(YamlFile1)},
			wantErr: false,
			rubric:  "a",
			cmd:     "b",
		},
		{
			name:    "test2",
			fields:  fields{map[string]string{}},
			args:    args{bytes.NewBufferString("ab")},
			wantErr: true,
			rubric:  "a",
			cmd:     "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &YamlRCFile{
				Commands: tt.fields.Commands,
			}
			if err := rc.Parse(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if cmd, exists := rc.GetCommand(tt.rubric); cmd != tt.cmd || !exists {
				t.Errorf("Parse()-get \"%v\":%v == wanted \"%v\"", cmd, exists, tt.cmd)
			}
		})
	}
}

const YamlFile2 = `
# This is a test comment
globals:
  - bob=77
  - test = thingy
wf_file:
  -
    rubric: a
    cmd: b {{.G.bob}}
  -
    rubric: c
    cmd: >
      This is a special
      test. tada!
  - rubric: b
    cmd: '{{.G.test}} {{.G.bob}}'
`

func TestYTRCFile_Parse(t *testing.T) {
	type fields struct {
		Globals  map[string]string
		Commands map[string]string
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		rubric  string
		cmd     string
	}{
		{
			name: "test1",
			fields: fields{Globals: make(map[string]string),
				Commands: make(map[string]string)},
			args:    args{bytes.NewBufferString(YamlFile2)},
			wantErr: false,
			rubric:  "a",
			cmd:     "b 77",
		},
		{
			name: "test2",
			fields: fields{Globals: make(map[string]string),
				Commands: make(map[string]string)},
			args:    args{bytes.NewBufferString(YamlFile2)},
			wantErr: false,
			rubric:  "b",
			cmd:     "thingy 77",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &YTRCFile{
				G:        tt.fields.Globals,
				Commands: tt.fields.Commands,
			}
			if err := rc.Parse(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if cmd, exists := rc.GetCommand(tt.rubric); cmd != tt.cmd || !exists {
				t.Errorf("Parse()-get \"%v\":%v == wanted \"%v\"", cmd, exists, tt.cmd)
			}
		})
	}
}
