// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
)

type synapseEndpointTransport func(*http.Request) (*http.Response, error)

func TestSynapseManagedPrivateEndpointFQDNRefresh(t *testing.T) {
	const fqdn = "examplestorageaccount.dfs.core.windows.net"
	const changedFQDN = "otherstorageaccount.dfs.core.windows.net"
	const workspaceID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Synapse/workspaces/example"
	const targetID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Storage/storageAccounts/examplestorageaccount"
	const endpointID = workspaceID + "/managedVirtualNetworks/default/managedPrivateEndpoints/example"
	tests := []struct {
		name          string
		prior         []string
		configured    []string
		responseField string
		expected      []string
		replacement   bool
	}{
		{name: "derived FQDN without configuration", responseField: `,"fqdns":["` + fqdn + `"]`, expected: []string{fqdn}},
		{name: "derived FQDN retained when omitted", prior: []string{fqdn}, expected: []string{fqdn}},
		{name: "omitted FQDN without prior state"},
		{name: "missing state cannot infer explicit FQDN", configured: []string{fqdn}, replacement: true},
		{name: "explicit FQDN returned", prior: []string{fqdn}, configured: []string{fqdn}, responseField: `,"fqdns":["` + fqdn + `"]`, expected: []string{fqdn}},
		{name: "explicit FQDN omitted", prior: []string{fqdn}, configured: []string{fqdn}, expected: []string{fqdn}},
		{name: "explicit FQDN null", prior: []string{fqdn}, configured: []string{fqdn}, responseField: `,"fqdns":null`, expected: []string{fqdn}},
		{name: "explicit FQDN authoritatively empty", prior: []string{fqdn}, configured: []string{fqdn}, responseField: `,"fqdns":[]`, replacement: true},
		{name: "derived FQDN authoritatively empty", prior: []string{fqdn}, responseField: `,"fqdns":[]`},
		{name: "explicit FQDN differs from API", prior: []string{fqdn}, configured: []string{fqdn}, responseField: `,"fqdns":["` + changedFQDN + `"]`, expected: []string{changedFQDN}, replacement: true},
		{name: "configured FQDN changed", prior: []string{fqdn}, configured: []string{changedFQDN}, responseField: `,"fqdns":["` + fqdn + `"]`, expected: []string{fqdn}, replacement: true},
		{name: "configured FQDN removed adopts API", prior: []string{fqdn}, responseField: `,"fqdns":["` + fqdn + `"]`, expected: []string{fqdn}},
		{name: "explicit empty configuration clears FQDN", prior: []string{fqdn}, configured: []string{}, responseField: `,"fqdns":["` + fqdn + `"]`, expected: []string{fqdn}, replacement: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resource := resourceSynapseManagedPrivateEndpoint()
			data := resource.TestResourceData()
			data.SetId(endpointID)
			if test.prior != nil {
				if err := data.Set("fully_qualified_domain_names", test.prior); err != nil {
					t.Fatal(err)
				}
			}
			// Start from persisted state, as refresh and import do, rather than configuration.
			data = resource.Data(data.State())
			api, err := managedprivateendpoints.NewManagedPrivateEndpointsClientUnconfigured()
			if err != nil {
				t.Fatal(err)
			}
			requests := 0
			// Authorization is local to this fake transport; no credentials are read.
			api.Client.AuthorizeRequest = func(context.Context, *http.Request, auth.Authorizer) error { return nil }
			api.Client.SetTransport(synapseEndpointTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				if request.Method != http.MethodGet || request.URL.Path != "/managedVirtualNetworks/default/managedPrivateEndpoints/example" {
					return nil, fmt.Errorf("unexpected request: %s %s", request.Method, request.URL.Path)
				}
				body := `{"properties":{"groupId":"dfs","privateLinkResourceId":"` + targetID + `"` + test.responseField + `}}`
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: request}, nil
			}))
			meta := &clients.Client{
				StopContext: context.Background(),
				Account:     &clients.ResourceManagerAccount{Environment: *environments.AzurePublic()},
				Synapse:     &client.Client{ManagedPrivateEndpointsClient: api},
			}
			for refresh := range 2 {
				if err := resourceSynapseManagedPrivateEndpointRead(data, meta); err != nil {
					t.Fatal(err)
				}
				got := data.Get("fully_qualified_domain_names").([]any)
				want := make([]any, len(test.expected))
				for i, value := range test.expected {
					want[i] = value
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("refresh %d FQDNs = %#v, want %#v", refresh, got, want)
				}
				config := map[string]any{"name": "example", "synapse_workspace_id": workspaceID, "target_resource_id": targetID, "subresource_name": "dfs"}
				if test.configured != nil {
					values := make([]any, len(test.configured))
					for i, value := range test.configured {
						values[i] = value
					}
					config["fully_qualified_domain_names"] = values
				}
				diff, err := resource.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(config), nil)
				if err != nil {
					t.Fatal(err)
				}
				if test.replacement {
					if diff == nil || !diff.RequiresNew() {
						t.Errorf("refresh %d: expected replacement, got %#v", refresh, diff)
					}
				} else if diff != nil && !diff.Empty() {
					t.Errorf("refresh %d: expected no changes, got %#v", refresh, diff.Attributes)
				}
				data = resource.Data(data.State())
			}
			if requests != 2 {
				t.Errorf("GET count = %d, want 2", requests)
			}
		})
	}
}

func (f synapseEndpointTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestSynapseManagedPrivateEndpointReadErrors(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			resource := resourceSynapseManagedPrivateEndpoint()
			data := resource.TestResourceData()
			id := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Synapse/workspaces/example/managedVirtualNetworks/default/managedPrivateEndpoints/example"
			data.SetId(id)
			if err := data.Set("fully_qualified_domain_names", []string{"example.dfs.core.windows.net"}); err != nil {
				t.Fatal(err)
			}
			api, err := managedprivateendpoints.NewManagedPrivateEndpointsClientUnconfigured()
			if err != nil {
				t.Fatal(err)
			}
			// Authorization is local to this fake transport; no credentials are read.
			api.Client.AuthorizeRequest = func(context.Context, *http.Request, auth.Authorizer) error { return nil }
			api.Client.SetTransport(synapseEndpointTransport(func(request *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"TestError","message":"synthetic response"}}`)), Request: request}, nil
			}))
			meta := &clients.Client{
				StopContext: context.Background(),
				Account:     &clients.ResourceManagerAccount{Environment: *environments.AzurePublic()},
				Synapse:     &client.Client{ManagedPrivateEndpointsClient: api},
			}
			err = resourceSynapseManagedPrivateEndpointRead(data, meta)
			if status == http.StatusNotFound {
				if err != nil {
					t.Fatal(err)
				}
				if data.Id() != "" {
					t.Errorf("404 did not remove resource from state: %s", data.Id())
				}
			} else {
				if err == nil {
					t.Fatal("expected forbidden error")
				}
				if data.Id() != id {
					t.Errorf("forbidden response changed resource ID: %s", data.Id())
				}
				if got := data.Get("fully_qualified_domain_names"); !reflect.DeepEqual(got, []any{"example.dfs.core.windows.net"}) {
					t.Errorf("forbidden response changed FQDNs: %#v", got)
				}
			}
		})
	}
}
