package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMDomainService_complete(t *testing.T) {
	dataSourceName := "data.azurerm_domain_service.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureDomainServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDomainService_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "filtered_sync", "Disabled"),
					resource.TestCheckResourceAttr(dataSourceName, "domain_controller_ip_address.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ldaps_settings.0.external_access_ip_address"),
				),
			},
		},
	})
}

func testAccDataSourceDomainService_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_domain_service" "test" {
  name                = "${azurerm_domain_service.test.name}"
  resource_group_name = "${azurerm_domain_service.test.resource_group_name}"
}
`, testAccAzureRMDomainService_complete(rInt, location))
}
