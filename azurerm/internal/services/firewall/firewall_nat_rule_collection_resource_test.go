package firewall_test

import (
	`context`
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure`
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils`
)

type FirewallNatRuleCollectionResource struct {
}

func TestAccFirewallNatRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatedName(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data2.CheckWithClient(r.doesNotExist),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}
	secondResourceName := "azurerm_firewall_nat_rule_collection.test_add"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleUpdate(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleRules(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_updateFirewallTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateFirewallTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_ipGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ipGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNatRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	r := FirewallNatRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.noSource(data),
			ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
		},
	})
}

func (FirewallNatRuleCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id, err = azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	firewallName := id.Path["azureFirewalls"]
	name := id.Path["natRuleCollections"]

	resp, err := clients.Firewall.AzureFirewallsClient.Get(ctx, id.ResourceGroup, firewallName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): %v", name, firewallName, id.ResourceGroup, err)
	}

	if resp.AzureFirewallPropertiesFormat == nil || resp.AzureFirewallPropertiesFormat.NatRuleCollections == nil {
		return nil, fmt.Errorf("retrieving Firewall  Nat Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", name, firewallName, id.ResourceGroup)
	}

	for _, rule := range *resp.AzureFirewallPropertiesFormat.NatRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name == name {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func (t FirewallNatRuleCollectionResource) doesNotExist (ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
	var id, err = azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return err
	}

	firewallName := id.Path["azureFirewalls"]
	name := id.Path["natRuleCollections"]

	exists, err := t.Exists(ctx, clients, state)
	if err != nil  {
		return err
	}

	if *exists {
		return fmt.Errorf("Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): still exists", name, firewallName, id.ResourceGroup)
	}

	return nil
}

func (t FirewallNatRuleCollectionResource) disappears (ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
	client := clients.Firewall.AzureFirewallsClient
	var id, err = azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return err
	}

	firewallName := id.Path["azureFirewalls"]
	name := id.Path["natRuleCollections"]

	resp, err := client.Get(ctx, id.ResourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): %v", name, firewallName, id.ResourceGroup, err)
	}

	if resp.AzureFirewallPropertiesFormat == nil || resp.AzureFirewallPropertiesFormat.NatRuleCollections == nil {
		return fmt.Errorf("retrieving Firewall  Nat Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", name, firewallName, id.ResourceGroup)
	}

	rules := make([]network.AzureFirewallNatRuleCollection, 0)
	for _, collection := range *resp.AzureFirewallPropertiesFormat.NatRuleCollections {
		if *collection.Name != name {
			rules = append(rules, collection)
		}
	}

	resp.AzureFirewallPropertiesFormat.NatRuleCollections = &rules

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, firewallName, resp)
	if err != nil {
		return fmt.Errorf("removing Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): %v", name, firewallName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the removal of Firewall Nat Rule Collection %q (Firewall %q / Resource Group %q): %v", name, firewallName, id.ResourceGroup, err)
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

  rule {
    name = "acctestrule_add"

    source_addresses = [
      "192.168.0.1",
    ]

    destination_ports = [
      "8888",
    ]

    destination_addresses = [
      azurerm_public_ip.test.ip_address,
    ]

    protocols = [
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

resource "azurerm_ip_group" "test" {
  name                = "acctestIpGroupForFirewallNatRules"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
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
      azurerm_ip_group.test.id,
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
