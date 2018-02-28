[![Build Status](https://travis-ci.org/stoewer/go-strcase.svg?branch=master)](https://travis-ci.org/stoewer/go-strcase)
[![Coverage Status](https://coveralls.io/repos/github/stoewer/go-strcase/badge.svg?branch=master)](https://coveralls.io/github/stoewer/go-strcase?branch=master)
[![GoDoc](https://godoc.org/github.com/stoewer/go-strcase?status.svg)](https://godoc.org/github.com/stoewer/go-strcase)
---

# Go strcase

The package `strcase` converts between different kinds of naming formats such as camel case 
(`CamelCase`), snake case (`snake_case`) or kebab case (`kebab-case`).
The package is designed to work only with strings consisting of standard ASCII letters. 
Unicode is currently not supported.

## Versioning and stability

Although the master branch is supposed to remain always backward compatible, the repository
contains version tags in order to support vendoring tools such as `glide`.
The tag names follow semantic versioning conventions and have the following format `v1.0.0`.


## Install and use

```sh
go get -u github.com/stoewer/go-strcase
```

```go
import "github.com/stoewer/go-strcase"

var snake = strcase.SnakeCase("CamelCase")
```
