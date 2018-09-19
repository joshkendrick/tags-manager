// Author: Josh Kendrick
// Version: v0.1.0
// Do whatever you want with this code

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	argsRegEx := regexp.MustCompile("'?[[:graph:]]+'?")

	// read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// split the line into arguments
		args := argsRegEx.FindAllString(line, -1)

		// perform action based on the command (first arg)
		cmd := args[0]
		switch cmd {
		case "exit":
			fmt.Println("GOODBYE")
			return
		case "index":
			// TODO
		case "add":
			// TODO
			// TODO
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
