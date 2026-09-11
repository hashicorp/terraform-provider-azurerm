// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loadtestservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PlaywrightWorkspaceResource struct{}

func TestAccPlaywrightWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_playwright_workspace", "test")
	r := PlaywrightWorkspaceResource{}

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

func TestAccPlaywrightWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_playwright_workspace", "test")
	r := PlaywrightWorkspaceResource{}

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

func TestAccPlaywrightWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_playwright_workspace", "test")
	r := PlaywrightWorkspaceResource{}

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

func TestAccPlaywrightWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_playwright_workspace", "test")
	r := PlaywrightWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
	})
}

func (PlaywrightWorkspaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := playwrightworkspaces.ParsePlaywrightWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.LoadTestService.PlaywrightWorkspacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (PlaywrightWorkspaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-pww-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Secondary)
}

func (r PlaywrightWorkspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_playwright_workspace" "test" {
  name                = "acctest-pww-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomIntOfLength(8))
}

func (r PlaywrightWorkspaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_playwright_workspace" "import" {
  name                = azurerm_playwright_workspace.test.name
  resource_group_name = azurerm_playwright_workspace.test.resource_group_name
  location            = azurerm_playwright_workspace.test.location
}
`, r.basic(data))
}

func (r PlaywrightWorkspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_playwright_workspace" "test" {
  name                = "acctest-pww-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    Environment = "Sandbox"
    Label       = "Test"
  }
}
`, r.template(data), data.RandomIntOfLength(8))
}
