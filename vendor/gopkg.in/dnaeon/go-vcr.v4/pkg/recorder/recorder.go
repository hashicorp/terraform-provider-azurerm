// Copyright (c) 2015-2024 Marin Atanasov Nikolov <dnaeon@gmail.com>
// Copyright (c) 2016 David Jack <davars@gmail.com>
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

package recorder

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"go.yaml.in/yaml/v4"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

type MatcherFunc = cassette.MatcherFunc

// ErrNoCassetteName is an error, which is returned when the recorder was
// created without specifying a cassette name.
var ErrNoCassetteName = errors.New("no cassette name specified")

// Mode represents the mode of operation of the recorder
type Mode int

// Recorder states
const (
	// ModeRecordOnly specifies that VCR will run in recording mode
	// only. HTTP interactions will be recorded for each interaction. If the
	// cassette file is present, it will be overwritten.
	ModeRecordOnly Mode = iota

	// ModeReplayOnly specifies that VCR will only replay interactions from
	// previously recorded cassette. If an interaction is missing from the
	// cassette it will return ErrInteractionNotFound error. If the cassette
	// file is missing it will return ErrCassetteNotFound error.
	ModeReplayOnly

	// ModeReplayWithNewEpisodes starts the recorder in replay mode, where
	// existing interactions are returned from the cassette, and missing
	// ones will be recorded and added to the cassette. This mode is useful
	// in cases where you need to update an existing cassette with new
	// interactions, but don't want to wipe out previously recorded
	// interactions. If the cassette file is missing it will create a new
	// one.
	ModeReplayWithNewEpisodes

	// ModeRecordOnce will record new HTTP interactions once only. This mode
	// is useful in cases where you need to record a set of interactions
	// once only and replay only the known interactions. Unknown/missing
	// interactions will cause the recorder to return an
	// ErrInteractionNotFound error. If the cassette file is missing, it
	// will be created.
	ModeRecordOnce

	// ModePassthrough specifies that VCR will not record any interactions
	// at all. In this mode all HTTP requests will be forwarded to the
	// endpoints using the real HTTP transport. In this mode no cassette
	// will be created.
	ModePassthrough
)

// ErrInvalidMode is returned when attempting to start the recorder with invalid
// mode
var ErrInvalidMode = errors.New("invalid recorder mode")

// HookFunc represents a function, which will be invoked in different stages of
// the playback. The hook functions allow for plugging in to the playback and
// transform an interaction, if needed. For example a hook function might redact
// or remove sensitive data from a request/response before it is added to the
// in-memory cassette, or before it is saved on disk. Another use case would be
// to transform the HTTP response before it is returned to the client during
// replay mode.
type HookFunc func(i *cassette.Interaction) error

// Hook kinds
type HookKind int

const (
	// AfterCaptureHook represents a hook, which will be invoked after
	// capturing a request/response pair.
	AfterCaptureHook HookKind = iota

	// BeforeSaveHook represents a hook, which will be invoked right before
	// the cassette is saved on disk.
	BeforeSaveHook

	// BeforeResponseReplayHook represents a hook, which will be invoked
	// before replaying a previously recorded response to the client.
	BeforeResponseReplayHook

	// OnRecorderStopHook is a hook, which will be invoked when the recorder
	// is about to be stopped. This hook is useful for performing any
	// post-actions such as cleanup or reporting.
	OnRecorderStopHook
)

// Hook represents a function hook of a given kind. Depending on the hook kind,
// the function will be invoked in different stages of the playback.
type Hook struct {
	// Handler is the function which will be invoked
	Handler HookFunc

	// Kind represents the hook kind
	Kind HookKind
}

// NewHook creates a new hook.
func NewHook(handler HookFunc, kind HookKind) *Hook {
	hook := &Hook{
		Handler: handler,
		Kind:    kind,
	}

	return hook
}

// PassthroughFunc is a predicate which determines whether a specific HTTP
// request is to be forwarded to the original endpoint. It should return true
// when a request needs to be passed through, and false otherwise.
type PassthroughFunc func(req *http.Request) bool

// ErrUnsafeRequestMethod is returned when the [Recorder] was configured to
// block unsafe methods, and an attempt to use such was invoked. Safe Methods
// are defined as part of RFC 9110, section 9.2.1.
var ErrUnsafeRequestMethod = errors.New("request uses an unsafe method")

type blockUnsafeMethodsRoundTripper struct {
	RoundTripper http.RoundTripper
}

func (r *blockUnsafeMethodsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	safeMethods := map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}
	if _, ok := safeMethods[req.Method]; !ok {
		return nil, ErrUnsafeRequestMethod
	}
	return r.RoundTripper.RoundTrip(req)
}

// Recorder represents a type used to record and replay client and server
// interactions.
type Recorder struct {
	// Cassette used by the recorder
	cassette *cassette.Cassette

	// cassetteName is the name of the cassette to be used by the recorder.
	cassetteName string

	// mode is the mode of the recorder
	mode Mode

	// RealTransport is the underlying http.RoundTripper to make
	// the real requests
	realTransport http.RoundTripper

	// blockUnsafeMethods specifies whether to block requests when making
	// HTTP requests which are not safe. The "Safe Methods" are defined as
	// part of RFC 9110, section 9.2.1, and SHOULD NOT have side effects on
	// the server.
	blockUnsafeMethods bool

	// skipRequestLatency specifies whether to simulate the latency of the
	// recorded interaction.
	skipRequestLatency bool

	// Passthrough handlers
	passthroughs []PassthroughFunc

	// hooks is a list of hooks, which are invoked in different
	// stages of the playback.
	hooks []*Hook

	// matcher is the [MatcherFunc] predicate used to match HTTP requests
	// against recorded interactions.
	matcher MatcherFunc

	// replayableInteractions specifies whether to allow interactions to be
	// replayed multiple times.
	replayableInteractions bool

	// fs specifies custom filesystem ([cassette.FS]) implementation.
	fs cassette.FS

	marshalFunc cassette.MarshalFunc
}

// Option is a function which configures the [Recorder].
type Option func(r *Recorder)

// WithMode is an [Option], which configures the [Recorder] to run in the
// specified mode.
func WithMode(mode Mode) Option {
	opt := func(r *Recorder) {
		r.mode = mode
	}

	return opt
}

// WithRealTransport is an [Option], which configures the [Recorder] to use the
// specified [http.RoundTripper] when making actual HTTP requests.
func WithRealTransport(rt http.RoundTripper) Option {
	opt := func(r *Recorder) {
		r.realTransport = rt
	}

	return opt
}

// WithBlockUnsafeMethods is an [Option], which configures the [Recorder] to
// block HTTP requests, which are not considered "Safe Methods", according to
// RFC 9110, section 9.2.1.
func WithBlockUnsafeMethods(val bool) Option {
	opt := func(r *Recorder) {
		r.blockUnsafeMethods = val
	}

	return opt
}

// WithSkipRequestLatency is an [Option], which configures the [Recorder] whether
// to simulate the latency of the recorded interaction. When set to false it
// will block for the period of time taken by the original request to simulate
// the latency between the recorder and the remote endpoints.
func WithSkipRequestLatency(val bool) Option {
	opt := func(r *Recorder) {
		r.skipRequestLatency = val
	}

	return opt
}

// WithPassthrough is an [Option], which configures the [Recorder] to
// passthrough requests for requests which satisfy the provided
// [PassthroughFunc] predicate.
func WithPassthrough(passfunc PassthroughFunc) Option {
	opt := func(r *Recorder) {
		r.passthroughs = append(r.passthroughs, passfunc)
	}

	return opt
}

// WithHook is an [Option], which configures the [Recorder] to invoke the
// provided hook at the specified playback stage.
func WithHook(handler HookFunc, kind HookKind) Option {
	opt := func(r *Recorder) {
		hook := NewHook(handler, kind)
		r.hooks = append(r.hooks, hook)
	}

	return opt
}

// WithMatchers is an [Option], which configures the [Recorder] to use the
// provided [MatcherFunc] predicate when matching HTTP requests against record
// interactions.
func WithMatcher(matcher MatcherFunc) Option {
	opt := func(r *Recorder) {
		r.matcher = matcher
	}

	return opt
}

// WithReplayableInteractions is an [Option], which configures the [Recorder] to
// allow replaying interactions multiple times. This is useful in situations
// when you need to hit the same endpoint multiple times and want to replay the
// interaction from the cassette each time.
func WithReplayableInteractions(val bool) Option {
	opt := func(r *Recorder) {
		r.replayableInteractions = val
	}

	return opt
}

// WithFS is an [Option], which configures the [Recorder] to use
// custom filesystem ([cassette.FS]) implementation. This allows the [Recorder] to use any
// FS-compatible backend (e.g., local disk, in-memory, or mock) for reading and writing files.
func WithFS(fs cassette.FS) Option {
	opt := func(r *Recorder) {
		r.fs = fs
	}

	return opt
}

// WithMarshalFunc is an [Option], which configures the [Recorder] to use
// custom YAML marshal func. This allows customization of the YAML encoding
// process, such as setting string literal style, etc.
func WithMarshalFunc(marshalFunc cassette.MarshalFunc) Option {
	return func(r *Recorder) {
		r.marshalFunc = marshalFunc
	}
}

// New creates a new [Recorder] and configures it using the provided options.
func New(cassetteName string, opts ...Option) (*Recorder, error) {
	r := &Recorder{
		cassetteName:           cassetteName,
		mode:                   ModeRecordOnce,
		realTransport:          http.DefaultTransport,
		passthroughs:           make([]PassthroughFunc, 0),
		hooks:                  make([]*Hook, 0),
		blockUnsafeMethods:     false,
		skipRequestLatency:     false,
		matcher:                cassette.DefaultMatcher,
		replayableInteractions: false,
		fs:                     cassette.NewDiskFS(),
		marshalFunc:            yaml.Marshal,
	}

	for _, opt := range opts {
		opt(r)
	}

	// Configure the cassette based on the recorder configuration
	c, err := r.getCassette()
	if err != nil {
		return nil, err
	}
	r.cassette = c
	r.cassette.Matcher = r.matcher
	r.cassette.ReplayableInteractions = r.replayableInteractions
	r.cassette.MarshalFunc = r.marshalFunc

	return r, nil
}

// getCassette creates a new [*cassette.Cassette], or loads an already existing
// one depending on the mode of the recorder.
func (rec *Recorder) getCassette() (*cassette.Cassette, error) {
	if rec.cassetteName == "" {
		return nil, ErrNoCassetteName
	}

	// Create or the cassette depending on the mode we are operating in.
	cassetteFile := cassette.New(rec.cassetteName).File
	cassetteExists := rec.fs.IsFileExists(cassetteFile)

	switch {
	case rec.mode == ModeRecordOnly:
		return cassette.New(rec.cassetteName), nil
	case rec.mode == ModeReplayOnly && !cassetteExists:
		return nil, fmt.Errorf("%w: %s", cassette.ErrCassetteNotFound, cassetteFile)
	case rec.mode == ModeReplayOnly && cassetteExists:
		return cassette.LoadWithFS(rec.cassetteName, rec.fs)
	case rec.mode == ModeReplayWithNewEpisodes && !cassetteExists:
		return cassette.New(rec.cassetteName), nil
	case rec.mode == ModeReplayWithNewEpisodes && cassetteExists:
		return cassette.LoadWithFS(rec.cassetteName, rec.fs)
	case rec.mode == ModeRecordOnce && !cassetteExists:
		return cassette.New(rec.cassetteName), nil
	case rec.mode == ModeRecordOnce && cassetteExists:
		return cassette.LoadWithFS(rec.cassetteName, rec.fs)
	case rec.mode == ModePassthrough:
		return cassette.New(rec.cassetteName), nil
	default:
		return nil, ErrInvalidMode
	}
}

// getRoundTripper returns the [http.RoundTripper] used by the recorder.
func (rec *Recorder) getRoundTripper() http.RoundTripper {
	if rec.blockUnsafeMethods {
		return &blockUnsafeMethodsRoundTripper{
			RoundTripper: rec.realTransport,
		}
	}

	return rec.realTransport
}

// requestHandler proxies requests to their original destination
// If serverResponse is provided, this is used for the recording instead of using RoundTrip
func (rec *Recorder) requestHandler(r *http.Request, serverResponse *http.Response) (*cassette.Interaction, error) {
	if err := r.Context().Err(); err != nil {
		return nil, err
	}

	switch {
	case rec.mode == ModeReplayOnly:
		return rec.cassette.GetInteraction(r)
	case rec.mode == ModeReplayWithNewEpisodes:
		interaction, err := rec.cassette.GetInteraction(r)
		if err == nil {
			// Interaction found, return it
			return interaction, nil
		} else if errors.Is(err, cassette.ErrInteractionNotFound) {
			// Interaction not found, we have a new episode
			break
		} else {
			// Any other error is an error
			return nil, err
		}
	case rec.mode == ModeRecordOnce && !rec.cassette.IsNew:
		// We've got an existing cassette, return what we've got
		return rec.cassette.GetInteraction(r)
	case rec.mode == ModePassthrough:
		// Passthrough requests always hit the original endpoint
		break
	case (rec.mode == ModeRecordOnly || rec.mode == ModeRecordOnce) && rec.cassette.ReplayableInteractions:
		// When running with replayable interactions look for existing
		// interaction first, so we avoid hitting multiple times the
		// same endpoint.
		interaction, err := rec.cassette.GetInteraction(r)
		if err == nil {
			// Interaction found, return it
			return interaction, nil
		} else if errors.Is(err, cassette.ErrInteractionNotFound) {
			// Interaction not found, we have to record it
			break
		} else {
			// Any other error is an error
			return nil, err
		}
	default:
		// Anything else hits the original endpoint
		break
	}

	// Copy the original request, so we can read the form values
	reqBytes, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, err
	}

	reqBuffer := bytes.NewBuffer(reqBytes)
	copiedReq, err := http.ReadRequest(bufio.NewReader(reqBuffer))
	if err != nil {
		return nil, err
	}

	err = copiedReq.ParseForm()
	if err != nil {
		return nil, err
	}

	reqBody := &bytes.Buffer{}
	if r.Body != nil && r.Body != http.NoBody {
		// Record the request body so we can add it to the cassette
		r.Body = io.NopCloser(io.TeeReader(r.Body, reqBody))
		if serverResponse != nil {
			// when serverResponse is provided by middleware, it has to be read in order
			// for reqBody buffer to be populated
			_, _ = io.ReadAll(r.Body)
		}
	}

	// Perform request to it's original destination and record the interactions
	// If serverResponse is provided, use it instead
	var start time.Time
	start = time.Now()
	resp := serverResponse
	if resp == nil {
		resp, err = rec.getRoundTripper().RoundTrip(r)
		if err != nil {
			return nil, err
		}
	}
	requestDuration := time.Since(start)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Add interaction to the cassette
	interaction := &cassette.Interaction{
		Request: cassette.Request{
			Proto:            r.Proto,
			ProtoMajor:       r.ProtoMajor,
			ProtoMinor:       r.ProtoMinor,
			ContentLength:    r.ContentLength,
			TransferEncoding: r.TransferEncoding,
			Trailer:          r.Trailer,
			Host:             r.Host,
			RemoteAddr:       r.RemoteAddr,
			RequestURI:       r.RequestURI,
			Body:             reqBody.String(),
			Form:             copiedReq.Form,
			Headers:          r.Header,
			URL:              r.URL.String(),
			Method:           r.Method,
		},
		Response: cassette.Response{
			Status:           resp.Status,
			Code:             resp.StatusCode,
			Proto:            resp.Proto,
			ProtoMajor:       resp.ProtoMajor,
			ProtoMinor:       resp.ProtoMinor,
			TransferEncoding: resp.TransferEncoding,
			Trailer:          resp.Trailer,
			ContentLength:    resp.ContentLength,
			Uncompressed:     resp.Uncompressed,
			Body:             string(respBody),
			Headers:          resp.Header,
			Duration:         requestDuration,
		},
	}

	// Apply after-capture hooks before we add the interaction to
	// the in-memory cassette.
	if err := rec.applyHooks(interaction, AfterCaptureHook); err != nil {
		return nil, err
	}

	rec.cassette.AddInteraction(interaction)

	return interaction, nil
}

// Stop is used to stop the recorder and save any recorded
// interactions if running in one of the recording modes. When
// running in ModePassthrough no cassette will be saved on disk.
func (rec *Recorder) Stop() error {
	cassetteFile := rec.cassette.File
	cassetteExists := rec.fs.IsFileExists(cassetteFile)

	// Nothing to do for ModeReplayOnly and ModePassthrough here
	switch {
	case rec.mode == ModeRecordOnly || rec.mode == ModeReplayWithNewEpisodes:
		if err := rec.persistCassette(); err != nil {
			return err
		}

	case rec.mode == ModeRecordOnce && !cassetteExists:
		if err := rec.persistCassette(); err != nil {
			return err
		}
	}

	// Apply on-recorder-stop hooks
	for _, interaction := range rec.cassette.Interactions {
		if err := rec.applyHooks(interaction, OnRecorderStopHook); err != nil {
			return err
		}
	}

	return nil
}

// persisteCassette persists the cassette on disk for future re-use
func (rec *Recorder) persistCassette() error {
	// Apply any before-save hooks
	for _, interaction := range rec.cassette.Interactions {
		if err := rec.applyHooks(interaction, BeforeSaveHook); err != nil {
			return err
		}
	}

	return rec.cassette.SaveWithFS(rec.fs)
}

// applyHooks applies the registered hooks of the given kind with the
// specified interaction
func (rec *Recorder) applyHooks(i *cassette.Interaction, kind HookKind) error {
	for _, hook := range rec.hooks {
		if hook.Kind == kind {
			if err := hook.Handler(i); err != nil {
				return err
			}
		}
	}

	return nil
}

// RoundTrip implements the [http.RoundTripper] interface
func (rec *Recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	return rec.executeAndRecord(req, nil)
}

// executeAndRecord is used internally by the HTTPMiddleware to allow recording a response on the server side
func (rec *Recorder) executeAndRecord(req *http.Request, serverResponse *http.Response) (*http.Response, error) {
	// Passthrough mode, use real transport
	if rec.mode == ModePassthrough {
		return rec.getRoundTripper().RoundTrip(req)
	}

	// Apply passthrough handler functions
	for _, passthroughFunc := range rec.passthroughs {
		if passthroughFunc(req) {
			return rec.getRoundTripper().RoundTrip(req)
		}
	}

	interaction, err := rec.requestHandler(req, serverResponse)
	if err != nil {
		return nil, err
	}

	// Apply before-response-replay hooks
	if err := rec.applyHooks(interaction, BeforeResponseReplayHook); err != nil {
		return nil, err
	}

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	default:
		// Apply the duration defined in the interaction
		if !rec.skipRequestLatency {
			<-time.After(interaction.Response.Duration)
		}

		return interaction.GetHTTPResponse()
	}
}

// Mode returns recorder state
func (rec *Recorder) Mode() Mode {
	return rec.mode
}

// GetDefaultClient returns an HTTP client with a pre-configured
// transport
func (rec *Recorder) GetDefaultClient() *http.Client {
	client := &http.Client{
		Transport: rec,
	}

	return client
}

// IsNewCassette returns true, if the recorder was started with a
// new/empty cassette. Returns false, if it was started using an
// existing cassette, which was loaded.
func (rec *Recorder) IsNewCassette() bool {
	return rec.cassette.IsNew
}

// IsRecording returns true, if the recorder is recording
// interactions, returns false otherwise. Note, that in some modes
// (e.g. ModeReplayWithNewEpisodes and ModeRecordOnce) the recorder
// might be recording new interactions. For example in ModeRecordOnce,
// we are replaying interactions only if there was an existing
// cassette, and we are recording it, if the cassette is a new one.
// ModeReplayWithNewEpisodes would replay interactions, if they are
// present in the cassette, but will also record new ones, if they are
// not part of the cassette already. In these cases the recorder is
// considered to be recording for these modes.
func (rec *Recorder) IsRecording() bool {
	switch {
	case rec.mode == ModeRecordOnly || rec.mode == ModeReplayWithNewEpisodes:
		return true
	case rec.mode == ModeReplayOnly || rec.mode == ModePassthrough:
		return false
	case rec.mode == ModeRecordOnce && rec.IsNewCassette():
		return true
	default:
		return false
	}
}
