package chatroom

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func init() {
	var err error

	db, err = bolt.Open("mychat.db", 0600, nil)
	if err != nil {
		log.Fatalf("crash!: %v", err)
	}
}

var bucketName = []byte("posts")

func Log(ts time.Time, msg string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return b.Put([]byte(ts.Format(time.RFC3339)), []byte(msg))
	})
	return err
}

func Retrieve(begin, end time.Time) ([]string, error) {
	// TODO: should it take a channel?
	var res []string
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()

		min := []byte(begin.Format(time.RFC3339))
		max := []byte(end.Format(time.RFC3339))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			// TODO: k? (timestamp)
			res = append(res, fmt.Sprintf("%s: %s", k, v))
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
