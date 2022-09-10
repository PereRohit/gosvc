# gosvc
Cli to generate skeleton HTTP server in Go with basic Health endpoints.

The HTTP server will be generated with [gorilla mux](https://github.com/gorilla/mux).

Gos [text/template](https://pkg.go.dev/text/template) is used to generate basic endpoints.

## Installation
```go
go install github.com/PereRohit/gosvc/cmd/gosvc@latest
```

## Usage
```bash
gosvc --init <module-name>
```

### Example: 
```bash
gosvc --init github.com/PereRohit/test-server
```
#### Output: will also create the folder `test-server`
```bash
tree test-server

test-server
├── cmd
│   └── test-server
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── handler
│   │   └── handler.go
│   └── router
│       └── router.go
└── pkg
    └── mocks
        └── internal
            └── mock
                └── mock_handler.go

9 directories, 6 files
```

### Help: only one option is available at the moment
```bash
gosvc --help
```

## Server Usage
### Start server on `localhost` with default port `80`
```bash
go run cmd/test-server/main.go
```
#### Output:
```bash
INF: 2022-09-10 15:52:27.902881 +0000 UTC | Starting server(:80)
```
### Health Check (cURL)
```bash
curl -v http://localhost/health

*   Trying 127.0.0.1:80...
* Connected to localhost (127.0.0.1) port 80 (#0)
> GET /health HTTP/1.1
> Host: localhost
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 10 Sep 2022 15:56:21 GMT
< Content-Length: 42
< 
{"status":200,"message":"OK","data":null}
* Connection #0 to host localhost left intact
```
#### Server Logs
```bash
INF: 2022-09-10 15:52:27.902881 +0000 UTC | Starting server(:80)
INF: 2022-09-10 15:55:22.11881 +0000 UTC |      127.0.0.1:50053 |   GET |              /health | 200 |  563.583µs | {"status":200,"message":"OK","data":null}
```

### Stop the Server
The server will be terminated when it receives the following signals:
- Interrupt
- SIGHUP
- SIGINT
- SIGTERM
- SIGQUIT

Kindly check the keyboard combinations for your OS for the above signals.

#### Server logs
`Command` + `C` on MacOS gracefully shuts down the server with the following logs
```bash
INF: 2022-09-10 15:55:22.11881 +0000 UTC |      127.0.0.1:50053 |   GET |              /health | 200 |  563.583µs | {"status":200,"message":"OK","data":null}

^CINF: 2022-09-10 16:01:37.83321 +0000 UTC | Closing Server
```