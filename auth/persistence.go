package main

import (
	"github.com/boltdb/bolt"
)

type boltDb struct {
	*bolt.DB
}

func (db boltDb) findPassword(user string) (password *string, ok *bool) {
	_ = db.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("user"))
			s := string(b.Get([]byte(user)))
			*password = s
			*ok = s == ""
			return nil
		})
	return password, ok
}
