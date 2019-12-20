package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataSourceHealthCareService_basic(t *testing.T) {
	dataSourceName := "data.azurerm_healthcare_service.test"
	ri := tf.AccRandTimeInt() / 10
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
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
	resource := testAccAzureRMHealthCareService_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_service" "test" {
  name                = "${azurerm_healthcare_service.test.name}"
  resource_group_name = "${azurerm_healthcare_service.test.resource_group_name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, resource)
}
