// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesAzurePolicyPayload(t *testing.T) {
	for _, before := range []bool{false, true} {
		for _, after := range []bool{false, true} {
			t.Run(fmt.Sprintf("%t_to_%t", before, after), func(t *testing.T) {
				resource := &pluginsdk.Resource{Schema: schemaKubernetesAddOns()}
				old := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{"azure_policy_enabled": before})
				old.SetId("existing-cluster")
				config := terraform.NewResourceConfigRaw(map[string]interface{}{"azure_policy_enabled": after})
				diff, err := resource.SimpleDiff(t.Context(), old.State(), config, nil)
				if err != nil {
					t.Fatal(err)
				}
				data, err := schema.InternalMap(resource.Schema).Data(old.State(), diff)
				if err != nil {
					t.Fatal(err)
				}
				input := map[string]interface{}{}
				for key := range resource.Schema {
					input[key] = data.Get(key)
				}
				profiles, err := expandKubernetesAddOns(data, input, *environments.AzurePublic())
				if err != nil {
					t.Fatal(err)
				}
				profile, ok := (*profiles)[azurePolicyKey]
				if !ok || profile.Enabled != after {
					t.Fatalf("policy must be synchronized even without an attribute change: %#v", profiles)
				}
				wire, err := json.Marshal(profile)
				if err != nil {
					t.Fatal(err)
				}
				var payload map[string]interface{}
				if err := json.Unmarshal(wire, &payload); err != nil {
					t.Fatal(err)
				}
				if _, ok := payload["config"]; ok {
					t.Fatalf("policy must not send a hard-coded version: %s", wire)
				}
				if got := flattenKubernetesAddOns(*profiles)["azure_policy_enabled"]; got != after {
					t.Fatalf("policy readback = %v, want %t", got, after)
				}
			})
		}
	}
}
