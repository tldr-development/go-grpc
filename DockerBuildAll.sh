#!/bin/bash

TAG_BASE=${TAG_BASE}
VERSION=1.0.10
TARGETS=(account apns inspire)

for TARGET in "${TARGETS[@]}"; do
cat << EOF > Dockerfile_$TARGET
FROM golang:1.21 as build-go
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app /app/$TARGET

FROM alpine:latest
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
COPY --from=build-go /bin/app ./
EXPOSE 3000
ENTRYPOINT ["./app"]

EOF
        echo "Building $TAG_BASE$TARGET:$VERSION"
        docker build -t $TAG_BASE$TARGET:$VERSION -f Dockerfile_$TARGET .
        gcloud builds submit . --tag=$TAG_BASE$TARGET:$VERSION
        rm Dockerfile_$TARGET
done