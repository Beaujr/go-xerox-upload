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
COPY --from=builder /golang/src/app/xerox/bin/beaujr/go-upload-xerox-linux_amd64 go-upload-xerox-linux_amd64
ENTRYPOINT ["./go-upload-xerox-linux_amd64"]



git init
git add README.md
git commit -m "first commit"
git remote add origin git@github.com:Beaujr/go-xerox-upload.git
git push -u origin master