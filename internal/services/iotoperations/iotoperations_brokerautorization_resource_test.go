package iotoperations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccBrokerAuthorization_basic(t *testing.T) {
    resourceName := "azurerm_iotoperations_broker_authorization.test"
    rgName := acctest.RandomResourceGroupName("acctestRG", 16)
    instanceName := acctest.RandomResourceName("acctestInstance", 16)
    brokerName := acctest.RandomResourceName("acctestBroker", 16)
    authName := acctest.RandomResourceName("acctestAuth", 16)

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { acctest.PreCheck(t) },
        Providers:    acctest.Providers,
        CheckDestroy: acctest.CheckResourceDestroy(resourceName),
        Steps: []resource.TestStep{
            {
                Config: testAccBrokerAuthorization_basic(rgName, instanceName, brokerName, authName),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "name", authName),
                    resource.TestCheckResourceAttr(resourceName, "resource_group_name", rgName),
                    resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
                    resource.TestCheckResourceAttr(resourceName, "broker_name", brokerName),
                    resource.TestCheckResourceAttrSet(resourceName, "provisioning_state"),
                ),
            },
        },
    })
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
  name                = "resource-name123"
  resource_group_name = azurerm_resource_group.test.name
  location           = azurerm_resource_group.test.location

  extended_location {
    name = "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ExtendedLocation/customLocations/location1"
    type = "CustomLocation"
  }
}

resource "azurerm_iotoperations_broker" "test" {
  name                = "resource-name123"
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
  name                = "resource-name123"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = azurerm_iotoperations_instance.test.name
  broker_name         = azurerm_iotoperations_broker.test.name

  authorization_policies {
    cache = "Enabled"
    rules {
      broker_resources {
        method     = "Connect"
        client_ids = ["nlc"]
        topics     = ["wvuca"]
      }
      principals {
        attributes = [{ key5526 = "nydhzdhbldygqcn" }]
        client_ids = ["smopeaeddsygz"]
        usernames  = ["iozngyqndrteikszkbasinzdjtm"]
      }
      state_store_resources {
        key_type = "Pattern"
        keys     = ["tkounsqtwvzyaklxjqoerpu"]
      }
    }
  }

  extended_location {
    name = azurerm_iotoperations_instance.test.extended_location[0].name
    type = azurerm_iotoperations_instance.test.extended_location[0].type
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Client().SubscriptionID, data.RandomStringOfLength(10))
}