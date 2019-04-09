package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMFirewallNatRuleCollection_basic(t *testing.T) {
	resourceName := "azurerm_firewall_nat_rule_collection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
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

func TestAccAzureRMFirewallNatRuleCollection_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_firewall_nat_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMFirewallNatRuleCollection_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_firewall_nat_rule_collection"),
			},
		},
	})
}

func TestAccAzureRMFirewallNatRuleCollection_updatedName(t *testing.T) {
	resourceName := "azurerm_firewall_nat_rule_collection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.3765122797.name", "rule1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNatRuleCollection_updatedName(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1700340761.name", "rule2"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNatRuleCollection_multipleRuleCollections(t *testing.T) {
	firstRule := "azurerm_firewall_nat_rule_collection.test"
	secondRule := "azurerm_firewall_nat_rule_collection.test_add"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Dnat"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNatRuleCollection_multiple(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Dnat"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
					testCheckAzureRMFirewallNatRuleCollectionExists(secondRule),
					resource.TestCheckResourceAttr(secondRule, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondRule, "priority", "200"),
					resource.TestCheckResourceAttr(secondRule, "action", "Snat"),
					resource.TestCheckResourceAttr(secondRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Dnat"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
					testCheckAzureRMFirewallNatRuleCollectionDoesNotExist("azurerm_firewall.test", "acctestnrc_add"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNatRuleCollection_update(t *testing.T) {
	firstResourceName := "azurerm_firewall_nat_rule_collection.test"
	secondResourceName := "azurerm_firewall_nat_rule_collection.test_add"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_multiple(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(firstResourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(firstResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallNatRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "200"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Snat"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNatRuleCollection_multipleUpdate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "300"),
					resource.TestCheckResourceAttr(firstResourceName, "action", "Snat"),
					resource.TestCheckResourceAttr(firstResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallNatRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "400"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNatRuleCollection_disappears(t *testing.T) {
	resourceName := "azurerm_firewall_nat_rule_collection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testCheckAzureRMFirewallNatRuleCollectionDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMFirewallNatRuleCollection_multipleRules(t *testing.T) {
	resourceName := "azurerm_firewall_nat_rule_collection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNatRuleCollection_multipleRules(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
				),
			},
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNatRuleCollection_updateFirewallTags(t *testing.T) {
	resourceName := "azurerm_firewall_nat_rule_collection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallNatRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallNatRuleCollection_updateFirewallTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallNatRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Dnat"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMFirewallNatRuleCollectionExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMFirewallNatRuleCollectionDoesNotExist(resourceName string, collectionName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		firewallName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testCheckAzureRMFirewallNatRuleCollectionDisappears(resourceName string) resource.TestCheckFunc {
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

func testAccAzureRMFirewallNatRuleCollection_basic(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, template)
}

func testAccAzureRMFirewallNatRuleCollection_requiresImport(rInt int, location string) string {
	template := testAccAzureRMFirewallNatRuleCollection_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "import" {
  name                = "${azurerm_firewall_nat_rule_collection.test.name}"
  azure_firewall_name = "${azurerm_firewall_nat_rule_collection.test.azure_firewall_name}"
  resource_group_name = "${azurerm_firewall_nat_rule_collection.test.resource_group_name}"
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
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, template)
}

func testAccAzureRMFirewallNatRuleCollection_updatedName(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, template)
}

func testAccAzureRMFirewallNatRuleCollection_multiple(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}

resource "azurerm_firewall_nat_rule_collection" "test_add" {
  name                = "acctestnrc_add"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  priority            = 200
  action              = "Snat"

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

func testAccAzureRMFirewallNatRuleCollection_multipleUpdate(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  priority            = 300
  action              = "Snat"

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

resource "azurerm_firewall_nat_rule_collection" "test_add" {
  name                = "acctestnrc_add"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
      "8.8.4.4",
    ]

    protocols = [
      "TCP",
    ]
  }
}
`, template)
}

func testAccAzureRMFirewallNatRuleCollection_multipleRules(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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

func testAccAzureRMFirewallNatRuleCollection_updateFirewallTags(rInt int, location string) string {
	template := testAccAzureRMFirewall_withTags(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_nat_rule_collection" "test" {
  name                = "acctestnrc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
      "8.8.8.8",
    ]

    protocols = [
      "Any",
    ]
  }
}
`, template)
}
