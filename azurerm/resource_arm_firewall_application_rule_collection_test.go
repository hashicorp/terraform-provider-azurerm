package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMFirewallApplicationRuleCollection_basic(t *testing.T) {
	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.source_addresses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.target_fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.type", "Https"),
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

func TestAccAzureRMFirewallApplicationRuleCollection_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMFirewallApplicationRuleCollection_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_firewall_application_rule_collection"),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_updatedName(t *testing.T) {
	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "rule1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_updatedName(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "rule2"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_multipleRuleCollections(t *testing.T) {
	firstRule := "azurerm_firewall_application_rule_collection.test"
	secondRule := "azurerm_firewall_application_rule_collection.test_add"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestarc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Allow"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multiple(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestarc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Allow"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionExists(secondRule),
					resource.TestCheckResourceAttr(secondRule, "name", "acctestarc_add"),
					resource.TestCheckResourceAttr(secondRule, "priority", "200"),
					resource.TestCheckResourceAttr(secondRule, "action", "Deny"),
					resource.TestCheckResourceAttr(secondRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(firstRule),
					resource.TestCheckResourceAttr(firstRule, "name", "acctestarc"),
					resource.TestCheckResourceAttr(firstRule, "priority", "100"),
					resource.TestCheckResourceAttr(firstRule, "action", "Allow"),
					resource.TestCheckResourceAttr(firstRule, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionDoesNotExist("azurerm_firewall.test", "acctestarc_add"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_update(t *testing.T) {
	firstResourceName := "azurerm_firewall_application_rule_collection.test"
	secondResourceName := "azurerm_firewall_application_rule_collection.test_add"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multiple(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(firstResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(firstResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestarc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "200"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleUpdate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(firstResourceName, "priority", "300"),
					resource.TestCheckResourceAttr(firstResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(firstResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestarc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "400"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_disappears(t *testing.T) {
	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_multipleRules(t *testing.T) {
	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleRules(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(t *testing.T) {
	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.type", "Http"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.port", "8001"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.type", "Https"),
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

func TestAccAzureRMFirewallApplicationRuleCollection_updateProtocols(t *testing.T) {
	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.type", "Http"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.port", "8001"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.type", "Https"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocolsUpdate(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.port", "9000"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.type", "Https"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.port", "9001"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.type", "Http"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.0.type", "Http"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.port", "8001"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.protocol.1.type", "Https"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_updateFirewallTags(t *testing.T) {
	resourceName := "azurerm_firewall_application_rule_collection.test"
	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "rule1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_updateFirewallTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
					resource.TestCheckResourceAttr(resourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "rule1"),
				),
			},
		},
	})
}

func testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		firewallName := rs.Primary.Attributes["azure_firewall_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.AzureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		read, err := client.Get(ctx, resourceGroup, firewallName)
		if err != nil {
			return err
		}

		found := false
		for _, collection := range *read.AzureFirewallPropertiesFormat.ApplicationRuleCollections {
			if *collection.Name == name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Expected Application Rule Collection %q (Firewall %q / Resource Group %q) to exist but it didn't", name, firewallName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMFirewallApplicationRuleCollectionDoesNotExist(resourceName string, collectionName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		firewallName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.AzureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		read, err := client.Get(ctx, resourceGroup, firewallName)
		if err != nil {
			return err
		}

		for _, collection := range *read.AzureFirewallPropertiesFormat.ApplicationRuleCollections {
			if *collection.Name == collectionName {
				return fmt.Errorf("Application Rule Collection %q exists in Firewall %q: %+v", collectionName, firewallName, collection)
			}
		}

		return nil
	}
}

func testCheckAzureRMFirewallApplicationRuleCollectionDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		firewallName := rs.Primary.Attributes["azure_firewall_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.AzureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		read, err := client.Get(ctx, resourceGroup, firewallName)
		if err != nil {
			return err
		}

		rules := make([]network.AzureFirewallApplicationRuleCollection, 0)
		for _, collection := range *read.AzureFirewallPropertiesFormat.ApplicationRuleCollections {
			if *collection.Name != name {
				rules = append(rules, collection)
			}
		}

		read.AzureFirewallPropertiesFormat.ApplicationRuleCollections = &rules
		future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, read)
		if err != nil {
			return fmt.Errorf("Error removing Application Rule Collection from Firewall: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for the removal of Application Rule Collection from Firewall: %+v", err)
		}

		_, err = client.Get(ctx, resourceGroup, firewallName)
		return err
	}
}

func testAccAzureRMFirewallApplicationRuleCollection_basic(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_requiresImport(rInt int, location string) string {
	template := testAccAzureRMFirewallApplicationRuleCollection_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "import" {
  name                = "${azurerm_firewall_application_rule_collection.test.name}"
  azure_firewall_name = "${azurerm_firewall_application_rule_collection.test.azure_firewall_name}"
  resource_group_name = "${azurerm_firewall_application_rule_collection.test.resource_group_name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_updatedName(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multiple(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multipleUpdate(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multipleRules(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multipleProtocolsUpdate(rInt int, location string) string {
	template := testAccAzureRMFirewall_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_updateFirewallTags(rInt int, location string) string {
	template := testAccAzureRMFirewall_withTags(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_application_rule_collection" "test" {
  name                = "acctestarc"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
`, template)
}
