package loganalytics_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMLogAnalyticsLinkedService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsLinkedService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsLinkedServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestlaw-%d/Automation", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "workspace_name", fmt.Sprintf("acctestlaw-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "linked_service_name", "automation"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsLinkedService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsLinkedService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsLinkedServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestlaw-%d/Automation", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "workspace_name", fmt.Sprintf("acctestlaw-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "linked_service_name", "automation"),
				),
			},
			{
				Config:      testAccAzureRMLogAnalyticsLinkedService_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_linked_service"),
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsLinkedService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsLinkedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsLinkedService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsLinkedServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "linked_service_name", "automation"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMLogAnalyticsLinkedServiceDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.LinkedServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_linked_service" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		lsName := rs.Primary.Attributes["linked_service_name"]

		resp, err := conn.Get(ctx, resourceGroup, workspaceName, lsName)
		if err != nil {
			return nil
		}
		if resp.ID == nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Log Analytics Linked Service still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMLogAnalyticsLinkedServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.LinkedServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		lsName := rs.Primary.Attributes["linked_service_name"]
		name := rs.Primary.Attributes["name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Log Analytics Linked Service: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, workspaceName, lsName)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Linked Service Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Linked Service '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogAnalyticsLinkedService_basic(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsLinkedService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  resource_id         = azurerm_automation_account.test.id
}
`, template)
}

func testAccAzureRMLogAnalyticsLinkedService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsLinkedService_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "import" {
  resource_group_name = azurerm_log_analytics_linked_service.test.resource_group_name
  workspace_name      = azurerm_log_analytics_linked_service.test.workspace_name
  resource_id         = azurerm_log_analytics_linked_service.test.resource_id
}
`, template)
}

func testAccAzureRMLogAnalyticsLinkedService_complete(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsLinkedService_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  linked_service_name = "automation"
  resource_id         = azurerm_automation_account.test.id
}
`, template)
}

func testAccAzureRMLogAnalyticsLinkedService_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutomation-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Basic"

  tags = {
    Environment = "Test"
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
