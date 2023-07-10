// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fluidrelay_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FluidRelayResource struct{}

func TestAccFluidRelay_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	f := FluidRelayResource{}

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
				check.That(data.ResourceName).Key("frs_tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("service_endpoints.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFluidRelay_storageBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	f := FluidRelayResource{}

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.storageBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
				check.That(data.ResourceName).Key("frs_tenant_id").IsUUID(),
			),
		},
		data.ImportStep("storage_sku"),
	})
}

func TestAccFluidRelay_ami(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	f := FluidRelayResource{}

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.userAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.systemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.systemAndUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFluidRelayServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	var f FluidRelayResource

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.RequiresImportErrorStep(f.requiresImport),
	})
}

func TestAccFluidRelayServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	var f FluidRelayResource

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
	})
}

func (f FluidRelayResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fluidrelayservers.ParseFluidRelayServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.FluidRelay.FluidRelayServers.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (f FluidRelayResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fluidrelay-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayResource) templateWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestRG-userAssignedIdentity-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, f.template(data), data.RandomInteger)
}

func (f FluidRelayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  tags = {
    foo = "bar"
  }
}
`, f.template(data), data.RandomInteger, data.Locations.Primary)
}

// basic storage sku only work with east asia and south-ease-asia
func (f FluidRelayResource) storageBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "SouthEastAsia"
  storage_sku         = "basic"
  tags = {
    foo = "bar"
  }
}
`, f.template(data), data.RandomInteger)
}

func (f FluidRelayResource) userAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, f.templateWithIdentity(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayResource) systemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  identity {
    type = "SystemAssigned"
  }
  tags = {
    foo = "bar"
  }
}
`, f.templateWithIdentity(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayResource) systemAndUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  tags = {
    foo = "bar"
  }
}
`, f.templateWithIdentity(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_fluid_relay_server" "import" {
  name                = azurerm_fluid_relay_server.test.name
  resource_group_name = azurerm_fluid_relay_server.test.resource_group_name
  location            = azurerm_fluid_relay_server.test.location
}
`, f.basic(data))
}

func (f FluidRelayResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  tags = {
    foo = "bar2"
  }
}
`, f.template(data), data.RandomInteger, data.Locations.Primary)
}
