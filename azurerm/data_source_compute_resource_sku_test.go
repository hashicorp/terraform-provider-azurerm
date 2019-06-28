package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMComputeResourceSku_basic(t *testing.T) {
	dataSourceName := "data.azurerm_compute_resource_sku.test"
	location := testLocation()
	name := "Standard_DS2_v2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMComputeResourceSku_basic(name, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "family"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMComputeResourceSku_basic(name string, location string) string {
	return fmt.Sprintf(`
data "azurerm_compute_resource_sku" "test" {
	name     = "%s"
	location = "%s"
}
`, name, location)
}
