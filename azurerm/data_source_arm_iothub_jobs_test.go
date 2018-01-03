package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMIotHubJobs_basic(t *testing.T) {
	dataSourceName := "data.azurerm_iothub_jobs.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIotHubJobsConfig_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "job_id", "test"),
				),
			},
		},
	})
}

func testAccDataSourceIotHubJobsConfig_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "foo" {
	job_id = "acctestIot-%d"
	location = "%s"
}

resource "azurerm_iothub" "bar" {
	job_id = "acctestiothub-%d"
	location = "${azurerm_resource_group.foo.location}"
	resource_group_name = "${azurerm_resource_group.foo.name}"
	job {
		name = "S1"
		tier = "Standard"
		capacity = "1"
	}

	tags {
		"purpose" = "testing"
	}
}

data "azurerm_iothub_jobs" "test" {
	job_id = "test"
	resource_group_name = "${azurerm_resource_group.foo.name}"
	iot_hub_name = "${azurerm_iothub.bar.name}"
}

	`, rInt, location, rInt)
}
