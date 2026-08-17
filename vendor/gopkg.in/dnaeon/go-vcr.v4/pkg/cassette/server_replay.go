package cassette

import (
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

// ReplayAssertFunc is used to assert the results of replaying a recorded request against a handler.
// It receives the current Interaction and the httptest.ResponseRecorder.
type ReplayAssertFunc func(t *testing.T, expected *Interaction, actual *httptest.ResponseRecorder)

// DefaultReplayAssertFunc compares the response status code, body, and headers.
// It can be overridden for more specific tests or to use your preferred assertion libraries
var DefaultReplayAssertFunc ReplayAssertFunc = func(t *testing.T, expected *Interaction, actual *httptest.ResponseRecorder) {
	if expected.Response.Code != actual.Result().StatusCode {
		t.Errorf("status code does not match: expected=%d actual=%d", expected.Response.Code, actual.Result().StatusCode)
	}

	if expected.Response.Body != actual.Body.String() {
		t.Errorf("body does not match: expected=%s actual=%s", expected.Response.Body, actual.Body.String())
	}

	if !headersEqual(expected.Response.Headers, actual.Header()) {
		t.Errorf("header values do not match. expected=%v actual=%v", expected.Response.Headers, actual.Header())
	}
}

// TestServerReplay loads a Cassette and replays each Interaction with the provided Handler, then compares the response
func TestServerReplay(t *testing.T, cassetteName string, handler http.Handler) {
	t.Helper()
	TestServerReplayWithFS(t, cassetteName, NewDiskFS(), handler)
}

// TestServerReplayWithFS loads a Cassette and replays each Interaction with the provided Handler, then compares the response.
// Function reads replay from abstract file system.
func TestServerReplayWithFS(t *testing.T, cassetteName string, fs FS, handler http.Handler) {
	t.Helper()

	c, err := LoadWithFS(cassetteName, fs)
	if err != nil {
		t.Errorf("unexpected error loading Cassette: %v", err)
	}

	if len(c.Interactions) == 0 {
		t.Error("no interactions in Cassette")
	}

	for _, interaction := range c.Interactions {
		t.Run(
			fmt.Sprintf("Interaction_%d", interaction.ID), func(t *testing.T) {
				TestInteractionReplay(t, handler, interaction)
			},
		)
	}
}

// TestInteractionReplay replays an Interaction with the provided Handler and compares the response
func TestInteractionReplay(t *testing.T, handler http.Handler, interaction *Interaction) {
	t.Helper()

	req, err := interaction.GetHTTPRequest()
	if err != nil {
		t.Errorf("unexpected error getting interaction request: %v", err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	DefaultReplayAssertFunc(t, interaction, w)
}

func headersEqual(expected, actual http.Header) bool {
	return maps.EqualFunc(
		expected, actual,
		func(v1, v2 []string) bool {
			slices.Sort(v1)
			slices.Sort(v2)

			if !slices.Equal(v1, v2) {
				return false
			}

			return true
		},
	)
}
