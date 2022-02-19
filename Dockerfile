FROM golang:1.17-alpine as builder

RUN apk update && apk add curl make git gcc cmake g++ ca-certificates

RUN mkdir -p /go/src/github.com/beaujr/go-xerox-upload

ENV GOPATH=/go

WORKDIR /go/src/github.com/beaujr/go-xerox-upload

COPY . .

ARG GOOS
ARG GOARCH
RUN make go_mod
RUN make go_upload_xerox GOOS=${GOOS} GOARCH=${GOARCH}

RUN mv bin/beaujr/go-xerox-upload-${GOOS}_${GOARCH} bin/beaujr/go-xerox-upload

FROM scratch
WORKDIR /

ENV PGID=1000
ENV GID=1000
COPY --from=builder /go/src/github.com/beaujr/go-xerox-upload/bin/beaujr/go-xerox-upload go-xerox-upload
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

ENTRYPOINT ["./go-xerox-upload"]
ARG VCS_REF
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/beaujr/go-xerox-upload" \
      org.label-schema.license="Apache-2.0"