package healthcare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccHealthCareServiceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_healthcare_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHealthCareServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceHealthcareService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kind"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "cosmosdb_throughput"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceHealthcareService_basic(data acceptance.TestData) string {
	resource := testAccHealthCareService_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_service" "test" {
  name                = azurerm_healthcare_service.test.name
  resource_group_name = azurerm_healthcare_service.test.resource_group_name
  location            = azurerm_resource_group.test.location
}
`, resource)
}
