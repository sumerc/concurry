package main

import (
	"bufio"
	"container/ring"
	"context"
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
	displayOutput    *bool
	verbose          *bool
	failFast         *bool
	colorize         *bool
	repeatCount      *uint
	repeatConcurrent *bool
	commandTimeout   *uint
}

// TODO:
// type ColorRing struct {
// 	colorReset
// }

var config Config
var colors = ring.New(6)
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

func GetNextColor() string {
	colors = colors.Next()
	return colors.Value.(string)
}

type TaskLogger struct {
	taskId int
	color  string
}

func (t TaskLogger) Sprintf(format string, args ...interface{}) string {
	if *config.colorize {
		format = fmt.Sprintf("%s(Task-%d) %s%s", t.color, t.taskId, format, colorReset)
	} else {
		format = fmt.Sprintf("(Task-%d) %s", t.taskId, format)
	}
	return fmt.Sprintf(format, args...)
}

// RunCmd TODO: Comment
// Note: log.Println() functions are goroutine safe. There is mutex involved when
// write() is called.
func RunCmd(command string, taskId int, wg *sync.WaitGroup, color string) {
	defer wg.Done()

	startTime := time.Now()
	taskLogger := TaskLogger{taskId: taskId, color: color}

	commandArr := strings.Split(command, " ")
	if *config.verbose {
		log.Println(taskLogger.Sprintf("Executing '%s'", command))
	}
	//cmd := exec.Command(commandArr[0], commandArr[1:]...)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*config.commandTimeout)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, commandArr[0], commandArr[1:]...)

	stdoutReader, _ := cmd.StdoutPipe()
	stderrReader, _ := cmd.StderrPipe()
	// TODO: err handling of above
	stdoutScanner := bufio.NewScanner(stdoutReader)
	stderrScanner := bufio.NewScanner(stderrReader)
	init := make(chan bool)

	// stderr shall be read everytime because if an error happens, we would like
	// to print it out.
	go func() {
		init <- true
		for stderrScanner.Scan() {
			log.Println(taskLogger.Sprintf(stderrScanner.Text()))
		}
	}()
	<-init

	cmd.Start()
	if *config.displayOutput {
		for stdoutScanner.Scan() {
			log.Println(taskLogger.Sprintf(stdoutScanner.Text()))
		}
	}
	err := cmd.Wait()

	if err != nil {
		errStr := fmt.Sprintf("%s", err)
		if ctx.Err() == context.DeadlineExceeded {
			errStr = "Task Timeout"
		}
		failure := taskLogger.Sprintf("'%s' failed. [%s] [%s]", command,
			errStr, time.Since(startTime))

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
			results = append(results, taskLogger.Sprintf("'%s' succeeded. [%s]",
				command, time.Since(startTime)))
			resultsMutex.Unlock()
		}
	}
}

func main() {
	var wg sync.WaitGroup

	startTime := time.Now()

	initColorRing()

	config.displayOutput = flag.Bool("o", false, "display stdout")
	//config.bufferIO = flag.Bool("b", false, "buffer stdout/stderr") // TODO
	config.verbose = flag.Bool("v", true, "show executed command and return values")
	config.repeatCount = flag.Uint("n", 1, "repeat command N times (synchronously)")
	config.failFast = flag.Bool("f", true, "fail if any concurrent command fails")
	config.colorize = flag.Bool("c", true, "colorize the command outputs")
	config.repeatConcurrent = flag.Bool("rc", false, "run repeated commands concurrently")
	config.commandTimeout = flag.Uint("t", 60, "timeout for executed command (secs)")
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

	taskId := 0
	for i := uint(0); i < *config.repeatCount; i++ {
		for _, command := range commands {
			command = strings.TrimSpace(command)
			if len(command) > 0 {
				wg.Add(1)
				taskId++
				go RunCmd(command, taskId, &wg, GetNextColor())
			}
		}
		if !*config.repeatConcurrent {
			wg.Wait()
		}
	}

	wg.Wait()

	// when we come here, all commands finish executing, so it is safe to read
	// results without a Lock
	for _, result := range results {
		log.Println(result)
	}

	log.Println(fmt.Sprintf("Total elapsed: %s", time.Since(startTime)))
}
