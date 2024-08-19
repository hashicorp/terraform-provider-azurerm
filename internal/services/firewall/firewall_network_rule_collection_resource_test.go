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

type FirewallNetworkRuleCollectionResource struct{}

func TestAccFirewallNetworkRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNetworkRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_firewall_network_rule_collection"),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_updatedName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatedName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNetworkRuleCollection_multipleRuleCollections(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}
	secondRule := "azurerm_firewall_network_rule_collection.test_add"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondRule).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(secondRule, "name", "acctestnrc_add"),
				acceptance.TestCheckResourceAttr(secondRule, "priority", "200"),
				acceptance.TestCheckResourceAttr(secondRule, "action", "Deny"),
				acceptance.TestCheckResourceAttr(secondRule, "rule.#", "1"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				data.CheckWithClient(r.checkFirewallNetworkRuleCollectionDoesNotExist("acctestnrc_add")),
			),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}
	secondResourceName := "azurerm_firewall_network_rule_collection.test_add"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
				acceptance.TestCheckResourceAttr(secondResourceName, "priority", "200"),
				acceptance.TestCheckResourceAttr(secondResourceName, "action", "Deny"),
				acceptance.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
			),
		},
		{
			Config: r.multipleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("300"),
				check.That(data.ResourceName).Key("action").HasValue("Deny"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
				acceptance.TestCheckResourceAttr(secondResourceName, "priority", "400"),
				acceptance.TestCheckResourceAttr(secondResourceName, "action", "Allow"),
				acceptance.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
			),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccFirewallNetworkRuleCollection_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.multipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("2"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_updateFirewallTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.updateFirewallTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_serviceTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceTag(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNetworkRuleCollection_ipGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNetworkRuleCollection_fqdns(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fqdns(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNetworkRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.noSource(data),
			ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_noDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.noDestination(data),
			ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q, %q and %q must be specified", "destination_addresses", "destination_ip_groups", "destination_fqdns")),
		},
	})
}

func (FirewallNetworkRuleCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallNetworkRuleCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	resp, err := clients.Network.AzureFirewalls.Get(ctx, firewallId)
	if err != nil {
		return nil, fmt.Errorf("retrieving Firewall  Network Rule Collection %q (Firewall %q / Resource Group %q): %v", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.NetworkRuleCollections == nil {
		return nil, fmt.Errorf("retrieving Firewall  Network Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	for _, rule := range *resp.Model.Properties.NetworkRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name == id.NetworkRuleCollectionName {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func (r FirewallNetworkRuleCollectionResource) checkFirewallNetworkRuleCollectionDoesNotExist(collectionName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
		defer cancel()
		// Ensure we have enough information in state to look up in API
		id, err := parse.FirewallNetworkRuleCollectionID(state.ID)
		if err != nil {
			return err
		}

		firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

		read, err := clients.Network.AzureFirewalls.Get(ctx, firewallId)
		if err != nil {
			return err
		}

		if read.Model == nil || read.Model.Properties == nil || read.Model.Properties.NetworkRuleCollections == nil {
			return fmt.Errorf("one of model/properties/networkRuleCollections was nil for %s", firewallId)
		}

		for _, collection := range *read.Model.Properties.NetworkRuleCollections {
			if *collection.Name == collectionName {
				return fmt.Errorf("network Rule Collection %q exists in Firewall %q: %+v", collectionName, firewallId.AzureFirewallName, collection)
			}
		}

		return nil
	}
}

func (FirewallNetworkRuleCollectionResource) Destroy(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallNetworkRuleCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	read, err := clients.Network.AzureFirewalls.Get(ctx, firewallId)
	if err != nil {
		return utils.Bool(false), err
	}

	if read.Model == nil || read.Model.Properties == nil || read.Model.Properties.NetworkRuleCollections == nil {
		return utils.Bool(false), fmt.Errorf("one of model/properties/networkRuleCollections was nil for %s", firewallId)
	}

	rules := make([]azurefirewalls.AzureFirewallNetworkRuleCollection, 0)
	for _, collection := range *read.Model.Properties.NetworkRuleCollections {
		if *collection.Name != id.NetworkRuleCollectionName {
			rules = append(rules, collection)
		}
	}

	read.Model.Properties.NetworkRuleCollections = &rules

	if err = clients.Network.AzureFirewalls.CreateOrUpdateThenPoll(ctx, firewallId, *read.Model); err != nil {
		return utils.Bool(false), fmt.Errorf("removing Network Rule Collection from Firewall: %+v", err)
	}

	return utils.Bool(true), nil
}

func (FirewallNetworkRuleCollectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

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
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (r FirewallNetworkRuleCollectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "import" {
  name                = azurerm_firewall_network_rule_collection.test.name
  azure_firewall_name = azurerm_firewall_network_rule_collection.test.azure_firewall_name
  resource_group_name = azurerm_firewall_network_rule_collection.test.resource_group_name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, r.basic(data))
}

func (FirewallNetworkRuleCollectionResource) updatedName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule2"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallNetworkRuleCollectionResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "acctestrule"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}

resource "azurerm_firewall_network_rule_collection" "test_add" {
  name                = "acctestnrc_add"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 200
  action              = "Deny"

  rule {
    name = "acctestruleadd"

    source_addresses = [
      "10.0.0.0/8",
    ]

    destination_ports = [
      "8080",
    ]

    destination_addresses = [
      "8.8.4.4",
    ]

    protocols = [
      "TCP",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallNetworkRuleCollectionResource) multipleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 300
  action              = "Deny"

  rule {
    name = "acctestrule"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}

resource "azurerm_firewall_network_rule_collection" "test_add" {
  name                = "acctestnrc_add"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 400
  action              = "Allow"

  rule {
    name = "acctestruleadd"

    source_addresses = [
      "10.0.0.0/8",
    ]

    destination_ports = [
      "8080",
    ]

    destination_addresses = [
      "8.8.4.4",
    ]

    protocols = [
      "TCP",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallNetworkRuleCollectionResource) multipleRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "acctestrule"

    source_addresses = [
      "10.0.0.0/16",
      "192.0.0.0/16",
    ]

    destination_ports = [
      "53",
      "64",
    ]

    destination_addresses = [
      "8.8.8.8",
      "1.1.1.1",
    ]

    protocols = [
      "UDP",
      "TCP",
    ]
  }

  rule {
    name = "acctestrule_add"

    source_addresses = [
      "192.168.0.1",
      "10.0.0.0/16",
    ]

    destination_ports = [
      "8888",
      "9999",
    ]

    destination_addresses = [
      "1.1.1.1",
      "8.8.8.8",
    ]

    protocols = [
      "TCP",
      "UDP",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallNetworkRuleCollectionResource) updateFirewallTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.withTags(data))
}

func (r FirewallNetworkRuleCollectionResource) serviceTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "ApiManagement",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallNetworkRuleCollectionResource) ipGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group" "test_source1" {
  name                = "acctestIpGroupForFirewallNetworkRulesSource1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}
resource "azurerm_ip_group" "test_source2" {
  name                = "acctestIpGroupForFirewallNetworkRulesSource2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["4.3.2.1/32", "65.43.21.0/24"]
}

resource "azurerm_ip_group" "test_destination1" {
  name                = "acctestIpGroupForFirewallNetworkRulesDestination1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
}

resource "azurerm_ip_group" "test_destination2" {
  name                = "acctestIpGroupForFirewallNetworkRulesDestination2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
}

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_ip_groups = [
      azurerm_ip_group.test_source1.id,
      azurerm_ip_group.test_source2.id,
    ]

    destination_ports = [
      "53",
    ]

    destination_ip_groups = [
      azurerm_ip_group.test_destination1.id,
      azurerm_ip_group.test_destination2.id,
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallNetworkRuleCollectionResource) fqdns(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_fqdns = [
      "time.windows.com",
      "time.linux.com"
    ]

    destination_ports = [
      "8080",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.enableDNS(data, true, "1.1.1.1", "8.8.8.8"))
}

func (r FirewallNetworkRuleCollectionResource) noSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallNetworkRuleCollectionResource) noDestination(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.basic(data))
}
