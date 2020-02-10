package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSchedulerJobCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_scheduler_job_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSchedulerJobCollection_basic(data),
				Check:  checkAccAzureRMSchedulerJobCollection_basic(data.ResourceName),
			},
		},
	})
}

func TestAccDataSourceAzureRMSchedulerJobCollection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_scheduler_job_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSchedulerJobCollection_complete(data),
				Check:  checkAccAzureRMSchedulerJobCollection_complete(data.ResourceName),
			},
		},
	})
}

func testAccDataSourceSchedulerJobCollection_basic(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJobCollection_basic(data, "")
	return fmt.Sprintf(`
%s

data "azurerm_scheduler_job_collection" "test" {
  name                = azurerm_scheduler_job_collection.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceSchedulerJobCollection_complete(data acceptance.TestData) string {
	template := testAccAzureRMSchedulerJobCollection_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_scheduler_job_collection" "test" {
  name                = azurerm_scheduler_job_collection.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
