package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAzureRMFirewallApplicationRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.source_addresses.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.target_fqdns.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.port", "443"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.type", "Https"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMFirewallApplicationRuleCollection_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_firewall_application_rule_collection"),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_updatedName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_updatedName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule2"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_multipleRuleCollections(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	secondRule := "azurerm_firewall_application_rule_collection.test_add"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionExists(secondRule),
					resource.TestCheckResourceAttr(secondRule, "name", "acctestarc_add"),
					resource.TestCheckResourceAttr(secondRule, "priority", "200"),
					resource.TestCheckResourceAttr(secondRule, "action", "Deny"),
					resource.TestCheckResourceAttr(secondRule, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionDoesNotExist("azurerm_firewall.test", "acctestarc_add"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")
	secondResourceName := "azurerm_firewall_application_rule_collection.test_add"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multiple(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "name", "acctestarc_add"),
					resource.TestCheckResourceAttr(secondResourceName, "priority", "200"),
					resource.TestCheckResourceAttr(secondResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(secondResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "300"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
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
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					testCheckAzureRMFirewallApplicationRuleCollectionDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "2"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.port", "8000"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.type", "Http"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.port", "8001"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.type", "Https"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_updateProtocols(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.port", "8000"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.type", "Http"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.port", "8001"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.type", "Https"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocolsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.port", "9000"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.type", "Https"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.port", "9001"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.type", "Http"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.port", "8000"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.0.type", "Http"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.port", "8001"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.protocol.1.type", "Https"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_updateFirewallTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
				),
			},
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_updateFirewallTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctestarc"),
					resource.TestCheckResourceAttr(data.ResourceName, "priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
				),
			},
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_ipGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallApplicationRuleCollection_ipGroups(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallApplicationRuleCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallApplicationRuleCollection_noSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_application_rule_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMFirewallApplicationRuleCollection_noSource(data),
				ExpectError: regexp.MustCompile(fmt.Sprintf("at least one of %q and %q must be specified", "source_addresses", "source_ip_groups")),
			},
		},
	})
}

func testCheckAzureRMFirewallApplicationRuleCollectionExists(resourceName string) resource.TestCheckFunc {
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

func testAccAzureRMFirewallApplicationRuleCollection_basic(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFirewallApplicationRuleCollection_basic(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_updatedName(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multiple(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multipleUpdate(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multipleRules(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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

func testAccAzureRMFirewallApplicationRuleCollection_multipleProtocols(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_multipleProtocolsUpdate(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_updateFirewallTags(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_withTags(data)
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
`, template)
}

func testAccAzureRMFirewallApplicationRuleCollection_ipGroups(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group" "test" {
  name                = "acctestIpGroupForFirewallAppRules"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
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
      azurerm_ip_group.test.id,
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

func testAccAzureRMFirewallApplicationRuleCollection_noSource(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_basic(data)
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
`, template)
}
