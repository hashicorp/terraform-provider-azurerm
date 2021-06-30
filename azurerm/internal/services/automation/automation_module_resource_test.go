package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationModuleResource struct {
}

func TestAccAutomationModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_module", "test")
	r := AutomationModuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("module_link"),
	})
}

func TestAccAutomationModule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_module", "test")
	r := AutomationModuleResource{}

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

func TestAccAutomationModule_multipleModules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_module", "test")
	r := AutomationModuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleModules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("module_link"),
	})
}

func (t AutomationModuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["modules"]

	resp, err := clients.Automation.ModuleClient.Get(ctx, resGroup, accName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Module '%s' (resource group: '%s') does not exist", name, id.ResourceGroup)
	}

	return utils.Bool(resp.ModuleProperties != nil), nil
}

func (AutomationModuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_module" "test" {
  name                    = "xActiveDirectory"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AutomationModuleResource) multipleModules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_module" "test" {
  name                    = "xActiveDirectory"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}

resource "azurerm_automation_module" "second" {
  name                    = "AzureRmMinus"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  module_link {
    uri = "https://www.powershellgallery.com/api/v2/package/AzureRmMinus/0.3.0.0"
  }

  depends_on = [azurerm_automation_module.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AutomationModuleResource) requiresImport(data acceptance.TestData) string {
	template := AutomationModuleResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_module" "import" {
  name                    = azurerm_automation_module.test.name
  resource_group_name     = azurerm_automation_module.test.resource_group_name
  automation_account_name = azurerm_automation_module.test.automation_account_name

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}
`, template)
}
