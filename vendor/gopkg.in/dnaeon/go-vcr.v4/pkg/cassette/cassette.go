// Copyright (c) 2015-2024 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer
//    in this position and unchanged.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package cassette

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.yaml.in/yaml/v4"
)

const (
	// CassetteFormatVersion is the supported cassette version.
	CassetteFormatVersion = 2
)

var (
	// ErrInteractionNotFound indicates that a requested interaction was not
	// found in the cassette file.
	ErrInteractionNotFound = errors.New("requested interaction not found")

	// ErrCassetteNotFound indicates that a requested cassette doesn't exist.
	ErrCassetteNotFound = errors.New("requested cassette not found")

	// ErrUnsupportedCassetteFormat is returned when attempting to use an
	// older and potentially unsupported format of a cassette.
	ErrUnsupportedCassetteFormat = errors.New("unsupported cassette version format")
)

// Request represents a client request as recorded in the cassette file.
type Request struct {
	Proto            string      `yaml:"proto"`
	ProtoMajor       int         `yaml:"proto_major"`
	ProtoMinor       int         `yaml:"proto_minor"`
	ContentLength    int64       `yaml:"content_length"`
	TransferEncoding []string    `yaml:"transfer_encoding,omitempty"`
	Trailer          http.Header `yaml:"trailer,omitempty"`
	Host             string      `yaml:"host"`
	RemoteAddr       string      `yaml:"remote_addr,omitempty"`
	RequestURI       string      `yaml:"request_uri,omitempty"`

	// Body of request
	Body string `yaml:"body,omitempty"`

	// Form values
	Form url.Values `yaml:"form,omitempty"`

	// Request headers
	Headers http.Header `yaml:"headers,omitempty"`

	// Request URL
	URL string `yaml:"url"`

	// Request method
	Method string `yaml:"method"`
}

// Response represents a server response as recorded in the cassette file.
type Response struct {
	Proto            string      `yaml:"proto"`
	ProtoMajor       int         `yaml:"proto_major"`
	ProtoMinor       int         `yaml:"proto_minor"`
	TransferEncoding []string    `yaml:"transfer_encoding,omitempty"`
	Trailer          http.Header `yaml:"trailer,omitempty"`
	ContentLength    int64       `yaml:"content_length"`
	Uncompressed     bool        `yaml:"uncompressed,omitempty"`

	// Body of response
	Body string `yaml:"body"`

	// Response headers
	Headers http.Header `yaml:"headers"`

	// Response status message
	Status string `yaml:"status"`

	// Response status code
	Code int `yaml:"code"`

	// Response duration
	Duration time.Duration `yaml:"duration"`
}

// Interaction type contains a pair of request/response for a single HTTP
// interaction between a client and a server.
type Interaction struct {
	// ID is the id of the interaction
	ID int `yaml:"id"`

	// Request is the recorded request
	Request Request `yaml:"request"`

	// Response is the recorded response
	Response Response `yaml:"response"`

	// DiscardOnSave if set to true will discard the interaction as a whole
	// and it will not be part of the final interactions when saving the
	// cassette on disk.
	DiscardOnSave bool `yaml:"-"`

	// replayed is true when this interaction has been played already.
	replayed bool `yaml:"-"`
}

// WasReplayed returns a boolean indicating whether the given interaction was
// already replayed.
func (i *Interaction) WasReplayed() bool {
	return i.replayed
}

// GetHTTPRequest converts the recorded interaction request to http.Request
// instance.
func (i *Interaction) GetHTTPRequest() (*http.Request, error) {
	url, err := url.Parse(i.Request.URL)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Proto:            i.Request.Proto,
		ProtoMajor:       i.Request.ProtoMajor,
		ProtoMinor:       i.Request.ProtoMinor,
		ContentLength:    i.Request.ContentLength,
		TransferEncoding: i.Request.TransferEncoding,
		Trailer:          i.Request.Trailer,
		Host:             i.Request.Host,
		RemoteAddr:       i.Request.RemoteAddr,
		RequestURI:       i.Request.RequestURI,
		Body:             io.NopCloser(strings.NewReader(i.Request.Body)),
		Form:             i.Request.Form,
		Header:           i.Request.Headers,
		URL:              url,
		Method:           i.Request.Method,
	}

	return req, nil
}

// GetHTTPResponse converts the recorded interaction response to http.Response
// instance.
func (i *Interaction) GetHTTPResponse() (*http.Response, error) {
	req, err := i.GetHTTPRequest()
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		Status:           i.Response.Status,
		StatusCode:       i.Response.Code,
		Proto:            i.Response.Proto,
		ProtoMajor:       i.Response.ProtoMajor,
		ProtoMinor:       i.Response.ProtoMinor,
		TransferEncoding: i.Response.TransferEncoding,
		Trailer:          i.Response.Trailer,
		ContentLength:    i.Response.ContentLength,
		Uncompressed:     i.Response.Uncompressed,
		Body:             io.NopCloser(strings.NewReader(i.Response.Body)),
		Header:           i.Response.Headers,
		Close:            true,
		Request:          req,
	}

	return resp, nil
}

// MatcherFunc is a predicate, which returns true when the actual request
// matches an interaction from the cassette.
type MatcherFunc func(*http.Request, Request) bool

// MarshalFunc is a function which marshals an object to a byte slice.
type MarshalFunc func(any) ([]byte, error)

// defaultMatcher is the default matcher used to match HTTP requests with
// recorded interactions.
type defaultMatcher struct {
	// If set, the default matcher will ignore matching on any of the
	// defined headers.
	ignoreHeaders []string
}

// DefaultMatcherOption is a function which configures the default matcher.
type DefaultMatcherOption func(m *defaultMatcher)

// WithIgnoreUserAgent is a [DefaultMatcherOption], which configures the default
// matcher to ignore matching on the User-Agent HTTP header.
func WithIgnoreUserAgent() DefaultMatcherOption {
	opt := func(m *defaultMatcher) {
		m.ignoreHeaders = append(m.ignoreHeaders, "User-Agent")
	}

	return opt
}

// WithIgnoreAuthorization is a [DefaultMatcherOption], which configures the default
// matcher to ignore matching on the Authorization HTTP header.
func WithIgnoreAuthorization() DefaultMatcherOption {
	opt := func(m *defaultMatcher) {
		m.ignoreHeaders = append(m.ignoreHeaders, "Authorization")
	}

	return opt
}

// WithIgnoreHeaders is a [DefaultMatcherOption], which configures the default
// matcher to ignore matching on the defined HTTP headers.
func WithIgnoreHeaders(val ...string) DefaultMatcherOption {
	opt := func(m *defaultMatcher) {
		m.ignoreHeaders = append(m.ignoreHeaders, val...)
	}

	return opt
}

// NewDefaultMatcher returns the default matcher.
func NewDefaultMatcher(opts ...DefaultMatcherOption) MatcherFunc {
	m := &defaultMatcher{}
	for _, opt := range opts {
		opt(m)
	}

	return m.matcher
}

// Similar to reflect.DeepEqual, but considers the contents of collections, so
// {} and nil would be considered equal. works with Array, Map, Slice, or
// pointer to Array.
func (m *defaultMatcher) deepEqualContents(x, y any) bool {
	if reflect.ValueOf(x).IsNil() {
		if reflect.ValueOf(y).IsNil() {
			return true
		} else {
			return reflect.ValueOf(y).Len() == 0
		}
	} else {
		if reflect.ValueOf(y).IsNil() {
			return reflect.ValueOf(x).Len() == 0
		} else {
			return reflect.DeepEqual(x, y)
		}
	}
}

// bodyMatches is a predicate which tests whether the bodies of the given HTTP
// request and interaction request match.
func (m *defaultMatcher) bodyMatches(r *http.Request, i Request) bool {
	if r.Body != nil {
		var buffer bytes.Buffer
		if _, err := buffer.ReadFrom(r.Body); err != nil {
			return false
		}

		r.Body = io.NopCloser(bytes.NewBuffer(buffer.Bytes()))
		if buffer.String() != i.Body {
			return false
		}
	} else {
		if len(i.Body) != 0 {
			return false
		}
	}

	return true
}

// matcher is a predicate which matches the provided HTTP request again a
// recorded interaction request.
func (m *defaultMatcher) matcher(r *http.Request, i Request) bool {
	if r.Method != i.Method {
		return false
	}

	if r.URL.String() != i.URL {
		return false
	}

	if r.Proto != i.Proto {
		return false
	}

	if r.ProtoMajor != i.ProtoMajor {
		return false
	}

	if r.ProtoMinor != i.ProtoMinor {
		return false
	}

	requestHeader := r.Header.Clone()
	cassetteRequestHeaders := i.Headers.Clone()

	for _, header := range m.ignoreHeaders {
		delete(requestHeader, header)
		delete(cassetteRequestHeaders, header)
	}

	if !m.deepEqualContents(requestHeader, cassetteRequestHeaders) {
		return false
	}

	if !m.bodyMatches(r, i) {
		return false
	}

	if r.ContentLength != i.ContentLength {
		return false
	}

	if !m.deepEqualContents(r.TransferEncoding, i.TransferEncoding) {
		return false
	}

	if r.Host != i.Host {
		return false
	}

	if err := r.ParseForm(); err != nil {
		return false
	}

	if !m.deepEqualContents(r.Form, i.Form) {
		return false
	}

	if !m.deepEqualContents(r.Trailer, i.Trailer) {
		return false
	}

	if r.RemoteAddr != i.RemoteAddr {
		return false
	}

	if r.RequestURI != i.RequestURI {
		return false
	}

	return true
}

// DefaultMatcher is the default matcher used to match HTTP requests with
// recorded interactions
var DefaultMatcher = NewDefaultMatcher()

// Cassette represents a cassette containing recorded interactions.
type Cassette struct {
	sync.Mutex `yaml:"-"`

	// Name of the cassette
	Name string `yaml:"-"`

	// File name of the cassette as written on disk
	File string `yaml:"-"`

	// Cassette format version
	Version int `yaml:"version"`

	// Interactions between client and server
	Interactions []*Interaction `yaml:"interactions"`

	// ReplayableInteractions defines whether to allow
	// interactions to be replayed or not
	ReplayableInteractions bool `yaml:"-"`

	// Matches actual request with interaction requests.
	Matcher MatcherFunc `yaml:"-"`

	// IsNew specifies whether this is a newly created cassette.
	// Returns false, when the cassette was loaded from an
	// existing source, e.g. a file.
	IsNew bool `yaml:"-"`

	nextInteractionId int `yaml:"-"`

	// MarshalFunc is a custom marshal func.
	MarshalFunc MarshalFunc `yaml:"-"`
}

// New creates a new empty cassette
func New(name string) *Cassette {
	c := &Cassette{
		Name:                   name,
		File:                   fmt.Sprintf("%s.yaml", name),
		Version:                CassetteFormatVersion,
		Interactions:           make([]*Interaction, 0),
		Matcher:                DefaultMatcher,
		ReplayableInteractions: false,
		IsNew:                  true,
		nextInteractionId:      0,
	}

	return c
}

// Load reads a cassette file from disk
func Load(name string) (*Cassette, error) {
	return LoadWithFS(name, NewDiskFS())
}

// Load reads a cassette file from disk
func LoadWithFS(name string, fs FS) (*Cassette, error) {
	c := New(name)
	data, err := fs.ReadFile(c.File)
	if err != nil {
		return nil, err
	}

	c.IsNew = false
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	if c.Version != CassetteFormatVersion {
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedCassetteFormat, CassetteFormatVersion)
	}
	c.nextInteractionId = len(c.Interactions)

	return c, err
}

// AddInteraction appends a new interaction to the cassette
func (c *Cassette) AddInteraction(i *Interaction) {
	c.Lock()
	defer c.Unlock()
	i.ID = c.nextInteractionId
	c.nextInteractionId += 1
	c.Interactions = append(c.Interactions, i)
}

// GetInteraction retrieves a recorded request/response interaction
func (c *Cassette) GetInteraction(r *http.Request) (*Interaction, error) {
	return c.getInteraction(r)
}

// getInteraction searches for the interaction corresponding to the given HTTP
// request, by using the configured [MatcherFunc].
func (c *Cassette) getInteraction(r *http.Request) (*Interaction, error) {
	c.Lock()
	defer c.Unlock()
	if r.Body == nil {
		// causes an error in the matcher when we try to do r.ParseForm if r.Body is nil
		// r.ParseForm returns missing form body error
		r.Body = http.NoBody
	}
	replayed := 0
	for _, i := range c.Interactions {
		if i.replayed {
			replayed++
		}
		if (c.ReplayableInteractions || !i.replayed) && c.Matcher(r, i.Request) {
			i.replayed = true
			return i, nil
		}
	}
	return nil, ErrInteractionNotFound
}

// Save writes the cassette data on disk for future re-use
func (c *Cassette) Save() error {
	return c.SaveWithFS(NewDiskFS())
}

// SaveWithFS writes the cassette data on abstract filesystem for future re-use
func (c *Cassette) SaveWithFS(fs FS) error {
	c.Lock()
	defer c.Unlock()

	// Filter out interactions which should be discarded. While discarding
	// interactions we should also fix the interaction IDs, so that we don't
	// introduce gaps in the final results.
	nextId := 0
	interactions := make([]*Interaction, 0)
	for _, i := range c.Interactions {
		if !i.DiscardOnSave {
			i.ID = nextId
			interactions = append(interactions, i)
			nextId += 1
		}
	}
	c.Interactions = interactions

	// Marshal to YAML and save interactions
	data, err := c.MarshalFunc(c)
	if err != nil {
		return err
	}

	// Honor the YAML structure specification
	// http://www.yaml.org/spec/1.2/spec.html#id2760395
	return fs.WriteFile(c.File, append([]byte("---\n"), data...))
}
