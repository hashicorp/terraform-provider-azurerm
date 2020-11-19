package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMonitorActivityLogAlert_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActivityLogAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActivityLogAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActivityLogAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActivityLogAlert_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActivityLogAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActivityLogAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActivityLogAlertExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorActivityLogAlert_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_monitor_activity_log_alert"),
			},
		},
	})
}

func TestAccAzureRMMonitorActivityLogAlert_singleResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActivityLogAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActivityLogAlert_singleResource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActivityLogAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operation_name", "Microsoft.Storage/storageAccounts/write"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "criteria.0.resource_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActivityLogAlert_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActivityLogAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActivityLogAlert_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActivityLogAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This is just a test resource."),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operation_name", "Microsoft.Storage/storageAccounts/write"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.resource_provider", "Microsoft.Storage"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.resource_type", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "criteria.0.resource_group"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "criteria.0.resource_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorActivityLogAlert_basicAndCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorActivityLogAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActivityLogAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.resource_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.caller", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.level", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.status", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "0"),
				),
			},
			{
				Config: testAccAzureRMMonitorActivityLogAlert_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActivityLogAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This is just a test resource."),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operation_name", "Microsoft.Storage/storageAccounts/write"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.resource_provider", "Microsoft.Storage"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.resource_type", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "criteria.0.resource_group"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "criteria.0.resource_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "2"),
				),
			},
			{
				Config: testAccAzureRMMonitorActivityLogAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorActivityLogAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.resource_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.caller", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.level", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.status", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "0"),
				),
			},
		},
	})
}

func testAccAzureRMMonitorActivityLogAlert_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_resource_group.test.id]

  criteria {
    category = "Recommendation"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMonitorActivityLogAlert_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorActivityLogAlert_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_activity_log_alert" "import" {
  name                = azurerm_monitor_activity_log_alert.test.name
  resource_group_name = azurerm_monitor_activity_log_alert.test.resource_group_name
  scopes              = [azurerm_resource_group.test.id]

  criteria {
    category = "Recommendation"
  }
}
`, template)
}

func testAccAzureRMMonitorActivityLogAlert_singleResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_resource_group.test.id]

  criteria {
    operation_name = "Microsoft.Storage/storageAccounts/write"
    category       = "Recommendation"
    resource_id    = azurerm_storage_account.test.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func testAccAzureRMMonitorActivityLogAlert_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag2"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  enabled             = true
  description         = "This is just a test resource."

  scopes = [
    azurerm_resource_group.test.id,
    azurerm_storage_account.test.id,
  ]

  criteria {
    operation_name          = "Microsoft.Storage/storageAccounts/write"
    category                = "Recommendation"
    resource_provider       = "Microsoft.Storage"
    resource_type           = "Microsoft.Storage/storageAccounts"
    resource_group          = azurerm_resource_group.test.name
    resource_id             = azurerm_storage_account.test.id
    recommendation_category = "OperationalExcellence"
    recommendation_impact   = "High"
  }

  action {
    action_group_id = azurerm_monitor_action_group.test1.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test2.id

    webhook_properties = {
      from = "terraform test"
      to   = "microsoft azure"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func testCheckAzureRMMonitorActivityLogAlertDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActivityLogAlertsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_activity_log_alert" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Activity log alert still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMMonitorActivityLogAlertExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ActivityLogAlertsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Activity Log Alert Instance: %s", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorActivityLogAlertsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Activity Log Alert Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}
