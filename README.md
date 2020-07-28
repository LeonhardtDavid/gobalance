# gobalance

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://david-leonhardt.mit-license.org/)
[![Tests status](https://github.com/LeonhardtDavid/gobalance/workflows/Tests/badge.svg)](https://github.com/LeonhardtDavid/gobalance/actions?query=workflow%3ATests)

Simple load balancer. Exercise took from  https://github.com/thewildcast/golang-primeros-pasos

## Running the project

### Run in development mode

```sh
go run main.go
```

## Run tests

```sh
go test ./...
```

### Run tests with coverage and coverage report

```sh
go test -coverprofile=c.out ./...
go tool cover -html=c.out -o coverage.html
```

Or just use the project's script

```sh
./test-coverage.sh
```

## Configurations

Check [this configuration template](/config.template.yml) to see configurations examples.

If a `config.yml` file exists in the same folder as the project, or the executable binary, or
if it exists in `/etc/gobalance`, then that file is used. Otherwise the applications will try
to create a configuration file from `config.template.yml` in the root folder.
