# go-ipify

The official client library for [ipify][]: *A Simple IP Address API*.

[![GitHub License](https://img.shields.io/badge/license-Unlicense-blue.svg)](https://raw.githubusercontent.com/rdegges/go-ipify/master/UNLICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/rdegges/go-ipify.svg)](https://github.com/rdegges/go-ipify/stargazers)
[![GoDoc](https://godoc.org/github.com/rdegges/go-ipify?status.svg)](https://godoc.org/github.com/rdegges/go-ipify)
[![Build Status](https://travis-ci.org/rdegges/go-ipify.svg?branch=master)](https://travis-ci.org/rdegges/go-ipify)


## Meta

- Author: Randall Degges
- Email: r@rdegges.com
- Site: http://www.rdegges.com
- Status: maintained, active


## Purpose

[ipify][] is the best IP address lookup service on the internet.  It's fast,
simple, scalable, open source, and well-funded (*by me!*).

In short: if you need a way to pragmatically get your public IP address, ipify
is the best possible choice!

This library will retrieve your public IP address from ipify's API service, and
return it as a string.  It can't get any simpler than that.

This library also has some other nice features you might care about:

- If a request fails for any reason, it is re-attempted 3 times using an
  exponential backoff algorithm for maximum effectiveness.
- This library handles errors properly, and usage examples below show you
  how to deal with errors in a foolproof way.
- This library only makes API requests over HTTPS.


## Installation

To install `ipify`, simply run:

```console
$ go get github.com/rdegges/go-ipify
```

This will install the latest version of the library automatically.


## Usage

Using this library is very simple.  Here's a simple example:

```go
package main

import (
    "fmt"
    "github.com/rdegges/go-ipify"
)

func main() {
    ip, err := ipify.GetIp()
    if err != nil {
        fmt.Println("Couldn't get my IP address:", err)
    } else {
        fmt.Println("My IP address is:", ip)
    }
}
```

Now, in regards to error handling, there are several ways this can fail:

- The ipify service is down (*not likely*), or:
- Your machine is unable to get the request to ipify because of a network error
  of some sort (*DNS, no internet, etc.*).

The library will output an informative error message in the event anything
fails.

One thing to keep in mind: regardless of how you decide to handle errors, the
ipify library will retry any failed requests 3 times before ever failing -- so
if you *do* need to handle errors, just remember that retry logic has already
been attempted.


## Contributing

This project is only possible due to the amazing contributors who work on it!

If you'd like to improve this library, please send me a pull request!  I'm
happy to review and merge pull requests.

The standard contribution workflow should look something like this:

- Fork this project on Github.
- Make some changes in the master branch (*this project is simple, so no need to
  complicate things*).
- Send a pull request when ready.

Also, if you're making changes, please write tests for your changes -- this
project has a full test suite you can easily modify / test.

To run the test suite, you can use the `go test` command after forking this
repository.


  [ipify]: http://www.ipify.org/ "ipify - A Simple IP Address API"
