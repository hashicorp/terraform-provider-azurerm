// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestSreAgentStableAPIResourcePresence(t *testing.T) {
	for _, exists := range []bool{false, true} {
		t.Run(fmt.Sprintf("exists=%t", exists), func(t *testing.T) {
			gets, puts := 0, 0
			created := false
			md, id, ctx := sreAgentTestMetadata(t, func(request *http.Request) (*http.Response, error) {
				assertSreAgentRequestVersion(t, request, sreAgentExpectedStableAPIVersion)
				switch request.Method {
				case http.MethodGet:
					gets++
					if !exists && !created {
						return sreAgentHTTPResponse(request, http.StatusNotFound, `{"error":{"code":"ResourceNotFound","message":"absent"}}`), nil
					}
					return sreAgentHTTPResponse(request, http.StatusOK, `{"location":"eastus","properties":{"provisioningState":"Succeeded"}}`), nil
				case http.MethodPut:
					puts++
					created = true
					return sreAgentHTTPResponse(request, http.StatusCreated, `{"location":"eastus","properties":{"provisioningState":"Succeeded"}}`), nil
				default:
					return nil, fmt.Errorf("unexpected resource request %s", request.Method)
				}
			})
			for key, value := range map[string]interface{}{
				"name": id.AgentName, "resource_group_name": id.ResourceGroupName, "location": "eastus",
				"identity": sreAgentSystemIdentityConfig(),
			} {
				if err := md.ResourceData.Set(key, value); err != nil {
					t.Fatal(err)
				}
			}
			err := (SreAgentResource{}).Create().Func(ctx, md)
			if exists {
				if err == nil || !strings.Contains(err.Error(), "import") || puts != 0 || gets != 1 {
					t.Fatalf("existing resource must require import: gets=%d puts=%d err=%v", gets, puts, err)
				}
				return
			}
			if err != nil || puts != 1 || gets < 1 || md.ResourceData.Id() != id.ID() {
				t.Fatalf("create must preserve presence check and ID: gets=%d puts=%d err=%v", gets, puts, err)
			}
		})
	}
}

func TestSreAgentStableAPIImporter(t *testing.T) {
	for _, valid := range []bool{false, true} {
		t.Run(fmt.Sprintf("valid_id=%t", valid), func(t *testing.T) {
			gets := 0
			md, id, ctx := sreAgentTestMetadata(t, func(request *http.Request) (*http.Response, error) {
				gets++
				assertSreAgentRequestVersion(t, request, sreAgentExpectedStableAPIVersion)
				if request.Method != http.MethodGet {
					return nil, fmt.Errorf("import refresh must only read, got %s", request.Method)
				}
				return sreAgentHTTPResponse(request, http.StatusOK, `{"name":"agent","location":"eastus","identity":{"type":"SystemAssigned"},"properties":{"provisioningState":"Succeeded"}}`), nil
			})
			importID := id.ID()
			if !valid {
				importID = strings.Replace(importID, "/agents/", "/notAgents/", 1)
			}
			md.ResourceData.SetId(importID)
			wrapped := sdk.WrappedResource(SreAgentResource{})
			imported, err := wrapped.Importer.StateContext(ctx, md.ResourceData, md.Client)
			if !valid {
				if err == nil || gets != 0 {
					t.Fatalf("invalid resource ID must fail before reading: requests=%d err=%v", gets, err)
				}
				return
			}
			if err != nil || len(imported) != 1 || imported[0].Id() != id.ID() || gets != 0 {
				t.Fatalf("import must return the validated ID before refresh: requests=%d err=%v", gets, err)
			}
			refreshed, diags := wrapped.RefreshWithoutUpgrade(ctx, imported[0].State(), md.Client)
			if diags.HasError() || refreshed == nil || refreshed.ID != id.ID() || gets != 1 {
				t.Fatalf("imported state refresh failed: requests=%d diagnostics=%v", gets, diags)
			}
			if refreshed.Attributes["name"] != id.AgentName || refreshed.Attributes["resource_group_name"] != id.ResourceGroupName {
				t.Fatalf("imported resource identity fields differ: %#v", refreshed.Attributes)
			}
		})
	}
}

func TestSreAgentStableAPILegacyState(t *testing.T) {
	gets := 0
	md, id, ctx := sreAgentTestMetadata(t, func(request *http.Request) (*http.Response, error) {
		gets++
		assertSreAgentRequestVersion(t, request, sreAgentExpectedStableAPIVersion)
		if request.Method != http.MethodGet {
			return nil, fmt.Errorf("state refresh must only read, got %s", request.Method)
		}
		body := fmt.Sprintf(`{"name":"agent","location":"eastus","identity":{"type":"SystemAssigned","principalId":"11111111-1111-1111-1111-111111111111","tenantId":"22222222-2222-2222-2222-222222222222"},"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet","vnetConfiguration":{"usePrivateDnsResolution":true}}},"provisioningState":"Succeeded"}}`, testSreSubnetID)
		return sreAgentHTTPResponse(request, http.StatusOK, body), nil
	})
	// Synthetic state from the January experiment; no retained fixture or state file is used.
	legacy := &terraform.InstanceState{
		ID: id.ID(),
		Attributes: map[string]string{
			"id": id.ID(), "name": id.AgentName, "resource_group_name": id.ResourceGroupName, "location": "eastus",
			"identity.#": "1", "identity.0.type": "SystemAssigned", "identity.0.identity_ids.#": "0",
			"identity.0.principal_id": "11111111-1111-1111-1111-111111111111",
			"identity.0.tenant_id":    "22222222-2222-2222-2222-222222222222",
			"networking.#":            "1", "networking.0.egress_mode": "AzureVNet", "networking.0.subnet_id": testSreSubnetID,
		},
		Meta: map[string]interface{}{"schema_version": "0"},
	}
	expectedAttributes := make(map[string]string, len(legacy.Attributes))
	for key, value := range legacy.Attributes {
		expectedAttributes[key] = value
	}
	wrapped := sdk.WrappedResource(SreAgentResource{})
	network := wrapped.Schema["networking"]
	if wrapped.SchemaVersion != 0 || network.ForceNew || network.Elem.(*pluginsdk.Resource).Schema["subnet_id"].ForceNew {
		t.Fatal("API version transition must preserve state schema and in-place networking")
	}
	refreshed, diags := wrapped.RefreshWithoutUpgrade(ctx, legacy, md.Client)
	if diags.HasError() || refreshed == nil || refreshed.ID != legacy.ID || gets != 1 {
		t.Fatalf("legacy state refresh failed: requests=%d diagnostics=%v", gets, diags)
	}
	for key, expected := range expectedAttributes {
		if refreshed.Attributes[key] != expected {
			t.Fatalf("legacy attribute %q = %q, want %q", key, refreshed.Attributes[key], expected)
		}
	}
	if _, exposed := network.Elem.(*pluginsdk.Resource).Schema["private_dns"]; exposed {
		t.Fatal("legacy response DNS must remain outside Terraform ownership")
	}
	diff, err := wrapped.Diff(ctx, refreshed, terraform.NewResourceConfigRaw(map[string]interface{}{
		"name": id.AgentName, "resource_group_name": id.ResourceGroupName, "location": "eastus",
		"identity": sreAgentSystemIdentityConfig(), "networking": sreAgentNetworkConfig("AzureVNet", testSreSubnetID),
	}), md.Client)
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil && !diff.Empty() {
		t.Fatalf("version-only state refresh must not propose changes or replacement: %#v", diff.Attributes)
	}
}
