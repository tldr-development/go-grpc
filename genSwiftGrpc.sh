protoc ./wscan/proto/wscan.proto \
  --proto_path=. \
  --swift_out=./Generated \
  --plugin=protoc-gen-grpc-swift=./bin/protoc-gen-grpc-swift \
  --grpc-swift_out=./Generated