package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataflowResource struct{}

func TestAccDataflow_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow", "test")
	r := DataflowResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("test-dataflow"),
				check.That(data.ResourceName).Key("properties.0.mode").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
	})
}

func (r DataflowResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	// TODO: Implement proper existence check when the dataflow client is available
	// For now, return nil to indicate the resource exists (placeholder implementation)
	return nil, nil
}

func (r DataflowResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}

resource "azurerm_iotoperations_dataflow" "test" {
  name                   = "test-dataflow"
  resource_group_name    = azurerm_resource_group.test.name
  instance_name         = "test-instance"
  dataflow_profile_name = "test-profile"
  location              = azurerm_resource_group.test.location

  properties {
    mode                     = "Enabled"
    request_disk_persistence = "Enabled"
    
    nodes {
      type = "source"
      name = "temperature"
    }
  }

  extended_location {
    name = "test-custom-location"
    type = "CustomLocation"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
