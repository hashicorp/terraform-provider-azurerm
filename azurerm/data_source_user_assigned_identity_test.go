package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMUserAssignedIdentity_basic(t *testing.T) {
	dataSourceName := "data.azurerm_user_assigned_identity.test"
	resourceName := "azurerm_user_assigned_identity.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMUserAssignedIdentity_basic(ri, acceptance.Location(), rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s-uai", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", fmt.Sprintf("acctest%d-rg", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "location", azure.NormalizeLocation(location)),
					resource.TestMatchResourceAttr(dataSourceName, "principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(dataSourceName, "client_id", validate.UUIDRegExp),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					testEqualResourceAttr(dataSourceName, resourceName, "principal_id"),
					testEqualResourceAttr(dataSourceName, resourceName, "client_id"),
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

func testAccDataSourceAzureRMUserAssignedIdentity_basic(rInt int, location string, rString string) string {
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
`, rInt, location, rString)
}
