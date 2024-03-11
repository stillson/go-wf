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

func preProcCmd(cmd string) (string, []string, error) {

	splitCmd := strings.Split(cmd, " ")

	outCmd, err := exec.LookPath(splitCmd[0])

	if err != nil {
		return splitCmd[0], splitCmd[1:], err
	}
	return outCmd, splitCmd[1:], nil

}

func main() {

	flag.Parse()
	rubric := flag.Arg(0)
	red := color.New(color.FgHiRed)

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
	green := color.New(color.FgHiGreen)
	_, _ = green.Printf("cmd: %v\t\targs: %v\n\n", splitCmd, splitArgs)

	// Because this executes arbitrary commands from and external file
	// there is no way to make this "safe", short of whitelisting
	// hence the nolint:gosec
	ecmd := exec.Command(splitCmd, splitArgs...) //nolint:gosec
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr

	err = ecmd.Run()
	if err != nil {
		_, _ = red.Printf("%v", err)
	}

	exitVal := ecmd.ProcessState.ExitCode()

	_, _ = green.Printf("\nProcess exited with %v\n", exitVal)

	os.Exit(ecmd.ProcessState.ExitCode())
}
