package blueprints_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type BlueprintPublishedVersionDataSource struct {
}

// lintignore:AT001
func TestAccBlueprintPublishedVersionDataSource_atSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_published_version", "test")
	r := BlueprintPublishedVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.atSubscription(data, "testAcc_basicSubscription", "v0.1_testAcc"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("target_scope").HasValue("subscription"),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test stub for Blueprints at Subscription"),
				check.That(data.ResourceName).Key("time_created").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
			),
		},
	})
}

func TestAccBlueprintPublishedVersionDataSource_atRootManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_published_version", "test")
	r := BlueprintPublishedVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.atRootManagementGroup("testAcc_basicRootManagementGroup", "v0.1_testAcc"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("target_scope").HasValue("subscription"),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test stub for Blueprints at Root Management Group"),
				check.That(data.ResourceName).Key("time_created").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
			),
		},
	})
}

func TestAccBlueprintPublishedVersionDataSource_atChildManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_blueprint_published_version", "test")
	r := BlueprintPublishedVersionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.atChildManagementGroup("testAcc_staticStubGroup", "testAcc_staticStubManagementGroup", "v0.1_testAcc"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("target_scope").HasValue("subscription"),
				check.That(data.ResourceName).Key("description").HasValue("Acceptance Test stub for Blueprints at Child Management Group"),
				check.That(data.ResourceName).Key("time_created").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
			),
		},
	})
}

func (BlueprintPublishedVersionDataSource) atSubscription(data acceptance.TestData, bpName string, version string) string {
	subscription := data.Client().SubscriptionIDAlt

	return fmt.Sprintf(`
provider "azurerm" {
  subscription_id = "%s"
  features {}
}

data "azurerm_subscription" "current" {}

data "azurerm_blueprint_published_version" "test" {
  scope_id       = data.azurerm_subscription.current.id
  blueprint_name = "%s"
  version        = "%s"
}
`, subscription, bpName, version)
}

func (BlueprintPublishedVersionDataSource) atRootManagementGroup(bpName, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_management_group" "root" {
  name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_blueprint_published_version" "test" {
  scope_id       = data.azurerm_management_group.root.id
  blueprint_name = "%s"
  version        = "%s"
}
`, bpName, version)
}

func (BlueprintPublishedVersionDataSource) atChildManagementGroup(mg, bpName, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_management_group" "test" {
  name = "%s"
}

data "azurerm_blueprint_published_version" "test" {
  scope_id       = data.azurerm_management_group.test.id
  blueprint_name = "%s"
  version        = "%s"
}
`, mg, bpName, version)
}
