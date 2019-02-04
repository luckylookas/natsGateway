package main

import (
	"github.com/BillD00r/natsGateway/common"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"log"
	"net/http"
	"strings"
	"time"
)

var natsConnection *nats.EncodedConn

func Route(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	path := r.URL.Path
	if username, password, ok := r.BasicAuth(); !ok {
		w.WriteHeader(401)
		return
	} else {
		message := &common.Request{
			Method: http.MethodGet,
			Headers: []*common.Header{
				{Key: "x-username", Value: username},
				{Key: "x-password", Value: password},
			}}

		response := new(common.Response)
		err := natsConnection.Request("api.auth", message, response, 250*time.Millisecond)

		if err != nil {
			w.WriteHeader(500)
		} else if response.Status == "200" {
			path = strings.Trim(strings.Replace(path, "/", ".", -1), ".")
			_ = natsConnection.Request(path, message, response, 250*time.Millisecond)
			_, _ = w.Write([]byte(response.Content))
		} else {
			w.WriteHeader(401)
		}
	}
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	natsConnection, _ = nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	defer natsConnection.Close()

	router := httprouter.New()
	router.GET("/api/user/info", Route)

	log.Fatal(http.ListenAndServe(":8080", router))
}
