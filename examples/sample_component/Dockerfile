FROM golang:1.17 AS build

WORKDIR /workspace/github.com/sustglobal/gost

COPY . /workspace/github.com/sustglobal/gost

RUN go mod download

RUN CGO_ENABLED=0 go build -o /sample_component ./examples/sample_component/cmd



FROM alpine:latest AS package

# address the glibc dependency
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=build /sample_component /usr/local/bin/sample_component

CMD ["/usr/local/bin/sample_component"]
