// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	databasemigrationclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/resources/database_migration_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestUnitDatabaseMigrationServiceRead_networkInterfaceId(t *testing.T) {
	id := serviceresource.NewServiceID("00000000-0000-0000-0000-000000000000", "test-rg", "test-dms")
	subnetId := commonids.NewSubnetID(id.SubscriptionId, id.ResourceGroupName, "test-vnet", "test-subnet").ID()
	networkInterfaceId := commonids.NewNetworkInterfaceID(id.SubscriptionId, id.ResourceGroupName, "test-nic").ID()
	resource := database_migration_service.RegisterResource()

	field, ok := resource.Schema["network_interface_id"]
	if !ok {
		t.Fatal("network_interface_id is missing from the resource schema")
	}
	if field.Type != pluginsdk.TypeString || !field.Computed || field.Optional || field.Required || field.ForceNew || field.Default != nil {
		t.Fatal("network_interface_id must be a computed-only string without a default or ForceNew")
	}

	testCases := []struct {
		name         string
		virtualNicId *string
		expected     string
	}{
		{
			name:         "populated",
			virtualNicId: pointer.To(networkInterfaceId),
			expected:     networkInterfaceId,
		},
		{
			name: "omitted clears previous value",
		},
		{
			name:         "empty clears previous value",
			virtualNicId: pointer.To(""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || req.URL.Path != id.ID() || req.URL.Query().Get("api-version") != "2021-06-30" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(serviceresource.DataMigrationService{
					Location: "westeurope",
					Properties: &serviceresource.DataMigrationServiceProperties{
						VirtualNicId:    testCase.virtualNicId,
						VirtualSubnetId: subnetId,
					},
					Sku: &serviceresource.ServiceSku{Name: pointer.To("Standard_1vCores")},
				}); err != nil {
					t.Errorf("encoding response: %v", err)
				}
			}))
			defer server.Close()

			client, err := serviceresource.NewServiceResourceClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
			if err != nil {
				t.Fatalf("building client: %v", err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true

			data := resource.Data(nil)
			data.SetId(id.ID())
			if err := data.Set("network_interface_id", networkInterfaceId+"-old"); err != nil {
				t.Fatalf("setting previous network interface ID: %v", err)
			}

			if err := resource.Read(data, &clients.Client{
				StopContext:       t.Context(),
				DatabaseMigration: &databasemigrationclient.Client{ServicesClient: client},
			}); err != nil {
				t.Fatalf("reading service: %v", err)
			}

			for key, expected := range map[string]string{
				"network_interface_id": testCase.expected,
				"name":                 id.ServiceName,
				"resource_group_name":  id.ResourceGroupName,
				"location":             "westeurope",
				"subnet_id":            subnetId,
				"sku_name":             "Standard_1vCores",
			} {
				if actual := data.Get(key); actual != expected {
					t.Errorf("expected %s to be %q, got %q", key, expected, actual)
				}
			}
		})
	}
}
