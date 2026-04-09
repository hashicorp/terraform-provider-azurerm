package vcr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

type fallbackTransport struct {
	r        *recorder.Recorder
	testName string
}

func GetVCRMode(rt http.RoundTripper) (recorder.Mode, bool) {
	if ft, ok := rt.(*fallbackTransport); ok {
		return ft.Mode(), true
	}

	return -1, false
}

func (f *fallbackTransport) Mode() recorder.Mode {
	return f.r.Mode()
}

func (f *fallbackTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Attempt to extract dynamic values from URL
	extracted := make(map[string]string)
	urlStr := req.URL.String()
	for placeholder, extractorRe := range dynamicExtractors {
		if match := extractorRe.FindStringSubmatch(urlStr); len(match) > 1 {
			extracted[placeholder] = match[1]
		}
	}

	if req.Body != nil && req.Body != http.NoBody {
		if bodyBytes, err := io.ReadAll(req.Body); err == nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // restore original
			bodyStr := string(bodyBytes)
			for placeholder, extractorRe := range dynamicExtractors {
				if match := extractorRe.FindStringSubmatch(bodyStr); len(match) > 1 {
					extracted[placeholder] = match[1]
				}
			}
		}
	}

	if len(extracted) > 0 {
		updateDynamicSidecar(f.testName, extracted)
	}

	resp, err := f.r.RoundTrip(req)

	if err != nil && errors.Is(err, cassette.ErrInteractionNotFound) {
		body := fmt.Sprintf("VCR Interaction not found for request: %s %s\nOriginal error: %v", req.Method, req.URL.String(), err)
		return &http.Response{
			Status:        "418 I'm a teapot (VCR Interaction Not Found)",
			StatusCode:    418,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          io.NopCloser(strings.NewReader(body)),
			Request:       req,
			Header:        make(http.Header),
			ContentLength: int64(len(body)),
		}, nil
	}

	// Intercept and rewrite response to inject the dynamic timestamp if we recorded one for this URL.
	if err == nil && resp != nil && resp.Body != nil && resp.Body != http.NoBody {
		mode := f.Mode()
		if mode == recorder.ModeReplayOnly || mode == recorder.ModeReplayWithNewEpisodes {
			dynamicMap := readDynamicSidecar(f.testName)
			if len(dynamicMap) > 0 {
				if respBodyBytes, err := io.ReadAll(resp.Body); err == nil {
					respString := string(respBodyBytes)
					modified := false
					for placeholder, dynamicVal := range dynamicMap {
						if strings.Contains(respString, placeholder) {
							respString = strings.ReplaceAll(respString, placeholder, dynamicVal)
							modified = true
						}
					}
					if modified {
						resp.Body = io.NopCloser(strings.NewReader(respString))
						resp.ContentLength = int64(len(respString))
					} else {
						resp.Body = io.NopCloser(bytes.NewReader(respBodyBytes))
					}
				}
			}
		}
	}

	return resp, err
}
