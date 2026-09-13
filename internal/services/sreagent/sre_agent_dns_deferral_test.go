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
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestSreAgentDNSDeferredSchema(t *testing.T) {
	fields := SreAgentResource{}.Arguments()["networking"].Elem.(*pluginsdk.Resource).Schema
	if len(fields) != 2 || fields["private_dns"] != nil || fields["use_vnet_dns"] != nil {
		t.Fatal("networking must expose only egress_mode and subnet_id; DNS control is deferred")
	}
}

func sreAgentOpaqueDNSResponse(mode, subnet, dnsMember string) string {
	egress := dnsMember
	if mode != "" {
		egress = fmt.Sprintf(`"mode":%q,`, mode) + egress
	}
	vnet := "null"
	if subnet != "" {
		vnet = fmt.Sprintf(`{"subnetResourceId":%q}`, subnet)
	}
	return fmt.Sprintf(`{"location":"eastus","identity":{"type":"SystemAssigned"},"properties":{"provisioningState":"Succeeded","vnetConfiguration":%s,"sandboxConfiguration":{"egress":{%s}}}}`, vnet, egress)
}

func TestSreAgentDNSDeferredReadCreate(t *testing.T) {
	for _, dns := range []string{
		`"vnetConfiguration":{"usePrivateDnsResolution":true}`,
		`"vnetConfiguration":{"usePrivateDnsResolution":false}`,
		`"vnetConfiguration":{"usePrivateDnsResolution":null}`,
		`"vnetConfiguration":{}`,
		`"vnetConfiguration":null`,
		`"unmodeledSetting":{"preserve":true}`,
	} {
		t.Run(dns, func(t *testing.T) {
			r := SreAgentResource{}
			md := sdk.NewResourceMetaData(nil, r)
			id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
			md.SetID(id)
			var remote agents.Agent
			if err := json.Unmarshal([]byte(sreAgentOpaqueDNSResponse("AzureVNet", testSreSubnetID, dns)), &remote); err != nil {
				t.Fatal(err)
			}
			if err := r.flatten(md, &id, &remote); err != nil {
				t.Fatal(err)
			}
			network := md.ResourceData.Get("networking").([]interface{})[0].(map[string]interface{})
			if len(network) != 2 {
				t.Fatalf("unowned DNS became Terraform state/control: %#v", network)
			}
			wrapped := sdk.WrappedResource(r)
			config := sreAgentNetworkBaseConfig()
			config["networking"] = sreAgentNetworkConfig("AzureVNet", testSreSubnetID)
			diff, err := wrapped.Diff(context.Background(), md.ResourceData.State(), terraform.NewResourceConfigRaw(config), &clients.Client{})
			if err != nil || (diff != nil && (len(diff.Attributes) != 0 || diff.RequiresNew())) {
				t.Fatalf("unowned DNS must not cause drift: diff=%#v error=%v", diff, err)
			}
			var model SreAgentModel
			if err := md.Decode(&model); err != nil {
				t.Fatal(err)
			}
			payload, err := r.expandCreate(model)
			if err != nil {
				t.Fatal(err)
			}
			body, err := json.Marshal(payload.Properties)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, body, fmt.Sprintf(`{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"}}}`, testSreSubnetID))
		})
	}
}

func TestSreAgentDNSDeferredOnlyRemoteFields(t *testing.T) {
	r := SreAgentResource{}
	md := sdk.NewResourceMetaData(nil, r)
	id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
	md.SetID(id)
	for _, dns := range []string{
		`"vnetConfiguration":{"usePrivateDnsResolution":true}`,
		`"vnetConfiguration":{"usePrivateDnsResolution":false}`,
		`"vnetConfiguration":{"usePrivateDnsResolution":null}`,
	} {
		var remote agents.Agent
		if err := json.Unmarshal([]byte(sreAgentOpaqueDNSResponse("", "", dns)), &remote); err != nil {
			t.Fatal(err)
		}
		if err := r.flatten(md, &id, &remote); err != nil {
			t.Fatal(err)
		}
		if md.ResourceData.Get("networking.#") != 0 {
			t.Fatal("legacy DNS-only response must not fabricate a networking block")
		}
	}
}

func TestSreAgentDNSDeferredSelectiveUpdates(t *testing.T) {
	for _, tc := range []struct {
		name, beforeMode, beforeSubnet, mode, subnet, want string
		tagsOnly                                           bool
	}{
		{"attach", "Limited", "", "AzureVNet", testSreSubnetID, fmt.Sprintf(`{"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"}}}}`, testSreSubnetID), false},
		{"replace", "AzureVNet", testSreSubnetID, "AzureVNet", testSreSubnetID + "-b", fmt.Sprintf(`{"properties":{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":"AzureVNet"}}}}`, testSreSubnetID+"-b"), false},
		{"detach", "AzureVNet", testSreSubnetID, "Limited", "", `{"properties":{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":"Limited"}}}}`, false},
		{"tag-only", "AzureVNet", testSreSubnetID, "AzureVNet", testSreSubnetID, `{"tags":{"stage":"updated"}}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const dns = `"vnetConfiguration":{"usePrivateDnsResolution":true,"futureSetting":"preserve"}`
			patches := 0
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
				return sreAgentHTTPResponse(req, http.StatusOK, sreAgentOpaqueDNSResponse(tc.mode, tc.subnet, dns)), nil
			})
			md.SetID(id)
			r := SreAgentResource{}
			var remote agents.Agent
			if err := json.Unmarshal([]byte(sreAgentOpaqueDNSResponse(tc.beforeMode, tc.beforeSubnet, dns)), &remote); err != nil {
				t.Fatal(err)
			}
			if err := r.flatten(md, &id, &remote); err != nil {
				t.Fatal(err)
			}
			config := sreAgentNetworkBaseConfig()
			if tc.tagsOnly {
				config["tags"] = map[string]interface{}{"stage": "updated"}
			} else {
				config["networking"] = sreAgentNetworkConfig(tc.mode, tc.subnet)
			}
			wrapped := sdk.WrappedResource(r)
			diff, err := wrapped.Diff(ctx, md.ResourceData.State(), terraform.NewResourceConfigRaw(config), md.Client)
			if err != nil {
				t.Fatal(err)
			}
			result, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
			if diags.HasError() || patches != 1 {
				t.Fatalf("selective update failed: patches=%d diagnostics=%v", patches, diags)
			}
			if result.ID != id.ID() || result.Attributes["networking.0.egress_mode"] != tc.mode ||
				result.Attributes["networking.0.subnet_id"] != tc.subnet || !reflect.DeepEqual(result.Identity, md.ResourceData.State().Identity) {
				t.Fatal("network update changed resource identity or lost attachment state")
			}
		})
	}
}
