package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main() {
	db, err := bolt.Open("auth.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_ = db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucket([]byte("user"))
		fmt.Println("inserting admin user")
		_ = b.Put([]byte("admin"), []byte("password"))
		return nil
	})
}
