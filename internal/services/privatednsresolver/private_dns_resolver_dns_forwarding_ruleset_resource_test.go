// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateDNSResolverDnsForwardingRulesetResource struct{}

func TestAccPrivateDNSResolverDnsForwardingRuleset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_dns_forwarding_ruleset", "test")
	r := PrivateDNSResolverDnsForwardingRulesetResource{}
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

func TestAccPrivateDNSResolverDnsForwardingRuleset_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_dns_forwarding_ruleset", "test")
	r := PrivateDNSResolverDnsForwardingRulesetResource{}
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

func TestAccPrivateDNSResolverDnsForwardingRuleset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_dns_forwarding_ruleset", "test")
	r := PrivateDNSResolverDnsForwardingRulesetResource{}
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

func TestAccPrivateDNSResolverDnsForwardingRuleset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_dns_forwarding_ruleset", "test")
	r := PrivateDNSResolverDnsForwardingRulesetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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

func (r PrivateDNSResolverDnsForwardingRulesetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PrivateDnsResolver.DnsForwardingRulesetsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-rg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "outbounddns"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.64/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_subnet" "test2" {
  name                 = "outbounddns2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_private_dns_resolver" "test" {
  name                = "acctest-dr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
}

resource "azurerm_private_dns_resolver_outbound_endpoint" "test" {
  name                    = "acctest-droe-%[2]d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  subnet_id               = azurerm_subnet.test.id
}

resource "azurerm_private_dns_resolver_outbound_endpoint" "test2" {
  name                    = "acctest-droe2-%[2]d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  subnet_id               = azurerm_subnet.test2.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                                       = "acctest-drdfr-%d"
  resource_group_name                        = azurerm_resource_group.test.name
  location                                   = azurerm_resource_group.test.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.test.id]
}
`, template, data.RandomInteger)
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "import" {
  name                                       = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.name
  resource_group_name                        = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.resource_group_name
  location                                   = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.location
  private_dns_resolver_outbound_endpoint_ids = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.private_dns_resolver_outbound_endpoint_ids
}
`, config)
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                                       = "acctest-drdfr-%d"
  resource_group_name                        = azurerm_resource_group.test.name
  location                                   = azurerm_resource_group.test.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.test.id]
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger)
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                = "acctest-drdfr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  private_dns_resolver_outbound_endpoint_ids = [
    azurerm_private_dns_resolver_outbound_endpoint.test.id,
    azurerm_private_dns_resolver_outbound_endpoint.test2.id,
  ]
  tags = {
    key = "updated value"
  }
}
`, template, data.RandomInteger)
}
