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
	"reflect"
	"testing"
)

func TestParseCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		want1   []string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    "This is a test",
			want:    "This",
			want1:   []string{"is", "a", "test"},
			wantErr: false,
		},
		{
			name:    "test2",
			args:    `This "is a test""`,
			want:    "This",
			want1:   []string{"is a test"},
			wantErr: false,
		},
		{
			name:    "test3",
			args:    `This 'is a test'`,
			want:    "This",
			want1:   []string{"is a test"},
			wantErr: false,
		},
		{
			name:    "test4",
			args:    `This is "{{ fooo balsi}}"`,
			want:    "This",
			want1:   []string{"is", "{{ fooo balsi}}"},
			wantErr: false,
		},
		{
			name:    "test5",
			args:    `This`,
			want:    "This",
			want1:   []string{},
			wantErr: false,
		},
		{
			name:    "test6",
			args:    `   a     b      c       d     e    "f 'g' h"`,
			want:    "a",
			want1:   []string{"b", "c", "d", "e", "f 'g' h"},
			wantErr: false,
		},
		{
			name:    "test7",
			args:    `   'a'     b      c       d     e    "f 'g' h"`,
			want:    "a",
			want1:   []string{"b", "c", "d", "e", "f 'g' h"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseCmd(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseCmd() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseCmd() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
