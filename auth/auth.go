package main

import (
	"fmt"
	"github.com/BillD00r/natsGateway/common"
	"github.com/boltdb/bolt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	defer c.Close()
	conn, _ := bolt.Open("auth.db", 0600, nil)
	boltDb := boltDb{conn}
	defer boltDb.Close()
	sigs := make(chan os.Signal, 1)

	sub, _ := c.QueueSubscribe("api.auth", "auth-service", func(subject, replySubject string, request common.Request) {
		user, ok_user := request.HeaderByName("x-username")
		password, ok_pass := request.HeaderByName("x-password")

		if ok_pass && ok_user {
			if storedPassword, ok := boltDb.findPassword(user); *ok && *storedPassword == password {
				_ = c.Publish(replySubject, &common.Response{Status: "200"})
			} else {
				_ = c.Publish(replySubject, &common.Response{Status: "401"})
			}
		} else {
			_ = c.Publish(replySubject, &common.Response{Status: "401"})
		}
	})

	defer sub.Unsubscribe()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("listening on nats connection...")
	<-sigs
}
