package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSentinelAlertRuleAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_action", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleAction_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleActionExists(data.ResourceName),
				),
			},
			data.ImportStep("logic_app_trigger_name"),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleAction_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_action", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleAction_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleActionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSentinelAlertRuleAction_requiresImport),
		},
	})
}

func testCheckAzureRMSentinelAlertRuleActionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sentinel.AlertRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Sentinel Alert Rule Action not found: %s", resourceName)
		}

		id, err := parse.SentinelAlertRuleActionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.GetAction(ctx, id.ResourceGroup, id.Workspace, id.Rule, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Sentinel Alert Rule Action %q (Resource Group %q / Workspace: %q / Alert Rule: %q) does not exist", id.Name, id.ResourceGroup, id.Workspace, id.Rule)
			}
			return fmt.Errorf("Getting on Sentinel.AlertRules: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSentinelAlertRuleActionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Sentinel.AlertRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sentinel_alert_rule_action" {
			continue
		}

		id, err := parse.SentinelAlertRuleActionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.GetAction(ctx, id.ResourceGroup, id.Workspace, id.Rule, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Getting on Sentinel.AlertRules: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMSentinelAlertRuleAction_basic(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleAction_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_action" "test" {
  name                   = "acctest-AlertRuleAction-%d"
  rule_id                = azurerm_sentinel_alert_rule_scheduled.test.id
  logic_app_id           = azurerm_logic_app_trigger_custom.test.logic_app_id
  logic_app_trigger_name = azurerm_logic_app_trigger_custom.test.name
  depends_on             = [azurerm_logic_app_trigger_custom.test]
}
`, template, data.RandomInteger)
}

func testAccAzureRMSentinelAlertRuleAction_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleAction_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_action" "import" {
  name                   = azurerm_sentinel_alert_rule_action.test.name
  rule_id                = azurerm_sentinel_alert_rule_action.test.rule_id
  logic_app_id           = azurerm_sentinel_alert_rule_action.test.logic_app_id
  logic_app_trigger_name = azurerm_sentinel_alert_rule_action.test.logic_app_trigger_name
}
`, template)
}

func testAccAzureRMSentinelAlertRuleAction_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%[1]d"
  location = "west europe"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-workspace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "pergb2018"
}

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%[1]d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  display_name               = "Some Rule"
  severity                   = "High"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctest-LogicApp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_logic_app_trigger_custom" "test" {
  name         = "acctest-LogicAppTrigger-%[1]d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "type": "ApiConnectionWebhook",
    "inputs": {
        "body": {
            "callback_url": "@{listCallbackUrl()}"
        },
        "host": {
            "connection": {
                "name": "@parameters('$connections')['azuresentinel']['connectionId']"
            }
        },
        "path": "/subscribe"
    }
}
BODY
}
`, data.RandomInteger)
}
