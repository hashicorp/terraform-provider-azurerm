package common

import (
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/go-uuid"

	"github.com/Azure/go-autorest/autorest"
)

func TestCorrelationRequestIDGenerated(t *testing.T) {
	first := correlationRequestID()

	if first == "" {
		t.Fatal("no correlation request ID generated")
	}

	second := correlationRequestID()
	if first != second {
		t.Fatal("subsequent correlation request ID not the same as the first")
	}
}

func TestCorrelationRequestIDSpecified(t *testing.T) {
	// Ensure the correlation request id is re-set in this test.
	msCorrelationRequestIDOnce = sync.Once{}

	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}
	if err := os.Setenv(envArmCorrelationRequestID, id); err != nil {
		t.Fatalf("failed to set environment variable for %s: %v", envArmCorrelationRequestID, err)
	}

	if correlationRequestID() != id {
		t.Fatalf("correlation id got is not the same as what is set to %q", envArmCorrelationRequestID)
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
