package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
)

func TestAccAzureRMApplicationInsightsAnalyticsItem_basic(t *testing.T) {
	resourceName := "azurerm_application_insights_analytics_item.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMApplicationInsightsAnalyticsItem_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightAnalyticsItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAnalyticsItemExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope", "shared"),
					resource.TestCheckResourceAttr(resourceName, "type", "query"),
					resource.TestCheckResourceAttr(resourceName, "content", "requests #test"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationInsightAnalyticsItemDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appInsights.AnalyticsItemsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights_analytics_item" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resGroup, name, insights.AnalyticsItems, "", "testquery")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Application Insights AnalyticsItem still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsAnalyticsItemExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		conn := testAccProvider.Meta().(*ArmClient).appInsights.AnalyticsItemsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resGroup, name, insights.AnalyticsItems, "", "testquery")
		if err != nil {
			return fmt.Errorf("Bad: Get on appInsightsAnalyticsItemsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Application Insights AnalyticsItem '%q' (resource group: '%q') does not exist", name, resGroup)
		}

		return nil
	}
}

func testAccAzureRMApplicationInsightsAnalyticsItem_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_application_insights_analytics_item" "test" {
  name                    = "testquery"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  application_insights_id = "${azurerm_application_insights.test.id}"
  content                 = "requests #test"
	scope                   = "shared"
	type                    = "query"
}
`, rInt, location, rInt)
}
