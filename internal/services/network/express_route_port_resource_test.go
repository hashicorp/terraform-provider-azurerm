// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExpressRoutePortResource struct{}

const ARMTestExpressRoutePortAdminState = "ARM_TEST_EXPRESS_ROUTE_PORT_ADMIN_STATE"

func TestAccExpressRoutePort_basic(t *testing.T) {
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
				check.That(data.ResourceName).Key("billing_type").Exists(),
				check.That(data.ResourceName).Key("guid").Exists(),
				check.That(data.ResourceName).Key("mtu").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccExpressRoutePort_adminState(t *testing.T) {
	if _, ok := os.LookupEnv(ARMTestExpressRoutePortAdminState); !ok {
		t.Skipf("Enabling admin state will cause high cost, please set environment variable %q if you want to test it.", ARMTestExpressRoutePortAdminState)
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

func TestAccExpressRoutePort_requiresImport(t *testing.T) {
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

func TestAccExpressRoutePort_userAssignedIdentity(t *testing.T) {
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

func TestAccExpressRoutePort_linkCipher(t *testing.T) {
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
	client := clients.Network.ExpressRoutePorts

	id, err := expressrouteports.ParseExpressRoutePortID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ExpressRoutePortResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_port" "test" {
  name                = "acctestERP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Equinix-London-LDS"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
  billing_type        = "MeteredData"
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
  peering_location    = "Equinix-Hong-Kong-HK1"
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
  purge_protection_enabled = false

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = [
      "Get",
    ]
  }
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge"
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
  peering_location    = "Airtel-Chennai2-CLS"
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
    macsec_sci_enabled            = true
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
