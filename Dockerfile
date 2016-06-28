FROM golang
ENV GO15VENDOREXPERIMENT 1
ADD . /go/src/github.com/crossdock/crossdock
RUN go install github.com/crossdock/crossdock
ENTRYPOINT /go/bin/crossdock
