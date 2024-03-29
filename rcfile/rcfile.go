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

// Package rcfile is all the code for finding
// the workflowrc file
package rcfile

import (
	"fmt"
	"os"
	"path"
)

// for now, just get full path to ./.workflow.yaml
// later, search parents, copy in a default if needed
// hmmm, how to test

func GetRCFile(fname string) (string, error) {

	fpath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fname = path.Clean(fname)

	for ; fpath != "/"; fpath = path.Dir(fpath) {
		rcCandidate := path.Join(fpath, fname)

		// does it exist?

		f, err := os.OpenFile(rcCandidate, os.O_RDONLY, 000) //nolint:gosec
		if err != nil {
			continue
		}

		fi, err := f.Stat()
		if err != nil {
			continue
		}

		if fi.IsDir() {
			continue
		}

		return rcCandidate, nil
	}

	return "", fmt.Errorf("workflow file not found")

}
