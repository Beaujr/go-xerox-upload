# go-xerox-upload

## Disclaimer
I have no connection with Xerox
Use this software at your own risk.
No Warranty and no support provided.

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

Example Docker Run with google drive json files

```
docker run --name xerox
 -e google=true
 -v $(pwd)/credentials.json:/credentials.json
 -v $(pwd)/token.json:/token.json
-p 8081:10000
beaujr/go-xerox-upload:latest_amd64_
```


Example Docker Run with google drive env vars

```
docker run --name xerox
 -e google=true
 -e ClientId=<ClientId>
 -e ProjectID=<ProjectID>
 -e ClientSecret=<ClientSecret>
 -e AccessToken=<AccessToken>
 -e TokenType=Bearer
 -e RefreshToken=<RefreshToken>
 -e expiry=<expiry>
-p 8081:10000
beaujr/go-xerox-upload:latest_amd64_
```

Download to local Filesystem
```
docker run --name xerox
-e PGID=1000
-e GID=1000
-v <volume_to_mount>:<volume_dest>
-p 8081:10000
beaujr/go-xerox-upload:latest_amd64_
```

Environment Variables for App Engine

These are the variables that are set in the env_variables.yaml for appengine
```
google=true
appengine=true
ClientId=<ClientId>
ProjectID=<ProjectID>
ClientSecret=<ClientSecret>
AccessToken=<AccessToken>
TokenType=Bearer
RefreshToken=<RefreshToken>
expiry=<expiry>
```


