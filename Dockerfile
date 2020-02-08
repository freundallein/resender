FROM golang:alpine AS intermediate

RUN apk update && \
    apk add --no-cache git make

RUN adduser -D -g '' resender

WORKDIR $GOPATH/src/

COPY . .

RUN go mod download
RUN go mod verify
RUN make build

FROM scratch


COPY --from=intermediate /go/src/bin/resender /go/bin/resender
COPY --from=intermediate /etc/passwd /etc/passwd

USER resender

WORKDIR /go/bin

CMD ["/go/bin/resender"]