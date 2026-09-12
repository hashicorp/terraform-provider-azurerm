// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func requireSreAgentPrivateDNS(t *testing.T) {
	t.Helper()
	fields := SreAgentResource{}.Arguments()["networking"].Elem.(*pluginsdk.Resource).Schema
	dns := fields["private_dns"]
	if dns == nil {
		t.Fatal("missing public43a private_dns schema")
	}
	if !dns.Optional || !dns.Computed || dns.MaxItems != 1 || dns.ForceNew {
		t.Fatal("private_dns must preserve omitted values as an optional/computed singleton")
	}
	enabled := dns.Elem.(*pluginsdk.Resource).Schema["enabled"]
	if enabled.Type != pluginsdk.TypeBool || !enabled.Required || enabled.Default != nil {
		t.Fatal("DNS bool must be explicit with no default")
	}
}

func sreAgentPrivateDNSResponse(mode, subnet string, value *bool) string {
	dns := `"vnetConfiguration":{"usePrivateDnsResolution":null}`
	if value != nil {
		dns = fmt.Sprintf(`"vnetConfiguration":{"usePrivateDnsResolution":%t}`, *value)
	}
	modeMember := ""
	if mode != "" {
		modeMember = fmt.Sprintf(`"mode":%q,`, mode)
	}
	subnetJSON := "null"
	if subnet != "" {
		subnetJSON = fmt.Sprintf("%q", subnet)
	}
	return fmt.Sprintf(`{"location":"eastus","identity":{"type":"SystemAssigned"},"properties":{"provisioningState":"Succeeded","vnetConfiguration":{"subnetResourceId":%s},"sandboxConfiguration":{"egress":{%s%s}}}}`, subnetJSON, modeMember, dns)
}

func TestSreAgentPrivateDNSReadPresence(t *testing.T) {
	requireSreAgentPrivateDNS(t)
	for _, tc := range []struct {
		name, mode, subnet string
		value              *bool
	}{
		{"true", "AzureVNet", testSreSubnetID, pointer.To(true)},
		{"false", "AzureVNet", testSreSubnetID, pointer.To(false)},
		{"null", "AzureVNet", testSreSubnetID, nil},
		{"mode-absent", "", "", pointer.To(true)},
		{"mode-future", "FutureMode", "", pointer.To(false)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := SreAgentResource{}
			md := sdk.NewResourceMetaData(nil, r)
			id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
			md.SetID(id)
			var model agents.Agent
			if err := json.Unmarshal([]byte(sreAgentPrivateDNSResponse(tc.mode, tc.subnet, tc.value)), &model); err != nil {
				t.Fatal(err)
			}
			if err := r.flatten(md, &id, &model); err != nil {
				t.Fatal(err)
			}
			wantCount := 0
			if tc.value != nil {
				wantCount = 1
				if md.ResourceData.Get("networking.0.private_dns.0.enabled") != *tc.value {
					t.Fatal("explicit DNS true/false was not preserved")
				}
			}
			if md.ResourceData.Get("networking.0.private_dns.#") != wantCount {
				t.Fatal("null DNS must not fabricate a false block")
			}
		})
	}
}

func TestSreAgentPrivateDNSCreateWire(t *testing.T) {
	requireSreAgentPrivateDNS(t)
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			r := SreAgentResource{}
			md := sdk.NewResourceMetaData(nil, r)
			config := sreAgentNetworkBaseConfig()
			network := sreAgentNetworkConfig("AzureVNet", testSreSubnetID)
			network[0].(map[string]interface{})["private_dns"] = []interface{}{map[string]interface{}{"enabled": enabled}}
			config["networking"] = network
			for key, value := range config {
				if err := md.ResourceData.Set(key, value); err != nil {
					t.Fatal(err)
				}
			}
			var model SreAgentModel
			if err := md.Decode(&model); err != nil {
				t.Fatal(err)
			}
			payload, err := r.expandCreate(model)
			if err != nil {
				t.Fatal(err)
			}
			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, body, fmt.Sprintf(`{"location":"eastus","identity":{"type":"SystemAssigned","userAssignedIdentities":null},"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet","vnetConfiguration":{"usePrivateDnsResolution":%t}}}}}`, testSreSubnetID, enabled))
		})
	}
}

func TestSreAgentPrivateDNSSelectiveUpdate(t *testing.T) {
	requireSreAgentPrivateDNS(t)
	for _, tc := range []struct {
		name, mode, subnet, want string
		before, configured       *bool
	}{
		{"disable", "AzureVNet", testSreSubnetID, `{"properties":{"sandboxConfiguration":{"egress":{"vnetConfiguration":{"usePrivateDnsResolution":false}}}}}`, pointer.To(true), pointer.To(false)},
		{"enable", "AzureVNet", testSreSubnetID, `{"properties":{"sandboxConfiguration":{"egress":{"vnetConfiguration":{"usePrivateDnsResolution":true}}}}}`, pointer.To(false), pointer.To(true)},
		{"unset-to-false", "AzureVNet", testSreSubnetID, `{"properties":{"sandboxConfiguration":{"egress":{"vnetConfiguration":{"usePrivateDnsResolution":false}}}}}`, nil, pointer.To(false)},
		{"subnet-preserves-dns", "AzureVNet", testSreSubnetID + "-b", fmt.Sprintf(`{"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"}}}}`, testSreSubnetID+"-b"), pointer.To(true), nil},
		{"detach-preserves-dns", "Limited", "", `{"properties":{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":"Limited"}}}}`, pointer.To(false), nil},
		{"detach-and-enable", "Limited", "", `{"properties":{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":"Limited","vnetConfiguration":{"usePrivateDnsResolution":true}}}}}`, pointer.To(false), pointer.To(true)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			patches := 0
			after := tc.configured
			if after == nil {
				after = tc.before
			}
			md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodPatch {
					patches++
					body, err := io.ReadAll(req.Body)
					if err != nil {
						return nil, err
					}
					assertJSONEqual(t, body, tc.want)
				} else if req.Method != http.MethodGet {
					return nil, fmt.Errorf("unexpected method %s", req.Method)
				}
				return sreAgentHTTPResponse(req, http.StatusOK, sreAgentPrivateDNSResponse(tc.mode, tc.subnet, after)), nil
			})
			r := SreAgentResource{}
			md.SetID(id)
			var initial agents.Agent
			if err := json.Unmarshal([]byte(sreAgentPrivateDNSResponse("AzureVNet", testSreSubnetID, tc.before)), &initial); err != nil {
				t.Fatal(err)
			}
			if err := r.flatten(md, &id, &initial); err != nil {
				t.Fatal(err)
			}
			config := sreAgentNetworkBaseConfig()
			network := sreAgentNetworkConfig(tc.mode, tc.subnet)
			if tc.configured != nil {
				network[0].(map[string]interface{})["private_dns"] = []interface{}{map[string]interface{}{"enabled": *tc.configured}}
			}
			config["networking"] = network
			wrapped := sdk.WrappedResource(r)
			diff, err := wrapped.Diff(ctx, md.ResourceData.State(), terraform.NewResourceConfigRaw(config), md.Client)
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil || diff.RequiresNew() {
				t.Fatal("DNS changes must plan an in-place update")
			}
			state, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
			if diags.HasError() || patches != 1 {
				t.Fatalf("DNS update failed: patches=%d diagnostics=%v", patches, diags)
			}
			if state.Attributes["networking.0.private_dns.0.enabled"] != fmt.Sprint(*after) {
				t.Fatal("DNS update readback mismatch")
			}
		})
	}
}

func TestSreAgentPrivateDNSOmissionAndUnknown(t *testing.T) {
	requireSreAgentPrivateDNS(t)
	r := SreAgentResource{}
	wrapped := sdk.WrappedResource(r)
	md := sdk.NewResourceMetaData(nil, r)
	id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
	md.SetID(id)
	var initial agents.Agent
	if err := json.Unmarshal([]byte(sreAgentPrivateDNSResponse("AzureVNet", testSreSubnetID, pointer.To(true))), &initial); err != nil {
		t.Fatal(err)
	}
	if err := r.flatten(md, &id, &initial); err != nil {
		t.Fatal(err)
	}
	config := sreAgentNetworkBaseConfig()
	config["networking"] = sreAgentNetworkConfig("AzureVNet", testSreSubnetID)
	diff, err := wrapped.Diff(context.Background(), md.ResourceData.State(), terraform.NewResourceConfigRaw(config), &clients.Client{})
	if err != nil || (diff != nil && len(diff.Attributes) > 0) {
		t.Fatalf("omitted DNS must preserve state without drift: diff=%#v err=%v", diff, err)
	}
	config["networking"].([]interface{})[0].(map[string]interface{})["private_dns"] = []interface{}{}
	if _, err := wrapped.Diff(context.Background(), md.ResourceData.State(), terraform.NewResourceConfigRaw(config), &clients.Client{}); err == nil {
		t.Fatal("raw DNS block removal must not silently clear the API flag")
	}
	for _, subnet := range []cty.Value{cty.UnknownVal(cty.String), cty.StringVal(testSreSubnetID)} {
		values := protocolValues(wrapped)
		values["networking"] = cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"egress_mode": cty.StringVal("AzureVNet"),
			"subnet_id":   subnet,
			"private_dns": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"enabled": cty.UnknownVal(cty.Bool),
			})}),
		})})
		if _, err := wrapped.Diff(context.Background(), nil, terraform.NewResourceConfigShimmed(cty.ObjectVal(values), wrapped.CoreConfigSchema()), &clients.Client{}); err != nil {
			t.Fatalf("unknown DNS with subnet known=%t must be plannable: %v", subnet.IsKnown(), err)
		}
	}
}

func TestSreAgentPrivateDNSReadNullObjects(t *testing.T) {
	requireSreAgentPrivateDNS(t)
	for _, sandbox := range []string{
		`null`, `{}`, `{"egress":null}`, `{"egress":{"mode":"AzureVNet"}}`,
		`{"egress":{"mode":"AzureVNet","vnetConfiguration":null}}`,
		`{"egress":{"mode":"AzureVNet","vnetConfiguration":{}}}`,
		`{"egress":{"mode":"AzureVNet","vnetConfiguration":{"usePrivateDnsResolution":null}}}`,
	} {
		t.Run(sandbox, func(t *testing.T) {
			r := SreAgentResource{}
			md := sdk.NewResourceMetaData(nil, r)
			id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
			md.SetID(id)
			for _, body := range []string{
				sreAgentPrivateDNSResponse("AzureVNet", testSreSubnetID, pointer.To(true)),
				fmt.Sprintf(`{"location":"eastus","identity":{"type":"SystemAssigned"},"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":%s}}`, testSreSubnetID, sandbox),
			} {
				var model agents.Agent
				if err := json.Unmarshal([]byte(body), &model); err != nil {
					t.Fatal(err)
				}
				if err := r.flatten(md, &id, &model); err != nil {
					t.Fatal(err)
				}
			}
			if md.ResourceData.Get("networking.0.private_dns.#") != 0 {
				t.Fatal("missing/null DNS must clear stale readback without fabricating false")
			}
		})
	}
}

func TestSreAgentPrivateDNSListFullState(t *testing.T) {
	requireSreAgentPrivateDNS(t)
	for _, tc := range []struct {
		name, member string
		want         *bool
	}{
		{"true", `,"vnetConfiguration":{"usePrivateDnsResolution":true}`, pointer.To(true)},
		{"false", `,"vnetConfiguration":{"usePrivateDnsResolution":false}`, pointer.To(false)},
		{"null", `,"vnetConfiguration":{"usePrivateDnsResolution":null}`, nil},
		{"missing-flag", `,"vnetConfiguration":{}`, nil},
		{"null-object", `,"vnetConfiguration":null`, nil},
		{"absent-object", "", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
			md, _, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				calls++
				if req.Method != http.MethodGet {
					return nil, fmt.Errorf("unexpected method %s", req.Method)
				}
				body := fmt.Sprintf(`{"value":[{"id":%q,"name":"agent","location":"eastus","identity":{"type":"SystemAssigned"},"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"%s}}}}]}`, id.ID(), testSreSubnetID, tc.member)
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
					t.Fatalf("incomplete DNS List result: %#v", result)
				}
				state, err := result.Resource.Unmarshal(sdk.WrappedResource(SreAgentResource{}).ProtoSchema(ctx)().ValueType())
				if err != nil {
					t.Fatal(err)
				}
				path := tftypes.NewAttributePath().WithAttributeName("networking").WithElementKeyInt(0).WithAttributeName("private_dns")
				if tc.want != nil {
					assertSreAgentListValue(t, state, path.WithElementKeyInt(0).WithAttributeName("enabled"), tftypes.NewValue(tftypes.Bool, *tc.want))
				} else {
					raw, _, err := tftypes.WalkAttributePath(state, path)
					if err != nil {
						t.Fatal(err)
					}
					value, ok := raw.(tftypes.Value)
					if !ok {
						t.Fatalf("unexpected DNS state type: %T", raw)
					}
					var blocks []tftypes.Value
					if !value.IsNull() {
						if err := value.As(&blocks); err != nil {
							t.Fatal(err)
						}
					}
					if len(blocks) != 0 {
						t.Fatal("List fabricated a DNS block for a missing/null flag")
					}
				}
			}
			if calls != 1 || results != 1 {
				t.Fatalf("expected one result without per-agent GETs: calls=%d results=%d", calls, results)
			}
		})
	}
}

func TestSreAgentPrivateDNSComposesIdentityAndDetach(t *testing.T) {
	requireSreAgentPrivateDNS(t)
	r := SreAgentResource{}
	a, b := testIdentityID, testIdentityID+"-b"
	var initial, after agents.Agent
	if err := json.Unmarshal([]byte(sreAgentPrivateDNSResponse("AzureVNet", testSreSubnetID, pointer.To(true))), &initial); err != nil {
		t.Fatal(err)
	}
	initial.Identity = testSreAgentIdentity(identity.TypeUserAssigned, []string{a, b})
	if err := json.Unmarshal([]byte(sreAgentPrivateDNSResponse("Limited", "", pointer.To(false))), &after); err != nil {
		t.Fatal(err)
	}
	after.Identity = testSreAgentIdentity(identity.TypeUserAssigned, []string{b})
	responseBody, err := json.Marshal(after)
	if err != nil {
		t.Fatal(err)
	}
	patches := 0
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodPatch {
			patches++
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			assertJSONEqual(t, body, fmt.Sprintf(`{"identity":{"type":"UserAssigned","userAssignedIdentities":{%q:null,%q:{}}},"properties":{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":"Limited","vnetConfiguration":{"usePrivateDnsResolution":false}}}}}`, a, b))
		} else if req.Method != http.MethodGet {
			return nil, fmt.Errorf("unexpected method %s", req.Method)
		}
		return sreAgentHTTPResponse(req, http.StatusOK, string(responseBody)), nil
	})
	md.SetID(id)
	if err := r.flatten(md, &id, &initial); err != nil {
		t.Fatal(err)
	}
	config := sreAgentNetworkBaseConfig()
	config["identity"] = []interface{}{map[string]interface{}{"type": "UserAssigned", "identity_ids": []interface{}{b}}}
	network := sreAgentNetworkConfig("Limited", "")
	network[0].(map[string]interface{})["private_dns"] = []interface{}{map[string]interface{}{"enabled": false}}
	config["networking"] = network
	wrapped := sdk.WrappedResource(r)
	diff, err := wrapped.Diff(ctx, md.ResourceData.State(), terraform.NewResourceConfigRaw(config), md.Client)
	if err != nil {
		t.Fatal(err)
	}
	state, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
	if diags.HasError() || patches != 1 || state.Attributes["identity.0.identity_ids.#"] != "1" ||
		state.Attributes["networking.0.private_dns.0.enabled"] != "false" || state.Attributes["networking.0.egress_mode"] != "Limited" {
		t.Fatalf("combined DNS/identity/detach failed: patches=%d diagnostics=%v state=%#v", patches, diags, state)
	}
}
