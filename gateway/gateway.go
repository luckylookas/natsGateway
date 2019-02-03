package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"strings"
	"time"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	defer c.Close()

	path := "/api/user/info"

	message := &Request{
		Method: "GET",
		Headers: []*Request_Header{
			{Key: "Authorization", Value: "Basic YWRtaW46cGFzc3dvcmQ="},
	}}

	path = strings.Trim(strings.Replace(path, "/", ".", -1), ".")
	fmt.Println(path)
	response := new(Request)
	auth := new(bool)

	_ = c.Request("api.auth", message, auth, 250 * time.Millisecond)

	if *auth {
		_ = c.Request(path, message, response, 250*time.Millisecond)
	} else {
		fmt.Println("401 unauthorized")
	}

	fmt.Println(response.Content)
}
