package automation_test

import (
	`context`
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure`
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationDscNodeConfiguration struct {
}

func TestAccAzureRMAutomationDscNodeConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_dsc_nodeconfiguration", "test")
	r := AutomationDscNodeConfiguration{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("configuration_name").HasValue("acctest"),
			),
		},
		data.ImportStep("content_embedded"),
	})
}

func TestAccAzureRMAutomationDscNodeConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_dsc_nodeconfiguration", "test")
	r := AutomationDscNodeConfiguration{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t AutomationDscNodeConfiguration) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["nodeConfigurations"]

	resp, err := clients.Automation.DscNodeConfigurationClient.Get(ctx, resGroup, accName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Dsc Node Configuration %q (resource group: %q): %+v", name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.DscNodeConfigurationProperties != nil), nil
}

func (AutomationDscNodeConfiguration) basic(data acceptance.TestData) string {
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

resource "azurerm_automation_dsc_configuration" "test" {
  name                    = "acctest"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  location                = azurerm_resource_group.test.location
  content_embedded        = "configuration acctest {}"
}

resource "azurerm_automation_dsc_nodeconfiguration" "test" {
  name                    = "acctest.localhost"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  depends_on              = [azurerm_automation_dsc_configuration.test]

  content_embedded = <<mofcontent
instance of MSFT_FileDirectoryConfiguration as $MSFT_FileDirectoryConfiguration1ref
{
  TargetResourceID = "[File]bla";
  Ensure = "Present";
  Contents = "bogus Content";
  DestinationPath = "c:\\bogus.txt";
  ModuleName = "PSDesiredStateConfiguration";
  SourceInfo = "::3::9::file";
  ModuleVersion = "1.0";
  ConfigurationName = "bla";
};
instance of OMI_ConfigurationDocument
{
  Version="2.0.0";
  MinimumCompatibleVersion = "1.0.0";
  CompatibleVersionAdditionalProperties= {"Omi_BaseResource:ConfigurationName"};
  Author="bogusAuthor";
  GenerationDate="06/15/2018 14:06:24";
  GenerationHost="bogusComputer";
  Name="acctest";
};
mofcontent

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AutomationDscNodeConfiguration) requiresImport(data acceptance.TestData) string {
	template := AutomationDscNodeConfiguration{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_dsc_nodeconfiguration" "import" {
  name                    = azurerm_automation_dsc_nodeconfiguration.test.name
  resource_group_name     = azurerm_automation_dsc_nodeconfiguration.test.resource_group_name
  automation_account_name = azurerm_automation_dsc_nodeconfiguration.test.automation_account_name
  content_embedded        = azurerm_automation_dsc_nodeconfiguration.test.content_embedded
}
`, template)
}
