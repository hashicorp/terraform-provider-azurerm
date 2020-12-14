package firewall_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFirewallPolicyRuleCollectionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyRuleCollectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyRuleCollectionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyRuleCollectionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicyRuleCollectionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyRuleCollectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyRuleCollectionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyRuleCollectionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicyRuleCollectionGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyRuleCollectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyRuleCollectionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyRuleCollectionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFirewallPolicyRuleCollectionGroup_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyRuleCollectionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFirewallPolicyRuleCollectionGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyRuleCollectionGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFirewallPolicyRuleCollectionGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallPolicyRuleCollectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewallPolicyRuleCollectionGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallPolicyRuleCollectionGroupExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFirewallPolicyRuleCollectionGroup_requiresImport),
		},
	})
}

func testCheckAzureRMFirewallPolicyRuleCollectionGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Firewall Policy Rule Collection Group not found: %s", resourceName)
		}

		id, err := parse.FirewallPolicyRuleCollectionGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Firewall Policy Rule Collection Group %q (Resource Group %q) does not exist", id.RuleCollectionGroupName, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Network.FirewallPolicyRuleGroups: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFirewallPolicyRuleCollectionGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_firewall_policy_rule_collection_group" {
			continue
		}

		id, err := parse.FirewallPolicyRuleCollectionGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
		if err == nil {
			return fmt.Errorf("Network.FirewallPolicyRuleGroups still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Network.FirewallPolicyRuleGroups: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMFirewallPolicyRuleCollectionGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RCG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RCG-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RCG-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMFirewallPolicyRuleCollectionGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RCG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RCG-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dns {
    network_rule_fqdn_enabled = false
  }
}

resource "azurerm_ip_group" "test_source" {
  name                = "acctestIpGroupForFirewallPolicySource"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}

resource "azurerm_ip_group" "test_destination" {
  name                = "acctestIpGroupForFirewallPolicyDest"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
}

resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RCG-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
  application_rule_collection {
    name     = "app_rule_collection1"
    priority = 500
    action   = "Deny"
    rule {
      name = "app_rule_collection1_rule1"
      protocols {
        type = "Http"
        port = 80
      }
      protocols {
        type = "Https"
        port = 443
      }
      source_addresses  = ["10.0.0.1"]
      destination_fqdns = [".microsoft.com"]
    }
    rule {
      name = "app_rule_collection1_rule2"
      protocols {
        type = "Http"
        port = 80
      }
      protocols {
        type = "Https"
        port = 443
      }
      source_ip_groups  = [azurerm_ip_group.test_source.id]
      destination_fqdns = [".microsoft.com"]
    }
    rule {
      name = "app_rule_collection1_rule3"
      protocols {
        type = "Http"
        port = 80
      }
      protocols {
        type = "Https"
        port = 443
      }
      source_addresses      = ["10.0.0.1"]
      destination_fqdn_tags = ["WindowsDiagnostics"]
    }
  }

  network_rule_collection {
    name     = "network_rule_collection1"
    priority = 400
    action   = "Deny"
    rule {
      name                  = "network_rule_collection1_rule1"
      protocols             = ["TCP", "UDP"]
      source_addresses      = ["10.0.0.1"]
      destination_addresses = ["192.168.1.1", "ApiManagement"]
      destination_ports     = ["80", "1000-2000"]
    }
    rule {
      name              = "network_rule_collection1_rule2"
      protocols         = ["TCP", "UDP"]
      source_addresses  = ["10.0.0.1"]
      destination_fqdns = ["time.windows.com"]
      destination_ports = ["80", "1000-2000"]
    }
    rule {
      name                  = "network_rule_collection1_rule3"
      protocols             = ["TCP", "UDP"]
      source_ip_groups      = [azurerm_ip_group.test_source.id]
      destination_ip_groups = [azurerm_ip_group.test_destination.id]
      destination_ports     = ["80", "1000-2000"]
    }
  }

  nat_rule_collection {
    name     = "nat_rule_collection1"
    priority = 300
    action   = "Dnat"
    rule {
      name                = "nat_rule_collection1_rule1"
      protocols           = ["TCP", "UDP"]
      source_addresses    = ["10.0.0.1", "10.0.0.2"]
      destination_address = "192.168.1.1"
      destination_ports   = ["80", "1000-2000"]
      translated_address  = "192.168.0.1"
      translated_port     = "8080"
    }
    rule {
      name                = "nat_rule_collection1_rule2"
      protocols           = ["TCP", "UDP"]
      source_ip_groups    = [azurerm_ip_group.test_source.id]
      destination_address = "192.168.1.1"
      destination_ports   = ["80", "1000-2000"]
      translated_address  = "192.168.0.1"
      translated_port     = "8080"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMFirewallPolicyRuleCollectionGroup_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RCG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RCG-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dns {
    network_rule_fqdn_enabled = false
  }
}

resource "azurerm_ip_group" "test_source" {
  name                = "acctestIpGroupForFirewallPolicySource"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}

resource "azurerm_ip_group" "test_destination" {
  name                = "acctestIpGroupForFirewallPolicyDest"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["192.168.0.0/25", "192.168.0.192/26"]
}

resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RCG-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
  application_rule_collection {
    name     = "app_rule_collection1"
    priority = 500
    action   = "Deny"
    rule {
      name = "app_rule_collection1_rule1"
      protocols {
        type = "Http"
        port = 80
      }
      protocols {
        type = "Https"
        port = 443
      }
      source_addresses  = ["10.0.0.1", "10.0.0.2"]
      destination_fqdns = [".microsoft.com"]
    }
    rule {
      name = "app_rule_collection1_rule2"
      protocols {
        type = "Http"
        port = 80
      }
      source_ip_groups  = [azurerm_ip_group.test_source.id]
      destination_fqdns = [".microsoft.com"]
    }
    rule {
      name = "app_rule_collection1_rule3"
      protocols {
        type = "Http"
        port = 80
      }
      protocols {
        type = "Https"
        port = 443
      }
      source_addresses      = ["10.0.0.1", "10.0.0.2"]
      destination_fqdn_tags = ["WindowsDiagnostics"]
    }
  }

  network_rule_collection {
    name     = "network_rule_collection1"
    priority = 400
    action   = "Deny"
    rule {
      name                  = "network_rule_collection1_rule1"
      protocols             = ["TCP", "UDP"]
      source_addresses      = ["10.0.0.1"]
      destination_addresses = ["192.168.1.2", "ApiManagement"]
      destination_ports     = ["80", "1000-2000"]
    }
    rule {
      name              = "network_rule_collection1_rule2"
      protocols         = ["TCP", "UDP"]
      source_addresses  = ["10.0.0.1", "10.0.0.2"]
      destination_fqdns = ["time.windows.com"]
      destination_ports = ["80", "1000-2000"]
    }
    rule {
      name                  = "network_rule_collection1_rule3"
      protocols             = ["TCP"]
      source_ip_groups      = [azurerm_ip_group.test_source.id]
      destination_ip_groups = [azurerm_ip_group.test_destination.id]
      destination_ports     = ["80", "1000-2000"]
    }
  }

  nat_rule_collection {
    name     = "nat_rule_collection1"
    priority = 300
    action   = "Dnat"
    rule {
      name                = "nat_rule_collection1_rule1"
      protocols           = ["TCP", "UDP"]
      source_addresses    = ["10.0.0.1", "10.0.0.2"]
      destination_address = "192.168.1.1"
      destination_ports   = ["80", "1000-2000"]
      translated_address  = "192.168.0.1"
      translated_port     = "8080"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMFirewallPolicyRuleCollectionGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFirewallPolicyRuleCollectionGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy_rule_collection_group" "import" {
  name               = azurerm_firewall_policy_rule_collection_group.test.name
  firewall_policy_id = azurerm_firewall_policy_rule_collection_group.test.firewall_policy_id
  priority           = azurerm_firewall_policy_rule_collection_group.test.priority
}
`, template)
}
