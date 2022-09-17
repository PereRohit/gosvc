# gosvc
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/PereRohit/gosvc) ![GitHub issues](https://img.shields.io/github/issues/PereRohit/gosvc?color=deep-green&label=Issues) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/PereRohit/gosvc?label=Release)

Cli to generate skeleton HTTP server in Go with basic Health endpoints.

The general [Go folder structure](https://github.com/golang-standards/project-layout) is followed.

The HTTP server will be generated with [gorilla mux](https://github.com/gorilla/mux).

Gos [text/template](https://pkg.go.dev/text/template) is used to generate basic endpoints.

## Installation
#### Prerequisites
- [go](https://go.dev/dl/)
- [mockgen](https://github.com/golang/mock/tree/main/mockgen)
```go
go install github.com/PereRohit/gosvc/cmd/gosvc@latest
```

## Usage
```bash
gosvc --init <module-name>
```

### Example: 
```bash
gosvc --init github.com/PereRohit/test-service
```
#### Output: will also create the folder `test-server`
```bash
tree test-service

test-service
├── README.md
├── build
│   └── ci
│       └── go.test.sh
├── cmd
│   └── test-service
│       └── main.go
├── codecov.yml
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
│   │   ├── datasource.go
│   │   └── request.go
│   ├── repo
│   │   └── datasource
│   │       ├── dummy.go
│   │       └── interface.go
│   └── router
│       └── router.go
└── pkg
    └── mock
        ├── mock_common.go
        ├── mock_datasource.go
        ├── mock_handler.go
        └── mock_logic.go

15 directories, 20 files
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
go run cmd/test-service/main.go
```
#### Output:
```bash
[2022-09-17 15:35:08] INF | service:testService | version:0.0.1 | Starting server(:80)
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
< Date: Sat, 17 Sep 2022 15:36:01 GMT
< Content-Length: 69
< 
{"status":200,"message":"OK","data":{"testService":{"status":"OK"}}}
* Connection #0 to host localhost left intact
```
#### Server Logs
```bash
[2022-09-17 15:35:08] INF | service:testService | version:0.0.1 | Starting server(:80)
[2022-09-17 15:36:01] DBG | service:testService | version:0.0.1 | caller:dummy.go:30 | this is a custom config
[2022-09-17 15:36:01] INF | service:testService | version:0.0.1 |      127.0.0.1:52106 | GET    | /v1/health                | 200 |  709.792µs | request_id:e27fde7c-ded8-40bb-aa76-b7e66b29273f | {"status":200,"message":"OK","data":{"testService":{"status":"OK"}}}
```

### Ping (cURL)
This is a demo endpoint to showcase the usage of registration of a new endpoint and how the development is meant to be  across the handler, logic & repo layer, following Go folder structure & CLEAN code principles.
```bash
curl -X POST http://localhost/v1/ping -H 'Content-Type: application/json' -d '{ "data": "hello world from test" }'
```
#### Server Logs
```bash
[2022-09-17 15:35:08] INF | service:testService | version:0.0.1 | Starting server(:80)
[2022-09-17 15:36:01] DBG | service:testService | version:0.0.1 | caller:dummy.go:30 | this is a custom config
[2022-09-17 15:36:01] INF | service:testService | version:0.0.1 |      127.0.0.1:52106 | GET    | /v1/health                | 200 |  709.792µs | request_id:e27fde7c-ded8-40bb-aa76-b7e66b29273f | {"status":200,"message":"OK","data":{"testService":{"status":"OK"}}}

[2022-09-17 15:37:22] INF | service:testService | version:0.0.1 | caller:dummy.go:22 | this is a custom config
[2022-09-17 15:37:22] INF | service:testService | version:0.0.1 |      127.0.0.1:52113 | POST   | /v1/ping                  | 200 |      323µs | request_id:5fba5e78-c335-4558-8ea0-af79bac587cf | {"status":200,"message":"Pong","data":{"Data":"hello world from test"}}
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
`Command` + `C` on MacOS, `Ctrl` + `c` on Windows or linux gracefully shuts down the server with the following logs
```bash
[2022-09-17 15:35:08] INF | service:testService | version:0.0.1 | Starting server(:80)
[2022-09-17 15:36:01] DBG | service:testService | version:0.0.1 | caller:dummy.go:30 | this is a custom config
[2022-09-17 15:36:01] INF | service:testService | version:0.0.1 |      127.0.0.1:52106 | GET    | /v1/health                | 200 |  709.792µs | request_id:e27fde7c-ded8-40bb-aa76-b7e66b29273f | {"status":200,"message":"OK","data":{"testService":{"status":"OK"}}}

[2022-09-17 15:37:22] INF | service:testService | version:0.0.1 | caller:dummy.go:22 | this is a custom config
[2022-09-17 15:37:22] INF | service:testService | version:0.0.1 |      127.0.0.1:52113 | POST   | /v1/ping                  | 200 |      323µs | request_id:5fba5e78-c335-4558-8ea0-af79bac587cf | {"status":200,"message":"Pong","data":{"Data":"hello world from test"}}

^C[2022-09-17 15:38:29] INF | service:testService | version:0.0.1 | Closing Server
[2022-09-17 15:38:29] INF | service:testService | version:0.0.1 | Server Closed!!
```

## Development
Once the skeleton project is generated for the given module following [Usage](##usage) section.

- `/internal/handler/handler.go` is meant to contain the raw HTTP handlers mainly responsible for:
  1. validating/sanitsing incoming requests for the endpoint(validator struct tags)
  2. call the business logic
  3. forward the response from logic to client.

    An interface will be dynamically generated for mocking in test cases.
    ```go
    type TestServiceHandler interface {
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
    type TestServiceLogicIer interface {
        Ping(*model.PingRequest) *respModel.Response
        HealthCheck() bool
        // Add API logic functions for internal calls from the handlers
    }
    ```
    Please use the `Ping` method as reference to add your custom logic for new functions.
    ```go
    func (l testServiceLogic) Ping(req *model.PingRequest) *respModel.Response {
        // add business logic here
        res, err := l.dummyDsSvc.Ping(&model.PingDs{
            Data: req.Data,
        })
        if err != nil {
            log.Error("datasource error", err)
            return &respModel.Response{
                Status:  http.StatusInternalServerError,
                Message: "",
                Data:    nil,
            }
        }
        return &respModel.Response{
            Status:  http.StatusOK,
            Message: "Pong",
            Data:    res,
        }
    }
    ```

- `internal/model/` is meant to contain the request structs. This will be used to parse the incoming request to it and forward the data to the logic or repo level, providing only the required data to be sent and separate control/logic/data layers.

  [go validator](https://github.com/go-playground/validator) struct tags can be used for out of the box struct validation.

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
            "name": "testService",  // service name taken during gosvc init
            "log_level": "info"    // log level to limit log type
        },
        "custom_svc": {           // custom config -- use your own to pass config
            "data": "this is a custom config"
        }
    }
    ```

- `internal/repo` provides implementation calls to dependant services. This can be initialised at router level and passed across to the logic level via the handler. Hence, it provides abstraction over infra level like datasources, etc. The idea is that, get underlying concrete implementation of the services from the router and pass it to the logic that requires it.
  
  *Example:* If the datasource wants to be changed from MySQl to MongoDB, implement the required `datasource` interface which will be the same for both. From router, call a init function like `NewMongoDs` and returns an implementation of the interface.

  A dummy implementation of datasource is provided.
  Similarly, for other services, create a directory under `repo` have an interface and corresponding file with concrete implementation.

This is part of the [github.com/PereRohit/util](https://github.com/PereRohit/util) library whose functions can be used.

*The demo service can be found here: [test-service](https://github.com/PereRohit/test-service)*