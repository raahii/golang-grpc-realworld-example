FROM golang:1.19-alpine

ENV ROOT=/go/src/app
ENV CGO_ENABLED 0
WORKDIR ${ROOT}

RUN apk update && apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 50051

CMD ["go", "run", "server.go"]
