package fluidrelay_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/sdk/2022-04-21/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay"
)

type FluidRelayResource struct{}

func (f FluidRelayResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fluidrelayservers.ParseFluidRelayServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.FluidRelay.ServerClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Ptr(false), nil
		}
		return nil, fmt.Errorf("retriving %s: %v", id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return utils.Ptr(false), nil
	}
	return utils.Ptr(true), nil
}

var s = fluidrelay.Server{}

func TestAccFluidRelay_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, s.ResourceType(), "test")
	f := FluidRelayResource{}

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
	})
}

func TestAccFluidRelayServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, s.ResourceType(), "test")
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

func (f FluidRelayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appServerDNSAlias-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestRG-userAssignedIdentity-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  identity_type = "SystemAssigned, UserAssigned"
  user_assigned_identity {
     identity_id= azurerm_user_assigned_identity.test.id
  }
  tags = {
    foo = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
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
