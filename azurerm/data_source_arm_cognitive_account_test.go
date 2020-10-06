package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMCognitiveAccount_basic(t *testing.T) {
	dataSourceName := "data.azurerm_cognitive_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccDataSourceAzureRMCognitiveAccount_basic(ri, rs, location)
	config := testAccDataSourceAzureRMCognitiveAccount_basicWithDataSource(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.Acceptance", "Test"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMCognitiveAccount_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Face"

  sku {
    name = "S0"
    tier = "Standard"
  }

  tags = {
    Acceptance = "Test"
  }
}
`, rInt, location, rString)
}

func testAccDataSourceAzureRMCognitiveAccount_basicWithDataSource(rInt int, rString string, location string) string {
	config := testAccDataSourceAzureRMCognitiveAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_cognitive_account" "test" {
  name                = "${azurerm_cognitive_account.test.name}"
  resource_group_name = "${azurerm_cognitive_account.test.resource_group_name}"
}
`, config)
}
