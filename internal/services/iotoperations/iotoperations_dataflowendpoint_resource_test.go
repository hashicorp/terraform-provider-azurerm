package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowendpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataflowEndpointResource struct{}

func TestAccIotOperationsDataflowEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_endpoint", "test")
	r := DataflowEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.endpoint_type").HasValue("DataExplorer"),
			),
		},
		data.ImportStep(),
	})
}

func (r DataflowEndpointResource) Exists(ctx context.Context, c *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataflowendpoint.ParseDataflowEndpointID(state.ID) // adjust if actual helper name differs
	if err != nil {
		return nil, fmt.Errorf("parsing ID %q: %w", state.ID, err)
	}

	resp, err := c.IoTOperations.DataflowEndpointClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r DataflowEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}

resource "azurerm_iotoperations_dataflow_endpoint" "test" {
  name                = "test-endpoint-%d"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "test-instance"
  location            = azurerm_resource_group.test.location

  properties {
    endpoint_type = "DataExplorer"
    data_explorer_settings {
      authentication {
        method = "SystemAssignedManagedIdentity"
        system_assigned_managed_identity_settings {
          audience = "psxomrfbhoflycm"
        }
      }
      database = "yqcdpjsifm"
      host     = "cluster.region.kusto.windows.net"
      batching {
        latency_seconds = 9312
        max_messages    = 9028
      }
    }
  }

  extended_location {
    name = "test-custom-location"
    type = "CustomLocation"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
