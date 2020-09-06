package natsvc

import (
	"time"

	nats "github.com/nats-io/nats.go"
)

type NatsClient struct {
	Client *nats.EncodedConn
}

func GetClient(url string) NatsClient {
	c := NatsClient{}

	nc, _ := nats.Connect(url)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	c.Client = ec

	return c
}

//Register a handler
func (self *NatsClient) RegisterTopicHandler(topic string, handler interface{}) {
	self.Client.Subscribe(topic, handler)
}

func (self *NatsClient) Publish(topic string, msg interface{}) {
	self.Client.Publish(topic, msg)
}

func (self *NatsClient) Request(topic string, msg interface{}, result interface{}) error {
	return self.Client.Request(topic, msg, &result, 30*time.Second)
}
