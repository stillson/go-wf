package rcparse

import (
	"bytes"
	"io"
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
