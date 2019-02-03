package main

import (
	"encoding/base64"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	defer c.Close()
	db, err := bolt.Open("auth.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sigs := make(chan os.Signal, 1)

	sub, _ := c.QueueSubscribe("api.auth", "auth-service", func (subject, replySubject string, request Request) {
		if request.Method == "GET" {
			if len(request.Headers) < 1 {
				_ = c.Publish(replySubject, &Response{Status: "401"})
			} else {
				authHeader := request.Headers[0].Value
				b64Credentials := strings.Split(authHeader, " ")[1]
				buf, _ := base64.StdEncoding.DecodeString(b64Credentials)
				credentails := strings.Split(string(buf), ":")

				_ = db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("user"))
					storedPassword := string(b.Get([]byte(credentails [0])))
					if storedPassword == credentails[1] {
						_ = c.Publish(replySubject, &Response{Status: "200"})
					} else {
						_ = c.Publish(replySubject, &Response{Status: "401"})
					}
					return nil
				})
			}
		}
	})
	defer sub.Unsubscribe()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("listening on nats connection...")
	<-sigs
}
