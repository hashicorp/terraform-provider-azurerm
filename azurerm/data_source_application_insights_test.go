package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceApplicationInsights_basic(t *testing.T) {
	dataSourceName := "data.azurerm_application_insights.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationInsights_complete(ri, location),
				Check:  resource.TestCheckResourceAttrSet(dataSourceName, "instrumentation_key"),
			},
		},
	})
}

func testAccResourceApplicationInsights_complete(rInt int, location string) string {
	resource := testAccAzureRMApplicationInsights_basic(rInt, location, "other")

	return fmt.Sprintf(`
	%s

data "azurerm_application_insights" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_application_insights.test.name}"
}
`, resource)
}
