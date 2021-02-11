package applicationinsights_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppInsightsSmartDetectionRule struct {
}

func TestAccApplicationInsightsSmartDetectionRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccApplicationInsightsSmartDetectionRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_smart_detection_rule", "test")
	r := AppInsightsSmartDetectionRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_application_insights_smart_detection_rule"),
		},
	})
}

func (t AppInsightsSmartDetectionRule) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func (AppInsightsSmartDetectionRule) requiresImport(data acceptance.TestData) string {
	template := AppInsightsSmartDetectionRule{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_smart_detection_rule" "import" {
  name                    = azurerm_application_insights_smart_detection_rule.test.name
  application_insights_id = azurerm_application_insights_smart_detection_rule.test.application_insights_id
}
`, template)
}
