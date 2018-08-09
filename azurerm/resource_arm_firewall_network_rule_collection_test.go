package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMFirewallNetworkRuleCollection_basic(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_addition(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	resourceNameAdd := "azurerm_firewall_network_rule_collection.test_add"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_basic(ri, testLocation())
	configAdd := testAccAzureRMFirewallNetworkRuleCollection_addition(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: configAdd,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc_add", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceNameAdd, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(resourceNameAdd, "priority", "200"),
					resource.TestCheckResourceAttr(resourceNameAdd, "action", "Deny"),
					resource.TestCheckResourceAttr(resourceNameAdd, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_removal(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	resourceNameAdd := "azurerm_firewall_network_rule_collection.test_add"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_addition(ri, testLocation())
	configRemove := testAccAzureRMFirewallNetworkRuleCollection_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc_add", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceNameAdd, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(resourceNameAdd, "priority", "200"),
					resource.TestCheckResourceAttr(resourceNameAdd, "action", "Deny"),
					resource.TestCheckResourceAttr(resourceNameAdd, "rule.#", "1"),
				),
			},
			{
				Config: configRemove,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_update(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	resourceNameAdd := "azurerm_firewall_network_rule_collection.test_add"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_addition(ri, testLocation())
	configUpdate := testAccAzureRMFirewallNetworkRuleCollection_update(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc_add", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceNameAdd, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(resourceNameAdd, "priority", "200"),
					resource.TestCheckResourceAttr(resourceNameAdd, "action", "Deny"),
					resource.TestCheckResourceAttr(resourceNameAdd, "rule.#", "1"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc_add", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "300"),
					resource.TestCheckResourceAttr(resourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceNameAdd, "name", "acctestnrc_add"),
					resource.TestCheckResourceAttr(resourceNameAdd, "priority", "400"),
					resource.TestCheckResourceAttr(resourceNameAdd, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceNameAdd, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_reapply(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_basic(ri, testLocation())
	deleteState := func(s *terraform.State) error {
		return s.Remove(resourceName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					deleteState,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_disappears(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testCheckAzureRMFirewallNetworkRuleCollectionDisappears("acctestnrc", &firewall),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_addrule(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_basic(ri, testLocation())
	configAddRule := testAccAzureRMFirewallNetworkRuleCollection_addRule(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: configAddRule,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallNetworkRuleCollection_removerule(t *testing.T) {
	var firewall network.AzureFirewall
	fwResourceName := "azurerm_firewall.test"
	resourceName := "azurerm_firewall_network_rule_collection.test"
	ri := acctest.RandInt()
	config := testAccAzureRMFirewallNetworkRuleCollection_basic(ri, testLocation())
	configAddRule := testAccAzureRMFirewallNetworkRuleCollection_addRule(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAddRule,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(fwResourceName, &firewall),
					testCheckAzureRMFirewallNetworkRuleCollectionExists("acctestnrc", &firewall),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestnrc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func testAccAzureRMFirewallNetworkRuleCollection_basic(rInt int, location string) string {
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
`, testAccAzureRMFirewall_basic(rInt, location))
}

func testAccAzureRMFirewallNetworkRuleCollection_addition(rInt int, location string) string {
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
`, testAccAzureRMFirewall_basic(rInt, location))
}

func testAccAzureRMFirewallNetworkRuleCollection_update(rInt int, location string) string {
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
`, testAccAzureRMFirewall_basic(rInt, location))
}

func testAccAzureRMFirewallNetworkRuleCollection_addRule(rInt int, location string) string {
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
`, testAccAzureRMFirewall_basic(rInt, location))
}

func testCheckAzureRMFirewallNetworkRuleCollectionExists(name string, firewall *network.AzureFirewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findArmFirewallNetworkRuleCollectionByName(firewall, name)
		if !exists {
			return fmt.Errorf("A Network Rule Collection with name %q cannot be found", name)
		}

		return nil
	}
}

func testCheckAzureRMFirewallNetworkRuleCollectionNotExists(name string, firewall *network.AzureFirewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findArmFirewallNetworkRuleCollectionByName(firewall, name)
		if exists {
			return fmt.Errorf("A Network Rule Collection with name %q has been found", name)
		}

		return nil
	}
}

func testCheckAzureRMFirewallNetworkRuleCollectionDisappears(name string, firewall *network.AzureFirewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, _, exists := findArmFirewallNetworkRuleCollectionByName(firewall, name)
		if !exists {
			return fmt.Errorf("A Network Rule Collection with name %q cannot be found", name)
		}

		updatedCollection := removeArmFirewallNetworkRuleCollectionByName(firewall, name)
		firewall.AzureFirewallPropertiesFormat.NetworkRuleCollections = updatedCollection

		id, err := parseAzureResourceID(*firewall.ID)
		if err != nil {
			return err
		}

		ipConfigs := fixArmFirewallIPConfiguration(firewall)
		firewall.AzureFirewallPropertiesFormat.IPConfigurations = &ipConfigs

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *firewall.Name, *firewall)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating Azure Firewall: %+v", err)
		}

		err = future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for completion for Azure Firewall: %+v", err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *firewall.Name)
		return err
	}
}
