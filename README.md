# gosvc
Cli to generate skeleton HTTP server in Go with basic Health endpoints.

The general [Go folder structure](https://github.com/golang-standards/project-layout) is followed.

The HTTP server will be generated with [gorilla mux](https://github.com/gorilla/mux).

Gos [text/template](https://pkg.go.dev/text/template) is used to generate basic endpoints.

## Installation
#### Prerequisites
- [mockgen](https://github.com/golang/mock/tree/main/mockgen)
- [go](https://go.dev/dl/)
```go
go install github.com/PereRohit/gosvc/cmd/gosvc@latest
```

## Issues
- Feel free to raise an issue and contribute
- Check existing issues [here](https://github.com/PereRohit/gosvc/issues)

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
├── configs
│   └── config.json
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── common.go
│   │   └── handler.go
│   ├── logic
│   │   └── logic.go
│   ├── model
│   │   └── model.go
│   └── router
│       └── router.go
└── pkg
    └── mock
        ├── mock_common.go
        ├── mock_handler.go
        └── mock_logic.go

11 directories, 14 files
```

### Help: `init` & `version` options are available at the moment
```bash
gosvc --help
```
#### Output
```bash
Usage of gosvc:
  -init string
        go module name
  -version
        version
```

## Server Usage
### Start server on `localhost` with default port `80`
```bash
go run cmd/test-server/main.go
```
#### Output:
```bash
INF: 2022-09-15 12:08:38.48473 +0000 UTC | name: testServer version: 0.0.1 | Starting server(:80)
```
### Health Check (cURL)
The baked in Health Check endpoint checks if the service is internally okay.

**NOTE:** The health check logic needs to be custom developed as every service has their own meaning for `"OK"`. This provides a template that may be used for the custom check, only for the newly introduced service.
```bash
curl -v http://localhost/v1/health
```
#### Output
```bash
*   Trying 127.0.0.1:80...
* Connected to localhost (127.0.0.1) port 80 (#0)
> GET /v1/health HTTP/1.1
> Host: localhost
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 15 Sep 2022 12:10:18 GMT
< Content-Length: 68
< 
{"status":200,"message":"OK","data":{"testServer":{"status":"OK"}}}
* Connection #0 to host localhost left intact
```
#### Server Logs
```bash
INF: 2022-09-15 12:08:38.48473 +0000 UTC | name: testServer version: 0.0.1 | Starting server(:80)
INF: 2022-09-15 12:10:18.315826 +0000 UTC | name: testServer version: 0.0.1 |      127.0.0.1:60538 |   GET |           /v1/health | 200 |   33.084µs | {"status":200,"message":"OK","data":{"testServer":{"status":"OK"}}}
```

### Ping (cURL)
This is a demo endpoint to showcase the usage of registration of a new endpoint and how the development is meant to be  across the handler and logic, following Go folder structure.
```bash
curl -X POST http://localhost/v1/ping -H 'Content-Type: application/json' -d '{ "data": "hello world from test" }'
```
#### Server Logs
```bash
INF: 2022-09-15 12:08:38.48473 +0000 UTC | name: testServer version: 0.0.1 | Starting server(:80)
INF: 2022-09-15 12:10:18.315826 +0000 UTC | name: testServer version: 0.0.1 |      127.0.0.1:60538 |   GET |           /v1/health | 200 |   33.084µs | {"status":200,"message":"OK","data":{"testServer":{"status":"OK"}}}

INF: 2022-09-15 12:11:37.195134 +0000 UTC | name: testServer version: 0.0.1 | this is a custom config
INF: 2022-09-15 12:11:37.195182 +0000 UTC | name: testServer version: 0.0.1 |      127.0.0.1:60541 |  POST |             /v1/ping | 200 |  362.667µs | {"status":200,"message":"Pong","data":"hello world from test"}
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
INF: 2022-09-15 12:08:38.48473 +0000 UTC | name: testServer version: 0.0.1 | Starting server(:80)
INF: 2022-09-15 12:10:18.315826 +0000 UTC | name: testServer version: 0.0.1 |      127.0.0.1:60538 |   GET |           /v1/health | 200 |   33.084µs | {"status":200,"message":"OK","data":{"testServer":{"status":"OK"}}}

INF: 2022-09-15 12:11:37.195134 +0000 UTC | name: testServer version: 0.0.1 | this is a custom config
INF: 2022-09-15 12:11:37.195182 +0000 UTC | name: testServer version: 0.0.1 |      127.0.0.1:60541 |  POST |             /v1/ping | 200 |  362.667µs | {"status":200,"message":"Pong","data":"hello world from test"}

^CINF: 2022-09-15 12:12:23.908422 +0000 UTC | name: testServer version: 0.0.1 | Closing Server
INF: 2022-09-15 12:12:23.909244 +0000 UTC | name: testServer version: 0.0.1 | Server Closed!!
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
        log.Info(l.dummySvc.Data)
        return &respModel.Response{
            Status:  http.StatusOK,
            Message: "Pong",
            Data:    req.Data,
        }
    }
    ```

- `internal/model/model.go` is meant to contain the request structs. This will be used to parse the incoming request to it and forward the data to the logic level, providing only the required data to be sent and separate control/logic.

- `configs/config.json` houses the config for internal services and server. Use this to pass configurations and initialise them in the code.

    For the added config in the `config.json`, add a field in the `Config` struct with appropriate json struct tag to populate it and use in code.
    > `Config` struct
    ```go
    type Config struct {
        ServiceRouteVersion string              `json:"service_route_version"`
        ServerConfig        config.ServerConfig `json:"server_config"`
        // add custom config structs below for any internal services
        DummyCfg DummySvcCfg `json:"custom_svc"`
    }
    ```
    > `config.json` file
    ```json
    {
        "service_route_version": "v1", // set sub route if not empty like /v1/health
        "server_config": {
            "host": "",
            "port": "80",     // port on which server is to be started, default is 80
            "version": "0.0.1",    // service version
            "name": "testServer",  // service name taken during gosvc init
            "log_level": "info"    // log level to limit log type
        },
        "custom_svc": {           // custom config -- use your own to pass config
            "data": "this is a custom config"
        }
    }
    ```

[go validator](https://github.com/go-playground/validator) struct tags can be used for out of the box struct validation.
This is part of the [github.com/PereRohit/util](https://github.com/PereRohit/util) library whose functions can be used.
