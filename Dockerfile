FROM golang:1.17.6 AS base

LABEL author="cpud1258@gmail.com"

ARG PORT

WORKDIR /src

COPY go.mod go.sum /src/

RUN  go mod download

COPY . /src/

FROM base AS build

ARG COMMIT_HASH

RUN CGO_ENABLED=0 GOOS=linux go build -a \
    -ldflags="-w -s -extldflags -static -X 'github.com/thienma1258/coding-challenge.CommitHash=$COMMIT_HASH'" \
    -installsuffix cgo -o /app/server cmd/server.go

FROM alpine:3.15.0 as release

WORKDIR /app

COPY --from=build /app/ /app/

EXPOSE $PORT

CMD [ "./server" ]
