package firewall_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FirewallNetworkRuleCollectionResource struct {
}

func TestAccFirewallNetworkRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.3765122797.name").HasValue("rule1"),
			),
		},
		{
			Config: r.updatedName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1700340761.name").HasValue("rule2"),
			),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_multipleRuleCollections(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}
	secondRule := "azurerm_firewall_network_rule_collection.test_add"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.multiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondRule).ExistsInAzure(r),
				resource.TestCheckResourceAttr(secondRule, "name", "acctestnrc_add"),
				resource.TestCheckResourceAttr(secondRule, "priority", "200"),
				resource.TestCheckResourceAttr(secondRule, "action", "Deny"),
				resource.TestCheckResourceAttr(secondRule, "rule.#", "1"),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				testCheckFirewallNetworkRuleCollectionDoesNotExist("azurerm_firewall.test", "acctestnrc_add"),
			),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}
	secondResourceName := "azurerm_firewall_network_rule_collection.test_add"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
				resource.TestCheckResourceAttr(secondResourceName, "priority", "200"),
				resource.TestCheckResourceAttr(secondResourceName, "action", "Deny"),
				resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
			),
		},
		{
			Config: r.multipleUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("300"),
				check.That(data.ResourceName).Key("action").HasValue("Deny"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(secondResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
				resource.TestCheckResourceAttr(secondResourceName, "priority", "400"),
				resource.TestCheckResourceAttr(secondResourceName, "action", "Allow"),
				resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
			),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				testCheckFirewallNetworkRuleCollectionDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccFirewallNetworkRuleCollection_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.multipleRules(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("2"),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestnrc"),
				check.That(data.ResourceName).Key("priority").HasValue("100"),
				check.That(data.ResourceName).Key("action").HasValue("Allow"),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
			),
		},
		{
			Config: r.updateFirewallTags(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.serviceTag(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ipGroup(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fqdns(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallNetworkRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.noSource(data),
			ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_noDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	r := FirewallNetworkRuleCollectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.noDestination(data),
			ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q, %q and %q must be specified", "destination_addresses", "destination_ip_groups", "destination_fqdns")),
		},
	})
}

func (FirewallNetworkRuleCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id, err = azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	firewallName := id.Path["azureFirewalls"]
	name := id.Path["networkRuleCollections"]

	resp, err := clients.Firewall.AzureFirewallsClient.Get(ctx, id.ResourceGroup, firewallName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Firewall  Network Rule Collection %q (Firewall %q / Resource Group %q): %v", name, firewallName, id.ResourceGroup, err)
	}

	if resp.AzureFirewallPropertiesFormat == nil || resp.AzureFirewallPropertiesFormat.NetworkRuleCollections == nil {
		return nil, fmt.Errorf("retrieving Firewall  Network Rule Collection %q (Firewall %q / Resource Group %q): properties or collections was nil", name, firewallName, id.ResourceGroup)
	}

	for _, rule := range *resp.AzureFirewallPropertiesFormat.NetworkRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name == name {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func testCheckFirewallNetworkRuleCollectionDoesNotExist(resourceName string, collectionName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.AzureFirewallsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		firewallName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		read, err := client.Get(ctx, resourceGroup, firewallName)
		if err != nil {
			return err
		}

		for _, collection := range *read.AzureFirewallPropertiesFormat.NetworkRuleCollections {
			if *collection.Name == collectionName {
				return fmt.Errorf("Network Rule Collection %q exists in Firewall %q: %+v", collectionName, firewallName, collection)
			}
		}

		return nil
	}
}

func testCheckFirewallNetworkRuleCollectionDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.AzureFirewallsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		firewallName := rs.Primary.Attributes["azure_firewall_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		read, err := client.Get(ctx, resourceGroup, firewallName)
		if err != nil {
			return err
		}

		rules := make([]network.AzureFirewallNetworkRuleCollection, 0)
		for _, collection := range *read.AzureFirewallPropertiesFormat.NetworkRuleCollections {
			if *collection.Name != name {
				rules = append(rules, collection)
			}
		}

		read.AzureFirewallPropertiesFormat.NetworkRuleCollections = &rules

		future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, read)
		if err != nil {
			return fmt.Errorf("Error removing Network Rule Collection from Firewall: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for the removal of Network Rule Collection from Firewall: %+v", err)
		}

		_, err = client.Get(ctx, resourceGroup, firewallName)
		return err
	}
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

  rule {
    name = "acctestrule_add"

    source_addresses = [
      "192.168.0.1",
    ]

    destination_ports = [
      "8888",
    ]

    destination_addresses = [
      "1.1.1.1",
    ]

    protocols = [
      "TCP",
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

resource "azurerm_ip_group" "test_source" {
  name                = "acctestIpGroupForFirewallNetworkRulesSource"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}

resource "azurerm_ip_group" "test_destination" {
  name                = "acctestIpGroupForFirewallNetworkRulesDestination"
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
      azurerm_ip_group.test_source.id,
    ]

    destination_ports = [
      "53",
    ]

    destination_ip_groups = [
      azurerm_ip_group.test_destination.id,
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
      "time.windows.com"
    ]

    destination_ports = [
      "8080",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, FirewallResource{}.enableDNS(data))
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
