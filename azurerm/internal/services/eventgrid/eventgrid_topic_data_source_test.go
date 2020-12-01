package eventgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMEventGridTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_topic", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMEventGridTopic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMEventGridTopic_basic(data acceptance.TestData) string {
	template := testAccAzureRMEventGridTopic_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_topic" "test" {
  name                = azurerm_eventgrid_topic.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
