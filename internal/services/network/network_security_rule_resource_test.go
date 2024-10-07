// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/securityrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkSecurityRuleResource struct{}

func TestAccNetworkSecurityRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test")
	r := NetworkSecurityRuleResource{}
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

func TestAccNetworkSecurityRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test")
	r := NetworkSecurityRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_security_rule"),
		},
	})
}

func TestAccNetworkSecurityRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test")
	r := NetworkSecurityRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccNetworkSecurityRule_addingRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test1")
	r := NetworkSecurityRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		{
			Config: r.updateExtraRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccNetworkSecurityRule_augmented(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test1")
	r := NetworkSecurityRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.augmented(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkSecurityRule_applicationSecurityGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_rule", "test1")
	r := NetworkSecurityRuleResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.applicationSecurityGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkSecurityRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securityrules.ParseSecurityRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.SecurityRules.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NetworkSecurityRuleResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securityrules.ParseSecurityRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	if err = client.Network.SecurityRules.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (NetworkSecurityRuleResource) basic(data acceptance.TestData) string {
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

resource "azurerm_network_security_rule" "test" {
  name                        = "test123"
  network_security_group_name = azurerm_network_security_group.test.name
  resource_group_name         = azurerm_resource_group.test.name
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkSecurityRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_rule" "import" {
  name                        = azurerm_network_security_rule.test.name
  network_security_group_name = azurerm_network_security_rule.test.network_security_group_name
  resource_group_name         = azurerm_network_security_rule.test.resource_group_name
  priority                    = azurerm_network_security_rule.test.priority
  direction                   = azurerm_network_security_rule.test.direction
  access                      = azurerm_network_security_rule.test.access
  protocol                    = azurerm_network_security_rule.test.protocol
  source_port_range           = azurerm_network_security_rule.test.source_port_range
  destination_port_range      = azurerm_network_security_rule.test.destination_port_range
  source_address_prefix       = azurerm_network_security_rule.test.source_address_prefix
  destination_address_prefix  = azurerm_network_security_rule.test.destination_address_prefix
}
`, r.basic(data))
}

func (NetworkSecurityRuleResource) updateBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test1.name
  network_security_group_name = azurerm_network_security_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityRuleResource) updateExtraRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test1.name
  network_security_group_name = azurerm_network_security_group.test1.name
}

resource "azurerm_network_security_rule" "test2" {
  name                        = "testing456"
  priority                    = 101
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "Icmp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.test1.name
  network_security_group_name = azurerm_network_security_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityRuleResource) augmented(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}

resource "azurerm_network_security_rule" "test1" {
  name                         = "test123"
  priority                     = 100
  direction                    = "Outbound"
  access                       = "Allow"
  protocol                     = "Tcp"
  source_port_ranges           = ["10000-40000"]
  destination_port_ranges      = ["80", "443", "8080", "8190"]
  source_address_prefixes      = ["10.0.0.0/8", "192.168.0.0/16"]
  destination_address_prefixes = ["172.16.0.0/20", "8.8.8.8"]
  resource_group_name          = azurerm_resource_group.test1.name
  network_security_group_name  = azurerm_network_security_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetworkSecurityRuleResource) applicationSecurityGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "source1" {
  name                = "acctest-source1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "source2" {
  name                = "acctest-source2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "destination1" {
  name                = "acctest-destination1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "destination2" {
  name                = "acctest-destination2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test1" {
  name                                       = "test123"
  resource_group_name                        = azurerm_resource_group.test.name
  network_security_group_name                = azurerm_network_security_group.test.name
  priority                                   = 100
  direction                                  = "Outbound"
  access                                     = "Allow"
  protocol                                   = "Tcp"
  source_application_security_group_ids      = [azurerm_application_security_group.source1.id, azurerm_application_security_group.source2.id]
  destination_application_security_group_ids = [azurerm_application_security_group.destination1.id, azurerm_application_security_group.destination2.id]
  source_port_ranges                         = ["10000-40000"]
  destination_port_ranges                    = ["80", "443", "8080", "8190"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
