// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FirewallNatRuleCollectionResource struct{}

func TestAccFirewallNatRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

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

func TestAccFirewallNatRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_firewall_nat_rule_collection"),
		},
	})
}

func TestAccFirewallNatRuleCollection_updatedName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatedName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_multipleRuleCollections(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test_add")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
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

func TestAccFirewallNatRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}
	secondResourceName := "azurerm_firewall_nat_rule_collection.test_add"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.disappears),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccFirewallNatRuleCollection_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleRules(data),
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

func TestAccFirewallNatRuleCollection_updateFirewallTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateFirewallTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_ipGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.noSource(data),
			ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
		},
	})
}

func (FirewallNatRuleCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallNatRuleCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	resp, err := clients.Network.AzureFirewalls.Get(ctx, firewallId)
	if err != nil {
		return nil, fmt.Errorf("retrieving Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): %v", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.NatRuleCollections == nil {
		return nil, fmt.Errorf("retrieving Firewall  Nat Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	for _, rule := range *resp.Model.Properties.NatRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name == id.NatRuleCollectionName {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func (t FirewallNatRuleCollectionResource) doesNotExist(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.FirewallNatRuleCollectionID(state.ID)
	if err != nil {
		return err
	}

	exists, err := t.Exists(ctx, clients, state)
	if err != nil {
		return err
	}

	if *exists {
		return fmt.Errorf("Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): still exists", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	return nil
}

func (t FirewallNatRuleCollectionResource) disappears(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	client := clients.Network.AzureFirewalls
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
	defer cancel()
	id, err := parse.FirewallNatRuleCollectionID(state.ID)
	if err != nil {
		return err
	}

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	resp, err := client.Get(ctx, firewallId)
	if err != nil {
		return fmt.Errorf("retrieving Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): %v", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.NatRuleCollections == nil {
		return fmt.Errorf("retrieving Firewall  Nat Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	rules := make([]azurefirewalls.AzureFirewallNatRuleCollection, 0)
	for _, collection := range *resp.Model.Properties.NatRuleCollections {
		if *collection.Name != id.NatRuleCollectionName {
			rules = append(rules, collection)
		}
	}

	resp.Model.Properties.NatRuleCollections = &rules

	if err := client.CreateOrUpdateThenPoll(ctx, firewallId, *resp.Model); err != nil {
		return fmt.Errorf("removing Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): %v", id.NatRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	return FirewallNatRuleCollectionResource{}.doesNotExist(ctx, clients, state)
}

func (FirewallNatRuleCollectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Dnat"

  rule {
    name        = "rule1"
    description = "test description"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "Any",
    ]

    translated_port    = 53
    translated_address = "8.8.8.8"
  }
}
`, FirewallResource{}.basic(data), data.RandomInteger)
}

func (r FirewallNatRuleCollectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "import" {
  name                = azurerm_firewall_nat_rule_collection.test.name
  azure_firewall_name = azurerm_firewall_nat_rule_collection.test.azure_firewall_name
  resource_group_name = azurerm_firewall_nat_rule_collection.test.resource_group_name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "Any",
    ]

    translated_port    = 53
    translated_address = "8.8.8.8"
  }
}
`, r.basic(data))
}

func (FirewallNatRuleCollectionResource) updatedName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "rule2"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "TCP",
    ]

    translated_port    = 53
    translated_address = "8.8.8.8"
  }
}
`, FirewallResource{}.basic(data), data.RandomInteger)
}

func (FirewallNatRuleCollectionResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "acctestrule"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "TCP",
    ]

    translated_port    = 53
    translated_address = "8.8.8.8"
  }
}

resource "azurerm_firewall_nat_rule_collection" "test_add" {
  name                = "acctestnrc_add-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 200
  action              = "Dnat"

  rule {
    name = "acctestruleadd"

    source_addresses = [
      "10.0.0.0/8",
    ]

    destination_ports = [
      "8080",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "TCP",
    ]

    translated_port    = 8080
    translated_address = "8.8.4.4"
  }
}
`, FirewallResource{}.basic(data), data.RandomInteger, data.RandomInteger)
}

func (FirewallNatRuleCollectionResource) multipleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 300
  action              = "Dnat"

  rule {
    name = "acctestrule"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "TCP",
    ]

    translated_port    = 53
    translated_address = "10.0.0.1"
  }
}

resource "azurerm_firewall_nat_rule_collection" "test_add" {
  name                = "acctestnrc_add-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 400
  action              = "Dnat"

  rule {
    name = "acctestruleadd"

    source_addresses = [
      "10.0.0.0/8",
    ]

    destination_ports = [
      "8080",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "TCP",
    ]

    translated_port    = 8080
    translated_address = "10.0.0.1"
  }
}
`, FirewallResource{}.basic(data), data.RandomInteger, data.RandomInteger)
}

func (FirewallNatRuleCollectionResource) multipleRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "acctestrule"

    source_addresses = [
      "10.0.0.0/16",
      "192.168.0.1",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "TCP",
      "UDP",
    ]

    translated_port    = 53
    translated_address = "10.0.0.1"
  }

  rule {
    name = "acctestrule_add"

    source_addresses = [
      "192.168.0.1",
      "10.0.0.0/16",
    ]

    destination_ports = [
      "8888",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "UDP",
      "TCP",
    ]

    translated_port    = 8888
    translated_address = "192.168.0.1"
  }
}
`, FirewallResource{}.basic(data), data.RandomInteger)
}

func (FirewallNatRuleCollectionResource) updateFirewallTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "TCP",
    ]

    translated_port    = 53
    translated_address = "10.0.0.1"
  }
}
`, FirewallResource{}.withTags(data), data.RandomInteger)
}

func (FirewallNatRuleCollectionResource) ipGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group" "test1" {
  name                = "acctestIpGroupForFirewallNatRules1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
}

resource "azurerm_ip_group" "test2" {
  name                = "acctestIpGroupForFirewallNatRules2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["193.168.0.0/25", "193.168.0.192/26"]
}

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "rule1"

    source_ip_groups = [
      azurerm_ip_group.test1.id,
      azurerm_ip_group.test2.id,
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "Any",
    ]

    translated_port    = 53
    translated_address = "8.8.8.8"
  }
}
`, FirewallResource{}.basic(data), data.RandomInteger)
}

func (FirewallNatRuleCollectionResource) noSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc-%d"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "rule1"

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
      "Any",
    ]

    translated_port    = 53
    translated_address = "8.8.8.8"
  }
}
`, FirewallResource{}.basic(data), data.RandomInteger)
}
