package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAzureRMUserAssignedIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_user_assigned_identity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMUserAssignedIdentity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctest%s-uai", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctest%d-rg", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
					resource.TestMatchResourceAttr(data.ResourceName, "principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "client_id", validate.UUIDRegExp),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					testEqualResourceAttr(data.ResourceName, "data."+data.ResourceName, "principal_id"),
					testEqualResourceAttr(data.ResourceName, "data."+data.ResourceName, "client_id"),
				),
			},
		},
	})
}

func testEqualResourceAttr(dataSourceName string, resourceName string, attrName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dsAttr := ds.Primary.Attributes[attrName]
		rsAttr := rs.Primary.Attributes[attrName]

		if dsAttr != rsAttr {
			return fmt.Errorf("Attributes not equal: %s, %s", dsAttr, rsAttr)
		}

		return nil
	}
}

func testAccDataSourceAzureRMUserAssignedIdentity_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s-uai"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  tags = {
    "foo" = "bar"
  }
}

data "azurerm_user_assigned_identity" "test" {
  name                = "${azurerm_user_assigned_identity.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
