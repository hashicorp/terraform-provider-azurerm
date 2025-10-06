package iotoperations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIotOperationsDataflowEndpoint_basic(t *testing.T) {
    resourceName := "azurerm_iotoperations_dataflow_endpoint.test"
    rgName := "acctestRG-dataflowendpoint"
    instanceName := "acctestInstance-dataflowendpoint"
    endpointName := "acctestEndpoint-dataflowendpoint"

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckIotOperationsDataflowEndpointDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccIotOperationsDataflowEndpointConfig(rgName, instanceName, endpointName),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "name", endpointName),
                    resource.TestCheckResourceAttr(resourceName, "resource_group_name", rgName),
                    resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
                    resource.TestCheckResourceAttr(resourceName, "properties.0.endpoint_type", "DataExplorer"),
                ),
            },
            {
                ResourceName:      resourceName,
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

func testAccIotOperationsDataflowEndpointConfig(rgName, instanceName, endpointName string) string {
    return fmt.Sprintf(`
resource "azurerm_iotoperations_dataflow_endpoint" "test" {
  name                = "%s"
  resource_group_name = "%s"
  instance_name       = "%s"

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
    name = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/%s/providers/Microsoft.ExtendedLocation/customLocations/resource-123"
    type = "CustomLocation"
  }
}
`, endpointName, rgName, instanceName, rgName)
}

func testAccCheckIotOperationsDataflowEndpointDestroy(s *terraform.State) error {
    // TODO: Implement check to verify resource destruction via Azure API
    return nil
}