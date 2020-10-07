package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStreamAnalyticsJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_stream_analytics_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStreamAnalyticsJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "job_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "streaming_units"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "transformation_query"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStreamAnalyticsJob_basic(data acceptance.TestData) string {
	config := testAccAzureRMStreamAnalyticsJob_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_stream_analytics_job" "test" {
  name                = azurerm_stream_analytics_job.test.name
  resource_group_name = azurerm_stream_analytics_job.test.resource_group_name
}
`, config)
}
