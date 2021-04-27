package firewall_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FirewallPolicyRuleCollectionGroupResource struct {
}

func TestAccFirewallPolicyRuleCollectionGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")
	r := FirewallPolicyRuleCollectionGroupResource{}

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

func TestAccFirewallPolicyRuleCollectionGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")
	r := FirewallPolicyRuleCollectionGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyRuleCollectionGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")
	r := FirewallPolicyRuleCollectionGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyRuleCollectionGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")
	r := FirewallPolicyRuleCollectionGroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (FirewallPolicyRuleCollectionGroupResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id, err = parse.FirewallPolicyRuleCollectionGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Firewall.FirewallPolicyRuleGroupClient.Get(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.FirewallPolicyRuleCollectionGroupProperties != nil), nil
}

func (FirewallPolicyRuleCollectionGroupResource) basic(data acceptance.TestData) string {
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

func (FirewallPolicyRuleCollectionGroupResource) complete(data acceptance.TestData) string {
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
    proxy_enabled             = true
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
      destination_fqdns = ["terraform.io"]
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
      destination_fqdns = ["terraform.io"]
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
    rule {
      name                  = "network_rule_collection1_rule4"
      protocols             = ["ICMP"]
      source_ip_groups      = [azurerm_ip_group.test_source.id]
      destination_ip_groups = [azurerm_ip_group.test_destination.id]
      destination_ports     = ["*"]
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

func (FirewallPolicyRuleCollectionGroupResource) update(data acceptance.TestData) string {
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
    proxy_enabled             = true
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
      destination_fqdns = ["terraform.io"]
    }
    rule {
      name = "app_rule_collection1_rule2"
      protocols {
        type = "Http"
        port = 80
      }
      source_ip_groups  = [azurerm_ip_group.test_source.id]
      destination_fqdns = ["terraform.io"]
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
      destination_ports     = ["80", "1-65535"]
    }
    rule {
      name              = "network_rule_collection1_rule2"
      protocols         = ["TCP", "UDP"]
      source_addresses  = ["10.0.0.1", "10.0.0.2"]
      destination_fqdns = ["time.windows.com"]
      destination_ports = ["80", "1-65535"]
    }
    rule {
      name                  = "network_rule_collection1_rule3"
      protocols             = ["TCP"]
      source_ip_groups      = [azurerm_ip_group.test_source.id]
      destination_ip_groups = [azurerm_ip_group.test_destination.id]
      destination_ports     = ["80", "1-65535"]
    }
    rule {
      name                  = "network_rule_collection1_rule4"
      protocols             = ["ICMP"]
      source_ip_groups      = [azurerm_ip_group.test_source.id]
      destination_ip_groups = [azurerm_ip_group.test_destination.id]
      destination_ports     = ["*"]
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

func (FirewallPolicyRuleCollectionGroupResource) requiresImport(data acceptance.TestData) string {
	template := FirewallPolicyRuleCollectionGroupResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall_policy_rule_collection_group" "import" {
  name               = azurerm_firewall_policy_rule_collection_group.test.name
  firewall_policy_id = azurerm_firewall_policy_rule_collection_group.test.firewall_policy_id
  priority           = azurerm_firewall_policy_rule_collection_group.test.priority
}
`, template)
}
