# gosvc
Cli to generate skeleton HTTP server in Go with basic Health endpoints.

The HTTP server will be generated with [gorilla mux](https://github.com/gorilla/mux).

Gos [text/template](https://pkg.go.dev/text/template) is used to generate basic endpoints.

## Installation
#### Prerequisites
- [mockgen](https://github.com/golang/mock/tree/main/mockgen)
- [go](https://go.dev/dl/)
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
├── README.md
├── cmd
│   └── test-server
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── handler
│   │   ├── common.go
│   │   └── handler.go
│   ├── logic
│   │   ├── logic.go
│   │   └── mockgen.go
│   ├── model
│   │   └── model.go
│   └── router
│       └── router.go
└── pkg
    └── mock
        ├── mock_common.go
        ├── mock_handler.go
        └── mock_logic.go

9 directories, 13 files
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
INF: 2022-09-11 11:04:13.010931 +0000 UTC | Starting server(:80)
```
### Health Check (cURL)
The baked in Health Check endpoint will internally calls checks if the service is internally okay.

**NOTE:** The health check logic needs to be custom developed as every service has their own meaning for `"OK"`. This provides a template that may be used for the custom check, only for the newly introduced service.
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
< Date: Sun, 11 Sep 2022 11:05:03 GMT
< Content-Length: 68
< 
{"status":200,"message":"OK","data":{"TestServer":{"status":"OK"}}}
* Connection #0 to host localhost left intact
```
#### Server Logs
```bash
INF: 2022-09-11 11:04:13.010931 +0000 UTC | Starting server(:80)
INF: 2022-09-11 11:05:03.383409 +0000 UTC |      127.0.0.1:50002 |   GET |              /health | 200 |  377.917µs | {"status":200,"message":"OK","data":{"TestServer":{"status":"OK"}}}
```

### Ping (cURL)
This is a demo endpoint to showcase the usage of registration of a new endpoint and how the development is meant to be  across the handler and logic, following Go folder structure.
```bash
curl -X POST http://localhost/ping -H 'Content-Type: application/json' -d '{ "data": "hello world from test" }'
```
#### Server Logs
```bash
INF: 2022-09-11 11:04:13.010931 +0000 UTC | Starting server(:80)
INF: 2022-09-11 11:05:03.383409 +0000 UTC |      127.0.0.1:50002 |   GET |              /health | 200 |  377.917µs | {"status":200,"message":"OK","data":{"TestServer":{"status":"OK"}}}

INF: 2022-09-11 11:12:25.013227 +0000 UTC |      127.0.0.1:50022 |  POST |                /ping | 200 |  980.375µs | {"status":200,"message":"Pong","data":"hello world from test"}
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
INF: 2022-09-11 11:04:13.010931 +0000 UTC | Starting server(:80)
INF: 2022-09-11 11:05:03.383409 +0000 UTC |      127.0.0.1:50002 |   GET |              /health | 200 |  377.917µs | {"status":200,"message":"OK","data":{"TestServer":{"status":"OK"}}}

INF: 2022-09-11 11:12:25.013227 +0000 UTC |      127.0.0.1:50022 |  POST |                /ping | 200 |  980.375µs | {"status":200,"message":"Pong","data":"hello world from test"}

^CINF: 2022-09-11 11:13:36.551407 +0000 UTC | Closing Server
INF: 2022-09-11 11:13:36.554467 +0000 UTC | Server Closed!!
```

## Development
Once the skeleton project is generated for the given module following [Usage](##usage) section.

- `/internal/handler/handler.go` is meant to contain the raw HTTP handlers mainly responsible for validating/sanitsing incoming requests for the endpoint --> call the business logic and forward the response back.

    An interface will be dynamically generated for mocking in test cases.
    ```go
    type TestServerHandler interface {
	    HealthChecker
	    Ping(w http.ResponseWriter, r *http.Request)
        // Add handler functions here that directly attaches to the endpoint
    }
    ```
    Add in the new handlers to the interface followed by their implementation for the service `struct`.

    [mockgen](https://github.com/golang/mock/tree/main/mockgen) is used to generate the mocks for testing using [gomock](https://github.com/golang/mock).

    The mocks need to be regenerated when the `interface` is updated.
    Please execute `go generate ./...` from project root for project wide mock regeneration.


- `/internal/logic/logic.go` is meant to contain the actual underlying API logic for the handler/endpoint.
    ```go
    type TestServerLogicIer interface {
        Ping(*model.PingRequest) *respModel.Response
        // Add API logic functions for internal calls from the handlers
    }
    ```
    Please use the `Ping` method as reference to add your custom logic for new functions.
    ```go
    func (l TestServerLogic) Ping(req *model.PingRequest) *respModel.Response {
        // add your business logic here
        return &respModel.Response{
            Status:  http.StatusOK,
            Message: "Pong",
            Data:    req.Data,
        }
    }
    ```

- `internal/model/model.go` is meant to contain the request structs. This will be used to parse the incoming request to it and forward the data to the logic level, providing only the required data to be sent and separate control/logic.

[go validator](https://github.com/go-playground/validator) struct tags can be used for out of the box struct validation.
This is part of the [github.com/PereRohit/util](https://github.com/PereRohit/util) library whose functions can be used.