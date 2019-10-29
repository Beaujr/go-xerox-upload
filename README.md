# go-xerox-upload

## Introduction

> A HTTP server to use with Xerox WorkCentre 3345DNi "scan to" functionality
Support for:
- Google Drive (credentials.json & token.json) requiredt.
- FileSystem upload

Default port is 10000


## Installation

Docker Image Available
https://hub.docker.com/r/beaujr/go-xerox-upload

AMD 64 Support and ARM Support

Example Docker Run

```
docker run --name xerox
 -e google=true
 -v $(pwd)/credentials.json:/credentials.json
 -v $(pwd)/token.json:/token.json
-p 8081:10000
beaujr/go-xerox-upload:0.1-amd64-2c6856d
```
