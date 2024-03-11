// Package rcfile is all the code for finding
// the workflowrc file
package rcfile

import (
	"os"
	"path"
)

// for now, just get full path to ./.workflowrc
// later, search parents, copy in a default if needed
// hmmm, how to test

func GetRCFile() (string, error) {

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(pwd, ".workflowrc"), nil
}
