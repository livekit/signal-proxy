FROM golang:latest

RUN go install github.com/pion/turn/v3/examples/turn-server/simple@latest


CMD ["simple", "-users", "foo=bar", "-public-ip", "127.0.0.1"] 