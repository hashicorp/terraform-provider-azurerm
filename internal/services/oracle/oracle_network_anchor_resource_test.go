// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkAnchorResource struct{}

// OracleNetworkAnchorResource matches the resource name inferred by the identity test generator.
type OracleNetworkAnchorResource = NetworkAnchorResource

func (a NetworkAnchorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networkanchors.ParseNetworkAnchorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.NetworkAnchors.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
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
			),
		},
		data.ImportStep(),
		{
			Config:   r.completeReordered(data),
			PlanOnly: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
		},
	})
}

func (r NetworkAnchorResource) completeReordered(data acceptance.TestData) string {
	config := strings.Replace(r.complete(data),
		`["abc.ocidelegated.ocinetworkanch.oraclevcn.com", "ghi.ocidelegated.ocinetworkanch.oraclevcn.com"]`,
		`["ghi.ocidelegated.ocinetworkanch.oraclevcn.com", "abc.ocidelegated.ocinetworkanch.oraclevcn.com"]`, 1)
	return strings.Replace(config, `["10.0.2.0/24", "10.0.3.0/24"]`, `["10.0.3.0/24", "10.0.2.0/24"]`, 1)
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (a NetworkAnchorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_oracle_network_anchor" "test" {
  location            = "%[3]s"
  name                = "acctestNA%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  zones               = local.zones

  resource_anchor_id = azurerm_oracle_resource_anchor.test.id
  subnet_id          = azurerm_subnet.virtual_network_subnet.id
}`, a.template(data), data.RandomString, data.Locations.Primary, data.Subscriptions.Primary)
}

func (a NetworkAnchorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_oracle_network_anchor" "test" {
  location            = "%[3]s"
  name                = "acctestNA%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  zones               = local.zones

  resource_anchor_id = azurerm_oracle_resource_anchor.test.id
  subnet_id          = azurerm_subnet.virtual_network_subnet.id

  backup_subnet_cidr                   = "10.0.0.0/24"
  dns_private_zone_replication_enabled = true

  dns_forwarding_rule {
    domain_names          = ["abc.ocidelegated.ocinetworkanch.oraclevcn.com", "ghi.ocidelegated.ocinetworkanch.oraclevcn.com"]
    forwarding_ip_address = "10.0.1.16"
  }
  dns_forwarding_rule {
    domain_names          = ["def.ocidelegated.ocinetworkanch.oraclevcn.com"]
    forwarding_ip_address = "10.0.1.24"
  }
  dns_listening_endpoint_allowed_cidrs = ["10.0.2.0/24", "10.0.3.0/24"]
  tags = {
    test = "testNA1"
  }
}`, a.template(data), data.RandomString, data.Locations.Primary)
}

func (a NetworkAnchorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_oracle_network_anchor" "test" {
  location            = "%[3]s"
  name                = "acctestNA%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  zones               = local.zones

  resource_anchor_id = azurerm_oracle_resource_anchor.test.id
  subnet_id          = azurerm_subnet.virtual_network_subnet.id

  backup_subnet_cidr = "10.0.2.0/24"
  tags = {
    test = "testNA1"
  }
}`, a.template(data), data.RandomString, data.Locations.Primary, data.Subscriptions.Primary)
}

func (a NetworkAnchorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_network_anchor" "import" {
  name                                 = azurerm_oracle_network_anchor.test.name
  location                             = azurerm_oracle_network_anchor.test.location
  resource_group_name                  = azurerm_oracle_network_anchor.test.resource_group_name
  zones                                = azurerm_oracle_network_anchor.test.zones
  resource_anchor_id                   = azurerm_oracle_network_anchor.test.resource_anchor_id
  subnet_id                            = azurerm_oracle_network_anchor.test.subnet_id
  dns_private_zone_replication_enabled = azurerm_oracle_network_anchor.test.dns_private_zone_replication_enabled
}
`, a.basic(data))
}

func (a NetworkAnchorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  zones = ["3"]
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "virtual_network" {
  name                = "acctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "virtual_network_subnet" {
  name                 = "acctest%[1]d_delegated"
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
  name                = "acctestRA%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomIntOfLength(15), data.Locations.Primary, data.RandomString)
}
