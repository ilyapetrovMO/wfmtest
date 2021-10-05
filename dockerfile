FROM golang:latest AS buildStage

WORKDIR /go/src/wfmtest/
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/api

FROM alpine:latest
WORKDIR /root/
COPY --from=buildStage /go/src/wfmtest/bin/app ./
CMD ["./app"]
