// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/dscnodeconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationDscNodeConfigurationResource struct{}

func TestAccAutomationDscNodeConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_dsc_nodeconfiguration", "test")
	r := AutomationDscNodeConfigurationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("configuration_name").HasValue("acctest"),
			),
		},
		data.ImportStep("content_embedded"),
	})
}

func TestAccAutomationDscNodeConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_dsc_nodeconfiguration", "test")
	r := AutomationDscNodeConfigurationResource{}

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

func (t AutomationDscNodeConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dscnodeconfiguration.ParseNodeConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Automation.DscNodeConfiguration.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (AutomationDscNodeConfigurationResource) basic(data acceptance.TestData) string {
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

func (AutomationDscNodeConfigurationResource) requiresImport(data acceptance.TestData) string {
	template := AutomationDscNodeConfigurationResource{}.basic(data)
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
