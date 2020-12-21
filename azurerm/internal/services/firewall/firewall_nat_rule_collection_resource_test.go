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

func TestAccFirewallNatRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNatRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccFirewallNatRuleCollection_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_firewall_nat_rule_collection"),
			},
		},
	})
}

func TestAccFirewallNatRuleCollection_updatedName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewallNatRuleCollection_updatedName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNatRuleCollection_multipleRuleCollections(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	secondRule := "azurerm_firewall_nat_rule_collection.test_add"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewallNatRuleCollection_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
					testCheckFirewallNatRuleCollectionExists(secondRule),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
					testCheckFirewallNatRuleCollectionDoesNotExist("azurerm_firewall.test", "acctestnrc_add"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNatRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")
	secondResourceName := "azurerm_firewall_nat_rule_collection.test_add"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
					testCheckFirewallNatRuleCollectionExists(secondResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewallNatRuleCollection_multipleUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
					testCheckFirewallNatRuleCollectionExists(secondResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNatRuleCollection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
					testCheckFirewallNatRuleCollectionDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccFirewallNatRuleCollection_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewallNatRuleCollection_multipleRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNatRuleCollection_updateFirewallTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccFirewallNatRuleCollection_updateFirewallTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNatRuleCollection_ipGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNatRuleCollection_ipGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckFirewallNatRuleCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccFirewallNatRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_nat_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallNatRuleCollection_noSource(data),
				ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
			},
		},
	})
}

func testCheckFirewallNatRuleCollectionExists(resourceName string) resource.TestCheckFunc {
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
		for _, collection := range *read.AzureFirewallPropertiesFormat.NatRuleCollections {
			if *collection.Name == name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Expected NAT Rule Collection %q (Firewall %q / Resource Group %q) to exist but it didn't", name, firewallName, resourceGroup)
		}

		return nil
	}
}

func testCheckFirewallNatRuleCollectionDoesNotExist(resourceName string, collectionName string) resource.TestCheckFunc {
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

		for _, collection := range *read.AzureFirewallPropertiesFormat.NatRuleCollections {
			if *collection.Name == collectionName {
				return fmt.Errorf("NAT Rule Collection %q exists in Firewall %q: %+v", collectionName, firewallName, collection)
			}
		}

		return nil
	}
}

func testCheckFirewallNatRuleCollectionDisappears(resourceName string) resource.TestCheckFunc {
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

		rules := make([]network.AzureFirewallNatRuleCollection, 0)
		for _, collection := range *read.AzureFirewallPropertiesFormat.NatRuleCollections {
			if *collection.Name != name {
				rules = append(rules, collection)
			}
		}

		read.AzureFirewallPropertiesFormat.NatRuleCollections = &rules

		future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, read)
		if err != nil {
			return fmt.Errorf("Error removing NAT Rule Collection from Firewall: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for the removal of NAT Rule Collection from Firewall: %+v", err)
		}

		_, err = client.Get(ctx, resourceGroup, firewallName)
		return err
	}
}

func testAccFirewallNatRuleCollection_basic(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template, data.RandomInteger)
}

func testAccFirewallNatRuleCollection_requiresImport(data acceptance.TestData) string {
	template := testAccFirewallNatRuleCollection_basic(data)
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
`, template)
}

func testAccFirewallNatRuleCollection_updatedName(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template, data.RandomInteger)
}

func testAccFirewallNatRuleCollection_multiple(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccFirewallNatRuleCollection_multipleUpdate(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccFirewallNatRuleCollection_multipleRules(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template, data.RandomInteger)
}

func testAccFirewallNatRuleCollection_updateFirewallTags(data acceptance.TestData) string {
	template := testAccFirewall_withTags(data)
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
`, template, data.RandomInteger)
}

func testAccFirewallNatRuleCollection_ipGroup(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template, data.RandomInteger)
}

func testAccFirewallNatRuleCollection_noSource(data acceptance.TestData) string {
	template := testAccFirewall_basic(data)
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
`, template, data.RandomInteger)
}
