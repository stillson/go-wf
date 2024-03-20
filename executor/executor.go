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

// Package executor is the system of executing a command.
// This is an interface to allow for other forms of execution,
// i.e. remote, maybe in a docker, etc.
package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/stillson/go-wf/rcparse"
)

// This is tricky to test. Depends on hidden variable and file system.
func preProcCmd(cmd string) (string, []string, error) {

	cmdStart, cmdRest, err := rcparse.ParseCmd(cmd)

	if err != nil {
		return "", nil, err
	}

	outCmd, err := exec.LookPath(cmdStart)
	if err != nil {
		return cmdStart, cmdRest, err
	}

	return outCmd, cmdRest, nil
}

type Executor interface {
	Run(rubric string, rcfile *rcparse.RCFile) (int, error)
	// RunWithContext(ctx context.Context, rubric string, rcfile *rcparse.RCFile) error
}

type LocalExecutor struct {
	name string
}

func NewLocalExec(name string) LocalExecutor {
	return LocalExecutor{name}
}

func (l *LocalExecutor) Run(rubric string, rcfile *rcparse.YRCfile) (int, error) {
	red := color.New(color.FgHiRed)
	green := color.New(color.FgHiGreen)

	cmd, env, exists := rcfile.GetCommandEnv(rubric)
	if !exists {
		_, _ = red.Printf("rubric does not exist\n")
		os.Exit(3)
	}

	splitCmd, splitArgs, err := preProcCmd(cmd)
	if err != nil {
		_, _ = red.Printf("cmd not found in path? %v\terr:%v\n", splitCmd, err)
		os.Exit(4)
	}

	_, _ = green.Printf("cmd: %v\t\targs: %#v\n", splitCmd, splitArgs)
	if env != nil {
		_, _ = green.Printf("Env : %+v\n\n", env)
	} else {
		fmt.Printf("\n")
	}

	// Because this executes arbitrary commands from and external file
	// there is no way to make this "safe", short of whitelisting
	// hence the nolint:gosec
	ecmd := exec.Command(splitCmd, splitArgs...) //nolint:gosec
	ecmd.Stdout, ecmd.Stderr = os.Stdout, os.Stderr

	for k, v := range env {
		line := fmt.Sprintf("%s=%s", k, v)
		ecmd.Env = append(ecmd.Env, line)
	}

	err = ecmd.Run()
	if err != nil {
		return -1, fmt.Errorf("error running command %s: %v", splitCmd, err)
	}

	exitVal := ecmd.ProcessState.ExitCode()

	return exitVal, nil
}
