FROM golang:1.13.3-alpine3.10 as builder

RUN apk update && apk add curl make git gcc cmake g++
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /golang/src/app/xerox

ENV GOPATH=/golang/

WORKDIR /golang/src/app/xerox

COPY . .

ARG GOOS
ARG GOARCH

RUN dep ensure

RUN make build GOOS=${GOOS} GOARCH=${GOARCH}

RUN mv bin/beaujr/go-xerox-upload-${GOOS}_${GOARCH} bin/beaujr/go-xerox-upload

FROM scratch

COPY --from=builder /golang/src/app/xerox/bin/beaujr/go-xerox-upload go-xerox-upload
ENTRYPOINT ["./go-xerox-upload"]