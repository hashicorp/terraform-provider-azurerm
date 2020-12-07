package iothub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMIotHubDpsSharedAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iothub_dps_shared_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDpsSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMIotHubDpsSharedAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDpsSharedAccessPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMIotHubDpsSharedAccessPolicy_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_dps_name     = azurerm_iothub_dps.test.name
  name                = "acctest"
  service_config      = true
}

data "azurerm_iothub_dps_shared_access_policy" "test" {
  name                = azurerm_iothub_dps_shared_access_policy.test.name
  iothub_dps_name     = azurerm_iothub_dps.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
