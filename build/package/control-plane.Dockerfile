FROM golang:1.26.3-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git build-base

RUN apk add --no-cache curl && \
    curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

COPY go.mod ./
RUN go mod download

COPY . .
 
CMD ["air", "-c", ".air.toml"]