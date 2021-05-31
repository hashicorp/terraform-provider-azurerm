// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ErrorUnmarshaller is the func to invoke when the endpoint returns an error response that requires unmarshalling.
type ErrorUnmarshaller func(*Response) error

// NewLROPoller creates an LROPoller based on the provided initial response.
// pollerID - a unique identifier for an LRO.  it's usually the client.Method string.
// NOTE: this is only meant for internal use in generated code.
func NewLROPoller(pollerID string, resp *Response, pl Pipeline, eu ErrorUnmarshaller) (*LROPoller, error) {
	// this is a back-stop in case the swagger is incorrect (i.e. missing one or more status codes for success).
	// ideally the codegen should return an error if the initial response failed and not even create a poller.
	if !lroStatusCodeValid(resp) {
		return nil, errors.New("the operation failed or was cancelled")
	}
	opLoc := resp.Header.Get(headerOperationLocation)
	loc := resp.Header.Get(headerLocation)
	// in the case of both headers, always prefer the operation-location header
	if opLoc != "" {
		return &LROPoller{
			lro:  newOpPoller(pollerID, opLoc, loc, resp),
			pl:   pl,
			eu:   eu,
			resp: resp,
		}, nil
	}
	if loc != "" {
		return &LROPoller{
			lro:  newLocPoller(pollerID, loc, resp.StatusCode),
			pl:   pl,
			eu:   eu,
			resp: resp,
		}, nil
	}
	return &LROPoller{lro: &nopPoller{}, resp: resp}, nil
}

// NewLROPollerFromResumeToken creates an LROPoller from a resume token string.
// pollerID - a unique identifier for an LRO.  it's usually the client.Method string.
// NOTE: this is only meant for internal use in generated code.
func NewLROPollerFromResumeToken(pollerID string, token string, pl Pipeline, eu ErrorUnmarshaller) (*LROPoller, error) {
	// unmarshal into JSON object to determine the poller type
	obj := map[string]interface{}{}
	err := json.Unmarshal([]byte(token), &obj)
	if err != nil {
		return nil, err
	}
	t, ok := obj["type"]
	if !ok {
		return nil, errors.New("missing type field")
	}
	tt, ok := t.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type format %T", t)
	}
	// the type is encoded as "pollerType;lroPoller"
	sem := strings.LastIndex(tt, ";")
	if sem < 0 {
		return nil, fmt.Errorf("invalid poller type %s", tt)
	}
	// ensure poller types match
	if received := tt[:sem]; received != pollerID {
		return nil, fmt.Errorf("cannot resume from this poller token.  expected %s, received %s", pollerID, received)
	}
	// now rehydrate the poller based on the encoded poller type
	var lro lroPoller
	switch pt := tt[sem+1:]; pt {
	case "opPoller":
		lro = &opPoller{}
	case "locPoller":
		lro = &locPoller{}
	default:
		return nil, fmt.Errorf("unhandled lroPoller type %s", pt)
	}
	if err = json.Unmarshal([]byte(token), lro); err != nil {
		return nil, err
	}
	return &LROPoller{lro: lro, pl: pl, eu: eu}, nil
}

// LROPoller encapsulates state and logic for polling on long-running operations.
// NOTE: this is only meant for internal use in generated code.
type LROPoller struct {
	lro  lroPoller
	pl   Pipeline
	eu   ErrorUnmarshaller
	resp *Response
	err  error
}

// Done returns true if the LRO has reached a terminal state.
func (l *LROPoller) Done() bool {
	if l.err != nil {
		return true
	}
	return l.lro.Done()
}

// Poll sends a polling request to the polling endpoint and returns the response or error.
func (l *LROPoller) Poll(ctx context.Context) (*http.Response, error) {
	if l.Done() {
		// the LRO has reached a terminal state, don't poll again
		if l.resp != nil {
			return l.resp.Response, nil
		}
		return nil, l.err
	}
	req, err := NewRequest(ctx, http.MethodGet, l.lro.URL())
	if err != nil {
		return nil, err
	}
	resp, err := l.pl.Do(req)
	if err != nil {
		// don't update the poller for failed requests
		return nil, err
	}
	if !lroStatusCodeValid(resp) {
		// the LRO failed.  unmarshall the error and update state
		l.err = l.eu(resp)
		l.resp = nil
		return nil, l.err
	}
	if err = l.lro.Update(resp); err != nil {
		return nil, err
	}
	l.resp = resp
	return l.resp.Response, nil
}

// ResumeToken returns a token string that can be used to resume a poller that has not yet reached a terminal state.
func (l *LROPoller) ResumeToken() (string, error) {
	if l.Done() {
		return "", errors.New("cannot create a ResumeToken from a poller in a terminal state")
	}
	b, err := json.Marshal(l.lro)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FinalResponse will perform a final GET request and return the final HTTP response for the polling
// operation and unmarshall the content of the payload into the respType interface that is provided.
func (l *LROPoller) FinalResponse(ctx context.Context, respType interface{}) (*http.Response, error) {
	if !l.Done() {
		return nil, errors.New("cannot return a final response from a poller in a non-terminal state")
	}
	// if there's nothing to unmarshall into just return the final response
	if respType == nil {
		return l.resp.Response, nil
	}
	u, err := l.lro.FinalGetURL(l.resp)
	if err != nil {
		return nil, err
	}
	if u != "" {
		req, err := NewRequest(ctx, http.MethodGet, u)
		if err != nil {
			return nil, err
		}
		resp, err := l.pl.Do(req)
		if err != nil {
			return nil, err
		}
		if !lroStatusCodeValid(resp) {
			return nil, l.eu(resp)
		}
		l.resp = resp
	}
	body, err := ioutil.ReadAll(l.resp.Body)
	l.resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, respType); err != nil {
		return nil, err
	}
	return l.resp.Response, nil
}

// PollUntilDone will handle the entire span of the polling operation until a terminal state is reached,
// then return the final HTTP response for the polling operation and unmarshal the content of the payload
// into the respType interface that is provided.
func (l *LROPoller) PollUntilDone(ctx context.Context, freq time.Duration, respType interface{}) (*http.Response, error) {
	logPollUntilDoneExit := func(v interface{}) {
		Log().Writef(LogLongRunningOperation, "END PollUntilDone() for %T: %v", l.lro, v)
	}
	Log().Writef(LogLongRunningOperation, "BEGIN PollUntilDone() for %T", l.lro)
	if l.resp != nil {
		// initial check for a retry-after header existing on the initial response
		if retryAfter := RetryAfter(l.resp.Response); retryAfter > 0 {
			Log().Writef(LogLongRunningOperation, "initial Retry-After delay for %s", retryAfter.String())
			if err := delay(ctx, retryAfter); err != nil {
				logPollUntilDoneExit(err)
				return nil, err
			}
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := l.Poll(ctx)
		if err != nil {
			logPollUntilDoneExit(err)
			return nil, err
		}
		if l.Done() {
			logPollUntilDoneExit(l.lro.Status())
			if !l.lro.Succeeded() {
				return nil, l.eu(&Response{resp})
			}
			return l.FinalResponse(ctx, respType)
		}
		d := freq
		if retryAfter := RetryAfter(resp); retryAfter > 0 {
			Log().Writef(LogLongRunningOperation, "Retry-After delay for %s", retryAfter.String())
			d = retryAfter
		} else {
			Log().Writef(LogLongRunningOperation, "delay for %s", d.String())
		}
		if err = delay(ctx, d); err != nil {
			logPollUntilDoneExit(err)
			return nil, err
		}
	}
}

var _ Poller = (*LROPoller)(nil)

// abstracts the differences between concrete poller types
type lroPoller interface {
	Done() bool
	Update(resp *Response) error
	FinalGetURL(resp *Response) (string, error)
	URL() string
	Status() string
	Succeeded() bool
}

// ====================================================================================================

// polls on the operation-location header
type opPoller struct {
	Type      string `json:"type"`
	ReqMethod string `json:"reqMethod"`
	ReqURL    string `json:"reqURL"`
	PollURL   string `json:"pollURL"`
	LocURL    string `json:"locURL"`
	status    string
}

func newOpPoller(pollerType, pollingURL, locationURL string, initialResponse *Response) *opPoller {
	return &opPoller{
		Type:      fmt.Sprintf("%s;opPoller", pollerType),
		ReqMethod: initialResponse.Request.Method,
		ReqURL:    initialResponse.Request.URL.String(),
		PollURL:   pollingURL,
		LocURL:    locationURL,
	}
}

func (p *opPoller) URL() string {
	return p.PollURL
}

func (p *opPoller) Done() bool {
	return strings.EqualFold(p.status, "succeeded") ||
		strings.EqualFold(p.status, "failed") ||
		strings.EqualFold(p.status, "cancelled")
}

func (p *opPoller) Succeeded() bool {
	return strings.EqualFold(p.status, "succeeded")
}

func (p *opPoller) Update(resp *Response) error {
	status, err := extractJSONValue(resp, "status")
	if err != nil {
		return err
	}
	if status == "" {
		return errors.New("no status found in body")
	}
	p.status = status
	// if the endpoint returned an operation-location header, update cached value
	if opLoc := resp.Header.Get(headerOperationLocation); opLoc != "" {
		p.PollURL = opLoc
	}
	return nil
}

func (p *opPoller) FinalGetURL(resp *Response) (string, error) {
	if !p.Done() {
		return "", errors.New("cannot return a final response from a poller in a non-terminal state")
	}
	resLoc, err := extractJSONValue(resp, "resourceLocation")
	if err != nil {
		return "", err
	}
	if resLoc != "" {
		return resLoc, nil
	}
	if p.ReqMethod == http.MethodPatch || p.ReqMethod == http.MethodPut {
		return p.ReqURL, nil
	}
	if p.ReqMethod == http.MethodPost && p.LocURL != "" {
		return p.LocURL, nil
	}
	return "", nil
}

func (p *opPoller) Status() string {
	return p.status
}

// ====================================================================================================

// polls on the location header
type locPoller struct {
	Type    string `json:"type"`
	PollURL string `json:"pollURL"`
	status  int
}

func newLocPoller(pollerType, pollingURL string, initialStatus int) *locPoller {
	return &locPoller{
		Type:    fmt.Sprintf("%s;locPoller", pollerType),
		PollURL: pollingURL,
		status:  initialStatus,
	}
}

func (p *locPoller) URL() string {
	return p.PollURL
}

func (p *locPoller) Done() bool {
	// a 202 means the operation is still in progress
	// zero-value indicates the poller was rehydrated from a token
	return p.status > 0 && p.status != http.StatusAccepted
}

func (p *locPoller) Succeeded() bool {
	// any 2xx status code indicates success
	return p.status >= 200 && p.status < 300
}

func (p *locPoller) Update(resp *Response) error {
	// if the endpoint returned a location header, update cached value
	if loc := resp.Header.Get(headerLocation); loc != "" {
		p.PollURL = loc
	}
	p.status = resp.StatusCode
	return nil
}

func (*locPoller) FinalGetURL(*Response) (string, error) {
	return "", nil
}

func (p *locPoller) Status() string {
	return strconv.Itoa(p.status)
}

// ====================================================================================================

// used if the endpoint didn't return any polling headers (synchronous completion)
type nopPoller struct{}

func (*nopPoller) URL() string {
	return ""
}

func (*nopPoller) Done() bool {
	return true
}

func (*nopPoller) Succeeded() bool {
	return true
}

func (*nopPoller) Update(*Response) error {
	return nil
}

func (*nopPoller) FinalGetURL(*Response) (string, error) {
	return "", nil
}

func (*nopPoller) Status() string {
	return "succeeded"
}

// returns true if the LRO response contains a valid HTTP status code
func lroStatusCodeValid(resp *Response) bool {
	return resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusNoContent)
}

// extracs a JSON value from the provided reader
func extractJSONValue(resp *Response, val string) (string, error) {
	if resp.ContentLength == 0 {
		return "", errors.New("the response does not contain a body")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// put the body back so it's available to our callers
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	// unmarshall the body to get the value
	var jsonBody map[string]interface{}
	if err = json.Unmarshal(body, &jsonBody); err != nil {
		return "", err
	}
	v, ok := jsonBody[val]
	if !ok {
		// it might be ok if the field doesn't exist, the caller must make that determination
		return "", nil
	}
	vv, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("the %s value %v was not in string format", val, v)
	}
	return vv, nil
}

func delay(ctx context.Context, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
