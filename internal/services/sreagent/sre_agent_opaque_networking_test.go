// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestSreAgentOpaqueDNSList(t *testing.T) {
	for _, member := range []string{
		`,"vnetConfiguration":{"usePrivateDnsResolution":true}`,
		`,"vnetConfiguration":{"usePrivateDnsResolution":false}`,
		`,"vnetConfiguration":{"usePrivateDnsResolution":null}`,
		`,"vnetConfiguration":{}`,
		`,"vnetConfiguration":null`,
		"",
	} {
		t.Run(member, func(t *testing.T) {
			calls := 0
			id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
			md, _, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				calls++
				if req.Method != http.MethodGet {
					return nil, fmt.Errorf("unexpected method %s", req.Method)
				}
				body := fmt.Sprintf(`{"value":[{"id":%q,"name":"agent","location":"eastus","identity":{"type":"SystemAssigned"},"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"%s}}}}]}`, id.ID(), testSreSubnetID, member)
				return sreAgentHTTPResponse(req, http.StatusOK, body), nil
			})
			server, request := sreAgentListTestServer(t, ctx, md)
			request.IncludeResource = true
			stream, err := server.ListResource(ctx, request)
			if err != nil {
				t.Fatal(err)
			}
			results := 0
			for result := range stream.Results {
				results++
				if len(result.Diagnostics) != 0 || result.Resource == nil || result.Identity == nil {
					t.Fatalf("incomplete opaque-field List result: %#v", result)
				}
				state, err := result.Resource.Unmarshal(sdk.WrappedResource(SreAgentResource{}).ProtoSchema(ctx)().ValueType())
				if err != nil {
					t.Fatal(err)
				}
				path := tftypes.NewAttributePath().WithAttributeName("networking").WithElementKeyInt(0)
				raw, _, err := tftypes.WalkAttributePath(state, path)
				if err != nil {
					t.Fatal(err)
				}
				value, ok := raw.(tftypes.Value)
				if !ok {
					t.Fatalf("unexpected networking value %T", raw)
				}
				var fields map[string]tftypes.Value
				if err := value.As(&fields); err != nil {
					t.Fatal(err)
				}
				if len(fields) != 2 {
					t.Fatal("List exposed unowned DNS as Terraform state/control")
				}
				assertSreAgentListValue(t, state, path.WithAttributeName("egress_mode"), tftypes.NewValue(tftypes.String, "AzureVNet"))
			}
			if calls != 1 || results != 1 {
				t.Fatalf("expected one result without per-agent GETs: calls=%d results=%d", calls, results)
			}
		})
	}
}

func TestSreAgentOpaqueDNSIdentityDetach(t *testing.T) {
	r := SreAgentResource{}
	a, b := testIdentityID, testIdentityID+"-b"
	const dns = `"vnetConfiguration":{"usePrivateDnsResolution":false,"futureSetting":"preserve"}`
	var initial agents.Agent
	if err := json.Unmarshal([]byte(sreAgentOpaqueDNSResponse("AzureVNet", testSreSubnetID, dns)), &initial); err != nil {
		t.Fatal(err)
	}
	initial.Identity = testSreAgentIdentity(identity.TypeUserAssigned, []string{a, b})
	responseBody := fmt.Sprintf(`{"location":"eastus","identity":{"type":"UserAssigned","userAssignedIdentities":{%q:{}}},"properties":{"provisioningState":"Succeeded","vnetConfiguration":{"subnetResourceId":null},"sandboxConfiguration":{"egress":{"mode":"Limited",%s}}}}`, b, dns)
	patches := 0
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodPatch {
			patches++
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			assertJSONEqual(t, body, fmt.Sprintf(`{"identity":{"type":"UserAssigned","userAssignedIdentities":{%q:null,%q:{}}},"properties":{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":"Limited"}}}}`, a, b))
		} else if req.Method != http.MethodGet {
			return nil, fmt.Errorf("unexpected method %s", req.Method)
		}
		return sreAgentHTTPResponse(req, http.StatusOK, responseBody), nil
	})
	md.SetID(id)
	if err := r.flatten(md, &id, &initial); err != nil {
		t.Fatal(err)
	}
	config := sreAgentNetworkBaseConfig()
	config["identity"] = []interface{}{map[string]interface{}{"type": "UserAssigned", "identity_ids": []interface{}{b}}}
	config["networking"] = sreAgentNetworkConfig("Limited", "")
	wrapped := sdk.WrappedResource(r)
	diff, err := wrapped.Diff(ctx, md.ResourceData.State(), terraform.NewResourceConfigRaw(config), md.Client)
	if err != nil {
		t.Fatal(err)
	}
	state, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
	if diags.HasError() || patches != 1 || state.Attributes["identity.0.identity_ids.#"] != "1" || state.Attributes["networking.0.egress_mode"] != "Limited" {
		t.Fatalf("identity/detach with unowned DNS failed: patches=%d diagnostics=%v state=%#v", patches, diags, state)
	}
}
