// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkSecurityGroupResource struct{}

func TestAccNetworkSecurityGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
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

func TestAccNetworkSecurityGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_security_group"),
		},
	})
}

func TestAccNetworkSecurityGroup_singleRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkSecurityGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				// The configuration for this step contains one security_rule
				// block, which should now be reflected in the state.
				check.That(data.ResourceName).Key("security_rule.#").HasValue("1"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				// The configuration for this step contains no security_rule
				// blocks at all, which means "ignore any existing security groups"
				// and thus the one from the previous step is preserved.
				check.That(data.ResourceName).Key("security_rule.#").HasValue("1"),
			),
		},
		{
			Config: r.rulesExplicitZero(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),

				// The configuration for this step assigns security_rule = []
				// to state explicitly that no rules are desired, so the
				// rule from the first step should now be removed.
				check.That(data.ResourceName).Key("security_rule.#").HasValue("0"),
			),
		},
	})
}

func TestAccNetworkSecurityGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccNetworkSecurityGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkSecurityGroup_addingExtraRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("1"),
			),
		},

		{
			Config: r.anotherRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("2"),
			),
		},
	})
}

func TestAccNetworkSecurityGroup_augmented(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.augmented(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkSecurityGroup_applicationSecurityGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.applicationSecurityGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkSecurityGroup_deleteRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	r := NetworkSecurityGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.deleteRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkSecurityGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networksecuritygroups.ParseNetworkSecurityGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.Client.NetworkSecurityGroups.Get(ctx, *id, networksecuritygroups.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkSecurityGroupResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networksecuritygroups.ParseNetworkSecurityGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := client.Network.Client.NetworkSecurityGroups.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (NetworkSecurityGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkSecurityGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_group" "import" {
  name                = azurerm_network_security_group.test.name
  location            = azurerm_network_security_group.test.location
  resource_group_name = azurerm_network_security_group.test.resource_group_name
}
`, r.basic(data))
}

func (NetworkSecurityGroupResource) rulesExplicitZero(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule = []
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityGroupResource) singleRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityGroupResource) anotherRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "testDeny"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Deny"
    protocol                   = "Udp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityGroupResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityGroupResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityGroupResource) augmented(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                         = "test123"
    priority                     = 100
    direction                    = "Inbound"
    access                       = "Allow"
    protocol                     = "Tcp"
    source_port_ranges           = ["10000-40000"]
    destination_port_ranges      = ["80", "443", "8080", "8190"]
    source_address_prefixes      = ["10.0.0.0/8", "192.168.0.0/16"]
    destination_address_prefixes = ["172.16.0.0/20", "8.8.8.8"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityGroupResource) applicationSecurityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "first" {
  name                = "acctest-first%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "second" {
  name                = "acctest-second%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                                       = "test123"
    priority                                   = 100
    direction                                  = "Inbound"
    access                                     = "Allow"
    protocol                                   = "Tcp"
    source_application_security_group_ids      = [azurerm_application_security_group.first.id]
    destination_application_security_group_ids = [azurerm_application_security_group.second.id]
    source_port_ranges                         = ["10000-40000"]
    destination_port_ranges                    = ["80", "443", "8080", "8190"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (NetworkSecurityGroupResource) deleteRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  security_rule       = []
}
`, data.RandomInteger, data.Locations.Primary)
}
