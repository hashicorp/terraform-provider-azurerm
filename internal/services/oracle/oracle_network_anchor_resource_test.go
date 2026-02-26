// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkAnchorResource struct{}

func (a NetworkAnchorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkanchors.ParseNetworkAnchorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.NetworkAnchors.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccNetworkAnchorResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.NetworkAnchorResource{}.ResourceType(), "test")
	r := NetworkAnchorResource{}
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

func TestAccNetworkAnchorResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.NetworkAnchorResource{}.ResourceType(), "test")
	r := NetworkAnchorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dns_forwarding_rule_url").IsNotEmpty(),
				check.That(data.ResourceName).Key("dns_forwarding_rule.#").HasValue("0"),
			),
		},
		data.ImportStep("dns_forwarding_rule"),
	})
}

func TestAccNetworkAnchorResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.NetworkAnchorResource{}.ResourceType(), "test")
	r := NetworkAnchorResource{}
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

func TestAccNetworkAnchorResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.NetworkAnchorResource{}.ResourceType(), "test")
	r := NetworkAnchorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
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
	})
}

func (a NetworkAnchorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_network_anchor" "test" {
  location            = "%[3]s"
  name                = "OFakeNA%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  zones               = local.zones

  resource_anchor_id = azurerm_oracle_resource_anchor.test.id
  subnet_id          = azurerm_subnet.virtual_network_subnet.id
}`, a.template(data), data.RandomString, data.Locations.Primary, data.Subscriptions.Primary)
}

func (a NetworkAnchorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_network_anchor" "test" {
  location            = "%[3]s"
  name                = "OFakeNA%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  zones               = local.zones

  resource_anchor_id = azurerm_oracle_resource_anchor.test.id
  subnet_id          = azurerm_subnet.virtual_network_subnet.id

  oci_backup_cidr_block                 = "10.0.0.0/24"
  oci_vcn_dns_label                     = "ociOFakeacctes"
  oracle_dns_listening_endpoint_enabled = true
  oracle_to_azure_dns_zone_sync_enabled = true

  oracle_dns_forwarding_endpoint_enabled = true
  dns_forwarding_rule {
    domain_names          = "abc.ocidelegated.ocinetworkanch.oraclevcn.com"
    forwarding_ip_address = "10.0.1.16"
  }
  dns_forwarding_rule {
    domain_names          = "def.ocidelegated.ocinetworkanch.oraclevcn.com"
    forwarding_ip_address = "10.0.1.24"
  }
  dns_listening_endpoint_allowed_cidrs = "10.0.2.0/24,10.0.3.0/24"
}`, a.template(data), data.RandomString, data.Locations.Primary, data.Subscriptions.Primary)
}

func (a NetworkAnchorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_network_anchor" "test" {
  location            = "%[3]s"
  name                = "OFakeNA%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  zones               = local.zones

  resource_anchor_id = azurerm_oracle_resource_anchor.test.id
  subnet_id          = azurerm_subnet.virtual_network_subnet.id

  oci_backup_cidr_block = "10.0.2.0/24"
  tags = {
    test = "testNA1"
  }
}`, a.template(data), data.RandomString, data.Locations.Primary, data.Subscriptions.Primary)
}

func (a NetworkAnchorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_network_anchor" "import" {
  name                                   = azurerm_oracle_network_anchor.test.name
  location                               = azurerm_oracle_network_anchor.test.location
  resource_group_name                    = azurerm_oracle_network_anchor.test.resource_group_name
  zones                                  = azurerm_oracle_network_anchor.test.zones
  resource_anchor_id                     = azurerm_oracle_network_anchor.test.resource_anchor_id
  subnet_id                              = azurerm_oracle_network_anchor.test.subnet_id
  oracle_dns_forwarding_endpoint_enabled = azurerm_oracle_network_anchor.test.oracle_dns_forwarding_endpoint_enabled
  oracle_dns_listening_endpoint_enabled  = azurerm_oracle_network_anchor.test.oracle_dns_listening_endpoint_enabled
  oracle_to_azure_dns_zone_sync_enabled  = azurerm_oracle_network_anchor.test.oracle_to_azure_dns_zone_sync_enabled
}
`, a.basic(data))
}

func (a NetworkAnchorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  zones = ["2"]
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "virtual_network" {
  name                = "OFakeacctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "virtual_network_subnet" {
  name                 = "OFakeacctest%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.virtual_network.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Oracle.Database/networkAttachments"
    }
  }
}

resource "azurerm_oracle_resource_anchor" "test" {
  name                = "tfRA%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}


`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
