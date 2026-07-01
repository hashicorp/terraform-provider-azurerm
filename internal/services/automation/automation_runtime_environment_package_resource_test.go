// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationRuntimeEnvironmentPackageResource struct{}

func TestAccAutomationRuntimeEnvironmentPackage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runtime_environment_package", "test")
	r := AutomationRuntimeEnvironmentPackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content_uri"),
	})
}

func TestAccAutomationRuntimeEnvironmentPackage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runtime_environment_package", "test")
	r := AutomationRuntimeEnvironmentPackageResource{}

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

func TestAccAutomationRuntimeEnvironmentPackage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runtime_environment_package", "test")
	r := AutomationRuntimeEnvironmentPackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content_uri", "content_version", "hash_algorithm", "hash_value"),
	})
}

func (AutomationRuntimeEnvironmentPackageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := packageresource.ParsePackageID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.PackageResource.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r AutomationRuntimeEnvironmentPackageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_runtime_environment_package" "test" {
  name                              = "acctest-package-%[2]d"
  automation_runtime_environment_id = azurerm_automation_runtime_environment.test.id
  content_uri                       = "https://www.powershellgallery.com/api/v2/package/Microsoft.Graph.Authentication/2.25.0"
}
`, r.template(data), data.RandomInteger)
}

func (r AutomationRuntimeEnvironmentPackageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_runtime_environment_package" "import" {
  name                              = azurerm_automation_runtime_environment_package.test.name
  automation_runtime_environment_id = azurerm_automation_runtime_environment_package.test.automation_runtime_environment_id
  content_uri                       = azurerm_automation_runtime_environment_package.test.content_uri
}
`, r.basic(data))
}

func (r AutomationRuntimeEnvironmentPackageResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_runtime_environment_package" "test" {
  name                              = "acctest-package-%[2]d"
  automation_runtime_environment_id = azurerm_automation_runtime_environment.test.id
  content_uri                       = "https://www.powershellgallery.com/api/v2/package/Microsoft.Graph.Authentication/2.25.0"
  content_version                   = "2.25.0"
  hash_algorithm                    = "SHA256"
  hash_value                        = "examplehashvalue"
}
`, r.template(data), data.RandomInteger)
}

func (r AutomationRuntimeEnvironmentPackageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runtime_environment" "test" {
  name                  = "acctest-rte-%[1]d"
  automation_account_id = azurerm_automation_account.test.id
  runtime_language      = "PowerShell"
  runtime_version       = "7.2"
  location              = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}
