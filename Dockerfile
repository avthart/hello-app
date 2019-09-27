# build stage
FROM golang:1.11-alpine AS build-env
RUN apk add --no-cache ca-certificates git

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /go/src/app
COPY go.mod .

RUN go mod download

COPY main.go .

RUN go build -a -installsuffix cgo -o /go/bin/hello

# final stage
FROM scratch
WORKDIR /app
COPY --from=build-env /go/bin/hello /app/hello
ENTRYPOINT ["/app/hello"]
EXPOSE 8080