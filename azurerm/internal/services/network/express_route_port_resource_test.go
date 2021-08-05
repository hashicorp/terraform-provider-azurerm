package network_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ExpressRoutePortResource struct{}

const ARMTestExpressRoutePortAdminState = "ARM_TEST_EXPRESS_ROUTE_PORT_ADMIN_STATE"

func TestAccAzureRMExpressRoutePort_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_port", "test")
	r := ExpressRoutePortResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("link1.0.id").Exists(),
				check.That(data.ResourceName).Key("link1.0.router_name").Exists(),
				check.That(data.ResourceName).Key("link1.0.interface_name").Exists(),
				check.That(data.ResourceName).Key("link1.0.patch_panel_id").Exists(),
				check.That(data.ResourceName).Key("link1.0.rack_id").Exists(),
				check.That(data.ResourceName).Key("link1.0.connector_type").Exists(),
				check.That(data.ResourceName).Key("link2.0.id").Exists(),
				check.That(data.ResourceName).Key("link2.0.router_name").Exists(),
				check.That(data.ResourceName).Key("link2.0.interface_name").Exists(),
				check.That(data.ResourceName).Key("link2.0.patch_panel_id").Exists(),
				check.That(data.ResourceName).Key("link2.0.rack_id").Exists(),
				check.That(data.ResourceName).Key("link2.0.connector_type").Exists(),
				check.That(data.ResourceName).Key("ethertype").Exists(),
				check.That(data.ResourceName).Key("guid").Exists(),
				check.That(data.ResourceName).Key("mtu").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMExpressRoutePort_adminState(t *testing.T) {
	if _, ok := os.LookupEnv(ARMTestExpressRoutePortAdminState); !ok {
		t.Skip(fmt.Sprintf("Enabling admin state will cause high cost, please set environment variable %q if you want to test it.", ARMTestExpressRoutePortAdminState))
	}
	data := acceptance.BuildTestData(t, "azurerm_express_route_port", "test")
	r := ExpressRoutePortResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.adminState(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMExpressRoutePort_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_port", "test")
	r := ExpressRoutePortResource{}

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

func TestAccAzureRMExpressRoutePort_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_port", "test")
	r := ExpressRoutePortResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMExpressRoutePort_linkCipher(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_port", "test")
	r := ExpressRoutePortResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linkCipher(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ExpressRoutePortResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Network.ExpressRoutePortsClient

	id, err := parse.ExpressRoutePortID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Express Route Port %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r ExpressRoutePortResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_port" "test" {
  name                = "acctestERP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Airtel-Chennai2-CLS"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r ExpressRoutePortResource) adminState(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_port" "test" {
  name                = "acctestERP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Area51-ERDirect"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
  link1 {
    admin_enabled = true
  }
}
`, template, data.RandomInteger)
}

func (r ExpressRoutePortResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_port" "import" {
  name                = azurerm_express_route_port.test.name
  resource_group_name = azurerm_express_route_port.test.resource_group_name
  location            = azurerm_express_route_port.test.location
  peering_location    = azurerm_express_route_port.test.peering_location
  bandwidth_in_gbps   = azurerm_express_route_port.test.bandwidth_in_gbps
  encapsulation       = azurerm_express_route_port.test.encapsulation
}
`, template)
}

func (r ExpressRoutePortResource) userAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest1%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_express_route_port" "test" {
  name                = "acctestERP-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "CDC-Canberra"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
	`, template, data.RandomIntOfLength(8))
}

func (r ExpressRoutePortResource) linkCipher(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest1%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestKv-%[2]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "premium"
  soft_delete_enabled      = true
  purge_protection_enabled = false

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = [
      "get",
    ]
  }
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    secret_permissions = [
      "get",
      "set",
      "delete",
      "purge"
    ]
  }
}

resource "azurerm_key_vault_secret" "cak" {
  name         = "cak"
  value        = "ead3664f508eb06c40ac7104cdae4ce5"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault_secret" "ckn" {
  name         = "ckn"
  value        = "dffafc8d7b9a43d5b9a3dfbbf6a30c16"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_express_route_port" "test" {
  name                = "acctestERP-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "CDC-Canberra2"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  link1 {
    macsec_cipher                 = "GcmAes256"
    macsec_ckn_keyvault_secret_id = azurerm_key_vault_secret.ckn.id
    macsec_cak_keyvault_secret_id = azurerm_key_vault_secret.cak.id
  }
  link2 {
    macsec_cipher                 = "GcmAes128"
    macsec_ckn_keyvault_secret_id = azurerm_key_vault_secret.ckn.id
    macsec_cak_keyvault_secret_id = azurerm_key_vault_secret.cak.id
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r ExpressRoutePortResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
