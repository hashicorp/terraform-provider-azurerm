package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestAccBrokerAuthorization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_broker_authorization", "test")
	r := BrokerAuthorizationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

type BrokerAuthorizationResource struct{}

func (BrokerAuthorizationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := brokerauthorization.ParseAuthorizationID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %+v", state.ID, err)
	}

	resp, err := clients.IoTOperations.BrokerAuthorizationClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r BrokerAuthorizationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iotoperations_instance" "test" {
  name                = "acctestinstance%d"
  resource_group_name = azurerm_resource_group.test.name
  location           = azurerm_resource_group.test.location

  extended_location {
    name = "/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.ExtendedLocation/customLocations/location1"
    type = "CustomLocation"
  }
}

resource "azurerm_iotoperations_broker" "test" {
  name                = "acctestbroker%d"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name

  properties {
    memory_profile = "Tiny"
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}

resource "azurerm_iotoperations_broker_authorization" "test" {
  name                = "acctestauth%d"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name
  broker_name         = azurerm_iotoperations_broker.test.name

  authorization_policies {
    cache = "Enabled"
    rules {
      broker_resources {
        method = "Connect"
        clients = ["test-client"]
        topics = ["test-topic"]
      }
      principals {
        clients = ["test-client"]
        usernames = ["test-user"]
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().SubscriptionID, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
