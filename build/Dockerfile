FROM golang:1.13.3-alpine3.10 as builder

RUN apk update && apk add curl make git gcc cmake g++ ca-certificates
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /go/src/github.com/beaujr/go-xerox-upload

ENV GOPATH=/go

WORKDIR /go/src/github.com/beaujr/go-xerox-upload

COPY . .

ARG GOOS
ARG GOARCH

RUN dep ensure

RUN make build GOOS=${GOOS} GOARCH=${GOARCH}

RUN mv bin/beaujr/go-xerox-upload-${GOOS}_${GOARCH} bin/beaujr/go-xerox-upload

FROM scratch

WORKDIR /
COPY --from=builder /go/src/github.com/beaujr/go-xerox-upload/bin/beaujr/go-xerox-upload go-xerox-upload
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
ENTRYPOINT ["./go-xerox-upload"]
ARG VCS_REF
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/beaujr/go-xerox-upload" \
      org.label-schema.license="Apache-2.0"