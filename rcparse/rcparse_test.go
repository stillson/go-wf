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
	"maps"
	"testing"
)

const YamlFile = `
# This is a test comment
globals:
  bob: BOBOB
  test: TESTEST
wf_file:
  -
    rule: alpha
    c:
     - bbb {{.G.bob}}
     - ddd
    env:
      FOO: BAR
      USER: me
  -
    rule: beta
    c: 
      - >
        This is a special
        test. tada!
      - foo
  - rule: delta
    c: 
      - '{{.G.test}} {{.G.bob}}'
`

func TestYTRCFile_Parse(t *testing.T) {
	type fields struct {
		Globals  map[string]string
		Commands map[string]CmdEnv
	}
	type args struct {
		r io.Reader
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		rule    string
		cmd     string
		env     map[string]string
	}{

		{
			name: "test1",
			fields: fields{Globals: make(map[string]string),
				Commands: make(map[string]CmdEnv)},
			args:    args{bytes.NewBufferString(YamlFile)},
			wantErr: false,
			rule:    "alpha",
			cmd:     "bbb BOBOB",
			env:     map[string]string{"FOO": "BAR", "USER": "me"},
		},
		{
			name: "test2",
			fields: fields{Globals: make(map[string]string),
				Commands: make(map[string]CmdEnv)},
			args:    args{bytes.NewBufferString(YamlFile)},
			wantErr: false,
			rule:    "delta",
			cmd:     "TESTEST BOBOB",
			env:     map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &YRCfile{
				G:        tt.fields.Globals,
				Commands: tt.fields.Commands,
			}
			if err := rc.Parse(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v\n", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			cmd, parsedEnv, exists := rc.GetCommandEnv(tt.rule)
			if !exists || cmd[0] != tt.cmd {
				t.Errorf("Parse()-get got: %v exists:%v -- wanted %v\n", cmd, exists, tt.cmd)
			}
			if !maps.Equal(parsedEnv, tt.env) {
				t.Errorf("Expected environment is incorrect \nexpected:\t%v\ngot:\t\t%v\n", tt.env, parsedEnv)
			}

		})
	}
}
