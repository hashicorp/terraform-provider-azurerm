package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this is a combined test rather than separate split out tests due to
// Azure only being happy about provisioning one per subscription at once
// (which our test suite can't easily workaround)

// this occasionally fails due to the rapid provisioning and deprovisioning,
// running the exact same test afterwards always results in a pass.

func TestAccAzureRMMonitorLogProfile(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":          testAccAzureRMMonitorLogProfile_basic,
			"requiresImport": testAccAzureRMMonitorLogProfile_requiresImport,
			"servicebus":     testAccAzureRMMonitorLogProfile_servicebus,
			"complete":       testAccAzureRMMonitorLogProfile_complete,
			"disappears":     testAccAzureRMMonitorLogProfile_disappears,
		},
		"datasource": {
			"eventhub":       testAccDataSourceAzureRMMonitorLogProfile_eventhub,
			"storageaccount": testAccDataSourceAzureRMMonitorLogProfile_storageaccount,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMMonitorLogProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMMonitorLogProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorLogProfile_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_monitor_log_profile"),
			},
		},
	})
}

func testAccAzureRMMonitorLogProfile_servicebus(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_servicebusConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccAzureRMMonitorLogProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_completeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccAzureRMMonitorLogProfile_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorLogProfile_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogProfileExists(data.ResourceName),
					testCheckAzureRMLogProfileDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLogProfileDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.LogProfilesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_log_profile" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Log Profile still exists:\n%#v", *resp.ID)
			}
		}
	}

	return nil
}

func testCheckAzureRMLogProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.LogProfilesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Log Profile %q does not exist", name)
			}

			return fmt.Errorf("Bad: Get on monitorLogProfilesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLogProfileDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.LogProfilesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]

		if _, err := client.Delete(ctx, name); err != nil {
			return fmt.Errorf("Error deleting Log Profile %q: %+v", name, err)
		}

		return nil
	}
}

func testAccAzureRMMonitorLogProfile_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_monitor_log_profile" "test" {
  name = "acctestlp-%d"

  categories = [
    "Action",
  ]

  locations = [
    "%s",
  ]

  storage_account_id = azurerm_storage_account.test.id

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMonitorLogProfile_requiresImportConfig(data acceptance.TestData) string {
	template := testAccAzureRMMonitorLogProfile_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_log_profile" "import" {
  name               = azurerm_monitor_log_profile.test.name
  categories         = azurerm_monitor_log_profile.test.categories
  locations          = azurerm_monitor_log_profile.test.locations
  storage_account_id = azurerm_monitor_log_profile.test.storage_account_id

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, template)
}

func testAccAzureRMMonitorLogProfile_servicebusConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestsbns-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = "acctestsbrule-%s"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = true
  manage = true
}

resource "azurerm_monitor_log_profile" "test" {
  name = "acctestlp-%d"

  categories = [
    "Action",
  ]

  locations = [
    "%s",
  ]

  servicebus_rule_id = azurerm_servicebus_namespace_authorization_rule.test.id

  retention_policy {
    enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMonitorLogProfile_completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctestehns-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_monitor_log_profile" "test" {
  name = "acctestlp-%d"

  categories = [
    "Action",
    "Delete",
    "Write",
  ]

  locations = [
    "%s",
    "%s",
  ]

  # RootManageSharedAccessKey is created by default with listen, send, manage permissions
  servicebus_rule_id = "${azurerm_eventhub_namespace.test.id}/authorizationrules/RootManageSharedAccessKey"
  storage_account_id = azurerm_storage_account.test.id

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
