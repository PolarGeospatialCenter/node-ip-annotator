FROM golang:stretch

WORKDIR /go/src/github.com/PolarGeospatialCenter/node-ip-annotator
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -vendor-only
COPY ./ .
RUN go build -o /bin/node-ip-annotator ./cmd/node-ip-annotator

FROM debian:stretch-slim
COPY --from=0 /bin/node-ip-annotator /bin/node-ip-annotator
ENTRYPOINT /bin/node-ip-annotator
