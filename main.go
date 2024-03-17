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
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/stillson/go-wf/rcfile"
	"github.com/stillson/go-wf/rcparse"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	VERSION = "0.0.1"
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
	// flags
	versionQ := flag.Bool("v", false, "Version of this program")
	verboseQ := flag.Bool("V", false, "Verbose output")
	timeQ := flag.Bool("t", false, "Time the command")
	dumpQ := flag.Bool("d", false, "Dump contents of workflow file")
	wfFile := flag.String("f", ".workflow.yaml", "Name of workflow file")
	flag.Parse()

	if *versionQ {
		fmt.Printf("wf version %v\n", VERSION)
		return
	}

	verbose := false
	if *verboseQ {
		verbose = true
		fmt.Println("Verbose is on")
	}

	timing := false
	if *timeQ {
		timing = true
		if verbose {
			fmt.Printf("Timing enabled\n")
		}
	}

	dump := false
	if *dumpQ {
		dump = true
		if verbose {
			fmt.Printf("Dumping workflow file\n")
		}
	}

	if verbose {
		fmt.Printf("File to search for: %s\n", *wfFile)
	}

	// set up colors
	red := color.New(color.FgHiRed)
	green := color.New(color.FgHiGreen)

	// get filename of rcfile
	f, err := rcfile.GetRCFile(*wfFile)
	if err != nil {
		_, _ = red.Printf("Error getting rcfile:%v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Actual file found: %s\n", f)
	}

	if dump {
		var fp, err = os.Open(f) //nolint:gosec
		if err != nil {
			_, _ = red.Printf("Error reading rcfile:%v\n", err)
			os.Exit(2)

		}

		defer func() {
			_ = fp.Close()
		}()

		if verbose {
			fmt.Printf("---\n")
		}
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			_, _ = red.Printf("error dumping file %v\n", err)
		}

		return
	}

	ourRcFile, err := rcparse.NewYTFile(f)
	if err != nil {
		_, _ = red.Printf("Error parsing rcfile:%v\n", err)
		os.Exit(2)
	}

	if verbose {
		fmt.Printf("\tRC: %v\n", ourRcFile)
	}

	rubric := flag.Arg(0)
	if verbose {
		fmt.Printf("rubric is: %s\n", rubric)
	}

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

	var now int64
	if timing {
		now = time.Now().UnixMicro()
		if verbose {
			fmt.Printf("Start time: %v\n", now)
		}
	}

	err = ecmd.Run()
	if err != nil {
		_, _ = red.Printf("%v\n", err)
	}

	if timing {
		end := time.Now().UnixMicro()
		if verbose {
			fmt.Printf("End time %v\n", end)
		}
		fmt.Printf("Total Time in µsecs: %v\n", end-now)
	}
	exitVal := ecmd.ProcessState.ExitCode()

	_, _ = green.Printf("\nProcess exited with %v\n", exitVal)

	os.Exit(ecmd.ProcessState.ExitCode())
}
