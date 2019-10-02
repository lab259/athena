# athena [![CircleCI](https://circleci.com/gh/lab259/athena/tree/master.svg?style=shield)](https://circleci.com/gh/lab259/athena/tree/master) [![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/lab259/athena) [![license](https://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/lab259/athena/master/LICENSE) [![Coverage](https://gocover.io/_badge/github.com/lab259/athena)](http://gocover.io/github.com/lab259/athena)

Wisely building web applications.

## Getting Started

### Graceful Applications

With `athena` you can easily run `net/http` applications that will gracefully be stopped once receive a `INTERRUPT` signal (<kbd>Ctrl</kbd>+<kbd>C</kbd>). For example:

```go
athena.GracefulHTTP(&http.Server{
    Addr:    ":3000",       // bind address
    Handler: app.Handler(), // your http handler
}, serviceStarter)          // from github.com/lab259/go-rscsrv
```

If you prefer `fasthttp`:

```go
athena.GracefulFastHTTP(
    app.NewServer(),    // your fasthttp server
    ":3000",            // bind address
    serviceStarter,     // from github.com/lab259/go-rscsrv
)
```

### Generate commands

`athena` ships with a CLI to generate models (for PostgreSQL and MongoDB) and services. Execute `athena --help` for more information:

```
Usage: athena [OPTIONS] COMMAND [arg...]

Wisely building web applications

Options:
  -v, --version   Show the version and exit

Commands:
  make:service    Generate a service
  make:model      Generate a model
  make:mgomodel   Generate a mgo model
```

### Additional packages

#### `athena/config`

A configuration loader implementation from YAML files.

#### `athena/pagination`

Parse pagination values (current page and page size) with sane defaults.

#### `athena/rscsrv`

Built-in `go-rscsrv`'s Service implementation for common services:

- [PostgreSQL](https://github.com/lab259/go-rscsrv-psql)
- [Redis](https://github.com/lab259/go-rscsrv-redigo)
- [MongoDB](https://github.com/lab259/go-rscsrv-mgo)

#### `athena/testing`

Utilities for testing applications:

- `envtest`: overwrite environment variables
- `ginkgtest`: initialize ginkgo with [macchiato](https://github.com/jamillosantos/macchiato) + reporters for CI
- `httptest`: wrappers for net/http, fasthttp, [hermes](https://github.com/lab259/hermes) using `httpexpect`
- `mgotest`: helpers for cleansing the database
- `psqltest`: overwrite the default PostgreSQL service implemenation to use transations
- `rscsrvtest`: initilize serveral services (with before and after hooks)

#### `athena/validator`

Validate structs and values using the [`go-playground/validator.v9`](https://gopkg.in/go-playground/validator.v9).

## Contributing

### Prerequisites

What things you need to setup the project:

- [go](https://golang.org/doc/install)
- [ginkgo](http://onsi.github.io/ginkgo/)

### Environment

Close the repository:

```bash
git clone git@github.com:lab259/athena.git
```

Now, the dependencies must be installed.

```
cd athena && go mod download
```

:wink: Finally, you are done to start developing.

### Running tests

In the root directory, execute:

```bash
make test
```

To enable coverage, execute:

```bash
make coverage
```

To generate the HTML coverage report, execute:

```bash
make coverage coverage-html
```
