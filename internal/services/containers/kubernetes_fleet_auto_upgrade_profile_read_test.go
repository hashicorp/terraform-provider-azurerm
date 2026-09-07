// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	containerclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

type fleetAutoUpgradeProfileTransport func(*http.Request) (*http.Response, error)

func (f fleetAutoUpgradeProfileTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestKubernetesFleetAutoUpgradeProfileRead(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		status    int
		wantError bool
		wantGone  bool
		update    bool
	}{
		{
			name:   "omitted optional properties",
			body:   `{"properties":{"channel":"Stable"}}`,
			status: http.StatusOK,
		},
		{
			name:      "missing properties",
			body:      `{}`,
			status:    http.StatusOK,
			wantError: true,
		},
		{
			name:      "null properties",
			body:      `{"properties":null}`,
			status:    http.StatusOK,
			wantError: true,
		},
		{
			name:      "update missing properties",
			body:      `{}`,
			status:    http.StatusOK,
			wantError: true,
			update:    true,
		},
		{
			name:      "update null properties",
			body:      `{"properties":null}`,
			status:    http.StatusOK,
			wantError: true,
			update:    true,
		},
		{
			name:     "not found",
			body:     `{"error":{"code":"ResourceNotFound","message":"not found"}}`,
			status:   http.StatusNotFound,
			wantGone: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := environments.NewApiEndpoint("fleet-test", "https://fleet.invalid", nil)
			client, err := autoupgradeprofiles.NewAutoUpgradeProfilesClientWithBaseURI(endpoint)
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			requests := 0
			client.Client.SetTransport(fleetAutoUpgradeProfileTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				if request.Method != http.MethodGet {
					t.Fatalf("unexpected request method: %s", request.Method)
				}
				return &http.Response{
					StatusCode: test.status,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(test.body)),
					Request:    request,
				}, nil
			}))

			resource := containers.KubernetesFleetAutoUpgradeProfileResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{
				Containers: &containerclient.Client{FleetAutoUpgradeProfilesClient: client},
			}, resource)
			id := autoupgradeprofiles.NewAutoUpgradeProfileID("00000000-0000-0000-0000-000000000000", "example", "fleet", "profile")
			metadata.SetID(id)
			previous := containers.KubernetesFleetAutoUpgradeProfileResourceModel{
				Name:                     id.AutoUpgradeProfileName,
				KubernetesFleetManagerId: commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID(),
				Channel:                  "Rapid",
				NodeImageSelectionType:   "Latest",
				UpdateStrategyId:         fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "strategy").ID(),
				Enabled:                  false,
			}
			if err := metadata.Encode(&previous); err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			if test.update {
				err = resource.Update().Func(ctx, metadata)
			} else {
				err = resource.Read().Func(ctx, metadata)
			}
			if requests != 1 {
				t.Fatalf("expected one GET request, got %d requests", requests)
			}
			if (err != nil) != test.wantError {
				t.Fatalf("want error = %t, got %v", test.wantError, err)
			}
			wantID := id.ID()
			if test.wantGone {
				wantID = ""
			}
			if got := metadata.ResourceData.Id(); got != wantID {
				t.Fatalf("unexpected state ID: want %q, got %q", wantID, got)
			}
			if test.wantError {
				if !strings.Contains(err.Error(), "properties") {
					t.Fatalf("expected missing properties error, got %v", err)
				}
				var state containers.KubernetesFleetAutoUpgradeProfileResourceModel
				if err := metadata.Decode(&state); err != nil {
					t.Fatal(err)
				}
				if state != previous {
					t.Fatalf("invalid response changed existing state: want %+v, got %+v", previous, state)
				}
				return
			}
			if test.wantGone {
				return
			}
			if got := metadata.ResourceData.Get("channel"); got != "Stable" {
				t.Fatalf("unexpected channel: %v", got)
			}
			if got := metadata.ResourceData.Get("node_image_selection_type"); got != "" {
				t.Fatalf("removed image selection was not cleared: %v", got)
			}
			if got := metadata.ResourceData.Get("update_strategy_id"); got != "" {
				t.Fatalf("removed update strategy was not cleared: %v", got)
			}
			if got := metadata.ResourceData.Get("enabled"); got != true {
				t.Fatalf("omitted disabled property must leave the profile enabled: %v", got)
			}
		})
	}
}
