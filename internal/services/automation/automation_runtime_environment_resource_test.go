// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutomationRuntimeEnvironmentResource struct{}

func (a AutomationRuntimeEnvironmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := runtimeenvironment.ParseRuntimeEnvironmentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.RuntimeEnvironment.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving runtime environment %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func TestAccAutomationRuntimeEnvironment_completePowerShell(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.AutomationRuntimeEnvironmentResource{}.ResourceType(), "test")
	r := AutomationRuntimeEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePowerShell(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationRuntimeEnvironment_basicPowerShell(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.AutomationRuntimeEnvironmentResource{}.ResourceType(), "test")
	r := AutomationRuntimeEnvironmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "PowerShell", "7.2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationRuntimeEnvironment_basicPython(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.AutomationRuntimeEnvironmentResource{}.ResourceType(), "test")
	r := AutomationRuntimeEnvironmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Python", "3.10"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationRuntimeEnvironment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.AutomationRuntimeEnvironmentResource{}.ResourceType(), "test")
	r := AutomationRuntimeEnvironmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePowerShell(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completePowerShellAddPackage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationRuntimeEnvironment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.AutomationRuntimeEnvironmentResource{}.ResourceType(), "test")
	r := AutomationRuntimeEnvironmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePowerShell(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (a AutomationRuntimeEnvironmentResource) requiresImport(data acceptance.TestData) string {
	template := a.completePowerShell(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_runtime_environment" "import" {
  name                    = azurerm_automation_runtime_environment.test.name
  resource_group_name     = azurerm_automation_runtime_environment.test.resource_group_name
  automation_account_name = azurerm_automation_runtime_environment.test.automation_account_name
  runtime_language        = azurerm_automation_runtime_environment.test.runtime_language
  runtime_version         = azurerm_automation_runtime_environment.test.runtime_version
  location                = azurerm_automation_runtime_environment.test.location
}
`, template)
}

func (a AutomationRuntimeEnvironmentResource) basic(data acceptance.TestData, runtime, version string) string {
	return fmt.Sprintf(`

%s

resource azurerm_automation_runtime_environment "test" {
  name                    = "acc_%[2]s_basic"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  runtime_language        = "%[2]s"
  runtime_version         = "%[3]s"

  location                = azurerm_resource_group.test.location
}

`, a.template(data), runtime, version)
}

func (a AutomationRuntimeEnvironmentResource) completePowerShell(data acceptance.TestData) string {
	return fmt.Sprintf(`

	%s

resource azurerm_automation_runtime_environment "test" {
  name                    = "powershell_complete"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  runtime_language        = "PowerShell"
  runtime_version         = "7.2"

  location                = azurerm_resource_group.test.location
  description             = "Test automation runtime environment"

  runtime_default_packages = {
    "az" = "11.2.0"
  }

  tags = {
    test = "value"
  }
}

`, a.template(data))
}

func (a AutomationRuntimeEnvironmentResource) completePowerShellAddPackage(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

  resource azurerm_automation_runtime_environment "test" {
    name                    = "powershell_complete"
    resource_group_name     = azurerm_resource_group.test.name
    automation_account_name = azurerm_automation_account.test.name

    runtime_language        = "PowerShell"
    runtime_version         = "7.2"

    location                = azurerm_resource_group.test.location
    description             = "Test automation runtime environment"

    runtime_default_packages = {
      "az" 			= "11.2.0",
      "azure cli" 	= "2.56.0",
    }

    tags = {
      test = "value"
    }
  }

`, a.template(data))
}

func (a AutomationRuntimeEnvironmentResource) template(data acceptance.TestData) string {
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

`, data.RandomInteger, data.Locations.Primary)
}
