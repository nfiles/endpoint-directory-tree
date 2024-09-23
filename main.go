package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func assertCommand(command string, expected int, args []string) {
	if expected != len(args) {
		log.Fatalf("command %s expects %d arg(s), but received %d\n",
			command, expected, len(args))
	}
}

const (
	CREATE = "CREATE"
	DELETE = "DELETE"
	MOVE   = "MOVE"
	LIST   = "LIST"
)

func parsePath(path string) []string {
	return strings.Split(path, SEPARATOR)
}

func main() {
	root := NewDirectory()

	// only echo the command if stdin has been redirected from a file
	// otherwise the user is typing commands in a REPL and doesn't need an echo
	stdinStat, _ := os.Stdin.Stat()
	echo := (stdinStat.Mode() & os.ModeCharDevice) != os.ModeCharDevice

	// read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		line := scanner.Text()

		if echo {
			fmt.Println(line)
		}

		args := strings.Split(line, " ")
		if len(args) == 0 {
			continue
		}

		var cmdErr error = nil
		switch args[0] {
		case CREATE:
			assertCommand(CREATE, 1, args[1:])
			cmdErr = root.Create(parsePath(args[1]))
		case DELETE:
			assertCommand(DELETE, 1, args[1:])
			cmdErr = root.Delete(parsePath(args[1]))
		case MOVE:
			assertCommand(MOVE, 2, args[1:])
			cmdErr = root.Move(parsePath(args[1]), parsePath(args[2]))
		case LIST:
			assertCommand(LIST, 0, args[1:])
			root.List()
		}

		if cmdErr != nil {
			fmt.Printf("%s\n", cmdErr.Error())
		}
	}
}
