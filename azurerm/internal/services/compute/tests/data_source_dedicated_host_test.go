package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAzureRMDedicatedHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "DSv3-Type1"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_replace_on_failure", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", string(compute.DedicatedHostLicenseTypesNone)),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedHost_basic(data acceptance.TestData) string {
	config := testAccAzureRMDedicatedHost_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_dedicated_host" "test" {
  name                = azurerm_dedicated_host.test.name
  resource_group_name = azurerm_dedicated_host.test.resource_group_name
  host_group_name     = azurerm_dedicated_host.test.host_group_name
}
`, config)
}
