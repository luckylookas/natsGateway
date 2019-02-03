package service

import (
	"fmt"
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
	sigs := make(chan os.Signal, 1)

	sub, _ := c.QueueSubscribe("api.user.info", "user-service", func (subject, replySubject string, request Request) {
		if request.Method == "GET" {
			_ = c.Publish(replySubject, &Request{Content: request.Headers[0].Value})
		}
	})
	defer sub.Unsubscribe()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("listening on nats connection...")
	<-sigs
}
