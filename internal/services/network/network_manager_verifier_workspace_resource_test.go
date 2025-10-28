// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/verifierworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagerVerifierWorkspaceResource struct{}

func testAccNetworkManagerVerifierWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_verifier_workspace", "test")
	r := ManagerVerifierWorkspaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerVerifierWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_verifier_workspace", "test")
	r := ManagerVerifierWorkspaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.complete(data),
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

func testAccNetworkManagerVerifierWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_verifier_workspace", "test")
	r := ManagerVerifierWorkspaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccNetworkManagerVerifierWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_verifier_workspace", "test")
	r := ManagerVerifierWorkspaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagerVerifierWorkspaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := verifierworkspaces.ParseVerifierWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.VerifierWorkspaces.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagerVerifierWorkspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_network_manager_verifier_workspace" "test" {
  name               = "acctest-vw-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerVerifierWorkspaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_verifier_workspace" "import" {
  name               = azurerm_network_manager_verifier_workspace.test.name
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
}
`, r.basic(data))
}

func (r ManagerVerifierWorkspaceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_network_manager_verifier_workspace" "test" {
  name               = "acctest-vw-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  description        = "This is another test verifier workspace"

  tags = {
    foo = "bar"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerVerifierWorkspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_network_manager_verifier_workspace" "test" {
  name               = "acctest-vw-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  description        = "This is a test verifier workspace"

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerVerifierWorkspaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-vw-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-vw-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity"]
}
`, data.RandomInteger, data.Locations.Primary)
}
