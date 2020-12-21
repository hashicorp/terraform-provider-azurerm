package firewall_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccFirewallNetworkRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccFirewallNetworkRuleCollection_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_firewall_network_rule_collection"),
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_updatedName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.3765122797.name", "rule1"),
				),
			},
			{
				Config: testAccFirewallNetworkRuleCollection_updatedName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1700340761.name", "rule2"),
				),
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_multipleRuleCollections(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	secondRule := "azurerm_firewall_network_rule_collection.test_add"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccFirewallNetworkRuleCollection_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckFirewallNetworkRuleCollectionExists(secondRule),
					resource.TestCheckResourceAttr(secondRule, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondRule, "priority", "200"),
					resource.TestCheckResourceAttr(secondRule, "action", "Deny"),
					resource.TestCheckResourceAttr(secondRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckFirewallNetworkRuleCollectionDoesNotExist("azurerm_firewall.test", "acctestnrc_add"),
				),
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")
	secondResourceName := "azurerm_firewall_network_rule_collection.test_add"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckFirewallNetworkRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "200"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccFirewallNetworkRuleCollection_multipleUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "300"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckFirewallNetworkRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "400"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckFirewallNetworkRuleCollectionDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccFirewallNetworkRuleCollection_multipleRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "2"),
				),
			},
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_updateFirewallTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccFirewallNetworkRuleCollection_updateFirewallTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_serviceTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_serviceTag(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_ipGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_ipGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_fqdns(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNetworkRuleCollection_fqdns(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNetworkRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNetworkRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallNetworkRuleCollection_noSource(data),
				ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
			},
		},
	})
}

func TestAccFirewallNetworkRuleCollection_noDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_network_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallNetworkRuleCollection_noDestination(data),
				ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q, %q and %q must be specified", "destination_addresses", "destination_ip_groups", "destination_fqdns")),
			},
		},
	})
}

func testCheckFirewallNetworkRuleCollectionExists(resourceName string) resource.TestCheckFunc {
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

		found := false
		for _, collection := range *read.AzureFirewallPropertiesFormat.NetworkRuleCollections {
			if *collection.Name == name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Expected Network Rule Collection %q (Firewall %q / Resource Group %q) to exist but it didn't", name, firewallName, resourceGroup)
		}

		return nil
	}
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

func testAccFirewallNetworkRuleCollection_basic(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_requiresImport(data acceptance.TestData) string {
	template := testAccFirewallNetworkRuleCollection_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_updatedName(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_multiple(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_multipleUpdate(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_multipleRules(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_updateFirewallTags(data acceptance.TestData) string {
	template := testAccFirewall_withTags(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_serviceTag(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_ipGroup(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_fqdns(data acceptance.TestData) string {
	template := testAccFirewall_enableDNS(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_noSource(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}

func testAccFirewallNetworkRuleCollection_noDestination(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template)
}
