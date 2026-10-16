// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesAzurePolicyPayload(t *testing.T) {
	resource := &pluginsdk.Resource{Schema: schemaKubernetesAddOns()}
	for _, enabled := range []bool{false, true} {
		data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{"azure_policy_enabled": enabled})
		input := map[string]any{}
		for key := range resource.Schema {
			input[key] = data.Get(key)
		}
		profiles, err := expandKubernetesAddOns(data, input, *environments.AzurePublic())
		if err != nil {
			t.Fatal(err)
		}
		profile, ok := (*profiles)[azurePolicyKey]
		if !ok || profile.Enabled != enabled {
			t.Fatalf("azure policy enabled = %t, want %t", profile.Enabled, enabled)
		}
		wire, err := json.Marshal(profile)
		if err != nil {
			t.Fatal(err)
		}
		var payload map[string]any
		if err := json.Unmarshal(wire, &payload); err != nil {
			t.Fatal(err)
		}
		if _, ok := payload["config"]; ok {
			t.Fatalf("azure policy must not send a config version: %s", wire)
		}
	}
}
