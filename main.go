package main

import (
	"flag"
	"fmt"
	"github.com/stillson/go-wf/rcfile"
	"github.com/stillson/go-wf/rcparse"
	"os"
	"os/exec"
	"strings"
)

func main() {

	flag.Parse()
	rub := flag.Arg(0)
	fmt.Println(rub)

	f, err := rcfile.GetRCFile()
	if err != nil {
		fmt.Printf("Error getting rcfile:%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\trcfile %v\n", f)

	rc, err := rcparse.NewPlainRcFile(f)

	if err != nil {
		fmt.Printf("Error parsing rcfile:%v\n", err)
		os.Exit(2)
	}

	fmt.Printf("\tRC: %v\n", rc)

	cmd, exists := rc.GetCommand(rub)

	if !exists {
		fmt.Printf("rubric does not exist\n")
		os.Exit(3)
	}

	fmt.Println(cmd)

	splitCmd := strings.Split(cmd, " ")
	// Because this executes arbitrary commands from and external file
	// there is no way to make this "safe", short of whitelisting
	// hence the nolint:gosec
	ecmd := exec.Command(splitCmd[0], splitCmd[1:]...) //nolint:gosec
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr

	err = ecmd.Run()
	if err != nil {
		fmt.Printf("%v", err)
	}

	os.Exit(ecmd.ProcessState.ExitCode())
}
