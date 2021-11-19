FROM golang:latest AS compiling_stage
RUN mkdir -p /linker
WORKDIR /linker
ADD . /linker
RUN go build ./cmd/server.go

FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="leejoys <test@test.test>"
WORKDIR /root/
COPY --from=compiling_stage /linker .
ENTRYPOINT ./linker