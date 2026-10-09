// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

const testOrderPath = "/providers/Microsoft.Capacity/reservationOrders/11111111-1111-1111-1111-111111111111"

type recordedRequest struct {
	Method     string
	Path       string
	APIVersion string
	Body       string
}

type fakeAzure struct {
	t        *testing.T
	server   *httptest.Server
	mu       sync.Mutex
	requests []recordedRequest

	// reservations returned by the list call, as name -> renew
	reservations []ReservationResponse
	// patchStatus is the status returned by the PATCH calls (defaults to 200)
	patchStatus int
	// polls counts the calls made to the long running operation URL
	polls int
}

func newFakeAzure(t *testing.T) *fakeAzure {
	f := &fakeAzure{t: t, patchStatus: http.StatusOK}
	f.server = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.server.Close)
	return f
}

func (f *fakeAzure) handle(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	f.mu.Lock()
	defer f.mu.Unlock()

	if r.URL.Path == "/operations/poll" {
		f.polls++
		w.Header().Set("Retry-After", "0")
		if f.polls < 2 {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	f.requests = append(f.requests, recordedRequest{
		Method:     r.Method,
		Path:       r.URL.Path,
		APIVersion: r.URL.Query().Get("api-version"),
		Body:       string(body),
	})

	switch {
	case r.Method == http.MethodGet && r.URL.Path == testOrderPath+"/reservations":
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ReservationListResponse{Value: f.reservations})
	case r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, testOrderPath+"/reservations/"):
		switch f.patchStatus {
		case http.StatusAccepted:
			w.Header().Set("Location", f.server.URL+"/operations/poll")
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusAccepted)
		case http.StatusOK:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(f.patchStatus)
			_, _ = w.Write([]byte(`{"error":{"code":"BadRequest","message":"nope"}}`))
		}
	default:
		f.t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}
}

func (f *fakeAzure) patches() []recordedRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []recordedRequest
	for _, r := range f.requests {
		if r.Method == http.MethodPatch {
			out = append(out, r)
		}
	}
	return out
}

func (f *fakeAzure) client(t *testing.T) *Client {
	t.Helper()
	rm, err := resourcemanager.NewClient(environments.NewApiEndpoint("test", f.server.URL, nil), "reservationOrders", ptuReservationAPIVersion)
	if err != nil {
		t.Fatalf("building client: %+v", err)
	}
	// requests go to a local fake, so there is nothing to authorize
	rm.AuthorizeRequest = func(context.Context, *http.Request, auth.Authorizer) error { return nil }
	return &Client{ReservationOrdersClient: rm}
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func reservation(name string, renew bool) ReservationResponse {
	return ReservationResponse{Name: name, Properties: &ReservationResponseProperties{Renew: renew}}
}

func TestGetReservations(t *testing.T) {
	f := newFakeAzure(t)
	f.reservations = []ReservationResponse{reservation("res-1", true), reservation("res-2", false)}

	got, err := f.client(t).GetReservations(testContext(t), testOrderPath)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if len(got) != 2 || got[0].Name != "res-1" || !got[0].Properties.Renew || got[1].Name != "res-2" || got[1].Properties.Renew {
		t.Fatalf("unexpected reservations: %+v", got)
	}
	if v := f.requests[0].APIVersion; v != ptuReservationAPIVersion {
		t.Fatalf("expected api-version %q, got %q", ptuReservationAPIVersion, v)
	}
}

func TestUpdateRenew_PatchesEveryReservation(t *testing.T) {
	for _, renew := range []bool{true, false} {
		t.Run(fmt.Sprintf("renew=%t", renew), func(t *testing.T) {
			f := newFakeAzure(t)
			f.reservations = []ReservationResponse{reservation("res-1", !renew), reservation("res-2", !renew)}

			if err := f.client(t).UpdateRenew(testContext(t), testOrderPath, renew); err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			patches := f.patches()
			if len(patches) != 2 {
				t.Fatalf("expected 2 PATCH requests, got %d: %+v", len(patches), patches)
			}
			for i, name := range []string{"res-1", "res-2"} {
				if want := testOrderPath + "/reservations/" + name; patches[i].Path != want {
					t.Errorf("PATCH %d: expected path %q, got %q", i, want, patches[i].Path)
				}
				if patches[i].APIVersion != ptuReservationAPIVersion {
					t.Errorf("PATCH %d: expected api-version %q, got %q", i, ptuReservationAPIVersion, patches[i].APIVersion)
				}
				if want := fmt.Sprintf(`{"properties":{"renew":%t}}`, renew); strings.TrimSpace(patches[i].Body) != want {
					t.Errorf("PATCH %d: expected body %s, got %s", i, want, patches[i].Body)
				}
			}
		})
	}
}

func TestUpdateRenew_PollsAcceptedResponse(t *testing.T) {
	f := newFakeAzure(t)
	f.patchStatus = http.StatusAccepted
	f.reservations = []ReservationResponse{reservation("res-1", true)}

	if err := f.client(t).UpdateRenew(testContext(t), testOrderPath, false); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if got := len(f.patches()); got != 1 {
		t.Fatalf("expected 1 PATCH request, got %d", got)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.polls < 2 {
		t.Fatalf("expected the long running operation to be polled until it completed, got %d polls", f.polls)
	}
}

func TestUpdateRenew_NoReservations(t *testing.T) {
	f := newFakeAzure(t)

	err := f.client(t).UpdateRenew(testContext(t), testOrderPath, false)
	if err == nil || !strings.Contains(err.Error(), "no reservations found") {
		t.Fatalf("expected a `no reservations found` error, got %v", err)
	}
	if got := len(f.patches()); got != 0 {
		t.Fatalf("expected no PATCH request, got %d", got)
	}
}

func TestUpdateRenew_PatchError(t *testing.T) {
	f := newFakeAzure(t)
	f.patchStatus = http.StatusBadRequest
	f.reservations = []ReservationResponse{reservation("res-1", true), reservation("res-2", true)}

	err := f.client(t).UpdateRenew(testContext(t), testOrderPath, false)
	if err == nil || !strings.Contains(err.Error(), "update renew") {
		t.Fatalf("expected an update renew error, got %v", err)
	}
	if got := len(f.patches()); got != 1 {
		t.Fatalf("expected to stop after the first failing PATCH, got %d PATCH requests", got)
	}
}
