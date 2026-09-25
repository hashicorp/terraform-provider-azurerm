// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	containersclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesMessageOfTheDayPlan(t *testing.T) {
	for _, standalone := range []bool{false, true} {
		name := "default pool"
		resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"default_node_pool": SchemaDefaultNodePool(),
		}}
		if standalone {
			name = "standalone pool"
			resource.Schema = map[string]*pluginsdk.Schema{
				"message_of_the_day": resourceKubernetesClusterNodePool().Schema["message_of_the_day"],
			}
		}
		t.Run(name, func(t *testing.T) {
			config := func(message interface{}, omit bool) map[string]interface{} {
				block := map[string]interface{}{}
				if !omit {
					block["message_of_the_day"] = message
				}
				if standalone {
					return block
				}
				block["name"] = "default"
				block["vm_size"] = "Standard_DS2_v2"
				block["node_count"] = 1
				block["temporary_name_for_rotation"] = "temp"
				return map[string]interface{}{"default_node_pool": []interface{}{block}}
			}
			for _, test := range []struct {
				name    string
				message interface{}
				omit    bool
			}{
				{name: "clear", message: ""},
				{name: "null", message: nil},
				{name: "omitted", omit: true},
				{name: "reset", message: "Reset message"},
			} {
				t.Run(test.name, func(t *testing.T) {
					data := schema.TestResourceDataRaw(t, resource.Schema, config("Original message", false))
					data.SetId("existing-pool")
					diff, err := resource.SimpleDiff(context.Background(), data.State(), terraform.NewResourceConfigRaw(config(test.message, test.omit)), nil)
					if err != nil {
						t.Fatal(err)
					}
					if diff.Empty() || diff.RequiresNew() != standalone {
						t.Fatalf("expected changed MOTD with replacement=%t, got %#v", standalone, diff)
					}
				})
			}
		})
	}
}

func TestKubernetesMessageOfTheDayExpand(t *testing.T) {
	for _, test := range []struct {
		name    string
		message interface{}
		omit    bool
		want    string
	}{
		{name: "omitted", omit: true},
		{name: "null", message: nil},
		{name: "empty", message: ""},
		{name: "text", message: "Welcome\n日本語", want: "Welcome\n日本語"},
		{name: "literal shell", message: "$(hostname)", want: "$(hostname)"},
	} {
		t.Run(test.name, func(t *testing.T) {
			block := map[string]interface{}{"name": "default", "vm_size": "Standard_DS2_v2", "node_count": 1}
			if !test.omit {
				block["message_of_the_day"] = test.message
			}
			data := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, map[string]interface{}{
				"default_node_pool": []interface{}{block},
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
	}{
		{name: "absent"},
		{name: "empty", encoded: pointer.To("")},
		{name: "text", encoded: pointer.To(base64.StdEncoding.EncodeToString([]byte("Welcome\n日本語"))), want: "Welcome\n日本語"},
		{name: "malformed", encoded: pointer.To("not base64!")},
		{name: "partially decodable", encoded: pointer.To("V2VsY29tZQ==!")},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("default pool", func(t *testing.T) {
				data := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, map[string]interface{}{
					"default_node_pool": []interface{}{map[string]interface{}{
						"name": "default", "vm_size": "Standard_DS2_v2", "node_count": 1, "message_of_the_day": "Previous message",
					}},
				})
				profiles := []managedclusters.ManagedClusterAgentPoolProfile{{
					Name: "default", Mode: pointer.To(managedclusters.AgentPoolModeSystem), MessageOfTheDay: test.encoded,
				}}
				flattened, err := FlattenDefaultNodePool(&profiles, data)
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
				data := schema.TestResourceDataRaw(t, resourceKubernetesClusterNodePool().Schema, map[string]interface{}{
					"message_of_the_day": "Previous message",
				})
				id := agentpools.NewAgentPoolID("00000000-0000-0000-0000-000000000000", "test", "test", "pool")
				data.SetId(id.ID())
				payload, err := json.Marshal(agentpools.AgentPool{Properties: &agentpools.ManagedClusterAgentPoolProfileProperties{MessageOfTheDay: test.encoded}})
				if err != nil {
					t.Fatal(err)
				}
				client, err := agentpools.NewAgentPoolsClientWithBaseURI(environments.AzurePublic().ResourceManager)
				if err != nil {
					t.Fatal(err)
				}
				client.Client.AuthorizeRequest = nil
				requests := 0
				client.Client.SetTransport(motdTestTransport(func(request *http.Request) (*http.Response, error) {
					requests++
					if request.Method != http.MethodGet || request.URL.Path != id.ID() || request.URL.Query().Get("api-version") != "2026-05-01" {
						t.Fatalf("unexpected request: %s %s", request.Method, request.URL)
					}
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(string(payload))), Request: request}, nil
				}))
				meta := &clients.Client{StopContext: context.Background(), Containers: &containersclient.Client{AgentPoolsClient: client}}
				if err := resourceKubernetesClusterNodePoolRead(data, meta); err != nil {
					t.Fatal(err)
				}
				if requests != 1 {
					t.Fatalf("expected one mocked GET, got %d", requests)
				}
				if got := data.Get("message_of_the_day"); got != test.want {
					t.Fatalf("MOTD = %q, want %q", got, test.want)
				}
			})
		})
	}
}

type motdTestTransport func(*http.Request) (*http.Response, error)

func (f motdTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
