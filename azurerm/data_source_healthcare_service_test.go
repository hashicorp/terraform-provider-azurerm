package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMDataSourceHealthcareService_basic(t *testing.T) {
	dataSourceName := "data.azurerm_healthcare_service.test"
	ri := tf.AccRandTimeInt() / 10
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHealthcareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceHealthcareService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kind"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cosmosdb_throughput"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "2"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceHealthcareService_basic(rInt int, location string) string {
	resource := testAccAzureRMHealthcareService_basic(rInt)
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_service" "test" {
  name                = "${azurerm_healthcare_service.test.name}"
  resource_group_name = "${azurerm_healthcare_service.test.resource_group_name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, resource)
}
