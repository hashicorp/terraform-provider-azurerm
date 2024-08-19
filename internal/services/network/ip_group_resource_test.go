// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IPGroupResource struct{}

func TestAccIpGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	r := IPGroupResource{}
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

func TestAccIpGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	r := IPGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_ip_group"),
		},
	})
}

func TestAccIpGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	r := IPGroupResource{}
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

func TestAccIpGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	r := IPGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIpGroup_updateWithAttachedPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test1")
	r := IPGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAzurePolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAzurePolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cidrs.#").HasValue("2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (t IPGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := ipgroups.ParseIPGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.Client.IPGroups.Get(ctx, *id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (IPGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r IPGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group" "import" {
  name                = azurerm_ip_group.test.name
  location            = azurerm_ip_group.test.location
  resource_group_name = azurerm_ip_group.test.resource_group_name
}
`, r.basic(data))
}

func (IPGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["192.168.0.1", "172.16.240.0/20", "10.48.0.0/12"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

resource "azurerm_ip_group" "test2" {
  name                = "acceptanceTestIpGroup2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["192.168.0.1", "172.16.240.0/20", "10.48.0.0/12"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

resource "azurerm_ip_group" "test3" {
  name                = "acceptanceTestIpGroup3"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["192.168.0.1", "172.16.240.0/20", "10.48.0.0/12"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IPGroupResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["172.16.240.0/20"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IPGroupResource) withAzurePolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test1" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["172.16.240.0/20"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

resource "azurerm_ip_group" "test2" {
  name                = "acceptanceTestIpGroup2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["172.17.240.0/20"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

resource "azurerm_firewall_policy" "test" {
  name                = "fwpol-test-policy"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "fwpol-test"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 100

  network_rule_collection {
    name     = "network-rule-collection1"
    priority = 100
    action   = "Allow"
    rule {
      name                  = "network-rule-collection1-rule1"
      protocols             = ["TCP"]
      source_ip_groups      = [azurerm_ip_group.test1.id]
      destination_ip_groups = [azurerm_ip_group.test2.id]
      destination_ports     = ["443"]
    }
  }

  network_rule_collection {
    name     = "network-rule-collection2"
    priority = 200
    action   = "Allow"
    rule {
      name                  = "network-rule-collection1-rule1"
      protocols             = ["TCP"]
      source_ip_groups      = [azurerm_ip_group.test2.id]
      destination_ip_groups = [azurerm_ip_group.test1.id]
      destination_ports     = ["443"]
    }
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "pip-fw"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "testfirewall"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  firewall_policy_id = azurerm_firewall_policy.test.id

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IPGroupResource) withAzurePolicyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test1" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["172.16.240.0/20", "172.18.240.0/20"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

resource "azurerm_ip_group" "test2" {
  name                = "acceptanceTestIpGroup2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["172.17.240.0/20", "172.19.240.0/20"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

resource "azurerm_firewall_policy" "test" {
  name                = "fwpol-test-policy"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "fwpol-test"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 100

  network_rule_collection {
    name     = "network-rule-collection1"
    priority = 100
    action   = "Allow"
    rule {
      name                  = "network-rule-collection1-rule1"
      protocols             = ["TCP"]
      source_ip_groups      = [azurerm_ip_group.test1.id]
      destination_ip_groups = [azurerm_ip_group.test2.id]
      destination_ports     = ["443"]
    }
  }

  network_rule_collection {
    name     = "network-rule-collection2"
    priority = 200
    action   = "Allow"
    rule {
      name                  = "network-rule-collection1-rule1"
      protocols             = ["TCP"]
      source_ip_groups      = [azurerm_ip_group.test2.id]
      destination_ip_groups = [azurerm_ip_group.test1.id]
      destination_ports     = ["443"]
    }
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "pip-fw"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "testfirewall"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  firewall_policy_id = azurerm_firewall_policy.test.id

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}




`, data.RandomInteger, data.Locations.Primary)
}
