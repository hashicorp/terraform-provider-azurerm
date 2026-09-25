# plural - Simple Go API for Pluralisation.

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/rickb777/plural)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/plural)](https://goreportcard.com/report/github.com/rickb777/plural)
[![Build](https://github.com/rickb777/plural/actions/workflows/go.yml/badge.svg)](https://github.com/rickb777/plural/actions)
[![Coverage](https://coveralls.io/repos/github/rickb777/plural/badge.svg?branch=main)](https://coveralls.io/github/rickb777/plural?branch=main)
[![Issues](https://img.shields.io/github/issues/rickb777/plural.svg)](https://github.com/rickb777/plural/issues)

Package plural provides simple support for localising plurals in a flexible range of different styles.

There are considerable differences around the world in the way plurals are handled. This is a simple
but useful API for catering with these differences when presenting to people formatted text with numbers.

This package is able to format **countable things** and **continuous values**. It can handle integers
and floating point numbers equally and this allows you to decide to what extent each is appropriate.

For example, `2 cars` might weigh `1.6 tonnes`; both categories are covered.

## Installation

```
    go get github.com/rickb777/plural
```

## Status

This library has been in reliable production use for some time. Versioning follows the well-known semantic version pattern.

## More info

This API is deliberately simple; it doesn't address the full gamut of internationalisation. If that's
what you need, you should consider other products instead, such as

 * [nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n)
 * [Seeingu/icu4xgo](https://github.com/Seeingu/icu4xgo) - from Unicode.org; requires Rust & Cgo
 * [x/text/feature/plural](https://pkg.go.dev/golang.org/x/text/feature/plural)
