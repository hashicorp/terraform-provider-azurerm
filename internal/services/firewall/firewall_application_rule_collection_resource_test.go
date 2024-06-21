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

type FirewallApplicationRuleCollectionResource struct{}

func TestAccFirewallApplicationRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.source_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.target_fqdns.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.protocol.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.port").HasValue("443"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.type").HasValue("Https"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallApplicationRuleCollection_fqdnTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fqdnTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallApplicationRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_firewall_application_rule_collection"),
		},
	})
}

func TestAccFirewallApplicationRuleCollection_updatedName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
			),
		},
		{
			Config: r.updatedName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule2"),
			),
		},
	})
}

func TestAccFirewallApplicationRuleCollection_multipleRuleCollections(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	secondRule := "azurerm_firewall_application_rule_collection.test_add"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondRule).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(secondRule, "name", "acctestarc_add"),
				acceptance.TestCheckResourceAttr(secondRule, "priority", "200"),
				acceptance.TestCheckResourceAttr(secondRule, "action", "Deny"),
				acceptance.TestCheckResourceAttr(secondRule, "rule.#", "1"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
	})
}

func TestAccFirewallApplicationRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}
	secondResourceName := "azurerm_firewall_application_rule_collection.test_add"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(secondResourceName).Key("name").HasValue("acctestarc_add"),
				check.That(secondResourceName).Key("priority").HasValue("200"),
				check.That(secondResourceName).Key("action").HasValue("Deny"),
				check.That(secondResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.multipleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("300"),
				check.That(data.ResourceName).Key("action").HasValue("Deny"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondResourceName).ExistsInAzure(r),
				check.That(secondResourceName).Key("name").HasValue("acctestarc_add"),
				check.That(secondResourceName).Key("priority").HasValue("400"),
				check.That(secondResourceName).Key("action").HasValue("Allow"),
				check.That(secondResourceName).Key("rule.#").HasValue("1"),
			),
		},
	})
}

func TestAccFirewallApplicationRuleCollection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				data.CheckWithClient(r.disappears),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccFirewallApplicationRuleCollection_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.multipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("2"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
	})
}

func TestAccFirewallApplicationRuleCollection_multipleProtocols(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleProtocols(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.protocol.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.port").HasValue("8000"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.type").HasValue("Http"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.port").HasValue("8001"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.type").HasValue("Https"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallApplicationRuleCollection_updateProtocols(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleProtocols(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.protocol.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.port").HasValue("8000"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.type").HasValue("Http"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.port").HasValue("8001"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.type").HasValue("Https"),
			),
		},
		{
			Config: r.multipleProtocolsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.protocol.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.port").HasValue("9000"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.type").HasValue("Https"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.port").HasValue("9001"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.type").HasValue("Http"),
			),
		},
		{
			Config: r.multipleProtocols(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.protocol.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.port").HasValue("8000"),
				check.That(data.ResourceName).Key("rule.0.protocol.0.type").HasValue("Http"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.port").HasValue("8001"),
				check.That(data.ResourceName).Key("rule.0.protocol.1.type").HasValue("Https"),
			),
		},
	})
}

func TestAccFirewallApplicationRuleCollection_updateFirewallTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
			),
		},
		{
			Config: r.updateFirewallTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestarc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
			),
		},
	})
}

func TestAccFirewallApplicationRuleCollection_ipGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallApplicationRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	r := FirewallApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.noSource(data),
			ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
		},
	})
}

func (FirewallApplicationRuleCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallApplicationRuleCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	resp, err := clients.Network.AzureFirewalls.Get(ctx, firewallId)
	if err != nil {
		return nil, fmt.Errorf("retrieving Firewall Application Rule Collection %q (Firewall %q / Resource Group %q): %v", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ApplicationRuleCollections == nil {
		return nil, fmt.Errorf("retrieving Firewall Application Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	for _, rule := range *resp.Model.Properties.ApplicationRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name == id.ApplicationRuleCollectionName {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func (t FirewallApplicationRuleCollectionResource) doesNotExist(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.FirewallApplicationRuleCollectionID(state.ID)
	if err != nil {
		return err
	}

	exists, err := t.Exists(ctx, clients, state)
	if err != nil {
		return err
	}

	if *exists {
		return fmt.Errorf("Firewall Application Rule Collection %q (Firewall %q / Resource Group %q): still exists", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	return nil
}

func (t FirewallApplicationRuleCollectionResource) disappears(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	client := clients.Network.AzureFirewalls
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(15*time.Minute))
	defer cancel()
	id, err := parse.FirewallApplicationRuleCollectionID(state.ID)
	if err != nil {
		return err
	}

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	resp, err := client.Get(ctx, firewallId)
	if err != nil {
		return fmt.Errorf("retrieving Firewall Application Rule Collection %q (Firewall %q / Resource Group %q): %v", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ApplicationRuleCollections == nil {
		return fmt.Errorf("retrieving Firewall  Application Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	rules := make([]azurefirewalls.AzureFirewallApplicationRuleCollection, 0)
	for _, collection := range *resp.Model.Properties.ApplicationRuleCollections {
		if *collection.Name != id.ApplicationRuleCollectionName {
			rules = append(rules, collection)
		}
	}

	resp.Model.Properties.ApplicationRuleCollections = &rules

	if err = client.CreateOrUpdateThenPoll(ctx, firewallId, *resp.Model); err != nil {
		return fmt.Errorf("removing Firewall Application Rule Collection %q (Firewall %q / Resource Group %q): %v", id.ApplicationRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	return FirewallApplicationRuleCollectionResource{}.doesNotExist(ctx, clients, state)
}

func (FirewallApplicationRuleCollectionResource) basic(data acceptance.TestData) string {
	template := FirewallResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
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

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}
`, template)
}

func (FirewallApplicationRuleCollectionResource) fqdnTags(data acceptance.TestData) string {
	template := FirewallResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name        = "rule1"
    description = "test description"
    fqdn_tags   = ["WindowsDiagnostics", "AzureBackup"]
    source_addresses = [
      "10.0.0.0/16",
    ]
  }
}
`, template)
}

func (FirewallApplicationRuleCollectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "import" {
  name                = azurerm_firewall_application_rule_collection.test.name
  azure_firewall_name = azurerm_firewall_application_rule_collection.test.azure_firewall_name
  resource_group_name = azurerm_firewall_application_rule_collection.test.resource_group_name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}
`, FirewallApplicationRuleCollectionResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) updatedName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule2"

    source_addresses = [
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}

resource "azurerm_firewall_application_rule_collection" "test_add" {
  name                = "acctestarc_add"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 200
  action              = "Deny"

  rule {
    name = "acctestruleadd"

    source_addresses = [
      "192.168.0.1",
    ]

    target_fqdns = [
      "*.microsoft.com",
    ]

    protocol {
      port = 80
      type = "Http"
    }
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) multipleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 300
  action              = "Deny"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}

resource "azurerm_firewall_application_rule_collection" "test_add" {
  name                = "acctestarc_add"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 400
  action              = "Allow"

  rule {
    name = "acctestruleadd"

    source_addresses = [
      "192.168.0.1",
    ]

    target_fqdns = [
      "*.microsoft.com",
    ]

    protocol {
      port = 80
      type = "Http"
    }
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) multipleRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "192.0.0.0/16",
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
      "*.microsoft.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }

  rule {
    name = "acctestruleadd"

    source_addresses = [
      "192.168.0.1",
      "10.0.0.1",
    ]

    target_fqdns = [
      "*.microsoft.com",
      "*.google.com",
    ]

    protocol {
      port = 80
      type = "Http"
    }
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) multipleProtocols(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 8000
      type = "Http"
    }

    protocol {
      port = 8001
      type = "Https"
    }
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) multipleProtocolsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 9000
      type = "Https"
    }

    protocol {
      port = 9001
      type = "Http"
    }
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) updateFirewallTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_addresses = [
      "10.0.0.0/16",
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}
`, FirewallResource{}.withTags(data))
}

func (FirewallApplicationRuleCollectionResource) ipGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group" "test1" {
  name                = "acctestIpGroupForFirewallAppRules1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
}

resource "azurerm_ip_group" "test2" {
  name                = "acctestIpGroupForFirewallAppRules2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["193.168.0.0/25", "193.168.0.192/26"]
}

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    source_ip_groups = [
      azurerm_ip_group.test1.id,
      azurerm_ip_group.test2.id,
    ]

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}
`, FirewallResource{}.basic(data))
}

func (FirewallApplicationRuleCollectionResource) noSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
  priority            = 100
  action              = "Allow"

  rule {
    name = "rule1"

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = 443
      type = "Https"
    }
  }
}
`, FirewallResource{}.basic(data))
}
