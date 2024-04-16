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

package main

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/fatih/color"
	"github.com/stillson/go-wf/rcparse"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		newArgs []string
		verbose bool
		time    bool
		dump    bool
		wfFile  string
		rules   bool
	}{
		{
			name:    "test1",
			newArgs: []string{"wf"},
			verbose: false,
			time:    false,
			dump:    false,
			wfFile:  ".workflow.yaml",
			rules:   false,
		},
		{
			name:    "test2",
			newArgs: []string{"wf", "-f", "TESTNAME"},
			verbose: false,
			time:    false,
			dump:    false,
			wfFile:  "TESTNAME",
			rules:   false,
		},
		{
			name:    "test3",
			newArgs: []string{"wf", "-V"},
			verbose: true,
			time:    false,
			dump:    false,
			wfFile:  ".workflow.yaml",
			rules:   false,
		},

		{
			name:    "test4",
			newArgs: []string{"wf", "-t"},
			verbose: false,
			time:    true,
			dump:    false,
			wfFile:  ".workflow.yaml",
			rules:   false,
		},
		{
			name:    "test5",
			newArgs: []string{"wf", "-d"},
			verbose: false,
			time:    false,
			dump:    true,
			wfFile:  ".workflow.yaml",
			rules:   false,
		},
		{
			name:    "test6",
			newArgs: []string{"wf", "-r"},
			verbose: false,
			time:    false,
			dump:    false,
			wfFile:  ".workflow.yaml",
			rules:   true,
		},
		{
			name:    "test7",
			newArgs: []string{"wf", "-V", "-t", "-d", "-r", "-f", "TESTNAME"},
			verbose: true,
			time:    true,
			dump:    true,
			wfFile:  "TESTNAME",
			rules:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.newArgs
			got, got1, got2, got3, got4 := ParseArgs()
			if !reflect.DeepEqual(*got, tt.verbose) {
				t.Errorf("ParseArgs() got = %v, verbose %v", got, tt.verbose)
			}
			if !reflect.DeepEqual(*got1, tt.time) {
				t.Errorf("ParseArgs() got1 = %v, verbose %v", got1, tt.time)
			}
			if !reflect.DeepEqual(*got2, tt.dump) {
				t.Errorf("ParseArgs() got2 = %v, verbose %v", got2, tt.dump)
			}
			if !reflect.DeepEqual(*got3, tt.wfFile) {
				t.Errorf("ParseArgs() got3 = %v, verbose %v", got3, tt.wfFile)
			}
			if !reflect.DeepEqual(*got4, tt.rules) {
				t.Errorf("ParseArgs() got4 = %v, verbose %v", got4, tt.rules)
			}

			// to reset flag module so it can be reused
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		})
	}
}

func Test_dumpRulesFile(t *testing.T) {
	type args struct {
		f    string
		verb bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dumpRulesFile(tt.args.f, tt.args.verb)
		})
	}
}

func Test_getColorPrints(t *testing.T) {
	tests := []struct {
		name  string
		want  *color.Color
		want1 *color.Color
	}{
		{
			name:  "test",
			want:  color.New(color.FgHiRed),
			want1: color.New(color.FgHiGreen),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getColorPrints()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getColorPrints() got = %v, verbose %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getColorPrints() got1 = %v, verbose %v", got1, tt.want1)
			}
		})
	}
}

func Test_printRules(t *testing.T) {
	type args struct {
		ourRcFile *rcparse.YRCfile
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printRules(tt.args.ourRcFile)
		})
	}
}

func Test_vprint(t *testing.T) {
	type args struct {
		verbose    bool
		printGreen bool
		format     string
		inputs     []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vprint(tt.args.verbose, tt.args.printGreen, tt.args.format, tt.args.inputs...)
		})
	}
}
