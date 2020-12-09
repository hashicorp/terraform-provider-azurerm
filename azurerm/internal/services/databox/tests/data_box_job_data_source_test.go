package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataBoxJobDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataBoxJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDataBoxJob_basic(data acceptance.TestData) string {
	config := testAccDataBoxJob_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_box_job" "test" {
  name                = azurerm_data_box_job.test.name
  resource_group_name = azurerm_data_box_job.test.resource_group_name
}
`, config)
}
