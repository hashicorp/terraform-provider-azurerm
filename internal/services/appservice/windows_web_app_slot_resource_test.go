package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsWebAppSlotResource struct{}

func TestAccWindowsWebAppSlot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

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

func TestAccWindowsWebAppSlot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

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

func TestAccWindowsWebAppSlot_autoSwap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoSwap(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// Complete

func TestAccWindowsWebAppSlot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

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

func TestAccWindowsWebAppSlot_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// Block Tests

func TestAccWindowsWebAppSlot_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_backupUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withConnectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withConnectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withConnectionStringsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withConnectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withConnectionStringsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withConnectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withIPRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withIPRestrictionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withIPRestrictionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withAuthSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAuthSettingsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withAutoHealRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccWindowsWebAppSlot_withAutoHealRulesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoHealRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoHealRulesUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("app_settings.secret").HasValue("sauce"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_identityKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssignedKeyVaultIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})

}

// Attributes

func TestAccWindowsWebAppSlot_loadBalancing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loadBalancing(data, "WeightedRoundRobin"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_detailedLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDetailedLogging(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withLoggingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withDetailedLogging(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLogsHttpBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLoggingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// App Stacks

func TestAccWindowsWebAppSlot_withDotNet3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "v3.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withDotNet4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "v4.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withDotNet5(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "v5.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withDotNet6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "v6.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withPhp74(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "7.4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withPython(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.4.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withNode12LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "12-LTS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withNode14LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "14-LTS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withNode16LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "16-LTS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withJava8Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "1.8", "JAVA", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withJava11Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11", "JAVA", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebAppSlot_withDocker(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.docker(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// Deployments
func TestAccWindowsWebAppSlot_zipDeploy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zipDeploy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("zip_deploy_file"),
	})
}

// Exists

func (r WindowsWebAppSlotResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebAppSlotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Windows %s: %+v", id, err)
	}

	if utils.ResponseWasNotFound(resp.Response) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

// Configs

func (r WindowsWebAppSlotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_windows_web_app_slot" "import" {
  name           = azurerm_windows_web_app_slot.test.name
  app_service_id = azurerm_windows_web_app_slot.test.app_service_id

  site_config {}

}
`, r.basic(data))
}

func (r WindowsWebAppSlotResource) appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) autoHealRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    auto_heal_enabled = true

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:01:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:05:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) autoHealRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    auto_heal_enabled = true

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:01:00"
          count             = 10
        }
        status_code {
          status_code_range = "400-404"
          interval          = "00:10:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:10:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) autoSwap(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    auto_swap_slot_name = "production"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) basicWithStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  app_settings = {
    foo = "bar"
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    facebook {
      app_id     = "facebookappid"
      app_secret = "facebookappsecret"

      oauth_scopes = [
        "facebookscope",
      ]
    }
  }

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 1
      frequency_unit     = "Day"
    }
  }

  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 2
      }
    }

    http_logs {
      azure_blob_storage {
        sas_url           = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
        retention_in_days = 3
      }
    }
  }

  client_affinity_enabled    = true
  client_certificate_enabled = true
  client_certificate_mode    = "Optional"

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  enabled    = false
  https_only = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  site_config {
    always_on        = true
    app_command_line = "/sbin/myserver -b 0.0.0.0"
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]
    http2_enabled               = true
    scm_use_main_ip_restriction = true
    local_mysql_enabled         = true
    managed_pipeline_mode       = "Integrated"
    remote_debugging_enabled    = true
    remote_debugging_version    = "VS2019"
    use_32_bit_worker           = true
    websockets_enabled          = true
    ftps_state                  = "FtpsOnly"
    health_check_path           = "/health"
    worker_count                = 1
    minimum_tls_version         = "1.1"
    scm_minimum_tls_version     = "1.1"
    cors {
      allowed_origins = [
        "http://www.contoso.com",
        "www.contoso.com",
      ]

      support_credentials = true
    }

    container_registry_use_managed_identity       = true
    container_registry_managed_identity_client_id = azurerm_user_assigned_identity.test.client_id

    auto_swap_slot_name = "Production"
    auto_heal_enabled   = true

    virtual_application {
      virtual_path  = "/"
      physical_path = "site\\wwwroot"
      preload       = true

      virtual_directory {
        virtual_path  = "/stuff"
        physical_path = "site\\stuff"
      }
    }

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:01:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:05:00"
      }
    }
  }

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "/mounts/files"
  }

  tags = {
    Environment = "AccTest"
    foo         = "bar"
  }
}
`, r.templateWithStorageAccount(data), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) withDetailedLogging(data acceptance.TestData, detailedErrorLogging bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  logs {
    detailed_error_messages = %t
  }
}
`, r.baseTemplate(data), data.RandomInteger, detailedErrorLogging)
}

func (r WindowsWebAppSlotResource) loadBalancing(data acceptance.TestData, loadBalancingMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    load_balancing_mode = "%s"
  }
}
`, r.baseTemplate(data), data.RandomInteger, loadBalancingMode)
}

func (r WindowsWebAppSlotResource) withAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = ["https://example.com"]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) withAuthSettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = ["https://example.com"]

    default_provider              = "AzureActiveDirectory"
    token_refresh_extension_hours = 24
    token_store_enabled           = true
    unauthenticated_client_action = "RedirectToLoginPage"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    facebook {
      app_id     = "facebookappid"
      app_secret = "facebookappsecret"

      oauth_scopes = [
        "facebookscope",
      ]
    }

    google {
      client_id     = "googleclientid"
      client_secret = "googleclientsecret"

      oauth_scopes = [
        "googlescope",
      ]
    }

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }

    twitter {
      consumer_key    = "twitterconsumerkey"
      consumer_secret = "twitterconsumersecret"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) withBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 7
      frequency_unit     = "Day"
    }
  }
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) withConnectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) withConnectionStringsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) withIPRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) withIPRestrictionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da", "6bde7211-57bc-4476-866a-c9676e22b9d7"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64", "9.9.9.8/32"]
        x_forwarded_host  = ["example.com", "anotherexample.com"]
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) withLogsHttpBlob(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 2
      }
    }

    http_logs {
      azure_blob_storage {
        sas_url           = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
        retention_in_days = 3
      }
    }
  }
}

`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) withLoggingComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  logs {
    detailed_error_messages = true
    failed_request_tracing  = true

    application_logs {
      file_system_level = "Warning"

      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 7
      }
    }

    http_logs {
      file_system {
        retention_in_days = 4
        retention_in_mb   = 25
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) dotNet(data acceptance.TestData, dotNetVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    application_stack {
      dotnet_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, dotNetVersion)
}

func (r WindowsWebAppSlotResource) php(data acceptance.TestData, phpVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    application_stack {
      php_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, phpVersion)
}

func (r WindowsWebAppSlotResource) python(data acceptance.TestData, pythonVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    application_stack {
      current_stack  = "python"
      python_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, pythonVersion)
}

func (r WindowsWebAppSlotResource) node(data acceptance.TestData, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    application_stack {
      node_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, nodeVersion)
}

func (r WindowsWebAppSlotResource) java(data acceptance.TestData, javaVersion string, javaContainer string, javaContainerVersion string) string {
	javaContainerStr := ""
	if javaContainer != "" {
		javaContainerStr = fmt.Sprintf("java_container = %q", javaContainer)
	}
	javaContainerVersionStr := ""
	if javaContainerVersion != "" {
		javaContainerVersionStr = fmt.Sprintf("java_container_version = %q", javaContainerVersion)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {
    application_stack {
      current_stack = "java"
      java_version  = "%s"
      %s
      %s
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, javaVersion, javaContainerStr, javaContainerVersionStr)
}

func (r WindowsWebAppSlotResource) docker(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL"          = "https://index.docker.io"
    "DOCKER_REGISTRY_SERVER_USERNAME"     = ""
    "DOCKER_REGISTRY_SERVER_PASSWORD"     = ""
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  }

  site_config {
    application_stack {
      docker_container_name = "%s"
      docker_container_tag  = "%s"
    }
  }
}
`, r.premiumV3PlanContainerTemplate(data), data.RandomInteger, "hello-world", "latest")
}

func (r WindowsWebAppSlotResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  identity {
    type = "SystemAssigned"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) identityUserAssignedKeyVaultIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "kv" {
  name                = "acctestUAI-kv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.kv.id]
  }

  key_vault_reference_identity_id = azurerm_user_assigned_identity.kv.id
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppSlotResource) zipDeploy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  app_settings = {
    WEBSITE_RUN_FROM_PACKAGE       = "1"
    SCM_DO_BUILD_DURING_DEPLOYMENT = "true"
  }

  site_config {
    application_stack {
      dotnet_version = "v6.0"
      current_stack  = "dotnet"
    }
  }

  zip_deploy_file = "./testdata/dotnet-zipdeploy.zip"
}

`, r.baseTemplate(data), data.RandomInteger)
}

// Templates

func (WindowsWebAppSlotResource) baseTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-WAS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "S1"
}

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r WindowsWebAppSlotResource) templateWithStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_share" "test" {
  name                 = "test"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-01"
  expiry = "2024-03-30"

  permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.RandomString)
}

func (WindowsWebAppSlotResource) premiumV3PlanContainerTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "P1v3"
  os_type             = "WindowsContainer"
}

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}
