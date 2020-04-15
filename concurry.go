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

	"github.com/mitchellh/go-ps"
)

type Config struct {
	//cmd            *string
	display_output *bool
	verbose        *bool
}

var config Config
var parentProcessName string

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
	} else {
		if *config.verbose {
			log.Println(fmt.Sprintf("Command %s succeeded.", command))
		}

		if *config.display_output {
			log.Println(outStr)
		}
	}

	return outStr
}

func main() {
	var wg sync.WaitGroup

	// TODO: Ability to set custom shell
	//config.cmd = flag.String("cmd", "", "command to be run")
	config.display_output = flag.Bool("o", false, "control displaying command output")
	config.verbose = flag.Bool("v", true, "control showing executed commands and return values")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	command, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Stdin could not be read. [%s]", err)
	}

	// get parent process name
	process, err := ps.FindProcess(os.Getppid())
	if err != nil {
		log.Fatalf("No Parent PID. [%s]", err)
	}
	parentProcessName = process.Executable()

	commands := strings.Split(command, ";")

	for _, command := range commands {
		command = strings.TrimSpace(command)
		if len(command) > 0 {
			wg.Add(1)
			go RunCmd(command, &wg)
		}
	}

	wg.Wait()
}
