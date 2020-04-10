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

// RunCmd TODO: Comment
func RunCmd(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()

	fmt.Println(fmt.Sprintf("Running command '%s %s'", name, strings.Join(arg, " ")))

	if err != nil {
		log.Fatalf("Fatal err: %s [%s]\n", out, err)
	}

	return string(out)
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

	// if len(os.Args) > 1 {
	// 	pyVersions = os.Args[1:]
	// } else {

	cmdPtr := flag.String("cmd", "", "pass the command that will run in Python interpreter. e.x: setup.py install")
	flag.Parse()

	if *cmdPtr == "" {
		log.Fatalf("Fatal err: cmd is not passed\n")
	}

	//fmt.Println("cmd=", *cmdPtr)

	pyVersions = GetPyVersions()

	// prepend local at the start of versions
	cmdSuffix := append([]string{"local"}, pyVersions...)

	out := RunCmd("pyenv", cmdSuffix...)
	fmt.Println(out)

	// clean first
	RunCmd("rm", "-Rf", "build/")
	RunCmd("rm", "-Rf", "dist/")

	wg.Add(len(pyVersions))
	for _, pyVersion := range pyVersions {
		go func(version string) {
			defer wg.Done()

			majorVersion := version[:3]
			pyExecutable := fmt.Sprintf("python%s", majorVersion)
			buildDir := fmt.Sprintf("/tmp/python%s", majorVersion)

			// out := RunCmd(pyExecutable, strings.Split(*cmdPtr, " ")...)
			// fmt.Println(out)

			out := RunCmd(pyExecutable, "setup.py", "clean")
			fmt.Println(out)

			out = RunCmd(pyExecutable, "setup.py", "build", "-b", buildDir, "install")
			fmt.Println(out)

			out = RunCmd(pyExecutable, "run_tests.py")
			fmt.Println(out)

		}(pyVersion)
	}

	wg.Wait()
}
