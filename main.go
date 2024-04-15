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
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/stillson/go-wf/executor"
	"github.com/stillson/go-wf/rcfile"
	"github.com/stillson/go-wf/rcparse"
)

const (
	VERSION = "0.0.2"
)

func vprint(verbose bool, format string, inputs ...any) {
	if !verbose {
		return
	}

	fmt.Printf(format, inputs...)
}

func main() {
	verboseQ, timeQ, dumpQ, wfFile, ruleQ := ParseArgs()

	vprint(*verboseQ, "Verbose is on\n")
	vprint(*timeQ && *verboseQ, "Timing enabled\n")

	dump := false
	if *dumpQ {
		dump = true
		if *verboseQ {
			fmt.Printf("Dumping workflow file\n")
		}
	}

	vprint(*verboseQ, "File to search for: %s\n", *wfFile)

	// set up colors
	red := color.New(color.FgHiRed)
	green := color.New(color.FgHiGreen)

	// get filename of rcfile
	f, err := rcfile.GetRCFile(*wfFile)
	if err != nil {
		_, _ = red.Printf("Error getting rcfile:%v\n", err)
		os.Exit(1)
	}

	vprint(*verboseQ, "Actual file found: %s\n", f)

	if dump {
		var fp, err = os.Open(f) //nolint:gosec
		if err != nil {
			_, _ = red.Printf("Error reading rcfile:%v\n", err)
			os.Exit(2)

		}

		defer func() {
			_ = fp.Close()
		}()

		if *verboseQ {
			vprint(*verboseQ, "---\n")
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

	ourRcFile, err := rcparse.NewYRCFile(f)
	if err != nil {
		_, _ = red.Printf("Error parsing rcfile:%v\n", err)
		os.Exit(2)
	}

	if *verboseQ {
		vprint(*verboseQ, "\tRC: %v\n", ourRcFile)
	}

	if *ruleQ {
		rules, err := ourRcFile.ListRules()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(7)
		}

		for _, rule := range rules {
			fmt.Printf("%s\n", rule)
		}

		return
	}

	rule := flag.Arg(0)
	vprint(*verboseQ, "rule is: %s\n", rule)

	var now int64
	if *timeQ {
		now = time.Now().UnixMicro()
		if *verboseQ {
			_, _ = green.Printf("Start time: %v\n", now)
		}
	}

	localExec := executor.NewLocalExec("main")
	rv, err := localExec.Run(rule, ourRcFile)

	if err != nil {
		_, _ = red.Printf("%v\n", err)
	}

	if *timeQ {
		end := time.Now().UnixMicro()
		if *verboseQ {
			_, _ = green.Printf("End time %v\n", end)
		}
		fmt.Printf("Total Time in µsecs: %v\n", end-now)
	}

	if rv != 0 {
		_, _ = red.Printf("\nProcess exited with %v\n", rv)
	} else {
		_, _ = green.Printf("\nProcess exited with %v\n", rv)
	}

	os.Exit(rv)
}

func ParseArgs() (*bool, *bool, *bool, *string, *bool) {
	// flags
	versionQ := flag.Bool("v", false, "Version of this program")
	verboseQ := flag.Bool("V", false, "Verbose output")
	timeQ := flag.Bool("t", false, "Time the command")
	dumpQ := flag.Bool("d", false, "Dump contents of workflow file")
	wfFile := flag.String("f", ".workflow.yaml", "Name of workflow file")
	ruleQ := flag.Bool("r", false, "Print available rules")

	flag.Parse()

	if *versionQ {
		fmt.Printf("wf version %v\n", VERSION)
		os.Exit(0)
	}
	return verboseQ, timeQ, dumpQ, wfFile, ruleQ
}
