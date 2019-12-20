package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStreamAnalyticsJob_basic(t *testing.T) {
	dataSourceName := "data.azurerm_stream_analytics_job.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStreamAnalyticsJob_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "job_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStreamAnalyticsJob_basic(rInt int, location string) string {
	config := testAccAzureRMStreamAnalyticsJob_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_stream_analytics_job" "test" {
  name                = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name = "${azurerm_stream_analytics_job.test.resource_group_name}"
}
`, config)
}
