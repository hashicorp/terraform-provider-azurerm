package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMonitorSmartDetectorAlertRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMMonitorSmartDetectorAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorSmartDetectorAlertRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMMonitorSmartDetectorAlertRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorSmartDetectorAlertRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMMonitorSmartDetectorAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorSmartDetectorAlertRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMMonitorSmartDetectorAlertRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMonitorSmartDetectorAlertRule_requiresImport),
		},
	})
}

func TestAccAzureRMMonitorSmartDetectorAlertRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMMonitorSmartDetectorAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorSmartDetectorAlertRule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMMonitorSmartDetectorAlertRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorSmartDetectorAlertRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMMonitorSmartDetectorAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorSmartDetectorAlertRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMMonitorSmartDetectorAlertRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorSmartDetectorAlertRule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMMonitorSmartDetectorAlertRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorSmartDetectorAlertRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMMonitorSmartDetectorAlertRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMMonitorSmartDetectorAlertRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.SmartDetectorAlertRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can not found Monitor Smart Detector Alert Rule: %s", resourceName)
		}
		id, err := parse.SmartDetectorAlertRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Monitor Smart Detector Alert Rule %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Monitor SmartDetectorAlertRulesClient: %+v", err)
		}
		return nil
	}
}

func testAccAzureRMMonitorSmartDetectorAlertRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_smart_detector_alert_rule" {
			continue
		}
		id, err := parse.SmartDetectorAlertRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Monitor SmartDetectorAlertRulesClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMMonitorSmartDetectorAlertRule_basic(data acceptance.TestData) string {
	template := testAccAzureRMMonitorSmartDetectorAlertRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_smart_detector_alert_rule" "test" {
  name                = "acctestSDAR-%d"
  resource_group_name = azurerm_resource_group.test.name
  severity            = "Sev0"
  scope_resource_ids  = [azurerm_application_insights.test.id]
  frequency           = "PT1M"
  detector_type       = "FailureAnomaliesDetector"

  action_group {
    ids = [azurerm_monitor_action_group.test.id]
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorSmartDetectorAlertRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorSmartDetectorAlertRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_smart_detector_alert_rule" "import" {
  name                = azurerm_monitor_smart_detector_alert_rule.test.name
  resource_group_name = azurerm_monitor_smart_detector_alert_rule.test.resource_group_name
  severity            = azurerm_monitor_smart_detector_alert_rule.test.severity
  scope_resource_ids  = azurerm_monitor_smart_detector_alert_rule.test.scope_resource_ids
  frequency           = azurerm_monitor_smart_detector_alert_rule.test.frequency
  detector_type       = azurerm_monitor_smart_detector_alert_rule.test.detector_type

  action_group {
    ids = [azurerm_monitor_action_group.test.id]
  }
}
`, template)
}

func testAccAzureRMMonitorSmartDetectorAlertRule_complete(data acceptance.TestData) string {
	template := testAccAzureRMMonitorSmartDetectorAlertRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_smart_detector_alert_rule" "test" {
  name                = "acctestSDAR-%d"
  resource_group_name = azurerm_resource_group.test.name
  severity            = "Sev0"
  scope_resource_ids  = [azurerm_application_insights.test.id]
  frequency           = "PT1M"
  detector_type       = "FailureAnomaliesDetector"

  description = "acctest"
  enabled     = false

  action_group {
    ids             = [azurerm_monitor_action_group.test.id]
    email_subject   = "acctest email subject"
    webhook_payload = <<BODY
{
    "msg": "Acctest payload body"
}
BODY
  }

  throttling_duration = "PT20M"
}
`, template, data.RandomInteger)
}

func testAccAzureRMMonitorSmartDetectorAlertRule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-monitor-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
