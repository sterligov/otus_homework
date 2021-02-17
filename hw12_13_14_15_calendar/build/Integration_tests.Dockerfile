FROM golang:1.15.2

WORKDIR /go/src

COPY . ${CODE_DIR}

RUN go test -i -tags integration ./tests/integration/...

CMD go test -v -tags integration ./tests/integration/...
