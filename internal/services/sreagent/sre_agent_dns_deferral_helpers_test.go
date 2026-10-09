// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestSreAgentDNSDeferredHelperSchema(t *testing.T) {
	fields := sreAgentNetworkingSchema().Elem.(*pluginsdk.Resource).Schema
	if len(fields) != 2 || fields["private_dns"] != nil {
		t.Fatal("production networking helper still exposes the deferred DNS control")
	}
}

func TestSreAgentDNSDeferredHelperRead(t *testing.T) {
	for _, flag := range []string{"true", "false", "null"} {
		var props agents.AgentProperties
		if err := json.Unmarshal([]byte(`{"sandboxConfiguration":{"egress":{"vnetConfiguration":{"usePrivateDnsResolution":`+flag+`}}}}`), &props); err != nil {
			t.Fatal(err)
		}
		if got := flattenSreAgentNetworking(&props); len(got) != 0 {
			t.Fatalf("legacy DNS-only fields became Terraform networking state: %#v", got)
		}
	}
}
func TestSreAgentDNSDeferredHelperPatch(t *testing.T) {
	subnet := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Network/virtualNetworks/test/subnets/agent"
	for _, mode := range []string{"AzureVNet", "Limited", "Unrestricted"} {
		network := SreAgentNetworking{EgressMode: mode}
		want := fmt.Sprintf(`{"vnetConfiguration":null,"sandboxConfiguration":{"egress":{"mode":%q}}}`, mode)
		if mode == "AzureVNet" {
			network.SubnetID = subnet
			want = fmt.Sprintf(`{"vnetConfiguration":{"subnetResourceId":%q},"sandboxConfiguration":{"egress":{"mode":%q}}}`, subnet, mode)
		}
		props := &agents.AgentPatchProperties{}
		if err := expandSreAgentNetworking(network, props); err != nil {
			t.Fatal(err)
		}
		got, err := marshalSreAgentPatchProperties(props)
		if err != nil {
			t.Fatal(err)
		}
		var actual, expected interface{}
		if err := json.Unmarshal(got, &actual); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(want), &expected); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("unexpected selective patch: %s", got)
		}
	}
}
