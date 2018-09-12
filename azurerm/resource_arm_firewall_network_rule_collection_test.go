package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func TestAccAzureRMFirewallNetworkRuleCollection_basic(t *testing.T) {
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_updatedName(t *testing.T) {
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_updatedName(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "updatedname"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_multipleRuleCollections(t *testing.T) {
	firstRule := "azurerm_firewall_network_rule_collection.test"
	secondRule := "azurerm_firewall_network_rule_collection.test_add"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Allow"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_multiple(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Allow"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
					testCheckAzureRMFirewallNetworkRuleCollectionExists(secondRule),
					resource.TestCheckResourceAttr(secondRule, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondRule, "priority", "200"),
					resource.TestCheckResourceAttr(secondRule, "action", "Deny"),
					resource.TestCheckResourceAttr(secondRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Allow"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
					testCheckAzureRMFirewallNetworkRuleCollectionDoesNotExist("azurerm_firewall.test", "acctestnrc_add"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_update(t *testing.T) {
	firstResourceName := "azurerm_firewall_network_rule_collection.test"
	secondResourceName := "azurerm_firewall_network_rule_collection.test_add"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_multiple(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(firstResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(firstResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallNetworkRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "200"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_multipleUpdate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "300"),
					resource.TestCheckResourceAttr(firstResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(firstResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallNetworkRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "400"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_disappears(t *testing.T) {
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testCheckAzureRMFirewallNetworkRuleCollectionDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_multipleRules(t *testing.T) {
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_multipleRules(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
				),
			},
			{
				Config: testAccAzureRMFirewallNetworkRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMFirewallNetworkRuleCollectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		firewallName := rs.Primary.Attributes["azure_firewall_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		read, err := client.Get(ctx, resourceGroup, name)
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

func testCheckAzureRMFirewallNetworkRuleCollectionDoesNotExist(resourceName string, collectionName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		read, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return err
		}

		for _, collection := range *read.AzureFirewallPropertiesFormat.NetworkRuleCollections {
			if *collection.Name == collectionName {
				return fmt.Errorf("Network Rule Collection %q exists: %+v", name, collection)
			}
		}

		return nil
	}
}

func testCheckAzureRMFirewallNetworkRuleCollectionDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		firewallName := rs.Primary.Attributes["azure_firewall_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
		ipConfigs, err := azure.FixFirewallIPConfiguration(read.AzureFirewallPropertiesFormat.IPConfigurations)
		if err != nil {
			return fmt.Errorf("Error fixing IP Configuration for Firewall: err")
		}

		read.AzureFirewallPropertiesFormat.IPConfigurations = ipConfigs

		future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, read)
		if err != nil {
			return fmt.Errorf("Error removing Network Rule Collection from Firewall: %+v", err)
		}

		err = future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for the removal of Network Rule Collection from Firewall: %+v", err)
		}

		_, err = client.Get(ctx, resourceGroup, name)
		return err
	}
}

func testAccAzureRMFirewallNetworkRuleCollection_basic(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallNetworkRuleCollection_updatedName(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  priority            = 100
  action              = "Allow"

  rule {
    name = "updatedname"

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

func testAccAzureRMFirewallNetworkRuleCollection_multiple(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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

func testAccAzureRMFirewallNetworkRuleCollection_multipleUpdate(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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

func testAccAzureRMFirewallNetworkRuleCollection_multipleRules(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
