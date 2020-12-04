package subscription_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type SubscriptionDataSource struct{}

func TestAccDataSourceSubscription_current(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription", "current")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: SubscriptionDataSource{}.currentConfig(),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureRMSubscriptionId(data.ResourceName),
				check.That(data.ResourceName).Key("subscription_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("tenant_id").Exists(),
				check.That(data.ResourceName).Key("state").HasValue("Enabled"),
			),
		},
	})
}

func TestAccDataSourceSubscription_specific(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subscription", "specific")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: SubscriptionDataSource{}.specificConfig(os.Getenv("ARM_SUBSCRIPTION_ID")),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureRMSubscriptionId(data.ResourceName),
				check.That(data.ResourceName).Key("subscription_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("tenant_id").Exists(),
				check.That(data.ResourceName).Key("location_placement_id").Exists(),
				check.That(data.ResourceName).Key("quota_id").Exists(),
				check.That(data.ResourceName).Key("spending_limit").Exists(),
			),
		},
	})
}

func testCheckAzureRMSubscriptionId(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		attributeName := "subscription_id"
		subscriptionId := rs.Primary.Attributes[attributeName]
		client := acceptance.AzureProvider.Meta().(*clients.Client)
		if subscriptionId != client.Account.SubscriptionId {
			return fmt.Errorf("%s: Attribute '%s' expected \"%s\", got \"%s\"", resourceName, attributeName, client.Account.SubscriptionId, subscriptionId)
		}

		return nil
	}
}

func (d SubscriptionDataSource) currentConfig() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}
`
}

func (d SubscriptionDataSource) specificConfig(subscriptionId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  subscription_id = "%s"
}

data "azurerm_subscription" "specific" {
  subscription_id = "%s"
}
`, subscriptionId, subscriptionId)
}
