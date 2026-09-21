// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	containersclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

func TestKubernetesNodeResourceGroupRestrictionPlan(t *testing.T) {
	resource := resourceKubernetesCluster()
	config := func(level interface{}, omit bool) map[string]interface{} {
		result := map[string]interface{}{
			"name": "test", "location": "westus2", "resource_group_name": "test", "dns_prefix": "test",
			"default_node_pool": []interface{}{map[string]interface{}{"name": "default", "node_count": 1, "vm_size": "Standard_DS2_v2"}},
			"identity":          []interface{}{map[string]interface{}{"type": "SystemAssigned"}},
			"network_profile":   []interface{}{map[string]interface{}{"network_plugin": "kubenet", "load_balancer_sku": "standard"}},
		}
		if !omit {
			result["node_resource_group_restriction_level"] = level
		}
		return result
	}
	for _, test := range []struct {
		name    string
		old     string
		new     interface{}
		omit    bool
		replace bool
	}{
		{name: "omitted to read only", new: "ReadOnly"},
		{name: "omitted to unrestricted", new: "Unrestricted"},
		{name: "restrict", old: "Unrestricted", new: "ReadOnly"},
		{name: "reverse", old: "ReadOnly", new: "Unrestricted"},
		{name: "read only to omitted", old: "ReadOnly", omit: true, replace: true},
		{name: "read only to null", old: "ReadOnly", replace: true},
		{name: "read only to empty", old: "ReadOnly", new: "", replace: true},
		{name: "unrestricted to omitted", old: "Unrestricted", omit: true, replace: true},
		{name: "unrestricted to null", old: "Unrestricted", replace: true},
		{name: "unrestricted to empty", old: "Unrestricted", new: "", replace: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, resource.Schema, config(test.old, test.old == ""))
			data.SetId("existing-cluster")
			diff, err := resource.SimpleDiff(context.Background(), data.State(), terraform.NewResourceConfigRaw(config(test.new, test.omit)), nil)
			if err != nil {
				t.Fatal(err)
			}
			attribute, ok := diff.Attributes["node_resource_group_restriction_level"]
			if !ok || attribute.RequiresNew != test.replace || diff.RequiresNew() != test.replace {
				t.Fatalf("expected restriction change with replacement=%t, got %#v", test.replace, diff)
			}
		})
	}
}

func TestKubernetesNodeResourceGroupRestrictionResponse(t *testing.T) {
	for _, prior := range []string{"", "ReadOnly"} {
		for _, test := range []struct {
			name    string
			profile *managedclusters.ManagedClusterNodeResourceGroupProfile
			want    string
		}{
			{name: "absent profile preserves prior state", want: prior},
			{name: "absent level preserves prior state", profile: &managedclusters.ManagedClusterNodeResourceGroupProfile{}, want: prior},
			{name: "read only", profile: &managedclusters.ManagedClusterNodeResourceGroupProfile{RestrictionLevel: pointer.To(managedclusters.RestrictionLevelReadOnly)}, want: "ReadOnly"},
			{name: "unrestricted", profile: &managedclusters.ManagedClusterNodeResourceGroupProfile{RestrictionLevel: pointer.To(managedclusters.RestrictionLevelUnrestricted)}, want: "Unrestricted"},
		} {
			t.Run(prior+"/"+test.name, func(t *testing.T) {
				data := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, map[string]interface{}{
					"node_resource_group_restriction_level": prior,
				})
				id := commonids.NewKubernetesClusterID("00000000-0000-0000-0000-000000000000", "test", "test")
				data.SetId(id.ID())
				payload, err := json.Marshal(managedclusters.ManagedCluster{
					Location: "westus2",
					Properties: &managedclusters.ManagedClusterProperties{
						NodeResourceGroupProfile: test.profile,
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				clusterClient, err := managedclusters.NewManagedClustersClientWithBaseURI(environments.AzurePublic().ResourceManager)
				if err != nil {
					t.Fatal(err)
				}
				maintenanceClient, err := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(environments.AzurePublic().ResourceManager)
				if err != nil {
					t.Fatal(err)
				}
				requests := 0
				transport := restrictionTestTransport(func(request *http.Request) (*http.Response, error) {
					requests++
					if request.URL.Query().Get("api-version") != "2026-05-01" {
						t.Fatalf("unexpected API version: %s", request.URL)
					}
					status, body := http.StatusOK, string(payload)
					switch {
					case request.Method == http.MethodGet && request.URL.Path == id.ID():
					case request.Method == http.MethodPost && request.URL.Path == id.ID()+"/listClusterUserCredential":
						body = `{"kubeconfigs":[]}`
					case request.Method == http.MethodGet && strings.HasPrefix(request.URL.Path, id.ID()+"/maintenanceConfigurations/"):
						status, body = http.StatusNotFound, `{"error":{"code":"NotFound","message":"not configured"}}`
					default:
						t.Fatalf("unexpected request: %s %s", request.Method, request.URL)
					}
					return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: request}, nil
				})
				clusterClient.Client.AuthorizeRequest = nil
				clusterClient.Client.SetTransport(transport)
				maintenanceClient.Client.AuthorizeRequest = nil
				maintenanceClient.Client.SetTransport(transport)
				meta := &clients.Client{StopContext: context.Background(), Containers: &containersclient.Client{
					KubernetesClustersClient: clusterClient, MaintenanceConfigurationsClient: maintenanceClient,
				}}
				if err := resourceKubernetesClusterRead(data, meta); err != nil {
					t.Fatal(err)
				}
				if requests != 5 {
					t.Fatalf("expected cluster, credentials and three maintenance requests, got %d", requests)
				}
				if got := data.Get("node_resource_group_restriction_level"); got != test.want {
					t.Fatalf("restriction = %q, want %q", got, test.want)
				}
			})
		}
	}
}

type restrictionTestTransport func(*http.Request) (*http.Response, error)

func (f restrictionTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
