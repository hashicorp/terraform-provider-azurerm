package applicationinsights_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppInsightsSmartDetectionRule struct {
}

func TestAccApplicationInsightsSmartDetectionRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccApplicationInsightsSmartDetectionRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccApplicationInsightsSmartDetectionRule_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_application_insights_smart_detection_rule.test2").ExistsInAzure(r),
			),
		},
	})
}

func TestAccApplicationInsightsSmartDetectionRule_longDependencyDuration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.longDependencyDuration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t AppInsightsSmartDetectionRule) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SmartDetectionRuleID(state.Attributes["id"])
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppInsights.SmartDetectionRuleClient.Get(ctx, id.ResourceGroup, id.ComponentName, id.SmartDetectionRuleName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Application Insights Smart Detection Rule '%q' does not exist", id.String())
	}

	return utils.Bool(resp.StatusCode != http.StatusNotFound), nil
}

func (AppInsightsSmartDetectionRule) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Slow page load time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsSmartDetectionRule) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Slow page load time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false

  send_emails_to_subscription_owners = false
  additional_email_recipients        = ["test@example.com", "test2@example.com"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsSmartDetectionRule) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Slow page load time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}

resource "azurerm_application_insights_smart_detection_rule" "test2" {
  name                    = "Slow server response time"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AppInsightsSmartDetectionRule) longDependencyDuration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_smart_detection_rule" "test" {
  name                    = "Long dependency duration"
  application_insights_id = azurerm_application_insights.test.id
  enabled                 = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
