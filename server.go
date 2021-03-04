package axonrpc

import (
	"context"
	"fmt"
	"github.com/Just4Ease/axon"
	"log"
	"reflect"
)

type Server struct {
	conn         axon.EventStore
	serveContext context.Context
	serve        bool
	services     map[string]*serviceInfo // service name -> service info
	serviceImpl  *ServiceDesc
}

// NewServer creates a axonRPC server which has no service registered and has not
// started to accept requests yet.
func NewServer(context context.Context, eventStore axon.EventStore) *Server {
	s := &Server{
		conn:         eventStore,
		services:     make(map[string]*serviceInfo),
		serveContext: context,
		serve:        false,
	}
	return s
}

// ServiceRegistrar wraps a single method that supports service registration. It
// enables users to pass concrete types other than axonRPC.Server to the service
// registration methods exported by the IDL generated code.
type ServiceRegistrar interface {
	// RegisterService registers a service and its implementation to the
	// concrete type implementing this interface.  It may not be called
	// once the server has started serving.
	// desc describes the service and its methods and handlers. impl is the
	// service implementation which is passed to the method handlers.
	RegisterService(desc *ServiceDesc, impl interface{})
}

// RegisterService registers a service and its implementation to the axonRPC
// server. It is called from the IDL generated code. This must be called before
// invoking Serve. If ss is non-nil (for legacy code), its type is checked to
// ensure it implements sd.HandlerType.
func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			log.Fatalf("axonRPC: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.register(sd, ss)
}

func (s *Server) register(sd *ServiceDesc, ss interface{}) {
	if s.serve {
		log.Fatalf("axonRPC: Server.RegisterService after Server.Serve for %q", sd.ServiceName)
	}
	if _, ok := s.services[sd.ServiceName]; ok {
		log.Fatalf("axonRPC: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	}
	info := &serviceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*MethodDesc),
		//streams:     make(map[string]*StreamDesc),
		//mdata:       sd.Metadata,
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	//for i := range sd.Streams {
	//	d := &sd.Streams[i]
	//	info.streams[d.StreamName] = d
	//}
	s.services[sd.ServiceName] = info
}

// axonRPCServerInterface contains methods from axonRPC.Server which are used by the
// axonRPCServer type here. This is useful for overriding in unit tests.
type tendonisServerInterface interface {
	RegisterService(*ServiceDesc, interface{})
	Serve() error
	Stop()
	GracefulStop()
}

func (s *Server) Serve() error {
	handlers := make([]axon.EventHandler, 0)
	for serviceName, svc := range s.services {
		for _, method := range svc.methods {
			handlers = append(handlers, func() error {
				endpoint := fmt.Sprintf("%s.%s", serviceName, method.MethodName)
				log.Printf("axonRPC endpoint registration: %s", endpoint)
				return s.conn.Reply(endpoint, func(input []byte) ([]byte, error) {
					return method.Handler(svc.serviceImpl, s.serveContext, input)
					//reply, appErr := method.Handler(, s.serveContext, )
				})
			})
		}
	}
	s.serve = true
	s.conn.Run(s.serveContext, handlers...)
	return nil
}

func (s *Server) Stop() {
	s.serveContext.Done()
}

func (s *Server) GracefulStop() {
	s.serveContext.Done()
}
