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

package executor

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stillson/go-wf/rcparse"
)

const YAMLFILE = `
# This is a test comment
globals:
wf_file:
  -
    rule: alpha
    c:
     - echo "TEST"
`

func TestLocalExecutor_Run(t *testing.T) {

	type args struct {
		rule   string
		rcfile *rcparse.YRCfile
	}

	rcfile, _ := rcparse.CreateYRCFile(strings.NewReader(YAMLFILE))

	tests := []struct {
		name    string
		fields  string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:   "test1",
			fields: "test",
			args: args{
				rule:   "alpha",
				rcfile: rcfile,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LocalExecutor{
				name: tt.fields,
			}
			got, err := l.Run(tt.args.rule, tt.args.rcfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalExecutor_displayCommand(t *testing.T) {
	type args struct {
		splitCmd  string
		splitArgs []string
		env       map[string]string
	}
	tests := []struct {
		name   string
		fields string
		args   args
	}{
		{
			name:   "test1",
			fields: "test",
			args: args{
				splitCmd:  "echo",
				splitArgs: []string{"TEST"},
				env:       map[string]string{"TEST": "TEST"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// need to capture and check stdout
			l := &LocalExecutor{
				name: tt.fields,
			}
			l.displayCommand(tt.args.splitCmd, tt.args.splitArgs, tt.args.env)
		})
	}
}

func TestLocalExecutor_getCommand(t *testing.T) {
	type args struct {
		splitCmd  string
		splitArgs []string
		env       map[string]string
	}
	tests := []struct {
		name   string
		fields string
		args   args
		want   string
	}{
		{
			name:   "test1",
			fields: "test",
			args: args{
				splitCmd:  "echo",
				splitArgs: []string{"TEST"},
				env:       map[string]string{"TEST": "TEST"},
			},
			want: "echo TEST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LocalExecutor{
				name: tt.fields,
			}
			if got := l.getCommand(tt.args.splitCmd, tt.args.splitArgs, tt.args.env).String(); !strings.Contains(got, tt.want) {
				t.Errorf("getCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalExecutor_subRun(t *testing.T) {
	type args struct {
		cmd string
		env map[string]string
	}
	tests := []struct {
		name    string
		fields  string
		args    args
		want    int
		wantErr bool
	}{

		{
			name:   "test1",
			fields: "test",
			args: args{
				cmd: "echo TEST",
				env: map[string]string{"TEST": "TEST"},
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LocalExecutor{
				name: tt.fields,
			}
			got, err := l.subRun(tt.args.cmd, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("subRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("subRun() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_preProcCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		want1   []string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    "echo TEST",
			want:    "/bin/echo",
			want1:   []string{"TEST"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := preProcCmd(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("preProcCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("preProcCmd() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("preProcCmd() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
