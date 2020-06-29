package common

import (
	"net/http"
	"testing"

	"github.com/Azure/go-autorest/autorest"
)

func TestCorrelationRequestID(t *testing.T) {
	first := correlationRequestID()

	if first == "" {
		t.Fatal("no correlation request ID generated")
	}

	second := correlationRequestID()
	if first != second {
		t.Fatal("subsequent correlation request ID not the same as the first")
	}
}

func TestWithCorrelationRequestID(t *testing.T) {
	uuid := correlationRequestID()
	req, _ := autorest.Prepare(&http.Request{}, withCorrelationRequestID(uuid))

	if req.Header.Get(HeaderCorrelationRequestID) != uuid {
		t.Fatalf("azure: withCorrelationRequestID failed to set %s -- expected %s, received %s",
			HeaderCorrelationRequestID, uuid, req.Header.Get(HeaderCorrelationRequestID))
	}
}
