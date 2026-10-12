// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	containersclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

func TestKubernetesClusterUpdateTags(t *testing.T) {
	testCases := []struct {
		name                 string
		newTags              map[string]string
		disableLocalAccounts bool
		method               string
		errorContains        string
	}{
		{
			name:          "change tag",
			newTags:       map[string]string{"phase": "updated"},
			method:        http.MethodPatch,
			errorContains: "updating tags for",
		},
		{
			name:          "remove all tags",
			newTags:       map[string]string{},
			method:        http.MethodPatch,
			errorContains: "updating tags for",
		},
		{
			name:                 "tags and cluster property",
			newTags:              map[string]string{"phase": "updated"},
			disableLocalAccounts: true,
			method:               http.MethodPut,
			errorContains:        "updating ",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			id := commonids.NewKubernetesClusterID("00000000-0000-0000-0000-000000000000", "test", "test")
			var writes atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.URL.Path != id.ID() {
					t.Errorf("expected path %q, got %q", id.ID(), req.URL.Path)
				}
				if req.Method == http.MethodGet {
					if err := json.NewEncoder(w).Encode(map[string]any{
						"location":   "eastus",
						"properties": map[string]any{},
						"tags":       map[string]string{"phase": "initial"},
					}); err != nil {
						t.Errorf("encoding cluster: %v", err)
					}
					return
				}

				writes.Add(1)
				if req.Method != testCase.method {
					t.Errorf("expected %s, got %s", testCase.method, req.Method)
				}
				var payload map[string]json.RawMessage
				if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
					t.Errorf("decoding update: %v", err)
				}
				var actualTags map[string]string
				if err := json.Unmarshal(payload["tags"], &actualTags); err != nil {
					t.Errorf("decoding tags: %v", err)
				}
				if !reflect.DeepEqual(actualTags, testCase.newTags) {
					t.Errorf("expected tags %#v, got %#v", testCase.newTags, actualTags)
				}
				if testCase.method == http.MethodPatch {
					if len(payload) != 1 {
						t.Errorf("PATCH must contain only tags, got %v", payload)
					}
				} else {
					var properties managedclusters.ManagedClusterProperties
					if err := json.Unmarshal(payload["properties"], &properties); err != nil {
						t.Errorf("decoding properties: %v", err)
					}
					if properties.DisableLocalAccounts == nil || !*properties.DisableLocalAccounts {
						t.Error("PUT did not include the cluster property change")
					}
				}

				// Stop before the final read and verify update errors are propagated.
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]string{"code": "TestFailure", "message": "sentinel update failure"},
				}); err != nil {
					t.Errorf("encoding update error: %v", err)
				}
			}))
			defer server.Close()

			client, err := managedclusters.NewManagedClustersClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
			if err != nil {
				t.Fatalf("building client: %v", err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true

			resource := resourceKubernetesCluster()
			prior := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{
				"tags": map[string]any{"phase": "initial"},
			})
			prior.SetId(id.ID())
			newValue, exists := testCase.newTags["phase"]
			diff := &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{
				"tags.phase": {Old: "initial", New: newValue, NewRemoved: !exists},
			}}
			if !exists {
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{
					Old: "1", New: "0",
				}
			}
			if testCase.disableLocalAccounts {
				diff.Attributes["local_account_disabled"] = &terraform.ResourceAttrDiff{Old: "false", New: "true"}
			}
			data, err := schema.InternalMap(resource.Schema).Data(prior.State(), diff)
			if err != nil {
				t.Fatalf("building resource data: %v", err)
			}
			err = resourceKubernetesClusterUpdate(data, &clients.Client{
				StopContext: t.Context(),
				Containers:  &containersclient.Client{KubernetesClustersClient: client},
			})
			if err == nil || !strings.Contains(err.Error(), testCase.errorContains) ||
				!strings.Contains(err.Error(), "sentinel update failure") {
				t.Errorf("expected propagated update error containing %q, got %v", testCase.errorContains, err)
			}
			if actual := writes.Load(); actual != 1 {
				t.Errorf("expected exactly one update, got %d", actual)
			}
		})
	}
}
