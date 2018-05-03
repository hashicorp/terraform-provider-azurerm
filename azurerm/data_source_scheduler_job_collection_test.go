package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMSchedulerJobCollection_basic(t *testing.T) {
	dataSourceName := "data.azurerm_scheduler_job_collection.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSchedulerJobCollection_basic(ri, testLocation()),
				Check:  checkAccAzureRMSchedulerJobCollection_basic(dataSourceName),
			},
		},
	})
}

func TestAccDataSourceAzureRMSchedulerJobCollection_complete(t *testing.T) {
	dataSourceName := "data.azurerm_scheduler_job_collection.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSchedulerJobCollection_complete(ri, testLocation()),
				Check:  checkAccAzureRMSchedulerJobCollection_complete(dataSourceName),
			},
		},
	})
}

func testAccDataSourceSchedulerJobCollection_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_scheduler_job_collection" "test" {
  name                = "${azurerm_scheduler_job_collection.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMSchedulerJobCollection_basic(rInt, location, ""))
}

func testAccDataSourceSchedulerJobCollection_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_scheduler_job_collection" "test" {
  name                = "${azurerm_scheduler_job_collection.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMSchedulerJobCollection_complete(rInt, location))
}
