package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestAccBroker_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker", "test")
	r := BrokerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("resource-name123"),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("instance_name").HasValue("resource-name123"),
				check.That(data.ResourceName).Key("properties.0.memory_profile").HasValue("Tiny"),
				check.That(data.ResourceName).Key("extended_location.0.type").HasValue("CustomLocation"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("properties.0.memory_profile").HasValue("Large"),
			),
		},
		data.ImportStep(),
	})
}

type BrokerResource struct{}

func (BrokerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := broker.ParseBrokerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTOperations.BrokerClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return acceptance.Exists(resp.Model), nil
}

func (BrokerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}

resource "azurerm_iotoperations_broker" "test" {
  name                = "resource-name123"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "resource-name123"

  properties {
    memory_profile = "Tiny"
  }

  extended_location {
    name = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/resource-123"
    type = "CustomLocation"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Client().SubscriptionID)
}

func (BrokerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iot-%d"
  location = "%s"
}

resource "azurerm_iotoperations_broker" "test" {
  name                = "resource-name123"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "resource-name123"

  properties {
    memory_profile = "Large"
  }

  extended_location {
    name = "/subscriptions/%s/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.ExtendedLocation/customLocations/resource-123"
    type = "CustomLocation"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Client().SubscriptionID)
}
