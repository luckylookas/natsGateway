package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"github.com/BillD00r/natsGateway/common"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	defer c.Close()
	sigs := make(chan os.Signal, 1)

	sub, _ := c.QueueSubscribe("api.user.info", "user-service", func (subject, replySubject string, request common.Request) {
		if request.Method == http.MethodGet {
			_ = c.Publish(replySubject, &common.Response{Content: "ok", Status: "200"})
		}
	})
	defer sub.Unsubscribe()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("listening on nats connection...")
	<-sigs
}
