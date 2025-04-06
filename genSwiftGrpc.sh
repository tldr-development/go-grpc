#!/bin/bash

# 생성된 Swift 파일이 들어갈 디렉토리
OUTPUT_DIR="./Generated"

# proto 파일이 위치한 루트 경로
PROTO_ROOT="."

# 디렉토리 생성
mkdir -p "$OUTPUT_DIR"

# protoc 실행
protoc \
  --proto_path="$PROTO_ROOT" \
  --swift_out="$OUTPUT_DIR" \
  --grpc-swift_out="$OUTPUT_DIR" \
  $(find "$PROTO_ROOT" -name '*.proto' \
    ! -path "*/swift-protobuf/*" \
    ! -path "*/PluginExamples/*")
