# Axon RPC

---
An RPC over event streams.

With Axon RPC, We can execute RPC calls like gRPC but this time, over event streams.

> It uses `github.com/Just4Ease/axon` as it's connection stream, while helping us with typed function signatures for RPC calls.

**AxonRPC** is also bundled with a code generator, where the same protobuf used to generate go code and gRPC handlers, is used to generate go code and axonRPC handlers, exactly the same as gRPC in go, so that there wouldn't be so much knowledge gaps.


### How to use

```shell
# install protoc plugin:

go install github.com/Just4Ease/axonrpc/protoc-gen-axonrpc

# add -axonrpc_out=/path/to/output 
# for example:
PROTO_SRC_DIR="${PWD}/proto" # Note: These variables here are simply paths to whatever directory your `*.proto` is located
PROTO_DST_DIR="${PWD}/proto"
mkdir -p ${PROTO_DST_DIR} && protoc -I=${PROTO_SRC_DIR}  --axonrpc_out=. ${PROTO_SRC_DIR}/service.proto
```


### TODO:
Things left to add to axonRPC

- [x] UnaryServer
- [x] UnaryClient
- [ ] UnaryServerInterceptor
- [ ] ServerStream
- [ ] ClientStream
- [ ] CallOptions

>Note, you could also just use `github.com/Just4Ease/axon` for streams with any broker of your preference.
Cheers.

You can join me to improve this package.
While this is experimental, I'm using it in a few systems, some have been modified to go to prod.

---

> @Just4Ease on all social platforms ( including linkedin )