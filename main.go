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
	"github.com/fatih/color"
	"github.com/stillson/go-wf/rcfile"
	"github.com/stillson/go-wf/rcparse"
	"os"
	"os/exec"
	"strings"
)

// This is tricky to test. Depends on hidden variable and file system.
func preProcCmd(cmd string) (string, []string, error) {
	splitCmd := strings.Split(cmd, " ")

	outCmd, err := exec.LookPath(splitCmd[0])
	if err != nil {
		return splitCmd[0], splitCmd[1:], err
	}

	return outCmd, splitCmd[1:], nil
}

func main() {

	// parse command line
	flag.Parse()
	rubric := flag.Arg(0)

	// set up colors
	red := color.New(color.FgHiRed)
	green := color.New(color.FgHiGreen)

	// get filename of rcfile
	f, err := rcfile.GetRCFile()
	if err != nil {
		_, _ = red.Printf("Error getting rcfile:%v\n", err)
		os.Exit(1)
	}

	ourRcFile, err := rcparse.NewPlainRcFile(f)
	if err != nil {
		_, _ = red.Printf("Error parsing rcfile:%v\n", err)
		os.Exit(2)
	}

	// fmt.Printf("\tRC: %v\n", ourRcFile)

	cmd, exists := ourRcFile.GetCommand(rubric)
	if !exists {
		_, _ = red.Printf("rubric does not exist\n")
		os.Exit(3)
	}

	splitCmd, splitArgs, err := preProcCmd(cmd)
	if err != nil {
		_, _ = red.Printf("cmd not found in path? %v\terr:%v\n", splitCmd, err)
		os.Exit(4)
	}

	_, _ = green.Printf("cmd: %v\t\targs: %v\n\n", splitCmd, splitArgs)

	// Because this executes arbitrary commands from and external file
	// there is no way to make this "safe", short of whitelisting
	// hence the nolint:gosec
	ecmd := exec.Command(splitCmd, splitArgs...) //nolint:gosec
	ecmd.Stdout, ecmd.Stderr = os.Stdout, os.Stderr

	err = ecmd.Run()
	if err != nil {
		_, _ = red.Printf("%v", err)
	}

	exitVal := ecmd.ProcessState.ExitCode()

	_, _ = green.Printf("\nProcess exited with %v\n", exitVal)

	os.Exit(ecmd.ProcessState.ExitCode())
}
