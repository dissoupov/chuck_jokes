# chuck_jokes

A demo web service

## Requirements

1. GoLang 1.15+

## Build

* `make all` initializes all dependencies, builds and tests.
* `make build` build the project
* `make test` run the tests
* `make testshort` runs the tests skipping the end-to-end tests and the code coverage reporting
* `make covtest` runs the tests with end-to-end and the code coverage reporting
* `make coverage` view the code coverage results from the last make test run.
* `make generate` runs go generate to update any code gen'd files (query_console.go in our case)
* `make fmt` runs go fmt on the project.
* `make lint` runs the go linter on the project.

run `make all` once, then run `make build` or `make test` as needed.

First run:

    make all

Subsequent builds:

    make build

Tests:

    make test

Optionally run golang race detector with test targets by setting RACE flag:

    make test RACE=true

Tests coverage:

    make covtest

Review coverage report:

    make coverage

Docker:

    make docker

## Debug with delve

```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "remotePath": "",
            "port": 2345,
            "host": "127.0.0.1",
            "program": "${workspaceRoot}/cmd/jokes",
            "env": {},
            "args": [
                "--std",
                "--cfg",
                "${workspaceRoot}/etc/dev/jokes-config.json"
            ],
            "showLog": true
        },
    ]
}
```

## Run locally

```
bin/jokes --std
```

## Run in docker

```
make docker
docker run -ti -p 5000:5000 -p 8080:8080 jokes:latest
```

## Service status

    curl http://localhost:8080/v1/status

## Configuration

```
etc/dev/jokes-config.json
```

## Security

1. To enable HTTPS, generate the certs and update `ServerTLS` configuration.
1. To enable mTLS, set `ClientCertAuth` config value and specify `CertMapper` configuration.
1. To enable APIKey-based roles, specify `APIKeyMapper` configuration.

## Logs

The logs are configured in `Logger.Directory`, by default `/tmp/jokes/logs/jokes.log`

## Metrics

At the moment is not configured and used with `in-memory` provider.

## Project structure

```
.
├── api
│   └── v1                    - Public API definition
├── cmd
│   └── jokes                 - jokes service main
├── etc                       - config folder
│   └── dev
├── internal
│   ├── config                - configuration schema and loader
│   └── version               - application version 
├── pkg
│   ├── printer               - helper to print responses
│   └── roles                 - Security roles
│       ├── apikeymapper
│       │   └── testdata
│       └── certmapper
│           └── testdata
└── service
    ├── jokes                 - jokes service
    └── status                - status service
```