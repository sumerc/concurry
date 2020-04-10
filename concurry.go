package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type Config struct {
	cmd            *string
	concurrent     *bool
	display_output *bool
	pyenv          *bool
}

var config Config

// RunCmd TODO: Comment
// Note: log.Println() functions are goroutine safe. There is mutex involved when
// write() is called.
func RunCmd(name string, arg ...string) string {

	realCmd := fmt.Sprintf("'%s %s'", name, strings.Join(arg, " "))
	log.Println("Executing ", realCmd)

	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()

	// if exitError, ok := err.(*exec.ExitError); ok {
	// 	fmt.Printf("Exit code is %d\n", exitError.ExitCode())
	// }

	outStr := string(out)
	if err != nil {
		log.Println(fmt.Sprintf("Command %s failed.", realCmd))
		log.Println(outStr)
	} else {
		log.Println(fmt.Sprintf("Command %s succeeded.", realCmd))

		if *config.display_output {
			log.Println(outStr)
		}
	}

	return outStr
}

func RunCmdConcurrent(wg *sync.WaitGroup, command string) string {
	defer wg.Done()

	commandArr := strings.Split(command, " ")

	return RunCmd(commandArr[0], commandArr[1:]...)
}

// GetPyVersions TODO: Comment
func GetPyVersions() []string {
	var result = []string{}

	out := RunCmd("pyenv", "versions", "--bare")

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		pyVersion := strings.Trim(scanner.Text(), " ")
		if strings.HasPrefix(pyVersion, "3") || strings.HasPrefix(pyVersion, "2") {
			result = append(result, pyVersion)
		}
	}

	return result
}

func main() {
	var wg sync.WaitGroup

	var pyVersions []string

	config.cmd = flag.String("cmd", "", "command to be run")
	config.concurrent = flag.Bool("concurrent", true, "control running command concurrently")
	config.display_output = flag.Bool("display-stdout", false, "control displaying command output")
	config.pyenv = flag.Bool("pyenv", false, "control running commands in supported pyenv interpreters")
	flag.Parse()

	if *config.cmd == "" {
		log.Fatalf("Fatal err: cmd is not passed\n")
	}

	commands := []string{}
	if *config.pyenv {
		// Set pyenv local interpreters
		pyVersions = GetPyVersions()
		cmdSuffix := append([]string{"local"}, pyVersions...) // prepend
		RunCmd("pyenv", cmdSuffix...)

		// Generate commands

		for _, pyVersion := range pyVersions {
			pyExecutable := fmt.Sprintf("python%s", pyVersion[:3])
			commands = append(commands, fmt.Sprintf("%s %s", pyExecutable, *config.cmd))
		}
	} else {
		commands = strings.Split(*config.cmd, "|")
	}

	// make sure no unnecessary whitespace exists
	for i := range commands {
		commands[i] = strings.TrimSpace(commands[i])
	}

	wg.Add(len(commands))
	for _, command := range commands {

		//fmt.Println(commandArr[0], commandArr[1:])
		if *config.concurrent {
			go RunCmdConcurrent(&wg, command)
		} else {
			RunCmdConcurrent(&wg, command)
		}
	}

	wg.Wait()
}
