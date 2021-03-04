#!/bin/zsh

PROTO_SRC_DIR="${PWD}/proto"
PROTO_DST_DIR="${PWD}/proto"
mkdir -p ${PROTO_DST_DIR} && protoc -I=${PROTO_SRC_DIR}  --axonrpc_out=${PROTO_DST_DIR} ${PROTO_SRC_DIR}/service.proto