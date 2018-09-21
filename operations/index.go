// Author: Josh Kendrick
// Version: v0.1.0
// Do whatever you want with this code

package operations

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// TagsByFiles is the name of the tagsByFiles bucket
const TagsByFiles = "tags_by_files"

// Index gets tag data from files/dirs and adds it to the buckets
func Index(boltDB *bolt.DB, path string) error {
	// create the files bucket if it doesnt exist
	boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tags_by_files"))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	// produce the files to the channel for the consumers
	filepaths := make(chan string, 300)
	producedCount := 0
	// number of processor(s)
	consumerCount := 20

	if info, err := os.Stat(path); err == nil && info.IsDir() {
		go func() {
			filepath.Walk(path, func(filepath string, f os.FileInfo, err error) error {
				filepaths <- filepath
				log.Printf("added path: %s", filepath)
				producedCount++

				return nil
			})
			close(filepaths)
		}()
	} else {
		filepaths <- path
		log.Printf("added path: %s", path)
		producedCount++
		close(filepaths)

		// only need one processor
		consumerCount = 1
	}

	// reporting channel
	done := make(chan int, consumerCount)

	// start the processor(s)
	for index := 0; index < consumerCount; index++ {
		go tagsProcessor(filepaths, boltDB, done, index+1)
	}

	// wait for processor(s) to finish
	consumedCount := 0
	for index := 0; index < consumerCount; index++ {
		consumedCount += <-done
	}

	log.Printf("produced: %d || consumed %d", producedCount, consumedCount)

	return nil
}

func tagsProcessor(filepaths <-chan string, boltDB *bolt.DB, done chan<- int, id int) {
	count := 0

	// get a filepath
	for {
		filepath, more := <-filepaths
		if !more {
			log.Printf("%4d consumed %d files", id, count)
			done <- count
			return
		}

		count++

		// try to get the file's tags
		metadata, err := extract(filepath)
		// if an error or no metadata, skip the file
		if err != nil || metadata == nil {
			continue
		}

		// try to pull the tags from the Subject field
		tags, exists := metadata["Subject"]

		// if still no tags found, try to pull from the Category field
		if !exists || tags == "" {
			tags, exists = metadata["Category"]
		}

		// if still no tags found, log and skip
		if !exists || tags == "" {
			log.Printf("******TAGS NOT FOUND****** %s", filepath)
			continue
		}

		// convert singleVal strings into an array
		// so all values in bolt database are the same format
		switch t := tags.(type) {
		case string:
			tags = []string{t}
		}

		log.Printf("%4d found tags - %s :: %v", id, filepath, tags)

		// for tag in tags
		// createbucketifnotexists
		// put filepath, true in bucket

		// marshal to json
		tagsAsJSON, err := json.Marshal(tags)
		if err != nil {
			log.Printf("%4d !!ERROR!! -- %v: %v", id, err, tags)
			continue
		}

		// save to bolt
		err = boltDB.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(TagsByFiles))
			err := bucket.Put([]byte(filepath), tagsAsJSON)
			return err
		})

		if err == nil {
			log.Printf("%4d saved tags - %s :: %s", id, filepath, tags)
		} else {
			log.Printf("%4d !!ERROR!! -- %v: %s", id, err, filepath)
		}
	}
}

// credit for below code goes to stale versions of:
// github.com/mostlygeek/go-exiftool
// extract data from files with exiftool
func extract(filename string) (map[string]interface{}, error) {
	// set up the command
	cmd := exec.Command("exiftool", "-json", "-binary", "--printConv", filename)
	var stdout, stderr bytes.Buffer

	// set receivers for the command output
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// run the command
	err := cmd.Run()

	container := make([]map[string]interface{}, 1, 1)

	// exiftool will exit and print valid output to stdout
	// if it exits with an unrecognized filetype, don't process
	// that situtation here
	if err != nil && stdout.Len() == 0 {
		return nil, errors.Errorf("%s", stderr.String())
	}

	// no exit error but also no output
	if stdout.Len() == 0 {
		return nil, errors.New("No output")
	}

	// try to unmarshal the bytes to json
	err = json.Unmarshal(stdout.Bytes(), &container)
	if err != nil {
		return nil, errors.Wrap(err, "JSON unmarshal failed")
	}

	// there should be at least one record
	if len(container) != 1 {
		return nil, errors.New("Expected one record")
	}

	// return the one record
	return container[0], nil
}
