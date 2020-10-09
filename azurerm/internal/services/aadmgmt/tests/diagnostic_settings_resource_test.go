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

func TestAccAzureADDiagnosticSettings_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aad_diagnostic_setting", "test")
	diagSettingName := fmt.Sprintf(`acctest-diagstng-%d`, data.RandomInteger)
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAADDiagnosticsSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAADDiagnosticSettings_basic(data, diagSettingName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAzureADDiagnosticSettingsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureADDiagnosticSettings_Complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aad_diagnostic_setting", "test")
	diagSettingName := fmt.Sprintf(`acctest-diagstng-%d`, data.RandomInteger)
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAADDiagnosticsSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAADDiagnosticSettings_complete(data, diagSettingName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAzureADDiagnosticSettingsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureADDiagnosticSettings_logAnalyticsLog(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aad_diagnostic_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAADDiagnosticsSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAADDiagnosticSettings_law(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAzureADDiagnosticSettingsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureADDiagnosticSettings_eventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aad_diagnostic_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAADDiagnosticsSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAADDiagnosticSettings_eventhub(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAzureADDiagnosticSettingsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureADDiagnosticSettings_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_aad_diagnostic_setting", "test")
	diagSettingName := fmt.Sprintf(`acctest-diagstng-%d`, data.RandomInteger)
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAADDiagnosticsSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAADDiagnosticSettings_complete(data, diagSettingName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAzureADDiagnosticSettingsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAADDiagnosticSettings_basic(data, diagSettingName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAzureADDiagnosticSettingsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAADDiagnosticSettings_complete(data, diagSettingName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAzureADDiagnosticSettingsExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAccAzureADDiagnosticSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).AADManagement.DiagnosticSettingsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("AAD Diagnostic settings resource not found: %s", resourceName)
		}

		diagSettingName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, diagSettingName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("AAD Diagnostic settings %s does not exist: %v", diagSettingName, err)
			}
			return fmt.Errorf("Bad request for get AAD diagnostic settings %s: %v", diagSettingName, err)
		}

		return nil
	}
}

func testAADDiagnosticsSettingsDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).AADManagement.DiagnosticSettingsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_aad_diagnostic_setting" {
			continue
		}

		diagSettingName := rs.Primary.Attributes["name"]
		if resp, err := client.Get(ctx, diagSettingName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return fmt.Errorf("Bad request for get AAD diagnostic settings %s: %v", diagSettingName, err)
		}

		return nil
	}

	return nil
}

func testAADDiagnosticSettings_basic(data acceptance.TestData, name string) string {
	template := testAccAADDiagnosticSettings_template(data)
	return fmt.Sprintf(`%s

resource "azurerm_aad_diagnostic_setting" "test" {
  name                = "%s"
  storage_account_id  = azurerm_storage_account.test.id
  logs  {
	category 	= "NonInteractiveUserSignInLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days =30
	  }
}
}
`, template, name)
}

func testAADDiagnosticSettings_eventhub(data acceptance.TestData) string {
	template := testAccAADDiagnosticSettings_template(data)

	return fmt.Sprintf(`%s
resource "azurerm_aad_diagnostic_setting" "test" {
  name                		= "acctestdiagsetng-%d"
  event_hub_name 			= azurerm_eventhub.test.name
  event_hub_auth_rule_id 	= "${azurerm_eventhub_namespace.test.id}/authorizationRules/RootManageSharedAccessKey"
  logs  {
	category 	= "ProvisioningLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days = 40
	  }
}
}
`, template, data.RandomInteger)
}

func testAADDiagnosticSettings_law(data acceptance.TestData) string {
	template := testAccAADDiagnosticSettings_template(data)
	return fmt.Sprintf(`%s

resource "azurerm_aad_diagnostic_setting" "test" {
  name          	= "acctest-diagsetng-%d"
  workspace_id 		= azurerm_log_analytics_workspace.test.id
  logs  {
	category 	= "ServicePrincipalSignInLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days = 50
	  }
}
}
`, template, data.RandomInteger)
}

func testAADDiagnosticSettings_complete(data acceptance.TestData, diagSettingName string) string {
	template := testAccAADDiagnosticSettings_template(data)
	return fmt.Sprintf(`%s
	
resource "azurerm_aad_diagnostic_setting" "test" {
  name                		= "%s"
  storage_account_id  		= azurerm_storage_account.test.id
  workspace_id 				= azurerm_log_analytics_workspace.test.id
  event_hub_name 			= azurerm_eventhub.test.name
  event_hub_auth_rule_id 	= "${azurerm_eventhub_namespace.test.id}/authorizationRules/RootManageSharedAccessKey"
  logs  {
	category = "AuditLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days =10
	  }
  }
  
  logs  {
	  category 	= "SignInLogs"
	  retention_policy {
		retention_policy_enabled = true
		retention_policy_days =10
	  }
  }

  logs  {
	category 	= "ManagedIdentitySignInLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days = 20
	  }
	}

logs  {
	category 	= "NonInteractiveUserSignInLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days =30
	  }
}

logs  {
	category 	= "ProvisioningLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days = 40
	  }
}

logs  {
	category 	= "ServicePrincipalSignInLogs"
	retention_policy {
		retention_policy_enabled = true
		retention_policy_days = 50
	  }
}
  
}
`, template, diagSettingName)
}

func testAccAADDiagnosticSettings_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubns%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 6
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomIntOfLength(8), data.RandomInteger, data.RandomInteger)
}
