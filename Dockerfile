FROM golang:latest AS compiling_stage
RUN mkdir -p /app
WORKDIR /app
ADD . /app
RUN go build ./cmd/server.go

FROM debian:stretch
LABEL version="1.0.0"
LABEL maintainer="leejoys <test@test.test>"
WORKDIR /root/
COPY --from=compiling_stage /app/server .
RUN chmod +x /root/server
ENTRYPOINT ./server