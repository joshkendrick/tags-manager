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
	consts "github.com/joshkendrick/tags-manager/utils"
)

func main() {
	fmt.Println("HELLO! TagsManager v0.1.0 I guess?")
	fmt.Println("type 'help' to see possible commands")

	boltDB := initDB()

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
			boltDB.Close()
			fmt.Println("GOODBYE")
			return

		case "index":
			if len(args) < 2 || args[1] == "" {
				fmt.Println("incorrect format, missing path: index <path>")
			} else {
				operations.Index(boltDB, args[1])
			}

		case "list":
			operations.List(boltDB, args)

		case "clear":
			// close current db
			boltDB.Close()

			// delete the file
			os.Remove(consts.DbName)

			// create a new db
			boltDB = initDB()

			// notify user
			fmt.Println("fresh, empty database!")

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

func initDB() *bolt.DB {
	// open db
	boltDB, err := bolt.Open(
		consts.DbName,
		0600,
		&bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	// create the files bucket if it doesnt exist
	boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(consts.TagsByFiles))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	return boltDB
}
