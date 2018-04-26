# FROM alpine

# RUN mkdir -p /app
# COPY ./assets /assets
# ADD ./cmd/app /app/app

# CMD ["/app/app"]

FROM golang:1.9

RUN mkdir -p /go/src/github.com/alikhil/quoridor-go-rpc
WORKDIR /go/src/github.com/alikhil/quoridor-go-rpc
RUN go get github.com/googollee/go-socket.io
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a cmd/main.go
RUN mkdir /app
RUN mv main /app/main

FROM scratch
COPY --from=0 /app /app
COPY ./assets /assets
CMD ["/app/main"]