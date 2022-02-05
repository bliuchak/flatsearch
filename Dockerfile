FROM golang:1.17.5-alpine3.15 as builder
WORKDIR /go/src/app
RUN go get github.com/cespare/reflex
COPY . .
RUN CGO_ENABLED=0 go build -o /flatsearch -v

FROM scratch
COPY --from=builder /flatsearch /
USER 9000
ENTRYPOINT [ "/flatsearch" ]
