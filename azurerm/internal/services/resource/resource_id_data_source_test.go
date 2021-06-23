package resource_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccDataSourceResourceId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_id", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: `
			provider "azurerm" {
				features {}
			}
			
			data "azurerm_resource_id" "test" {
				resource_id = "/subscriptions/c90e9ba4-9a69-49d6-be99-2110471ec1a4/resourceGroups/SomeResourceGroup/providers/Microsoft.ResourceProvider/resourceType/MyResource"
			}
			`,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue("c90e9ba4-9a69-49d6-be99-2110471ec1a4"),
				check.That(data.ResourceName).Key("resource_group_name").HasValue("SomeResourceGroup"),
				check.That(data.ResourceName).Key("provider_namespace").HasValue("Microsoft.ResourceProvider"),
				check.That(data.ResourceName).Key("resource_type").HasValue("resourceType"),
				check.That(data.ResourceName).Key("name").HasValue("MyResource"),
				check.That(data.ResourceName).Key("full_resource_type").HasValue("Microsoft.ResourceProvider/resourceType"),
			),
		},
	})
}
