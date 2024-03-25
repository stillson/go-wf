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

package rcfile

import (
	"os"
	"path/filepath"
	"testing"
)

func sameFile(a string, b string) (bool, error) {
	af, err := os.Open(a) //nolint:gosec
	if err != nil {
		return false, err
	}
	defer func() {
		_ = af.Close()
	}()

	bf, err := os.Open(b) //nolint:gosec
	if err != nil {
		return false, err
	}
	defer func() {
		_ = bf.Close()
	}()

	aInfo, err := af.Stat()
	if err != nil {
		return false, err
	}

	bInfo, err := bf.Stat()
	if err != nil {
		return false, err
	}

	return os.SameFile(aInfo, bInfo), nil
}

func TestGetRCFile(t *testing.T) {

	dir, err := os.MkdirTemp("", "rcfile_test*")
	if err != nil {
		t.Fatalf("Unable to creat tmp directory\n")
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	oldPwd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(oldPwd)
	}()

	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Unable to change directory to %s: %v\n", dir, err)
	}

	// make a .workflow.yaml
	file := filepath.Join(dir, ".workflow.yaml")
	if err = os.WriteFile(file, []byte("content"), 0600); err != nil {
		t.Fatalf("Unable to create test .workflow.yaml")
	}
	// make subdirectory TEST
	subDir := filepath.Join(dir, "TEST")
	if err = os.Mkdir(subDir, 0750); err != nil {
		t.Fatalf("Unable to create testing subDir")
	}
	if err = os.Mkdir(filepath.Join(subDir, ".workflow.yaml"), 0750); err != nil {
		t.Fatalf("Unable to create testing fake workflow file")
	}
	subDir2 := filepath.Join(subDir, "TEST2")
	if err = os.Mkdir(subDir2, 0750); err != nil {
		t.Fatalf("Unable to create testing fake workflow file")
	}
	file2 := filepath.Join(subDir2, ".workflow.yaml")
	if err = os.WriteFile(file2, []byte("content"), 0); err != nil {
		t.Fatalf("Unable to create second test .workflow.yaml")
	}

	tests := []struct {
		name    string
		dir     string
		fName   string
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			dir:     dir,
			fName:   ".workflow.yaml",
			want:    filepath.Join(dir, ".workflow.yaml"),
			wantErr: false,
		},
		{
			name:    "test2",
			dir:     subDir,
			fName:   ".workflow.yaml",
			want:    filepath.Join(dir, ".workflow.yaml"),
			wantErr: false,
		},
		{
			name:    "test3",
			dir:     subDir2,
			fName:   ".workflow.yaml",
			want:    filepath.Join(dir, ".workflow.yaml"),
			wantErr: false,
		},
		{
			name:    "test4",
			dir:     subDir2,
			fName:   "NOTFOUND.yaml",
			want:    filepath.Join(dir, "NOTFOUND.yaml"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = os.Chdir(tt.dir)
			if err != nil {
				t.Errorf("unable to change director to %v\n", tt.dir)
			}
			got, err := GetRCFile(tt.fName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRCFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				// we wanted and got an error, return now
				return
			}
			same, err := sameFile(got, tt.want)
			if err != nil {
				t.Errorf("sameFile error %v: %s  %s", err, got, tt.want)
			}
			if !same {
				t.Errorf("GetRCFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
