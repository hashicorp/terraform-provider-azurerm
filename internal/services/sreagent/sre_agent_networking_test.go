// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const testSreSubnetID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/test/subnets/agent"

func TestSreAgentNetworkingSchema(t *testing.T) {
	schema := SreAgentResource{}.Arguments()["networking"]
	if schema == nil {
		t.Fatal("missing public-minimum networking block")
	}
	if !schema.Optional || !schema.Computed || schema.ForceNew || schema.MaxItems != 1 {
		t.Fatal("networking must preserve omitted state and support in-place updates")
	}
	fields := schema.Elem.(*pluginsdk.Resource).Schema
	if len(fields) != 3 || fields["private_dns"] == nil || fields["use_vnet_dns"] != nil {
		t.Fatal("only public subnet_id, egress_mode and private_dns belong in this experiment")
	}
	if !fields["egress_mode"].Required || fields["egress_mode"].Default != nil || fields["subnet_id"].ForceNew {
		t.Fatal("egress selection must be explicit and subnet changes must not force replacement")
	}
}

func sreAgentNetworkConfig(mode, subnet string) []interface{} {
	return []interface{}{map[string]interface{}{"egress_mode": mode, "subnet_id": subnet}}
}

func sreAgentNetworkBaseConfig() map[string]interface{} {
	return map[string]interface{}{
		"name": "agent", "resource_group_name": "example", "location": "eastus",
		"identity": sreAgentSystemIdentityConfig(),
	}
}

func TestSreAgentNetworkingCreatePayload(t *testing.T) {
	for _, mode := range []string{"AzureVNet", "Limited", "Unrestricted"} {
		t.Run(mode, func(t *testing.T) {
			r := SreAgentResource{}
			if r.Arguments()["networking"] == nil {
				t.Fatal("missing networking schema")
			}
			subnet := ""
			network := `"sandboxConfiguration":{"egress":{"mode":"` + mode + `"}}`
			if mode == "AzureVNet" {
				subnet = testSreSubnetID
				network += fmt.Sprintf(`,"vnetConfiguration":{"subnetResourceId":%q}`, subnet)
			}
			config := sreAgentNetworkBaseConfig()
			config["networking"] = sreAgentNetworkConfig(mode, subnet)
			md := sdk.NewResourceMetaData(nil, r)
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
			assertJSONEqual(t, body, `{"location":"eastus","identity":{"type":"SystemAssigned","userAssignedIdentities":null},"properties":{`+network+`}}`)
		})
	}
}

func TestSreAgentNetworkingUpdateWire(t *testing.T) {
	for _, scenario := range []struct {
		name, mode, subnet, want string
		omit                     bool
	}{
		{"attach", "AzureVNet", testSreSubnetID, fmt.Sprintf(`{"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"}}}}`, testSreSubnetID), false},
		{"replace", "AzureVNet", testSreSubnetID + "-b", fmt.Sprintf(`{"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"}}}}`, testSreSubnetID+"-b"), false},
		{"detach", "Limited", "", `{"properties":{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":"Limited"}}}}`, false},
		{"tag-only-omission", "", "", `{"tags":{"stage":"updated"}}`, true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			r := SreAgentResource{}
			if r.Arguments()["networking"] == nil {
				t.Fatal("missing networking schema")
			}
			initialMode, initialSubnet := "AzureVNet", testSreSubnetID
			if scenario.name == "attach" {
				initialMode, initialSubnet = "Limited", ""
			}
			mode, subnet := scenario.mode, scenario.subnet
			if scenario.omit {
				mode, subnet = initialMode, initialSubnet
			}
			patches := 0
			md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				switch req.Method {
				case http.MethodPatch:
					patches++
					body, err := io.ReadAll(req.Body)
					if err != nil {
						return nil, err
					}
					assertJSONEqual(t, body, scenario.want)
				case http.MethodGet:
				default:
					return nil, fmt.Errorf("unexpected networking method: %s", req.Method)
				}
				return sreAgentHTTPResponse(req, http.StatusOK, sreAgentNetworkResponse(mode, subnet, "updated")), nil
			})
			md.SetID(id)
			var initial agents.Agent
			if err := json.Unmarshal([]byte(sreAgentNetworkResponse(initialMode, initialSubnet, "")), &initial); err != nil {
				t.Fatal(err)
			}
			if err := r.flatten(md, &id, &initial); err != nil {
				t.Fatal(err)
			}
			config := sreAgentNetworkBaseConfig()
			if scenario.omit {
				config["tags"] = map[string]interface{}{"stage": "updated"}
			} else {
				config["networking"] = sreAgentNetworkConfig(mode, subnet)
			}
			wrapped := sdk.WrappedResource(r)
			diff, err := wrapped.Diff(ctx, md.ResourceData.State(), terraform.NewResourceConfigRaw(config), md.Client)
			if err != nil {
				t.Fatal(err)
			}
			if diff.RequiresNew() {
				t.Fatal("network update must not recreate agent")
			}
			state, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
			if diags.HasError() {
				t.Fatal(diags)
			}
			if patches != 1 || state.Attributes["networking.0.egress_mode"] != mode || state.Attributes["networking.0.subnet_id"] != subnet {
				t.Fatalf("network update/readback mismatch: patches=%d state=%#v", patches, state.Attributes)
			}
		})
	}
}

func sreAgentNetworkResponse(mode, subnet, stage string) string {
	tags := "{}"
	if stage != "" {
		tags = fmt.Sprintf(`{"stage":%q}`, stage)
	}
	return fmt.Sprintf(`{"location":"eastus","identity":{"type":"SystemAssigned"},"tags":%s,"properties":{"provisioningState":"Succeeded","vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":%q,"unowned":{"keep":true}}}}}`, tags, subnet, mode)
}

func TestSreAgentNetworkingPlanValidation(t *testing.T) {
	wrapped := sdk.WrappedResource(SreAgentResource{})
	if wrapped.Schema["networking"] == nil {
		t.Fatal("missing networking schema")
	}
	for _, tc := range []struct {
		mode, subnet string
		valid        bool
	}{
		{"AzureVNet", testSreSubnetID, true}, {"Limited", "", true}, {"Unrestricted", "", true},
		{"AzureVNet", "", false}, {"Limited", testSreSubnetID, false}, {"", testSreSubnetID, false},
	} {
		t.Run(tc.mode+"/"+tc.subnet, func(t *testing.T) {
			config := sreAgentNetworkBaseConfig()
			config["networking"] = sreAgentNetworkConfig(tc.mode, tc.subnet)
			_, err := wrapped.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), &clients.Client{})
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%t error=%v", tc.valid, err)
			}
		})
	}
	values := protocolValues(wrapped)
	values["networking"] = cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
		"egress_mode": cty.StringVal("AzureVNet"), "subnet_id": cty.UnknownVal(cty.String),
		"private_dns": cty.NullVal(wrapped.CoreConfigSchema().ImpliedType().AttributeType("networking").ElementType().AttributeType("private_dns")),
	})})
	if _, err := wrapped.Diff(context.Background(), nil, terraform.NewResourceConfigShimmed(cty.ObjectVal(values), wrapped.CoreConfigSchema()), &clients.Client{}); err != nil {
		t.Fatalf("unknown subnet references must remain plannable: %v", err)
	}
}

func TestSreAgentNetworkingDetachComposesIdentityRemoval(t *testing.T) {
	r := SreAgentResource{}
	if r.Arguments()["networking"] == nil {
		t.Fatal("missing networking schema")
	}
	a, b := testIdentityID, testIdentityID+"-b"
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
			return nil, fmt.Errorf("unexpected request %s", req.Method)
		}
		return sreAgentHTTPResponse(req, http.StatusOK, fmt.Sprintf(`{"location":"eastus","identity":{"type":"UserAssigned","userAssignedIdentities":{%q:{}}},"properties":{"provisioningState":"Succeeded","vnetConfiguration":{"subnetResourceId":null},"sandboxConfiguration":{"egress":{"mode":"Limited"}}}}`, b)), nil
	})
	md.SetID(id)
	var model agents.Agent
	if err := json.Unmarshal([]byte(sreAgentNetworkResponse("AzureVNet", testSreSubnetID, "")), &model); err != nil {
		t.Fatal(err)
	}
	model.Identity = testSreAgentIdentity(identity.TypeUserAssigned, []string{a, b})
	if err := r.flatten(md, &id, &model); err != nil {
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
		t.Fatalf("composed detach/identity update failed: patches=%d diagnostics=%v state=%#v", patches, diags, state)
	}
}

func TestSreAgentNetworkingListFullState(t *testing.T) {
	r := SreAgentResource{}
	if r.Arguments()["networking"] == nil {
		t.Fatal("missing networking schema")
	}
	calls := 0
	id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
	md, _, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		calls++
		if req.Method != http.MethodGet || !strings.HasSuffix(req.URL.Path, "/providers/Microsoft.App/agents") {
			return nil, fmt.Errorf("unexpected List request %s %s", req.Method, req.URL)
		}
		model := sreAgentNetworkResponse("AzureVNet", testSreSubnetID, "")
		model = strings.Replace(model, "{", fmt.Sprintf(`{"id":%q,"name":"agent",`, id.ID()), 1)
		return sreAgentHTTPResponse(req, http.StatusOK, `{"value":[`+model+`]}`), nil
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
			t.Fatalf("incomplete network List result: %#v", result)
		}
		state, err := result.Resource.Unmarshal(sdk.WrappedResource(r).ProtoSchema(ctx)().ValueType())
		if err != nil {
			t.Fatal(err)
		}
		for key, want := range map[string]string{"egress_mode": "AzureVNet", "subnet_id": testSreSubnetID} {
			assertSreAgentListValue(t, state, tftypes.NewAttributePath().WithAttributeName("networking").WithElementKeyInt(0).WithAttributeName(key), tftypes.NewValue(tftypes.String, want))
		}
	}
	if results != 1 || calls != 1 {
		t.Fatalf("wanted one List result without per-agent GET: results=%d calls=%d", results, calls)
	}
}

func TestSreAgentNetworkingPollingFailure(t *testing.T) {
	r := SreAgentResource{}
	if r.Arguments()["networking"] == nil {
		t.Fatal("missing networking schema")
	}
	patches, polls := 0, 0
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == http.MethodPatch:
			patches++
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			assertJSONEqual(t, body, `{"properties":{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":"Limited"}}}}`)
			resp := sreAgentHTTPResponse(req, http.StatusAccepted, `{"properties":{"provisioningState":"Accepted"}}`)
			resp.Header.Set("Azure-AsyncOperation", "https://management.azure.com/operations/network-update?api-version=2026-01-01")
			return resp, nil
		case req.Method == http.MethodGet && req.URL.Path == "/operations/network-update":
			polls++
			return sreAgentHTTPResponse(req, http.StatusOK, `{"status":"Failed","error":{"code":"NetworkUpdateFailed","message":"synthetic network failure"}}`), nil
		case req.Method == http.MethodGet:
			return sreAgentHTTPResponse(req, http.StatusOK, sreAgentNetworkResponse("AzureVNet", testSreSubnetID, "")), nil
		default:
			return nil, fmt.Errorf("unexpected request %s", req.Method)
		}
	})
	md.SetID(id)
	var model agents.Agent
	if err := json.Unmarshal([]byte(sreAgentNetworkResponse("AzureVNet", testSreSubnetID, "")), &model); err != nil {
		t.Fatal(err)
	}
	if err := r.flatten(md, &id, &model); err != nil {
		t.Fatal(err)
	}
	config := sreAgentNetworkBaseConfig()
	config["networking"] = sreAgentNetworkConfig("Limited", "")
	wrapped := sdk.WrappedResource(r)
	diff, err := wrapped.Diff(ctx, md.ResourceData.State(), terraform.NewResourceConfigRaw(config), md.Client)
	if err != nil {
		t.Fatal(err)
	}
	state, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
	if !diags.HasError() || !strings.Contains(fmt.Sprint(diags), "NetworkUpdateFailed") || patches != 1 || polls != 1 {
		t.Fatalf("failed network LRO must surface: patches=%d polls=%d diagnostics=%v", patches, polls, diags)
	}
	if state == nil || state.ID != id.ID() {
		t.Fatal("failed update must retain the resource ID for cleanup")
	}
}

func TestSreAgentNetworkingOmissionNoDrift(t *testing.T) {
	r := SreAgentResource{}
	if r.Arguments()["networking"] == nil {
		t.Fatal("missing networking schema")
	}
	for _, mode := range []string{"AzureVNet", "Limited", "FutureMode"} {
		t.Run(mode, func(t *testing.T) {
			md := sdk.NewResourceMetaData(nil, r)
			id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
			md.SetID(id)
			var model agents.Agent
			subnet := ""
			if mode == "AzureVNet" {
				subnet = testSreSubnetID
			}
			if err := json.Unmarshal([]byte(sreAgentNetworkResponse(mode, subnet, "")), &model); err != nil {
				t.Fatal(err)
			}
			if err := r.flatten(md, &id, &model); err != nil {
				t.Fatal(err)
			}
			wrapped := sdk.WrappedResource(r)
			diff, err := wrapped.Diff(context.Background(), md.ResourceData.State(), terraform.NewResourceConfigRaw(sreAgentNetworkBaseConfig()), &clients.Client{})
			if err != nil {
				t.Fatal(err)
			}
			if diff != nil && (len(diff.Attributes) != 0 || diff.RequiresNew() || diff.Destroy || !reflect.DeepEqual(diff.Identity, md.ResourceData.State().Identity)) {
				t.Fatalf("omitted networking caused drift: %#v", diff)
			}
			config := sreAgentNetworkBaseConfig()
			config["networking"] = []interface{}{}
			_, err = wrapped.Diff(context.Background(), md.ResourceData.State(), terraform.NewResourceConfigRaw(config), &clients.Client{})
			if err == nil || !strings.Contains(err.Error(), "networking") {
				t.Fatalf("raw whole-block reset must not silently detach: %v", err)
			}
		})
	}
}
