package main

import (
	"bufio"
	"container/ring"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Config struct {
	//cmd            *string
	displayOutput *bool
	verbose       *bool
	failFast      *bool
	repeatCount   *uint
}

var config Config
var colors = ring.New(6)
var colorMutex = &sync.Mutex{}
var results = []string{}
var resultsMutex = &sync.Mutex{}
var colorReset = "\033[0m"

// initColorRing initializes an array containing color codes for terminals to
// output different colors for different Commands.
func initColorRing() {
	r := colors
	r.Value = "\033[31m" // red
	r = r.Next()
	r.Value = "\033[32m" // green
	r = r.Next()
	r.Value = "\033[33m" // yellow
	r = r.Next()
	r.Value = "\033[34m" // blue
	r = r.Next()
	r.Value = "\033[35m" // purple
	r = r.Next()
	r.Value = "\033[36m" // cyan
}

func getColor() string {
	colorMutex.Lock()
	defer colorMutex.Unlock()
	colors = colors.Next()
	return colors.Value.(string)
}

// RunCmd TODO: Comment
// Note: log.Println() functions are goroutine safe. There is mutex involved when
// write() is called.
func RunCmd(command string, wg *sync.WaitGroup) {
	defer wg.Done()

	color := getColor()
	startTime := time.Now()

	commandArr := strings.Split(command, " ")
	if *config.verbose {
		log.Println("Executing ", commandArr)
	}
	cmd := exec.Command(commandArr[0], commandArr[1:]...)

	stdoutReader, _ := cmd.StdoutPipe()
	stderrReader, _ := cmd.StderrPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	stderrScanner := bufio.NewScanner(stderrReader)
	init := make(chan bool)

	// stderr shall be read everytime because if an error happens, we would like
	// to print it out.
	go func() {
		init <- true
		for stderrScanner.Scan() {
			log.Println(color, stderrScanner.Text(), colorReset)
		}
	}()
	<-init

	cmd.Start()
	if *config.displayOutput {
		for stdoutScanner.Scan() {
			log.Println(color, stdoutScanner.Text(), colorReset)
		}
	}
	err := cmd.Wait()

	if err != nil {
		failure := fmt.Sprintf("%s'%s' failed. [%s] [%s] %s", color, command,
			err, time.Since(startTime), colorReset)

		if *config.failFast {
			log.Println(failure)
			os.Exit(1)
		}

		resultsMutex.Lock()
		results = append(results, failure)
		resultsMutex.Unlock()

	} else {
		if *config.verbose {
			resultsMutex.Lock()
			results = append(results, fmt.Sprintf("%s'%s' succeeded. [%s] %s", color, command,
				time.Since(startTime), colorReset))
			resultsMutex.Unlock()
		}
	}
}

func main() {
	var wg sync.WaitGroup

	startTime := time.Now()

	initColorRing()

	// TODO: Ability to set custom shell
	//config.cmd = flag.String("cmd", "", "command to be run")
	config.displayOutput = flag.Bool("o", false, "display command output")
	config.verbose = flag.Bool("v", true, "show executed command and return values")
	config.repeatCount = flag.Uint("n", 1, "repeat command N times (synchronously)")
	config.failFast = flag.Bool("f", true, "fail if any concurrent command fails")
	//config.waitTimeout
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	commands := []string{}
	for {
		command, _ := reader.ReadString('\n')

		// EOF?
		if len(command) == 0 {
			break
		}

		commands = append(commands, command)
	}

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

	// when we come here, all commands finish executing, so it is safe to read
	// results without a Lock
	for _, result := range results {
		log.Println(result)
	}

	log.Println(fmt.Sprintf("Total elapsed: %s", time.Since(startTime)))
}
