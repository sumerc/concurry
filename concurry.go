package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type Config struct {
	//cmd            *string
	displayOutput *bool
	verbose       *bool
	failFast      *bool
	repeatCount   *uint
}

var config Config

//var parentProcessName string

// RunCmd TODO: Comment
// Note: log.Println() functions are goroutine safe. There is mutex involved when
// write() is called.
func RunCmd(command string, wg *sync.WaitGroup) string {
	defer wg.Done()

	var cmd *exec.Cmd
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "windows" {

		commandArr := strings.Split(command, " ")
		if *config.verbose {
			log.Println("Executing '", commandArr)
		}
		cmd = exec.Command(commandArr[0], commandArr[1:]...)

	} else {
		log.Fatalf("Unsupported platform. [%s]", runtime.GOOS)
	}

	out, err := cmd.CombinedOutput()

	// if exitError, ok := err.(*exec.ExitError); ok {
	// 	fmt.Printf("Exit code is %d\n", exitError.ExitCode())
	// }

	outStr := string(out)
	if err != nil {
		log.Println(fmt.Sprintf("Command %s failed.", command))
		log.Println(outStr)

		if *config.failFast {
			os.Exit(1)
		}

	} else {
		if *config.verbose {
			log.Println(fmt.Sprintf("Command %s succeeded.", command))
		}

		if *config.displayOutput {
			log.Println(outStr)
		}
	}

	return outStr
}

func main() {
	var wg sync.WaitGroup

	// TODO: Ability to set custom shell
	//config.cmd = flag.String("cmd", "", "command to be run")
	config.displayOutput = flag.Bool("o", false, "display command output")
	config.verbose = flag.Bool("v", true, "show executed command and return values")
	config.repeatCount = flag.Uint("n", 1, "repeat command N times (synchronously)")
	config.failFast = flag.Bool("f", true, "fail if any concurrent command fails")
	//config.waitTimeout
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	command, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Stdin could not be read. [%s]", err)
	}

	// // get parent process name
	//process, err := ps.FindProcess(os.Getppid())
	// if err != nil {
	// 	log.Fatalf("No Parent PID. [%s]", err)
	// }
	// parentProcessName = process.Executable()

	commands := strings.Split(command, ";")

	for i := uint(0); i < *config.repeatCount; i++ {
		for _, command := range commands {
			command = strings.TrimSpace(command)
			if len(command) > 0 {
				wg.Add(1)
				go RunCmd(command, &wg)
			}
		}
		wg.Wait()
	}
}
