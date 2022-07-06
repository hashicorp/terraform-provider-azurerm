package appservice_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const windowsPythonVersionString = "3.4.0"

type WindowsWebAppResource struct{}

func TestAccWindowsWebApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_freeSkuAlwaysOnShouldFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.windowsFreeSku(data),
			ExpectError: regexp.MustCompile("always_on cannot be set to true when using Free, F1, D1 Sku"),
		},
	})
}

func TestAccWindowsWebApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_completeUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_backupUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_withConnectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withConnectionStringsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withLoggingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_withLoggingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_updateServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.secondServicePlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_loadBalancing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withIPRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withIPRestrictionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withAuthSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_withStorageAccountUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withStorageAccountUpdate(data),
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

func TestAccWindowsWebApp_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_identityKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

// Windows Specific
func TestAccWindowsWebApp_virtualDirectories(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualDirectories(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_virtualDirectoriesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.virtualDirectories(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

// App Stacks
func TestAccWindowsWebApp_withDotNet3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNetCore(data, "v3.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_withDotNet4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withDotNet5(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withPhp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withPython34(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, windowsPythonVersionString),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_withPythonUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.python(data, windowsPythonVersionString),
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

func TestAccWindowsWebApp_withJava8Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withJava11Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_basicDockerContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.docker(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.windows_fx_version").HasValue("DOCKER|hello-world:latest"),
			),
		},
		data.ImportStep(),
	})
}

// TODO: More Java matrix tests...

func TestAccWindowsWebApp_withNode12lts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withNode14lts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withMultiStack(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiStack(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_updateAppStack(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11", "JAVA", "9.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dotNet(data, "v5.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_removeAppStack(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11", "JAVA", "9.3"),
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

func TestAccWindowsWebApp_containerRegistryCredentials(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.acrCredentials(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_containerRegistryCredentialsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.container_registry_use_managed_identity").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.acrCredentials(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.container_registry_use_managed_identity").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.container_registry_managed_identity_client_id").HasValue(""),
			),
		},
		data.ImportStep(),
		{
			Config: r.acrCredentialsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.container_registry_use_managed_identity").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.container_registry_managed_identity_client_id").IsSet(),
			),
		},
		data.ImportStep(),
		{
			Config: r.acrCredentials(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.container_registry_use_managed_identity").HasValue("true"),
				check.That(data.ResourceName).Key("site_config.0.container_registry_managed_identity_client_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

// TODO - Needs more property tests for autoheal

func TestAccWindowsWebApp_withAutoHealRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_withAutoHealRulesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_withAutoHealRulesStatusCodeRange(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRulesStatusCodeRange(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_withAutoHealRulesSlowRequest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRulesSlowRequest(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_stickySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsWebApp_stickySettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings").DoesNotExist(),
				check.That(data.ResourceName).Key("sticky_settings").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettingsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettingsRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

// Deployments

func TestAccWindowsWebApp_zipDeploy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

// ASE based tests - Deliberately have longer prefix to make it possible to exclude from testing unrelated changes in the app resource
// as they take a significant amount of time to execute (anything up to 6h)

func TestAccWindowsWebAppASEV3_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withASEV3(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r WindowsWebAppResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.Get(ctx, id.ResourceGroup, id.SiteName)
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

func (r WindowsWebAppResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) windowsFreeSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, r.windowsFreeSkuTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) withBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 1
      frequency_unit     = "Day"
    }
  }

  site_config {}
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r WindowsWebAppResource) withConnectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) withConnectionStringsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (r WindowsWebAppResource) secondServicePlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_service_plan" "test2" {
  name                = "acctestASP2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "S1"

}

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test2.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) virtualDirectories(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    virtual_application {
      virtual_path  = "/"
      physical_path = "site\\wwwroot"
      preload       = true

      virtual_directory {
        virtual_path  = "/stuff"
        physical_path = "site\\stuff"
      }
    }

    virtual_application {
      virtual_path  = "/static-content"
      physical_path = "site\\static"
      preload       = true

      virtual_directory {
        virtual_path  = "/images"
        physical_path = "site\\static\\images"
      }

      virtual_directory {
        virtual_path  = "/css"
        physical_path = "site\\static\\css"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) withAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsWebAppResource) withStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "\\mounts\\files"
  }

  site_config {}
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r WindowsWebAppResource) withStorageAccountUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account {
    name         = "updatedfiles"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "\\mounts\\otherfiles"
  }

  site_config {}
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r WindowsWebAppResource) withAuthSettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsWebAppResource) withDetailedLogging(data acceptance.TestData, detailedErrorLogging bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  logs {
    detailed_error_messages = %t
  }

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, detailedErrorLogging)
}

func (r WindowsWebAppResource) withLoggingComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) withIPRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (r WindowsWebAppResource) withIPRestrictionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (r WindowsWebAppResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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
  //client_certificate_mode    = "Optional"

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
    always_on = true
    // api_management_config_id = // TODO
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
    use_32_bit_worker           = false
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

    // auto_swap_slot_name = // TODO
    auto_heal_enabled = true

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

  sticky_settings {
    app_setting_names       = ["foo"]
    connection_string_names = ["First", "Third"]
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

func (r WindowsWebAppResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  app_settings = {
    foo    = "bar"
    SECRET = "sauce"
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value_new"
    }

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecretNew"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    facebook {
      app_id     = "updatedfacebookappid"
      app_secret = "updatedfacebookappsecret"

      oauth_scopes = [
        "facebookscope",
        "facebookscope2"
      ]
    }
  }

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 12
      frequency_unit     = "Hour"
    }
  }

  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Warning"
        sas_url           = "http://x.com/"
        retention_in_days = 7
      }
    }

    http_logs {
      azure_blob_storage {
        sas_url           = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
        retention_in_days = 5
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

  enabled    = true
  https_only = true

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  site_config {
    always_on = true
    // api_management_config_id = // TODO
    app_command_line = "/sbin/myserver -b 0.0.0.0"
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]
    http2_enabled               = false
    scm_use_main_ip_restriction = false
    local_mysql_enabled         = false
    managed_pipeline_mode       = "Integrated"
    remote_debugging_enabled    = true
    remote_debugging_version    = "VS2017"
    websockets_enabled          = true
    ftps_state                  = "FtpsOnly"
    health_check_path           = "/health2"
    worker_count                = 2
    minimum_tls_version         = "1.2"
    scm_minimum_tls_version     = "1.2"
    cors {
      allowed_origins = [
        "http://www.contoso.com",
        "www.contoso.com",
        "contoso.com",
      ]

      support_credentials = true
    }

    container_registry_use_managed_identity = true

    auto_heal_enabled = true

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:05:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:05:00"
      }
    }
    // auto_swap_slot_name = // TODO - Not supported yet
  }

  tags = {
    foo = "bar"
  }
}
`, r.templateWithStorageAccount(data), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsWebAppResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_web_app" "import" {
  name                = azurerm_windows_web_app.test.name
  location            = azurerm_windows_web_app.test.location
  resource_group_name = azurerm_windows_web_app.test.resource_group_name
  service_plan_id     = azurerm_windows_web_app.test.service_plan_id

  site_config {}
}
`, r.basic(data))
}

func (r WindowsWebAppResource) loadBalancing(data acceptance.TestData, loadBalancingMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    load_balancing_mode = "%s"
  }
}
`, r.baseTemplate(data), data.RandomInteger, loadBalancingMode)
}

func (r WindowsWebAppResource) dotNet(data acceptance.TestData, dotNetVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      dotnet_version = "%s"
      current_stack  = "dotnet"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, dotNetVersion)
}

func (r WindowsWebAppResource) dotNetCore(data acceptance.TestData, dotNetVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      dotnet_version = "%s"
      current_stack  = "dotnetcore"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, dotNetVersion)
}

func (r WindowsWebAppResource) docker(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (r WindowsWebAppResource) node(data acceptance.TestData, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      node_version  = "%s"
      current_stack = "node"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, nodeVersion)
}

func (r WindowsWebAppResource) php(data acceptance.TestData, phpVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      php_version   = "%s"
      current_stack = "php"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, phpVersion)
}

func (r WindowsWebAppResource) python(data acceptance.TestData, pythonVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      python_version = "%s"
      current_stack  = "python"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, pythonVersion)
}

func (r WindowsWebAppResource) java(data acceptance.TestData, javaVersion string, javaContainer string, javaContainerVersion string) string {
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

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (r WindowsWebAppResource) multiStack(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      dotnet_version         = "%s"
      php_version            = "%s"
      python_version         = "%s"
      java_version           = "%s"
      java_container         = "%s"
      java_container_version = "%s"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, "v4.0", "7.4", windowsPythonVersionString, "1.8", "TOMCAT", "9.0")
}

func (r WindowsWebAppResource) withLogsHttpBlob(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

  site_config {}
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r WindowsWebAppResource) autoHealRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (r WindowsWebAppResource) autoHealRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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
        action_type                    = "LogEvent"
        minimum_process_execution_time = "00:10:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) autoHealRulesStatusCodeRange(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    auto_heal_enabled = true

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500-599"
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

func (r WindowsWebAppResource) autoHealRulesSlowRequest(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    auto_heal_enabled = true

    auto_heal_setting {
      trigger {
        slow_request {
          count      = "10"
          interval   = "00:10:00"
          time_taken = "00:00:10"
          path       = null
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

// Note - this test omits the features block as the referenced template is a complete test on Service Plan
func (r WindowsWebAppResource) withASEV3(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, ServicePlanResource{}.aseV3(data), data.RandomInteger)
}

func (r WindowsWebAppResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s
resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  identity {
    type = "SystemAssigned"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
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

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) identityUserAssigned(data acceptance.TestData) string {
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

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) identityUserAssignedKeyVaultIdentity(data acceptance.TestData) string {
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

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id
  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.kv.id]
  }

  key_vault_reference_identity_id = azurerm_user_assigned_identity.kv.id
}
`, r.baseTemplate(data), data.RandomInteger)
}

// Misc Properties

func (r WindowsWebAppResource) acrCredentials(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    container_registry_use_managed_identity = true
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) acrCredentialsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    container_registry_use_managed_identity       = true
    container_registry_managed_identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) stickySettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  app_settings = {
    foo                                     = "bar"
    secret                                  = "sauce"
    third                                   = "degree"
    "Special chars: !@#$%%^&*()_+-=' \";/?" = "Supported by the Azure portal"
  }

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

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "Special chars: !@#$%%^&*()_+-=' \";/?"
    value = "characters-supported-by-the-Azure-portal"
    type  = "Custom"
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret", "Special chars: !@#$%%^&*()_+-=' \";/?"]
    connection_string_names = ["First", "Third", "Special chars: !@#$%%^&*()_+-=' \";/?"]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) stickySettingsRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
    third  = "degree"
  }

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

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) stickySettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
    third  = "degree"
  }

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

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret", "third"]
    connection_string_names = ["First", "Second", "Third"]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r WindowsWebAppResource) zipDeploy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (WindowsWebAppResource) baseTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "S1"

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (WindowsWebAppResource) windowsFreeSkuTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "F1"

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (WindowsWebAppResource) premiumV3PlanContainerTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "P1v3"
  os_type             = "WindowsContainer"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WindowsWebAppResource) templateWithStorageAccount(data acceptance.TestData) string {
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
