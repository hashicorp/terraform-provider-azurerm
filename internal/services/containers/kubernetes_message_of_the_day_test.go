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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	containersclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

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
				data := schema.TestResourceDataRaw(t, resourceKubernetesClusterNodePool().Schema, map[string]any{
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
				err = resourceKubernetesClusterNodePoolRead(data, meta)
				if requests != 1 {
					t.Fatalf("expected one mocked GET, got %d", requests)
				}
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

type motdTestTransport func(*http.Request) (*http.Response, error)

func (f motdTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
