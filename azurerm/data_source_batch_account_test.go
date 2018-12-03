package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMBatchAccount_basic(t *testing.T) {
	dataSourceName := "data.azurerm_batch_account.test"
	ri := acctest.RandInt()
	name := fmt.Sprintf("acctestbatchaccount%d", ri)
	location := testLocation()
	config := testAccDataSourceAzureRMBatchAccount_basic(name, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "location", azureRMNormalizeLocation(location)),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMBatchAccount_basic(name string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_batch_account" "test" {
  name = "%s"
  location = "%s"

  tags {
    env = "test"
  }
}

data "azurerm_batch_account" "test" {
	name = "${azurerm_batch_account.test.name}"
  }
`, name, location)
}
