FROM golang:1.18.10 as build-go
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app .

FROM alpine:latest
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
COPY --from=build-go /bin/app ./
EXPOSE 3000
ENTRYPOINT ["./app"]