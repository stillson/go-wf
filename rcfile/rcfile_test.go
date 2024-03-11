package rcfile

import (
	"fmt"
	"os"
	"testing"
)

func TestGetRCFile(t *testing.T) {

	pwd, _ := os.Getwd()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			want:    fmt.Sprintf("%s/.workflowrc", pwd),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRCFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRCFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRCFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
