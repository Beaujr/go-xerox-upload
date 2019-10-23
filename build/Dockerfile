FROM golang:1.13.3-alpine3.10 as builder

RUN apk update && apk add curl make git gcc cmake g++
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /golang/src/app/xerox

ENV GOPATH=/golang/

WORKDIR /golang/src/app/xerox

COPY . .

RUN dep ensure

RUN make build

FROM scratch
ARG GOOS
ARG GOARCH
COPY --from=builder /golang/src/app/xerox/bin/beaujr/go-xerox-upload-${GOOS}_${GOARCH} go-xerox-upload-${GOOS}_${GOARCH}
ENTRYPOINT ["./go-upload-xerox-linux_amd64"]