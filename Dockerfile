FROM golang:alpine as builder
LABEL stage=go-builder
WORKDIR /app/
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -v -o /app/bin/openlist-proxy -ldflags="-w -s" .

FROM alpine:3
LABEL MAINTAINER="OpenList"
WORKDIR /app/
COPY --from=builder /app/bin/openlist-proxy ./
EXPOSE 5243
CMD [ "./openlist-proxy" ]
