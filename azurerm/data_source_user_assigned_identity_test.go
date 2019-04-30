package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMUserAssignedIdentity_basic(t *testing.T) {
	generatedUuidRegex := "^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$"
	dataSourceName := "data.azurerm_user_assigned_identity.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMUserAssignedIdentity_basic(ri, testLocation(), rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s-uai", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", fmt.Sprintf("acctest%d-rg", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "location", azureRMNormalizeLocation(location)),
					resource.TestMatchResourceAttr(dataSourceName, "principal_id", regexp.MustCompile(generatedUuidRegex)),
					resource.TestMatchResourceAttr(dataSourceName, "client_id", regexp.MustCompile(generatedUuidRegex)),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMUserAssignedIdentity_basic(rInt int, location string, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name = "acctest%d-rg"
	location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
	name = "acctest%s-uai"
	resource_group_name = "${azurerm_resource_group.test.name}"
	location = "${azurerm_resource_group.test.location}"
	tags = {
		"foo" = "bar"
	}
}

data "azurerm_user_assigned_identity" "test" {
	name = "${azurerm_user_assigned_identity.test.name}"
	resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rString)
}
