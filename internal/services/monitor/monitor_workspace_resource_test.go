// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkspaceTestResource struct{}

func TestMonitorWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace", "test")
	r := WorkspaceTestResource{}
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

func TestMonitorWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace", "test")
	r := WorkspaceTestResource{}
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

func TestMonitorWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace", "test")
	r := WorkspaceTestResource{}
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

func TestMonitorWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace", "test")
	r := WorkspaceTestResource{}
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

func TestMonitorWorkspace_publicNetworkAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace", "test")
	r := WorkspaceTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicNetworkAccessDisabled(data),
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

func (r WorkspaceTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azuremonitorworkspaces.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Monitor.WorkspacesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkspaceTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r WorkspaceTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-mamw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkspaceTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_workspace" "import" {
  name                = azurerm_monitor_workspace.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, config, data.Locations.Primary)
}

func (r WorkspaceTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-mamw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkspaceTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-mamw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  tags = {
    key2 = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkspaceTestResource) publicNetworkAccessDisabled(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-mamw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  public_network_access_enabled = false
}
`, template, data.RandomInteger, data.Locations.Primary)
}
