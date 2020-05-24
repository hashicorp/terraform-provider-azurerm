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

func TestAccAzureRMSentinelAlertRuleScheduled_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleScheduledDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleScheduled_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleScheduledExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleScheduled_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleScheduledDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleScheduled_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleScheduledExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleScheduled_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleScheduledDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleScheduled_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleScheduledExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSentinelAlertRuleScheduled_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleScheduledExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSentinelAlertRuleScheduled_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleScheduledExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleScheduled_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleScheduledDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleScheduled_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleScheduledExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSentinelAlertRuleScheduled_requiresImport),
		},
	})
}

func testCheckAzureRMSentinelAlertRuleScheduledExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sentinel.AlertRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Sentinel Alert Rule Scheduled not found: %s", resourceName)
		}

		id, err := parse.SentinelAlertRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Sentinel Alert Rule Scheduled %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Sentinel.AlertRules: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSentinelAlertRuleScheduledDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Sentinel.AlertRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sentinel_alert_rule_scheduled" {
			continue
		}

		id, err := parse.SentinelAlertRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.Workspace, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Getting on Sentinel.AlertRules: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMSentinelAlertRuleScheduled_basic(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleScheduled_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%d"
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
`, template, data.RandomInteger)
}

func testAccAzureRMSentinelAlertRuleScheduled_complete(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleScheduled_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  display_name               = "Updated Rule"
  description                = "Some Description"
  tactics                    = ["Collection", "CommandAndControl"]
  severity                   = "Low"
  enabled                    = false
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(3d), now(), 1d) by Caller
QUERY
  query_frequency            = "PT20M"
  query_period               = "PT40M"
  trigger_operator           = "Equal"
  trigger_threshold          = 5
  suppression_enabled        = true
  suppression_duration       = "PT40M"
}
`, template, data.RandomInteger)
}

func testAccAzureRMSentinelAlertRuleScheduled_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleScheduled_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "import" {
  name                       = azurerm_sentinel_alert_rule_scheduled.test.name
  log_analytics_workspace_id = azurerm_sentinel_alert_rule_scheduled.test.log_analytics_workspace_id
  display_name               = azurerm_sentinel_alert_rule_scheduled.test.display_name
  severity                   = azurerm_sentinel_alert_rule_scheduled.test.severity
  query                      = azurerm_sentinel_alert_rule_scheduled.test.query
}
`, template)
}

func testAccAzureRMSentinelAlertRuleScheduled_template(data acceptance.TestData) string {
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
