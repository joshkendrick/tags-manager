// Author: Josh Kendrick
// Version: v0.1.0
// Do whatever you want with this code

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/boltdb/bolt"

	"github.com/joshkendrick/tags-manager/operations"
)

const Version = "v0.1.0"

func main() {
	fmt.Println("HELLO! TagsManager v0.1.0 I guess?")
	fmt.Println("type 'help' to see possible commands")

	// open db
	boltDB, err := bolt.Open(
		"tags-manager.db",
		0600,
		&bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer boltDB.Close()

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
			if len(args) < 2 || args[1] == "" {
				fmt.Println("incorrect format, missing path: index <path>")
			} else {
				operations.Index(boltDB, args[1])
			}
		case "list":
			// TODO
		case "help":
			fmt.Println("things you can do")
			fmt.Println("index <path> -> gets file tags in that path added to the database")
			fmt.Println("list -> displays all data we have")
			fmt.Println("list <tag_or_absolute_filepath> -> displays data about that key")
			fmt.Println("exit -> PEACE")
		default:
			fmt.Println("no comprende")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
