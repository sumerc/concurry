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
func RunCmd(command string, wg *sync.WaitGroup) string {
	if wg != nil {
		defer wg.Done()
	}

	//arg := strings.Split(command, " ")

	//realCmd := fmt.Sprintf("'%s %s'", name, strings.Join(arg, " "))
	log.Println("Executing ", command)

	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()

	// if exitError, ok := err.(*exec.ExitError); ok {
	// 	fmt.Printf("Exit code is %d\n", exitError.ExitCode())
	// }

	outStr := string(out)
	if err != nil {
		log.Println(fmt.Sprintf("Command %s failed.", command))
		log.Println(outStr)
	} else {
		log.Println(fmt.Sprintf("Command %s succeeded.", command))

		if *config.display_output {
			log.Println(outStr)
		}
	}

	return outStr
}

// GetPyVersions TODO: Comment
func GetPyVersions() []string {
	var result = []string{}

	out := RunCmd("pyenv versions --bare", nil)

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
		flag.Usage()
		return
	}

	commands := []string{}
	if *config.pyenv {
		// Set pyenv local interpreters
		pyVersions = GetPyVersions()

		RunCmd(fmt.Sprintf("pyenv local %s", strings.Join(pyVersions, " ")), nil)

		// Generate commands
		for _, pyVersion := range pyVersions {
			pyExecutable := fmt.Sprintf("python%s", pyVersion[:3])
			cmd := strings.ReplaceAll(*config.cmd, "python", pyExecutable)
			commands = append(commands, cmd)
		}
	} else {
		commands = strings.Split(*config.cmd, ";")
	}

	wg.Add(len(commands))
	for _, command := range commands {
		// make sure no unnecessary whitespace exists
		command = strings.TrimSpace(command)

		if *config.concurrent {
			go RunCmd(command, &wg)
		} else {
			RunCmd(command, &wg)
		}
	}

	wg.Wait()
}
