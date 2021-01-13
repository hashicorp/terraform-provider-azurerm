package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SentinelAlertRuleFusionResource struct{}

func TestAccSentinelAlertRuleFusion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleFusion_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleFusion_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleFusion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_fusion", "test")
	r := SentinelAlertRuleFusionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r SentinelAlertRuleFusionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	alertRuleClient := client.Sentinel.AlertRulesClient
	id, err := parse.SentinelAlertRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := alertRuleClient.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.Workspace, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Sentinel Alert Rule Fusion (%q): %+v", state.String(), err)
	}

	return utils.Bool(resp.Value != nil), nil
}

func (r SentinelAlertRuleFusionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "Advanced Multistage Attack Detection"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}

resource "azurerm_sentinel_alert_rule_fusion" "test" {
  name                       = "acctest-SentinelAlertRule-Fusion-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleFusionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "Advanced Multistage Attack Detection"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}

resource "azurerm_sentinel_alert_rule_fusion" "test" {
  name                       = "acctest-SentinelAlertRule-Fusion-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test.name
  enabled                    = false
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleFusionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_fusion" "import" {
  name                       = azurerm_sentinel_alert_rule_fusion.test.name
  log_analytics_workspace_id = azurerm_sentinel_alert_rule_fusion.test.log_analytics_workspace_id
  alert_rule_template_guid   = azurerm_sentinel_alert_rule_fusion.test.alert_rule_template_guid
}
`, r.basic(data))
}

func (r SentinelAlertRuleFusionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
