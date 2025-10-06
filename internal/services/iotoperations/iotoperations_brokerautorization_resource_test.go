package iotoperations

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/acctest"
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

func testAccBrokerAuthorization_basic(rgName, instanceName, brokerName, authName string) string {
    return acctest.Config(`
resource "azurerm_resource_group" "test" {
  name     = "` + rgName + `"
  location = "eastus"
}

resource "azurerm_iotoperations_broker_authorization" "test" {
  name                = "` + authName + `"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "` + instanceName + `"
  broker_name         = "` + brokerName + `"

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
        method   = "Read"
      }
    }
  }

  extended_location {
    name = "/subscriptions/F8C729F9-DF9C-4743-848F-96EE433D8E53/resourceGroups/rgiotoperations/providers/Microsoft.ExtendedLocation/customLocations/resource-123"
    type = "CustomLocation"
  }
}
`)
}