package main

import "testing"

func Test_intIdent(t *testing.T) {
	type args struct {
		input int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "Normal", args: args{input: 17}, want: 17},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intIdent(tt.args.input); got != tt.want {
				t.Errorf("intIdent() = %v, want %v", got, tt.want)
			}
		})
	}
}
