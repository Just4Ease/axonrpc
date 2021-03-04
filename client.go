package axonrpc

import (
	"context"
	"encoding/json"
	"github.com/Just4Ease/axon"
	"log"
)

type ClientConnInterface interface {
	// Invoke performs a unary RPC and returns after the response is received
	// into reply.
	Invoke(ctx context.Context, method string, args interface{}, reply interface{}) error
	// NewStream begins a streaming RPC.
	//NewStream(ctx context.Context, desc *StreamDesc, method string, opts ...CallOption) (ClientStream, error)
}

type tClient struct {
	axon.EventStore
}

func (t tClient) Invoke(ctx context.Context, method string, args interface{}, reply interface{}) error {
	payload, err := json.Marshal(args) // TODO: Make this jsonable or protoable
	if err != nil {
		log.Printf("failed to marshal input in %s with the following erros: %s", method, err)
		return err
	}
	return t.EventStore.Request(method, payload, reply)
}

func NewClient(store axon.EventStore) ClientConnInterface {
	return &tClient{store}
}
