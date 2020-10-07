package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSecurityCenterAutomation_logicApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMSecurityCenterAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterAutomation_logicApp(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMSecurityCenterAutomationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSecurityCenterAutomation_logAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMSecurityCenterAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterAutomation_logAnalytics(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMSecurityCenterAutomationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSecurityCenterAutomation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccAzureRMSecurityCenterAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterAutomation_logicApp(data),
				Check: resource.ComposeTestCheckFunc(
					testAccAzureRMSecurityCenterAutomationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSecurityCenterAutomation_requiresImport),
		},
	})
}

func testAccAzureRMSecurityCenterAutomationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.AutomationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Security Center automation: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Security Center automation %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on Security Center automation: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSecurityCenterAutomationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.AutomationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_security_center_automation" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Security Center automation still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMSecurityCenterAutomation_logicApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "LogicApp"
		resource_id = azurerm_logic_app_workflow.test.id
		trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
		event_source = "Alerts"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMSecurityCenterAutomation_logAnalytics(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlogs-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Free"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_security_center_automation" "test" {
  name                = "acctestautomation-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "LogAnalytics"
    resource_id = azurerm_log_analytics_workspace.test.id
  }

  source {
		event_source = "Alerts"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMSecurityCenterAutomation_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_automation" "import" {
  name                = azurerm_security_center_automation.test.name
  location            = azurerm_security_center_automation.test.location
  resource_group_name = azurerm_security_center_automation.test.resource_group_name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "LogicApp"
		resource_id = azurerm_logic_app_workflow.test.id
		trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
		event_source = "Alerts"
  }
}
`, testAccAzureRMSecurityCenterAutomation_logicApp(data))
}
