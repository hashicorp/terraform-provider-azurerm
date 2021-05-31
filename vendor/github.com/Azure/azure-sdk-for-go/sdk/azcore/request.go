// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"
)

const (
	contentTypeAppJSON = "application/json"
	contentTypeAppXML  = "application/xml"
)

// Base64Encoding is usesd to specify which base-64 encoder/decoder to use when
// encoding/decoding a slice of bytes to/from a string.
type Base64Encoding int

const (
	// Base64StdFormat uses base64.StdEncoding for encoding and decoding payloads.
	Base64StdFormat Base64Encoding = 0

	// Base64URLFormat uses base64.RawURLEncoding for encoding and decoding payloads.
	Base64URLFormat Base64Encoding = 1
)

// Request is an abstraction over the creation of an HTTP request as it passes through the pipeline.
// Don't use this type directly, use NewRequest() instead.
type Request struct {
	*http.Request
	body     ReadSeekCloser
	policies []Policy
	values   opValues
}

type opValues map[reflect.Type]interface{}

// Set adds/changes a value
func (ov opValues) set(value interface{}) {
	ov[reflect.TypeOf(value)] = value
}

// Get looks for a value set by SetValue first
func (ov opValues) get(value interface{}) bool {
	v, ok := ov[reflect.ValueOf(value).Elem().Type()]
	if ok {
		reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v))
	}
	return ok
}

// JoinPaths concatenates multiple URL path segments into one path,
// inserting path separation characters as required.
func JoinPaths(paths ...string) string {
	if len(paths) == 0 {
		return ""
	}
	path := paths[0]
	for i := 1; i < len(paths); i++ {
		if path[len(path)-1] == '/' && paths[i][0] == '/' {
			// strip off trailing '/' to avoid doubling up
			path = path[:len(path)-1]
		} else if path[len(path)-1] != '/' && paths[i][0] != '/' {
			// add a trailing '/'
			path = path + "/"
		}
		path += paths[i]
	}
	return path
}

// NewRequest creates a new Request with the specified input.
func NewRequest(ctx context.Context, httpMethod string, endpoint string) (*Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if req.URL.Host == "" {
		return nil, errors.New("no Host in request URL")
	}
	if !(req.URL.Scheme == "http" || req.URL.Scheme == "https") {
		return nil, fmt.Errorf("unsupported protocol scheme %s", req.URL.Scheme)
	}
	return &Request{Request: req}, nil
}

// Next calls the next policy in the pipeline.
// If there are no more policies, nil and ErrNoMorePolicies are returned.
// This method is intended to be called from pipeline policies.
// To send a request through a pipeline call Pipeline.Do().
func (req *Request) Next() (*Response, error) {
	if len(req.policies) == 0 {
		return nil, ErrNoMorePolicies
	}
	nextPolicy := req.policies[0]
	nextReq := *req
	nextReq.policies = nextReq.policies[1:]
	return nextPolicy.Do(&nextReq)
}

// MarshalAsByteArray will base-64 encode the byte slice v, then calls SetBody.
// The encoded value is treated as a JSON string.
func (req *Request) MarshalAsByteArray(v []byte, format Base64Encoding) error {
	// send as a JSON string
	encode := fmt.Sprintf("\"%s\"", EncodeByteArray(v, format))
	return req.SetBody(NopCloser(strings.NewReader(encode)), contentTypeAppJSON)
}

// MarshalAsJSON calls json.Marshal() to get the JSON encoding of v then calls SetBody.
func (req *Request) MarshalAsJSON(v interface{}) error {
	v = cloneWithoutReadOnlyFields(v)
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("error marshalling type %T: %s", v, err)
	}
	return req.SetBody(NopCloser(bytes.NewReader(b)), contentTypeAppJSON)
}

// MarshalAsXML calls xml.Marshal() to get the XML encoding of v then calls SetBody.
func (req *Request) MarshalAsXML(v interface{}) error {
	b, err := xml.Marshal(v)
	if err != nil {
		return fmt.Errorf("error marshalling type %T: %s", v, err)
	}
	return req.SetBody(NopCloser(bytes.NewReader(b)), contentTypeAppXML)
}

// SetOperationValue adds/changes a mutable key/value associated with a single operation.
func (req *Request) SetOperationValue(value interface{}) {
	if req.values == nil {
		req.values = opValues{}
	}
	req.values.set(value)
}

// OperationValue looks for a value set by SetOperationValue().
func (req *Request) OperationValue(value interface{}) bool {
	if req.values == nil {
		return false
	}
	return req.values.get(value)
}

// SetBody sets the specified ReadSeekCloser as the HTTP request body.
func (req *Request) SetBody(body ReadSeekCloser, contentType string) error {
	// Set the body and content length.
	size, err := body.Seek(0, io.SeekEnd) // Seek to the end to get the stream's size
	if err != nil {
		return err
	}
	if size == 0 {
		body.Close()
		return nil
	}
	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// keep a copy of the original body.  this is to handle cases
	// where req.Body is replaced, e.g. httputil.DumpRequest and friends.
	req.body = body
	req.Request.Body = body
	req.Request.ContentLength = size
	req.Header.Set(HeaderContentType, contentType)
	req.Header.Set(HeaderContentLength, strconv.FormatInt(size, 10))
	return nil
}

// SetMultipartFormData writes the specified keys/values as multi-part form
// fields with the specified value.  File content must be specified as a ReadSeekCloser.
// All other values are treated as string values.
func (req *Request) SetMultipartFormData(formData map[string]interface{}) error {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	for k, v := range formData {
		if rsc, ok := v.(ReadSeekCloser); ok {
			// this is the body to upload, the key is its file name
			fd, err := writer.CreateFormFile(k, k)
			if err != nil {
				return err
			}
			// copy the data to the form file
			if _, err = io.Copy(fd, rsc); err != nil {
				return err
			}
			continue
		}
		// ensure the value is in string format
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprintf("%v", v)
		}
		if err := writer.WriteField(k, s); err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}
	req.body = NopCloser(bytes.NewReader(body.Bytes()))
	req.Body = req.body
	req.ContentLength = int64(body.Len())
	req.Header.Set(HeaderContentType, writer.FormDataContentType())
	req.Header.Set(HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))
	return nil
}

// SkipBodyDownload will disable automatic downloading of the response body.
func (req *Request) SkipBodyDownload() {
	req.SetOperationValue(bodyDownloadPolicyOpValues{skip: true})
}

// RewindBody seeks the request's Body stream back to the beginning so it can be resent when retrying an operation.
func (req *Request) RewindBody() error {
	if req.body != nil {
		// Reset the stream back to the beginning and restore the body
		_, err := req.body.Seek(0, io.SeekStart)
		req.Body = req.body
		return err
	}
	return nil
}

// Close closes the request body.
func (req *Request) Close() error {
	if req.Body == nil {
		return nil
	}
	return req.Body.Close()
}

// Telemetry adds telemetry data to the request.
// If telemetry reporting is disabled the value is discarded.
func (req *Request) Telemetry(v string) {
	req.SetOperationValue(requestTelemetry(v))
}

// clone returns a deep copy of the request with its context changed to ctx
func (req *Request) clone(ctx context.Context) *Request {
	r2 := Request{}
	r2 = *req
	r2.Request = req.Request.Clone(ctx)
	return &r2
}

// valid returns nil if the underlying http.Request is well-formed.
func (req *Request) valid() error {
	// check copied from Transport.roundTrip()
	for k, vv := range req.Header {
		if !httpguts.ValidHeaderFieldName(k) {
			req.Close()
			return fmt.Errorf("invalid header field name %q", k)
		}
		for _, v := range vv {
			if !httpguts.ValidHeaderFieldValue(v) {
				req.Close()
				return fmt.Errorf("invalid header field value %q for key %v", v, k)
			}
		}
	}
	return nil
}

// writes to a buffer, used for logging purposes
func (req *Request) writeBody(b *bytes.Buffer) error {
	if req.Body == nil {
		fmt.Fprint(b, "   Request contained no body\n")
		return nil
	}
	if ct := req.Header.Get(HeaderContentType); !shouldLogBody(b, ct) {
		return nil
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(b, "   Failed to read request body: %s\n", err.Error())
		return err
	}
	if err := req.RewindBody(); err != nil {
		return err
	}
	logBody(b, body)
	return nil
}

// EncodeByteArray will base-64 encode the byte slice v.
func EncodeByteArray(v []byte, format Base64Encoding) string {
	if format == Base64URLFormat {
		return base64.RawURLEncoding.EncodeToString(v)
	}
	return base64.StdEncoding.EncodeToString(v)
}

// returns a clone of the object graph pointed to by v, omitting values of all read-only
// fields. if there are no read-only fields in the object graph, no clone is created.
func cloneWithoutReadOnlyFields(v interface{}) interface{} {
	val := reflect.Indirect(reflect.ValueOf(v))
	if val.Kind() != reflect.Struct {
		// not a struct, skip
		return v
	}
	// first walk the graph to find any R/O fields.
	// if there aren't any, skip cloning the graph.
	if !recursiveFindReadOnlyField(val) {
		return v
	}
	return recursiveCloneWithoutReadOnlyFields(val)
}

// returns true if any field in the object graph of val contains the `azure:"ro"` tag value
func recursiveFindReadOnlyField(val reflect.Value) bool {
	t := val.Type()
	// iterate over the fields, looking for the "azure" tag.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		aztag := field.Tag.Get("azure")
		if azureTagIsReadOnly(aztag) {
			return true
		} else if reflect.Indirect(val.Field(i)).Kind() == reflect.Struct && recursiveFindReadOnlyField(reflect.Indirect(val.Field(i))) {
			return true
		}
	}
	return false
}

// clones the object graph of val.  all non-R/O properties are copied to the clone
func recursiveCloneWithoutReadOnlyFields(val reflect.Value) interface{} {
	clone := reflect.New(val.Type())
	t := val.Type()
	// iterate over the fields, looking for the "azure" tag.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		aztag := field.Tag.Get("azure")
		if azureTagIsReadOnly(aztag) {
			// omit from payload
		} else if reflect.Indirect(val.Field(i)).Kind() == reflect.Struct {
			// recursive case
			v := recursiveCloneWithoutReadOnlyFields(reflect.Indirect(val.Field(i)))
			if t.Field(i).Anonymous {
				// NOTE: this does not handle the case of embedded fields of unexported struct types.
				// this should be ok as we don't generate any code like this at present
				reflect.Indirect(clone).Field(i).Set(reflect.Indirect(reflect.ValueOf(v)))
			} else {
				reflect.Indirect(clone).Field(i).Set(reflect.ValueOf(v))
			}
		} else {
			// no azure RO tag, non-recursive case, include in payload
			reflect.Indirect(clone).Field(i).Set(val.Field(i))
		}
	}
	return clone.Interface()
}

// returns true if the "azure" tag contains the option "ro"
func azureTagIsReadOnly(tag string) bool {
	if tag == "" {
		return false
	}
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		if part == "ro" {
			return true
		}
	}
	return false
}

func logBody(b *bytes.Buffer, body []byte) {
	fmt.Fprintln(b, "   --------------------------------------------------------------------------------")
	fmt.Fprintln(b, string(body))
	fmt.Fprintln(b, "   --------------------------------------------------------------------------------")
}
