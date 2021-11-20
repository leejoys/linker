FROM golang:latest AS compiling_stage
RUN mkdir -p /app
WORKDIR /app
ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -a -v ./cmd/server.go

FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="leejoys <test@test.test>"
WORKDIR /root/
COPY --from=compiling_stage /app/server .
RUN chmod +x /root/server
ENTRYPOINT ["./server"]