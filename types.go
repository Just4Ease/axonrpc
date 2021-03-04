package axonrpc

import (
	"bytes"
	"context"
	"encoding/json"
)

type ServiceDesc struct {
	ServiceName string
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface{}
	Methods     []MethodDesc
	//Streams     []StreamDesc
	Metadata interface{}
}

// serviceInfo wraps information about a service. It is very similar to
// ServiceDesc and is constructed from it for internal purposes.
type serviceInfo struct {
	// Contains the implementation for the methods in this service.
	serviceImpl interface{}
	methods     map[string]*MethodDesc
	//streams     map[string]*StreamDesc
	mdata interface{}
}

type methodHandler func(srv interface{}, ctx context.Context, input []byte) ([]byte, error)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
}

func UnPack(in interface{}, target interface{}) error {
	var e1 error
	var b []byte
	switch in := in.(type) {
	case []byte:
		b = in
	// Do something.
	default:
		// Do the rest.
		b, e1 = json.Marshal(in)
		if e1 != nil {
			return e1
		}
	}

	buf := bytes.NewBuffer(b)
	enc := json.NewDecoder(buf)
	enc.UseNumber()
	if err := enc.Decode(&target); err != nil {
		return err
	}
	return nil
}
