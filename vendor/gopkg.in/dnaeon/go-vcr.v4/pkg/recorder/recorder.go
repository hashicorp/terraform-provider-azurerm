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
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"slices"
	"strconv"
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

// String implements the [fmt.Stringer] interface for [Mode].
func (m Mode) String() string {
	switch m {
	case ModeRecordOnly:
		return "RecordOnly"
	case ModeReplayOnly:
		return "ReplayOnly"
	case ModeReplayWithNewEpisodes:
		return "ReplayWithNewEpisodes"
	case ModeRecordOnce:
		return "RecordOnce"
	case ModePassthrough:
		return "Passthrough"
	default:
		return fmt.Sprintf("Mode(%d)", m)
	}
}

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

// String implements the [fmt.Stringer] interface for [HookKind].
func (hk HookKind) String() string {
	switch hk {
	case AfterCaptureHook:
		return "AfterCapture"
	case BeforeSaveHook:
		return "BeforeSave"
	case BeforeResponseReplayHook:
		return "BeforeResponseReplay"
	case OnRecorderStopHook:
		return "OnRecorderStop"
	default:
		return fmt.Sprintf("HookKind(%d)", hk)
	}
}

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
	debugLogger  *slog.Logger
}

// safeMethods enumerates the HTTP methods that are considered safe per
// RFC 9110, section 9.2.1.
var safeMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodTrace,
}

func (r *blockUnsafeMethodsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !slices.Contains(safeMethods, req.Method) {
		r.debugLogger.Debug("unsafe method blocked",
			"method", req.Method,
			"url", req.URL.String(),
		)
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

	// marshalFunc is the [cassette.MarshalFunc] which is called for
	// encoding the cassette.
	marshalFunc cassette.MarshalFunc

	// debugWriter and debugLogger are used for emitting debug events by the
	// recorder.
	debugWriter io.Writer
	debugLogger *slog.Logger
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

// WithFS is an [Option], which configures the [Recorder] to use custom
// filesystem ([cassette.FS]) implementation. This allows the [Recorder] to use
// any FS-compatible backend (e.g., local disk, in-memory, or mock) for reading
// and writing files.
func WithFS(fs cassette.FS) Option {
	opt := func(r *Recorder) {
		r.fs = fs
	}

	return opt
}

// WithMarshalFunc is an [Option], which configures the [Recorder] to use custom
// YAML marshal func. This allows customization of the YAML encoding process,
// such as setting string literal style, etc.
func WithMarshalFunc(marshalFunc cassette.MarshalFunc) Option {
	return func(r *Recorder) {
		r.marshalFunc = marshalFunc
	}
}

// WithDebugWriter is an [Option], which configures the [Recorder] to use the
// given [io.Writer] for recording debug events.
func WithDebugWriter(w io.Writer) Option {
	opt := func(r *Recorder) {
		r.debugWriter = w
	}

	return opt
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
		debugWriter:            io.Discard,
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

	// Configure debug logger to [os.Stderr] if `VCR_DEBUG' env var is set,
	// and the current debug writer is [io.Discard]. That way tests which
	// configure their own debug writers will not be affected, but tests
	// that don't explicitly set it will get one for free.
	ok, _ := strconv.ParseBool(os.Getenv("VCR_DEBUG"))
	if ok && r.debugWriter == io.Discard {
		r.debugWriter = os.Stderr
	}

	logHandler := slog.NewTextHandler(r.debugWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	r.debugLogger = slog.New(logHandler).With(
		"component", "recorder",
		"mode", r.mode.String(),
	)
	r.cassette.DebugLogger = slog.New(logHandler).With(
		"component", "cassette",
		"file", r.cassette.File,
	)

	if !r.cassette.IsNew {
		r.cassette.DebugLogger.Debug("cassette loaded",
			"version", r.cassette.Version,
			"interaction_count", len(r.cassette.Interactions),
		)
	}

	r.debug("recorder initialized",
		"cassette_name", r.cassette.Name,
		"cassette_file", r.cassette.File,
		"is_new_cassette", r.cassette.IsNew,
		"interaction_count", len(r.cassette.Interactions),
		"replayable_interactions", r.replayableInteractions,
		"block_unsafe_methods", r.blockUnsafeMethods,
		"skip_request_latency", r.skipRequestLatency,
		"passthroughs", len(r.passthroughs),
		"hooks", len(r.hooks),
	)

	return r, nil
}

// debug emits a debug event
func (r *Recorder) debug(msg string, args ...any) {
	r.debugLogger.Debug(msg, args...)
}

// getCassette creates a new [*cassette.Cassette], or loads an already existing
// one depending on the mode of the recorder.
func (r *Recorder) getCassette() (*cassette.Cassette, error) {
	if r.cassetteName == "" {
		return nil, ErrNoCassetteName
	}

	// Create or the cassette depending on the mode we are operating in.
	cassetteFile := cassette.New(r.cassetteName).File
	cassetteExists := r.fs.IsFileExists(cassetteFile)

	switch {
	case r.mode == ModeRecordOnly:
		return cassette.New(r.cassetteName), nil
	case r.mode == ModeReplayOnly && !cassetteExists:
		return nil, fmt.Errorf("%w: %s", cassette.ErrCassetteNotFound, cassetteFile)
	case r.mode == ModeReplayOnly && cassetteExists:
		return cassette.LoadWithFS(r.cassetteName, r.fs)
	case r.mode == ModeReplayWithNewEpisodes && !cassetteExists:
		return cassette.New(r.cassetteName), nil
	case r.mode == ModeReplayWithNewEpisodes && cassetteExists:
		return cassette.LoadWithFS(r.cassetteName, r.fs)
	case r.mode == ModeRecordOnce && !cassetteExists:
		return cassette.New(r.cassetteName), nil
	case r.mode == ModeRecordOnce && cassetteExists:
		return cassette.LoadWithFS(r.cassetteName, r.fs)
	case r.mode == ModePassthrough:
		return cassette.New(r.cassetteName), nil
	default:
		return nil, ErrInvalidMode
	}
}

// getRoundTripper returns the [http.RoundTripper] used by the recorder.
func (r *Recorder) getRoundTripper() http.RoundTripper {
	if r.blockUnsafeMethods {
		return &blockUnsafeMethodsRoundTripper{
			RoundTripper: r.realTransport,
			debugLogger:  r.debugLogger,
		}
	}

	return r.realTransport
}

// requestHandler proxies requests to their original destination
// If serverResponse is provided, this is used for the recording instead of using RoundTrip
func (rec *Recorder) requestHandler(r *http.Request, serverResponse *http.Response) (*cassette.Interaction, error) {
	if err := r.Context().Err(); err != nil {
		rec.debug("request cancelled before handling", "error", err)
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
			rec.debug("no match, will record new episode")
			break
		} else {
			// Any other error is an error
			return nil, err
		}
	case rec.mode == ModeRecordOnce && !rec.cassette.IsNew:
		// We've got an existing cassette, return what we've got
		rec.debug("replaying from existing cassette")
		return rec.cassette.GetInteraction(r)
	case rec.mode == ModePassthrough:
		// Passthrough requests always hit the original endpoint
		break
	case (rec.mode == ModeRecordOnly || rec.mode == ModeRecordOnce) && rec.cassette.ReplayableInteractions:
		// When running with replayable interactions look for existing
		// interaction first, so we avoid hitting multiple times the
		// same endpoint.
		rec.debug("checking cassette before recording (replayable interactions enabled)")
		interaction, err := rec.cassette.GetInteraction(r)
		if err == nil {
			// Interaction found, return it
			return interaction, nil
		} else if errors.Is(err, cassette.ErrInteractionNotFound) {
			// Interaction not found, we have to record it
			rec.debug("no match, will record")
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
	switch resp {
	case nil:
		rec.debug("forwarding to real transport",
			"method", r.Method,
			"url", r.URL.String(),
		)
		resp, err = rec.getRoundTripper().RoundTrip(r)
		if err != nil {
			rec.debug("real transport error", "error", err)
			return nil, err
		}
	default:
		rec.debug("got handler response",
			"method", r.Method,
			"url", r.URL.String(),
			"status", resp.StatusCode,
		)
	}
	requestDuration := time.Since(start)
	rec.debug("response received",
		"status", resp.StatusCode,
		"duration_ms", requestDuration.Milliseconds(),
		"content_length", resp.ContentLength,
	)
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
func (r *Recorder) Stop() error {
	r.debug("stopping recorder")
	cassetteFile := r.cassette.File
	cassetteExists := r.fs.IsFileExists(cassetteFile)

	// Nothing to do for ModeReplayOnly and ModePassthrough here
	switch {
	case r.mode == ModeRecordOnly || r.mode == ModeReplayWithNewEpisodes:
		r.debug("persisting cassette", "reason", "recording mode")
		if err := r.persistCassette(); err != nil {
			return err
		}

	case r.mode == ModeRecordOnce && !cassetteExists:
		r.debug("persisting cassette", "reason", "record_once with new cassette")
		if err := r.persistCassette(); err != nil {
			return err
		}
	default:
		r.debug("skipping persist", "reason", "replay or existing cassette")
	}

	// Apply on-recorder-stop hooks
	for _, interaction := range r.cassette.Interactions {
		if err := r.applyHooks(interaction, OnRecorderStopHook); err != nil {
			return err
		}
	}

	return nil
}

// persisteCassette persists the cassette on disk for future re-use
func (r *Recorder) persistCassette() error {
	// Apply any before-save hooks
	for _, interaction := range r.cassette.Interactions {
		if err := r.applyHooks(interaction, BeforeSaveHook); err != nil {
			return err
		}
	}

	return r.cassette.SaveWithFS(r.fs)
}

// applyHooks applies the registered hooks of the given kind with the
// specified interaction
func (r *Recorder) applyHooks(i *cassette.Interaction, kind HookKind) error {
	for _, hook := range r.hooks {
		if hook.Kind != kind {
			continue
		}
		r.debug("applying hook", "kind", kind, "id", i.ID)
		if err := hook.Handler(i); err != nil {
			r.debug("hook returned error", "kind", kind, "id", i.ID, "error", err)
			return err
		}
	}

	return nil
}

// RoundTrip implements the [http.RoundTripper] interface
func (r *Recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	r.debug("request received",
		"method", req.Method,
		"url", req.URL.String(),
		"host", req.Host,
	)

	return r.executeAndRecord(req, nil)
}

// executeAndRecord is used internally by the HTTPMiddleware to allow recording a response on the server side
func (r *Recorder) executeAndRecord(req *http.Request, serverResponse *http.Response) (*http.Response, error) {
	// Passthrough mode, use real transport
	if r.mode == ModePassthrough {
		r.debug("passthrough mode, bypassing cassette")
		return r.getRoundTripper().RoundTrip(req)
	}

	// Apply passthrough handler functions
	for i, passthroughFunc := range r.passthroughs {
		if passthroughFunc(req) {
			r.debug("passthrough func matched", "index", i)
			return r.getRoundTripper().RoundTrip(req)
		}
	}

	interaction, err := r.requestHandler(req, serverResponse)
	if err != nil {
		return nil, err
	}

	// Apply before-response-replay hooks
	if err := r.applyHooks(interaction, BeforeResponseReplayHook); err != nil {
		return nil, err
	}

	select {
	case <-req.Context().Done():
		err := req.Context().Err()
		r.debug("request cancelled during replay", "error", err)
		return nil, err
	default:
		// Apply the duration defined in the interaction
		if !r.skipRequestLatency {
			r.debug("simulating latency", "duration_ms", interaction.Response.Duration.Milliseconds())
			<-time.After(interaction.Response.Duration)
		}

		r.debug("replaying response",
			"id", interaction.ID,
			"status", interaction.Response.Code,
			"response", interaction.Response,
		)

		return interaction.GetHTTPResponse()
	}
}

// Mode returns recorder state
func (r *Recorder) Mode() Mode {
	return r.mode
}

// GetDefaultClient returns an HTTP client with a pre-configured
// transport
func (r *Recorder) GetDefaultClient() *http.Client {
	client := &http.Client{
		Transport: r,
	}

	return client
}

// IsNewCassette returns true, if the recorder was started with a
// new/empty cassette. Returns false, if it was started using an
// existing cassette, which was loaded.
func (r *Recorder) IsNewCassette() bool {
	return r.cassette.IsNew
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
func (r *Recorder) IsRecording() bool {
	switch {
	case r.mode == ModeRecordOnly || r.mode == ModeReplayWithNewEpisodes:
		return true
	case r.mode == ModeReplayOnly || r.mode == ModePassthrough:
		return false
	case r.mode == ModeRecordOnce && r.IsNewCassette():
		return true
	default:
		return false
	}
}
