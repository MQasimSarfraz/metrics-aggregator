FROM golang:1.11.4-alpine3.8 as go-builder

ENV PACKAGE github.com/MQasimSarfraz/metrics-aggregator
ENV CGO_ENABLED 0
ENV GO111MODULE=on

WORKDIR $GOPATH/src/$PACKAGE

# create directories for binary
RUN mkdir -p /out

# copy sources, test and build the application
COPY . ./
RUN go vet ./...
RUN go test --parallel=1 ./...
RUN go build -v -ldflags="-s -w" -o /out/metrics-aggregator ./cmd/metrics-aggregator


# build the final container image
FROM alpine:3.8

EXPOSE 3080

COPY --from=go-builder /out/metrics-aggregator /

ENTRYPOINT ["/metrics-aggregator"]