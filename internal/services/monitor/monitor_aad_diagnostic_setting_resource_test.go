// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorAADDiagnosticSettingResource struct{}

// NOTE: this is a combined test rather than separate split out tests due to
// Azure only being happy about provisioning five per subscription at once and
// there are existing resource in the test subscription hard to clear.
// (which our test suite can't easily workaround)
func TestAccMonitorAADDiagnosticSetting(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"eventhubDefault":       testAccMonitorAADDiagnosticSetting_eventhubDefault,
			"eventhub":              testAccMonitorAADDiagnosticSetting_eventhub,
			"requiresImport":        testAccMonitorAADDiagnosticSetting_requiresImport,
			"logAnalyticsWorkspace": testAccMonitorAADDiagnosticSetting_logAnalyticsWorkspace,
			"storageAccount":        testAccMonitorAADDiagnosticSetting_storageAccount,
			"storageAccountUpdate":  testAccMonitorAADDiagnosticSetting_updateToEnabledLog,
			"updateEnabledLog":      testAccMonitorAADDiagnosticSetting_updateEnabledLog,
			"updateToDisabled":      testAccMonitorAADDiagnosticSetting_updateToDisabled, // remove this test in 4.0 version
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

func testAccMonitorAADDiagnosticSetting_eventhubDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventhubDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorAADDiagnosticSetting_eventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorAADDiagnosticSetting_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_aad_diagnostic_setting"),
		},
	})
}

func testAccMonitorAADDiagnosticSetting_logAnalyticsWorkspace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.logAnalyticsWorkspace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorAADDiagnosticSetting_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorAADDiagnosticSetting_updateToEnabledLog(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiEnabledLog(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.retentionDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.storageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorAADDiagnosticSetting_updateEnabledLog(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleEnabledLog(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiEnabledLog(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleEnabledLog(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorAADDiagnosticSetting_updateToDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_aad_diagnostic_setting", "test")
	r := MonitorAADDiagnosticSettingResource{}

	if features.FourPointOhBeta() {
		t.Skip("remove this test in 4.0 version")
	}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.disabled(data),
			ExpectError: regexp.MustCompile("at least one of the `log`"),
		},
	})
}

func (t MonitorAADDiagnosticSettingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := diagnosticsettings.ParseDiagnosticSettingID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Monitor.AADDiagnosticSettingsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (MonitorAADDiagnosticSettingResource) eventhub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name
  enabled_log {
    category = "SignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "AuditLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "NonInteractiveUserSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "ServicePrincipalSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "RiskyUsers"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "UserRiskEvents"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "B2CRequestLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MonitorAADDiagnosticSettingResource) eventhubDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  enabled_log {
    category = "SignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "AuditLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "NonInteractiveUserSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "ServicePrincipalSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "RiskyUsers"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "UserRiskEvents"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "B2CRequestLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorAADDiagnosticSettingResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_aad_diagnostic_setting" "import" {
  name                           = azurerm_monitor_aad_diagnostic_setting.test.name
  eventhub_authorization_rule_id = azurerm_monitor_aad_diagnostic_setting.test.eventhub_authorization_rule_id
  eventhub_name                  = azurerm_monitor_aad_diagnostic_setting.test.eventhub_name

  enabled_log {
    category = "SignInLogs"
    retention_policy {}
  }
}
`, r.eventhub(data))
}

func (MonitorAADDiagnosticSettingResource) logAnalyticsWorkspace(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[1]d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  enabled_log {
    category = "SignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "AuditLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "NonInteractiveUserSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "ServicePrincipalSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "RiskyUsers"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "UserRiskEvents"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "B2CRequestLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MonitorAADDiagnosticSettingResource) storageAccount(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  log {
    category = "SignInLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "AuditLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "NonInteractiveUserSignInLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "ServicePrincipalSignInLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "ManagedIdentitySignInLogs"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "ProvisioningLogs"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "ADFSSignInLogs"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "RiskyUsers"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "UserRiskEvents"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "NetworkAccessTrafficLogs"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "RiskyServicePrincipals"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "ServicePrincipalRiskEvents"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "B2CRequestLogs"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  log {
    category = "EnrichedOffice365AuditLogs"
    enabled  = false
    retention_policy {}
  }
  log {
    category = "MicrosoftGraphActivityLogs"
    enabled  = false
    retention_policy {}
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  enabled_log {
    category = "SignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "AuditLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "NonInteractiveUserSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "ServicePrincipalSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "RiskyUsers"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "UserRiskEvents"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "B2CRequestLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))
}

func (MonitorAADDiagnosticSettingResource) multiEnabledLog(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  enabled_log {
    category = "AuditLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
  enabled_log {
    category = "SignInLogs"
    retention_policy {
      enabled = true
      days    = 2
    }
  }
  enabled_log {
    category = "NonInteractiveUserSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))
}

func (MonitorAADDiagnosticSettingResource) singleEnabledLog(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  enabled_log {
    category = "NonInteractiveUserSignInLogs"
    retention_policy {
      enabled = true
      days    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))
}

func (MonitorAADDiagnosticSettingResource) retentionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  enabled_log {
    category = "AuditLogs"
    retention_policy {}
  }
  enabled_log {
    category = "SignInLogs"
    retention_policy {}
  }
  enabled_log {
    category = "NonInteractiveUserSignInLogs"
    retention_policy {
      enabled = false
      days    = 3
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))
}

// remove this in 4.0 version
func (MonitorAADDiagnosticSettingResource) disabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_aad_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  log {
    category = "SignInLogs"
    enabled  = false
    retention_policy {
      enabled = true
      days    = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))

}
