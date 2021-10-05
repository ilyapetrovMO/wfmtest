FROM golang:latest AS build
WORKDIR /go/src/wfmtest/
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/api

FROM alpine:latest
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=build /go/src/wfmtest/bin/app /app/
WORKDIR /app
EXPOSE 8080 8080
ENTRYPOINT ["./app"]
