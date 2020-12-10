package sentinel_test

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

func TestAccSentinelAlertRuleMsSecurityIncident_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSentinelAlertRuleMsSecurityIncident_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentinelAlertRuleMsSecurityIncident_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSentinelAlertRuleMsSecurityIncident_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccSentinelAlertRuleMsSecurityIncident_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSentinelAlertRuleMsSecurityIncident_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccSentinelAlertRuleMsSecurityIncident_requiresImport),
		},
	})
}

func testCheckSentinelAlertRuleMsSecurityIncidentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sentinel.AlertRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Sentinel Alert Rule Ms Security Incident not found: %s", resourceName)
		}

		id, err := parse.SentinelAlertRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.Workspace, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Sentinel Alert Rule Ms Security Incident %q (Resource Group %q / Workspace: %q) does not exist", id.Name, id.ResourceGroup, id.Workspace)
			}
			return fmt.Errorf("Getting on Sentinel.AlertRules: %+v", err)
		}

		return nil
	}
}

func testCheckSentinelAlertRuleMsSecurityIncidentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Sentinel.AlertRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sentinel_alert_rule_ms_security_incident" {
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

func testAccSentinelAlertRuleMsSecurityIncident_basic(data acceptance.TestData) string {
	template := testAccSentinelAlertRuleMsSecurityIncident_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_ms_security_incident" "test" {
  name                       = "acctest-SentinelAlertRule-MSI-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  product_filter             = "Microsoft Cloud App Security"
  display_name               = "some rule"
  severity_filter            = ["High"]
}
`, template, data.RandomInteger)
}

func testAccSentinelAlertRuleMsSecurityIncident_complete(data acceptance.TestData) string {
	template := testAccSentinelAlertRuleMsSecurityIncident_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_ms_security_incident" "test" {
  name                       = "acctest-SentinelAlertRule-MSI-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  product_filter             = "Azure Security Center"
  display_name               = "updated rule"
  severity_filter            = ["High", "Low"]
  description                = "this is a alert rule"
  display_name_filter        = ["alert"]
}
`, template, data.RandomInteger)
}

func testAccSentinelAlertRuleMsSecurityIncident_requiresImport(data acceptance.TestData) string {
	template := testAccSentinelAlertRuleMsSecurityIncident_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_ms_security_incident" "import" {
  name                       = azurerm_sentinel_alert_rule_ms_security_incident.test.name
  log_analytics_workspace_id = azurerm_sentinel_alert_rule_ms_security_incident.test.log_analytics_workspace_id
  product_filter             = azurerm_sentinel_alert_rule_ms_security_incident.test.product_filter
  display_name               = azurerm_sentinel_alert_rule_ms_security_incident.test.display_name
  severity_filter            = azurerm_sentinel_alert_rule_ms_security_incident.test.severity_filter
}
`, template)
}

func testAccSentinelAlertRuleMsSecurityIncident_template(data acceptance.TestData) string {
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
