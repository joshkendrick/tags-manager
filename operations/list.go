// Author: Josh Kendrick
// Do whatever you want with this code

package operations

import (
	"fmt"

	"github.com/boltdb/bolt"

	consts "github.com/joshkendrick/tags-manager/utils"
)

// List gets tag data from files/dirs and adds it to the buckets
func List(boltDB *bolt.DB, args []string) {
	// if listing a term:
	if len(args) > 1 {
		searchTerm := args[1]
		listArg(boltDB, searchTerm)
	} else {
		listAll(boltDB)
	}
}

func listArg(boltDB *bolt.DB, arg string) {
	boltDB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(arg))

		// if arg was a tag and we have a bucket for that tag
		if b != nil {
			tagStats := b.Stats()
			fmt.Printf("%s : %d files", arg, tagStats.KeyN)

			// print all the files
			b.ForEach(func(k, v []byte) error {
				fmt.Println(string(k))
				return nil
			})
		} else {
			// it's a file? try to get tags for that file
			b = tx.Bucket([]byte(consts.TagsByFiles))
			v := b.Get([]byte(arg))
			fmt.Printf("result: %s\n", v)
		}

		return nil
	})
}

func listAll(boltDB *bolt.DB) {
	// loop through the tags buckets
	boltDB.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println(string(name))

			return nil
		})
		return nil
	})

	// then loop through and display tags_by_files
	boltDB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(consts.TagsByFiles))

		b.ForEach(func(k, v []byte) error {
			fmt.Println(string(k))
			return nil
		})
		return nil
	})
}
