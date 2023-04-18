package automanage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutoManageConfigurationProfileResource struct{}

func TestAccAutoManageConfigurationProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutoManageConfigurationProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutoManageConfigurationProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutoManageConfigurationProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutoManageConfigurationProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AutomanageConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Automanage.ConfigurationClient
	resp, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Response.Response != nil), nil
}

func (r AutoManageConfigurationProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  antimalware {
	enabled = true
	exclusions {
      extensions = "exe;dll"
	  paths = "C:\\Windows\\Temp;D:\\Temp"
      processes = "svchost.exe;notepad.exe"
	}
    real_time_protection = true
    scheduled_scan_enabled = true
    scheduled_scan_type = "Quick"
	scheduled_scan_day = 1
	scheduled_scan_time_in_minutes = 1339
  }
  automation_account_enabled = true
  boot_diagnostics_enabled = true
  change_tracking_and_inventory_enabled = true
  defender_for_cloud_enabled = true
  guest_configuration_enabled = true
  status_change_alert_enabled = true
  tags = {
	"env" = "test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "import" {
  name                = azurerm_automanage_configuration.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  configuration_json = jsonencode({
    "Antimalware/Enable" : false,
    "AzureSecurityCenter/Enable" : true,
    "Backup/Enable" : false,
    "BootDiagnostics/Enable" : true,
    "ChangeTrackingAndInventory/Enable" : true,
    "GuestConfiguration/Enable" : true,
    "LogAnalytics/Enable" : true,
    "UpdateManagement/Enable" : true,
    "VMInsights/Enable" : true
  })
}
`, config, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  antimalware {
	enabled = true
	exclusions {
      extensions = "exe;dll"
	  paths = "C:\\Windows\\Temp;D:\\Temp"
      processes = "svchost.exe;notepad.exe"
	}
    real_time_protection_enabled = true
    scheduled_scan_enabled = true
    scheduled_scan_type = "Quick"
	scheduled_scan_day = 1
	scheduled_scan_time_in_minutes = 1339
  }
  automation_account_enabled = true
  boot_diagnostics_enabled = true
  change_tracking_and_inventory_enabled = true
  defender_for_cloud_enabled = true
  guest_configuration_enabled = true
  status_change_alert_enabled = true
  tags = {
	"env" = "test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  configuration_json = jsonencode({
    "Antimalware/Enable" : false,
    "AzureSecurityCenter/Enable" : true,
    "Backup/Enable" : true,
    "BootDiagnostics/Enable" : true,
    "ChangeTrackingAndInventory/Enable" : true,
    "GuestConfiguration/Enable" : true,
    "LogAnalytics/Enable" : true,
    "UpdateManagement/Enable" : true,
    "VMInsights/Enable" : true
  })
  tags = {
    key2 = "value2"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
