// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesMessageOfTheDayPlan(t *testing.T) {
	rawConfig := func(nested bool, message *string) cty.Value {
		value := cty.NullVal(cty.String)
		if message != nil {
			value = cty.StringVal(*message)
		}
		if nested {
			return cty.ObjectVal(map[string]cty.Value{"default_node_pool": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"message_of_the_day": value})})})
		}
		return cty.ObjectVal(map[string]cty.Value{"message_of_the_day": value})
	}
	standaloneResource := &schema.Resource{Schema: map[string]*schema.Schema{
		"message_of_the_day": resourceKubernetesClusterNodePoolSchema()["message_of_the_day"],
	}}
	defaultPoolResource := &schema.Resource{Schema: map[string]*schema.Schema{"default_node_pool": {
		Type: schema.TypeList, Required: true, MaxItems: 1,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{"message_of_the_day": SchemaDefaultNodePool().Elem.(*schema.Resource).Schema["message_of_the_day"]}},
	}}}
	tests := []struct {
		name        string
		resource    *schema.Resource
		oldConfig   map[string]any
		newConfig   map[string]any
		rawConfig   cty.Value
		key         string
		wantChange  bool
		wantReplace bool
	}{
		{
			name:      "standalone omitted",
			resource:  standaloneResource,
			oldConfig: map[string]any{"message_of_the_day": "Existing message"},
			newConfig: map[string]any{},
			rawConfig: rawConfig(false, nil),
			key:       "message_of_the_day",
		},
		{
			name:        "standalone explicit empty",
			resource:    standaloneResource,
			oldConfig:   map[string]any{"message_of_the_day": "Existing message"},
			newConfig:   map[string]any{"message_of_the_day": ""},
			rawConfig:   rawConfig(false, pointer.To("")),
			key:         "message_of_the_day",
			wantChange:  true,
			wantReplace: true,
		},
		{
			name:      "default pool omitted",
			resource:  defaultPoolResource,
			oldConfig: map[string]any{"default_node_pool": []any{map[string]any{"message_of_the_day": "Existing message"}}},
			newConfig: map[string]any{"default_node_pool": []any{map[string]any{}}},
			rawConfig: rawConfig(true, nil),
			key:       "default_node_pool.0.message_of_the_day",
		},
		{
			name:       "default pool explicit empty",
			resource:   defaultPoolResource,
			oldConfig:  map[string]any{"default_node_pool": []any{map[string]any{"message_of_the_day": "Existing message"}}},
			newConfig:  map[string]any{"default_node_pool": []any{map[string]any{"message_of_the_day": ""}}},
			rawConfig:  rawConfig(true, pointer.To("")),
			key:        "default_node_pool.0.message_of_the_day",
			wantChange: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createDiff, err := test.resource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(test.oldConfig), nil)
			if err != nil {
				t.Fatal(err)
			}
			attributes, err := createDiff.Apply(nil, test.resource.CoreConfigSchema())
			if err != nil {
				t.Fatal(err)
			}
			state := &terraform.InstanceState{ID: "test", Attributes: attributes, RawConfig: test.rawConfig}
			diff, err := test.resource.Diff(context.Background(), state, terraform.NewResourceConfigRaw(test.newConfig), nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil {
				if test.wantChange || test.wantReplace {
					t.Fatal("expected a message change, got no diff")
				}
				return
			}
			if _, changed := diff.Attributes[test.key]; changed != test.wantChange {
				t.Fatalf("message change = %t, want %t: %#v", changed, test.wantChange, diff.Attributes)
			}
			if got := diff.RequiresNew(); got != test.wantReplace {
				t.Fatalf("replacement = %t, want %t: %#v", got, test.wantReplace, diff.Attributes)
			}
		})
	}
}

func TestKubernetesMessageOfTheDayExpand(t *testing.T) {
	for _, test := range []struct {
		name    string
		message any
		omit    bool
		want    string
	}{
		{name: "omitted", omit: true},
		{name: "null", message: nil},
		{name: "empty", message: ""},
		{name: "text", message: "Welcome\n日本語", want: "Welcome\n日本語"},
		{name: "literal shell", message: "$(hostname)", want: "$(hostname)"},
		{name: "literal base64", message: "V2VsY29tZQ==", want: "V2VsY29tZQ=="},
	} {
		t.Run(test.name, func(t *testing.T) {
			block := map[string]any{"name": "default", "vm_size": "Standard_DS2_v2", "node_count": 1}
			if !test.omit {
				block["message_of_the_day"] = test.message
			}
			data := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, map[string]any{
				"default_node_pool": []any{block},
			})
			profiles, err := ExpandDefaultNodePool(data)
			if err != nil {
				t.Fatal(err)
			}
			profile := (*profiles)[0]
			converted := ConvertDefaultNodePoolToAgentPool(profiles)
			for name, value := range map[string]*string{"cluster request": profile.MessageOfTheDay, "rotation request": converted.Properties.MessageOfTheDay} {
				if test.want == "" {
					if value != nil {
						t.Fatalf("%s must omit empty MOTD, got %q", name, *value)
					}
				} else if value == nil || *value != base64.StdEncoding.EncodeToString([]byte(test.want)) {
					t.Fatalf("%s did not preserve the base64-encoded MOTD: %v", name, value)
				}
			}
			wire, err := json.Marshal(converted)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(wire), `"messageOfTheDay"`) != (test.want != "") {
				t.Fatalf("unexpected MOTD wire presence: %s", wire)
			}
		})
	}
}

func TestKubernetesMessageOfTheDayResponse(t *testing.T) {
	for _, test := range []struct {
		name    string
		encoded *string
		want    string
		wantErr bool
	}{
		{name: "absent"},
		{name: "empty", encoded: pointer.To("")},
		{name: "text", encoded: pointer.To(base64.StdEncoding.EncodeToString([]byte("Welcome\n日本語"))), want: "Welcome\n日本語"},
		{name: "literal base64", encoded: pointer.To(base64.StdEncoding.EncodeToString([]byte("V2VsY29tZQ=="))), want: "V2VsY29tZQ=="},
		{name: "malformed", encoded: pointer.To("not base64!"), wantErr: true},
		{name: "partially decodable", encoded: pointer.To("V2VsY29tZQ==!"), wantErr: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("default pool", func(t *testing.T) {
				data := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, map[string]any{
					"default_node_pool": []any{map[string]any{
						"name": "default", "vm_size": "Standard_DS2_v2", "node_count": 1, "message_of_the_day": "Previous message",
					}},
				})
				profiles := []managedclusters.ManagedClusterAgentPoolProfile{{
					Name: "default", Mode: pointer.To(managedclusters.AgentPoolModeSystem), MessageOfTheDay: test.encoded,
				}}
				flattened, err := FlattenDefaultNodePool(&profiles, data)
				if test.wantErr {
					if err == nil || !strings.Contains(err.Error(), "decoding `message_of_the_day`") {
						t.Fatalf("expected MOTD decoding error, got %v", err)
					}
					if got := data.Get("default_node_pool.0.message_of_the_day"); got != "Previous message" {
						t.Fatalf("invalid MOTD overwrote prior state: %q", got)
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				if err := data.Set("default_node_pool", flattened); err != nil {
					t.Fatal(err)
				}
				if got := data.Get("default_node_pool.0.message_of_the_day"); got != test.want {
					t.Fatalf("MOTD = %q, want %q", got, test.want)
				}
			})
			t.Run("standalone pool", func(t *testing.T) {
				data := schema.TestResourceDataWithIdentityRaw(t, resourceKubernetesClusterNodePool().Schema, pluginsdk.GenerateIdentitySchema(&agentpools.AgentPoolId{})(), nil)
				if err := data.Set("message_of_the_day", "Previous message"); err != nil {
					t.Fatal(err)
				}
				id := agentpools.NewAgentPoolID("00000000-0000-0000-0000-000000000000", "test", "test", "pool")
				data.SetId(id.ID())
				err := resourceKubernetesClusterNodePoolFlatten(data, &id, &agentpools.AgentPool{Properties: &agentpools.ManagedClusterAgentPoolProfileProperties{MessageOfTheDay: test.encoded}})
				if test.wantErr {
					if err == nil || !strings.Contains(err.Error(), "decoding `message_of_the_day`") {
						t.Fatalf("expected MOTD decoding error, got %v", err)
					}
					if got := data.Get("message_of_the_day"); got != "Previous message" {
						t.Fatalf("invalid MOTD overwrote prior state: %q", got)
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				if got := data.Get("message_of_the_day"); got != test.want {
					t.Fatalf("MOTD = %q, want %q", got, test.want)
				}
			})
		})
	}
}
