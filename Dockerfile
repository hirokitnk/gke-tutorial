FROM golang:1.10-alpine
ADD  . /go/src/echo
RUN go install echo

FROM alpine:latest
COPY --from=0 /go/bin/echo .
RUN apk add --no-cache ca-certificates
ENV PORT 8080
CMD ["./echo"]
