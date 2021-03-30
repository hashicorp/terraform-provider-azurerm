// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
)

// Policy represents an extensibility point for the Pipeline that can mutate the specified
// Request and react to the received Response.
type Policy interface {
	// Do applies the policy to the specified Request.  When implementing a Policy, mutate the
	// request before calling req.Next() to move on to the next policy, and respond to the result
	// before returning to the caller.
	Do(req *Request) (*Response, error)
}

// PolicyFunc is a type that implements the Policy interface.
// Use this type when implementing a stateless policy as a first-class function.
type PolicyFunc func(*Request) (*Response, error)

// Do implements the Policy interface on PolicyFunc.
func (pf PolicyFunc) Do(req *Request) (*Response, error) {
	return pf(req)
}

// Transport represents an HTTP pipeline transport used to send HTTP requests and receive responses.
type Transport interface {
	// Do sends the HTTP request and returns the HTTP response or error.
	Do(req *http.Request) (*http.Response, error)
}

// TransportFunc is a type that implements the Transport interface.
// Use this type when implementing a stateless transport as a first-class function.
type TransportFunc func(*http.Request) (*http.Response, error)

// Do implements the Transport interface on TransportFunc.
func (tf TransportFunc) Do(req *http.Request) (*http.Response, error) {
	return tf(req)
}

// used to adapt a TransportPolicy to a Policy
type transportPolicy struct {
	trans Transport
}

func (tp transportPolicy) Do(req *Request) (*Response, error) {
	resp, err := tp.trans.Do(req.Request)
	if err != nil {
		return nil, err
	} else if resp == nil {
		// there was no response and no error (rare but can happen)
		// this ensures the retry policy will retry the request
		return nil, errors.New("received nil response")
	}
	return &Response{Response: resp}, nil
}

// Pipeline represents a primitive for sending HTTP requests and receiving responses.
// Its behavior can be extended by specifying policies during construction.
type Pipeline struct {
	policies []Policy
}

// NewPipeline creates a new Pipeline object from the specified Transport and Policies.
// If no transport is provided then the default *http.Client transport will be used.
func NewPipeline(transport Transport, policies ...Policy) Pipeline {
	if transport == nil {
		transport = defaultHTTPClient
	}
	// transport policy must always be the last in the slice
	policies = append(policies, PolicyFunc(httpHeaderPolicy), PolicyFunc(bodyDownloadPolicy), transportPolicy{trans: transport})
	return Pipeline{
		policies: policies,
	}
}

// Do is called for each and every HTTP request. It passes the request through all
// the Policy objects (which can transform the Request's URL/query parameters/headers)
// and ultimately sends the transformed HTTP request over the network.
func (p Pipeline) Do(req *Request) (*Response, error) {
	if err := req.valid(); err != nil {
		return nil, err
	}
	req.policies = p.policies
	return req.Next()
}

// ReadSeekCloser is the interface that groups the io.ReadCloser and io.Seeker interfaces.
type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) ReadSeekCloser {
	return nopCloser{rs}
}

// Poller provides operations for checking the state of a long-running operation.
// An LRO can be in either a non-terminal or terminal state.  A non-terminal state
// indicates the LRO is still in progress.  A terminal state indicates the LRO has
// completed successfully, failed, or was cancelled.
type Poller interface {
	// Done returns true if the LRO has reached a terminal state.
	Done() bool

	// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is update and the HTTP
	// response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is
	// updated and the error is returned.
	// If the LRO has not reached a terminal state, the poller's state is updated and
	// the latest HTTP response is returned.
	// If Poll fails, the poller's state is unmodified and the error is returned.
	// Calling Poll on an LRO that has reached a terminal state will return the final
	// HTTP response or error.
	Poll(context.Context) (*http.Response, error)

	// ResumeToken returns a value representing the poller that can be used to resume
	// the LRO at a later time. ResumeTokens are unique per service operation.
	ResumeToken() (string, error)
}

// Pager provides operations for iterating over paged responses.
type Pager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// Err returns the last error encountered while paging.
	Err() error
}

// holds sentinel values used to send nulls
var nullables map[reflect.Type]interface{} = map[reflect.Type]interface{}{}

// NullValue is used to send an explicit 'null' within a request.
// This is typically used in JSON-MERGE-PATCH operations to delete a value.
func NullValue(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		// t is not of pointer type, make it be of pointer type
		t = reflect.PtrTo(t)
	}
	v, found := nullables[t]
	if !found {
		o := reflect.New(t.Elem())
		v = o.Interface()
		nullables[t] = v
	}
	// return the sentinel object
	return v
}

// IsNullValue returns true if the field contains a null sentinel value.
// This is used by custom marshallers to properly encode a null value.
func IsNullValue(v interface{}) bool {
	// see if our map has a sentinel object for this *T
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		// v isn't a pointer type so it can never be a null
		return false
	}
	if o, found := nullables[t]; found {
		// we found it; return true if v points to the sentinel object
		return o == v
	}
	// no sentinel object for this *t
	return false
}
