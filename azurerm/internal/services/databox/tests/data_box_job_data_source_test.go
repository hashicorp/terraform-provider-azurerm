package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataBoxJobDataSource_basic(t *testing.T) {
	location, err := testGetLocationFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataBoxJob_basic(data, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDataBoxJob_basic(data acceptance.TestData, location string) string {
	config := testAccDataBoxJob_basic(data, location)
	return fmt.Sprintf(`
%s

data "azurerm_data_box_job" "test" {
  name                = azurerm_data_box_job.test.name
  resource_group_name = azurerm_data_box_job.test.resource_group_name
}
`, config)
}
