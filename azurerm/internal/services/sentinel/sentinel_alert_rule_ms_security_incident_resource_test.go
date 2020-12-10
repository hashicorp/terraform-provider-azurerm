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

func TestAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleMsSecurityIncident_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleMsSecurityIncident_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleMsSecurityIncident_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSentinelAlertRuleMsSecurityIncident_requiresImport),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleMsSecurityIncident_withAlertRuleTemplateName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_withAlertRuleTemplateName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSentinelAlertRuleMsSecurityIncident_withDisplayNameExcludeFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSentinelAlertRuleMsSecurityIncidentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_withDisplayNameExcludeFilter(data, "alert3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_withDisplayNameExcludeFilter(data, "alert4"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSentinelAlertRuleMsSecurityIncidentExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMSentinelAlertRuleMsSecurityIncidentDestroy(s *terraform.State) error {
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

func testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleMsSecurityIncident_template(data)
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

func testAccAzureRMSentinelAlertRuleMsSecurityIncident_complete(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleMsSecurityIncident_template(data)
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

func testAccAzureRMSentinelAlertRuleMsSecurityIncident_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data)
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

func testAccAzureRMSentinelAlertRuleMsSecurityIncident_withAlertRuleTemplateName(data acceptance.TestData) string {
	template := testAccAzureRMSentinelAlertRuleMsSecurityIncident_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_ms_security_incident" "test" {
  name                       = "acctest-SentinelAlertRule-MSI-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  product_filter             = "Microsoft Cloud App Security"
  display_name               = "some rule"
  severity_filter            = ["High"]
  alert_rule_template_name   = "b3cfc7c0-092c-481c-a55b-34a3979758cb"
}
`, template, data.RandomInteger)
}

func testAccAzureRMSentinelAlertRuleMsSecurityIncident_withDisplayNameExcludeFilter(data acceptance.TestData, displayNameExcludeFilter string) string {
	template := testAccAzureRMSentinelAlertRuleMsSecurityIncident_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_ms_security_incident" "test" {
  name                        = "acctest-SentinelAlertRule-MSI-%d"
  log_analytics_workspace_id  = azurerm_log_analytics_workspace.test.id
  product_filter              = "Microsoft Cloud App Security"
  display_name                = "some rule"
  severity_filter             = ["High"]
  display_name_filter         = ["alert1"]
  display_name_exclude_filter = ["%s"]
}
`, template, data.RandomInteger, displayNameExcludeFilter)
}

func testAccAzureRMSentinelAlertRuleMsSecurityIncident_template(data acceptance.TestData) string {
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
