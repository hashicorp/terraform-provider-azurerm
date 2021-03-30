// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

const okPage = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>Login Succeeded</title>
</head>
<body>
    <h4>You have logged into Microsoft Azure!</h4>
    <p>You can now close this window.</p>
</body>
</html>
`

const failPage = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>Login Failed</title>
</head>
<body>
    <h4>An error occurred during authentication</h4>
    <p>Please open an issue in the <a href="https://github.com/azure/go-autorest/issues">Go Autorest repo</a> for assistance.</p>
</body>
</html>
`

type server struct {
	wg   *sync.WaitGroup
	s    *http.Server
	code string
	err  error
}

// NewServer creates an object that satisfies the Server interface.
func newServer() *server {
	rs := &server{
		wg: &sync.WaitGroup{},
		s:  &http.Server{},
	}
	return rs
}

// Start starts the local HTTP server on a separate go routine.
// The return value is the full URL plus port number.
func (s *server) Start(reqState string, port int) string {
	if port == 0 {
		port = rand.Intn(600) + 8400
	}
	s.s.Addr = fmt.Sprintf(":%d", port)
	s.s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer s.wg.Done()
		qp := r.URL.Query()
		if respState, ok := qp["state"]; !ok {
			s.err = errors.New("missing OAuth state")
			return
		} else if respState[0] != reqState {
			s.err = errors.New("mismatched OAuth state")
			return
		}
		if err, ok := qp["error"]; ok {
			w.Write([]byte(failPage))
			s.err = fmt.Errorf("authentication error: %s; description: %s", err[0], qp.Get("error_description"))
			return
		}
		if code, ok := qp["code"]; ok {
			w.Write([]byte(okPage))
			s.code = code[0]
		} else {
			s.err = errors.New("authorization code missing in query string")
		}
	})
	s.wg.Add(1)
	go s.s.ListenAndServe()
	return fmt.Sprintf("http://localhost:%d", port)
}

// Stop will shut down the local HTTP server.
func (s *server) Stop() {
	s.s.Shutdown(context.Background())
}

// WaitForCallback will wait until Azure interactive login has called us back with an authorization code or error.
func (s *server) WaitForCallback() {
	s.wg.Wait()
}

// AuthorizationCode returns the authorization code or error result from the interactive login.
func (s *server) AuthorizationCode() (string, error) {
	return s.code, s.err
}
