FROM golang:1.14.4 as builder

RUN git config --global core.eol lf
RUN git config --global core.autocrlf input

COPY . /go/src/github.com/kucinghitam13/gcs-drive
WORKDIR /go/src/github.com/kucinghitam13/gcs-drive

# building binary file
RUN go mod download
RUN go mod verify
RUN go build -v -o ./build/gcs-drive .

FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates

FROM ubuntu

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/src/github.com/kucinghitam13/gcs-drive/build/gcs-drive ./build/
COPY --from=builder /go/src/github.com/kucinghitam13/gcs-drive/files /root/files/
RUN ls /root/
WORKDIR /build

ENTRYPOINT ["./gcs-drive"]
