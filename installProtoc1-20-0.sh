git clone https://github.com/grpc/grpc-swift.git
cd grpc-swift
git checkout 1.20.0

swift build -c release

# 바이너리 복사
cp .build/release/protoc-gen-grpc-swift /usr/local/bin/