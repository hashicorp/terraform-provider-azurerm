package firewall_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type FirewallPolicyApplicationRuleCollectionResource struct{}

func TestAccFirewallPolicyApplicationRuleCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_application_rule_collection", "test")
	r := FirewallPolicyApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyApplicationRuleCollection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_application_rule_collection", "test")
	r := FirewallPolicyApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyApplicationRuleCollection_completePremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_application_rule_collection", "test")
	r := FirewallPolicyApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyApplicationRuleCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_application_rule_collection", "test")
	r := FirewallPolicyApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyApplicationRuleCollection_updatePremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_application_rule_collection", "test")
	r := FirewallPolicyApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatePremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completePremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewallPolicyApplicationRuleCollection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_application_rule_collection", "test")
	r := FirewallPolicyApplicationRuleCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (FirewallPolicyApplicationRuleCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallPolicyRuleCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Firewall.FirewallPolicyRuleGroupClient.Get(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	idx := indexFunc(*resp.RuleCollections, func(policy network.BasicFirewallPolicyRuleCollection) bool {
		info, _ := policy.AsFirewallPolicyFilterRuleCollection()
		return *info.Name == id.RuleCollectionName
	})

	return utils.Bool(idx != -1), nil
}

func (FirewallPolicyApplicationRuleCollectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RC-%[1]d"
  location = "%[2]s"
}
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RC-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RC-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
}
resource "azurerm_firewall_policy_application_rule_collection" "test" {
  name                     = "acctest-fwpolicy-RC-%[1]d"
  rule_collection_group_id = azurerm_firewall_policy_rule_collection_group.test.id
  priority                 = 500
  action                   = "Allow"
  rule {
    name = "app_rule_collection_rule1"
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses  = ["10.0.0.1"]
    destination_fqdns = ["*.microsoft.com"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FirewallPolicyApplicationRuleCollectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RC-%[1]d"
  location = "%[2]s"
}
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RC-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dns {
    proxy_enabled = true
  }
}
resource "azurerm_ip_group" "test_source" {
  name                = "acctestIpGroupForFirewallPolicySource"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}
resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RC-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
}
resource "azurerm_firewall_policy_application_rule_collection" "test" {
  name                     = "acctest-fwpolicy-RC-%[1]d"
  rule_collection_group_id = azurerm_firewall_policy_rule_collection_group.test.id
  priority                 = 500
  action                   = "Deny"
  rule {
    name = "app_rule_collection_rule1"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses  = ["10.0.0.1"]
    destination_fqdns = ["pluginsdk.io"]
  }
  rule {
    name = "app_rule_collection_rule2"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_ip_groups  = [azurerm_ip_group.test_source.id]
    destination_fqdns = ["pluginsdk.io"]
  }
  rule {
    name = "app_rule_collection_rule3"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    protocols {
      type = "Mssql"
      port = 1443
    }
    source_addresses      = ["10.0.0.1"]
    destination_fqdn_tags = ["WindowsDiagnostics"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FirewallPolicyApplicationRuleCollectionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RC-%[1]d"
  location = "%[2]s"
}
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RC-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dns {
    proxy_enabled = true
  }
}
resource "azurerm_ip_group" "test_source" {
  name                = "acctestIpGroupForFirewallPolicySource"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}
resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RC-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
}
resource "azurerm_firewall_policy_application_rule_collection" "test" {
  name                     = "acctest-fwpolicy-RC-%[1]d"
  rule_collection_group_id = azurerm_firewall_policy_rule_collection_group.test.id
  priority                 = 500
  action                   = "Deny"
  rule {
    name = "app_rule_collection_rule1"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses  = ["10.0.0.1", "10.0.0.2"]
    destination_fqdns = ["pluginsdk.io"]
  }
  rule {
    name = "app_rule_collection_rule2"
    protocols {
      type = "Http"
      port = 80
    }
    source_ip_groups  = [azurerm_ip_group.test_source.id]
    destination_fqdns = ["pluginsdk.io"]
  }
  rule {
    name = "app_rule_collection_rule3"
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
`, data.RandomInteger, data.Locations.Primary)
}

func (FirewallPolicyApplicationRuleCollectionResource) completePremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RC-%[1]d"
  location = "%[2]s"
}
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RC-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  dns {
    proxy_enabled = true
  }
}
resource "azurerm_ip_group" "test_source1" {
  name                = "acctestIpGroupForFirewallPolicySource1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}
resource "azurerm_ip_group" "test_source2" {
  name                = "acctestIpGroupForFirewallPolicySource2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["4.3.2.1/32", "87.65.43.0/24"]
}
resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RC-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
}
resource "azurerm_firewall_policy_application_rule_collection" "test" {
  name                     = "acctest-fwpolicy-RC-%[1]d"
  rule_collection_group_id = azurerm_firewall_policy_rule_collection_group.test.id
  priority                 = 500
  action                   = "Deny"
  rule {
    name        = "app_rule_collection_rule1"
    description = "app_rule_collection_rule1"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses      = ["10.0.0.1"]
    destination_addresses = ["10.0.0.1"]
    destination_urls      = ["www.google.com/en"]
    terminate_tls         = true
    web_categories        = ["News"]
  }
  rule {
    name        = "app_rule_collection_rule2"
    description = "app_rule_collection_rule2"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_ip_groups      = [azurerm_ip_group.test_source1.id, azurerm_ip_group.test_source2.id]
    destination_addresses = ["10.0.0.1"]
    destination_fqdns     = ["pluginsdk.io"]
    terminate_tls         = true
    web_categories        = ["News"]
  }
  rule {
    name        = "app_rule_collection_rule3"
    description = "app_rule_collection_rule3"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses      = ["10.0.0.1"]
    destination_addresses = ["10.0.0.1"]
    destination_urls      = ["www.google.com/en"]
    terminate_tls         = true
    web_categories        = ["News"]
  }
  rule {
    name        = "app_rule_collection_rule4"
    description = "app_rule_collection_rule4"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses      = ["10.0.0.1", "10.0.0.2"]
    destination_addresses = ["10.0.0.1", "10.0.0.2"]
    destination_urls      = ["www.google.com/en", "www.google.com/cn"]
    terminate_tls         = true
    web_categories        = ["News", "Arts"]
  }
  rule {
    name        = "app_rule_collection_rule5"
    description = "app_rule_collection_rule5"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_ip_groups      = [azurerm_ip_group.test_source1.id, azurerm_ip_group.test_source2.id]
    destination_addresses = ["10.0.0.1", "10.0.0.2"]
    destination_fqdns     = ["pluginsdk.io", "pluginframework.io"]
    terminate_tls         = true
    web_categories        = ["News", "Arts"]
  }
  rule {
    name        = "app_rule_collection_rule6"
    description = "app_rule_collection_rule6"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses      = ["10.0.0.1", "10.0.0.2"]
    destination_addresses = ["10.0.0.1", "10.0.0.2"]
    destination_urls      = ["www.google.com/en", "www.google.com/cn"]
    terminate_tls         = true
    web_categories        = ["News", "Arts"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FirewallPolicyApplicationRuleCollectionResource) updatePremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fwpolicy-RC-%[1]d"
  location = "%[2]s"
}
resource "azurerm_firewall_policy" "test" {
  name                = "acctest-fwpolicy-RC-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dns {
    proxy_enabled = true
  }
}
resource "azurerm_ip_group" "test_source1" {
  name                = "acctestIpGroupForFirewallPolicySource1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["1.2.3.4/32", "12.34.56.0/24"]
}
resource "azurerm_ip_group" "test_source2" {
  name                = "acctestIpGroupForFirewallPolicySource2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidrs               = ["4.3.2.1/32", "87.65.43.0/24"]
}
resource "azurerm_firewall_policy_rule_collection_group" "test" {
  name               = "acctest-fwpolicy-RC-%[1]d"
  firewall_policy_id = azurerm_firewall_policy.test.id
  priority           = 500
}
resource "azurerm_firewall_policy_application_rule_collection" "test" {
  name                     = "acctest-fwpolicy-RC-%[1]d"
  rule_collection_group_id = azurerm_firewall_policy_rule_collection_group.test.id
  priority                 = 500
  action                   = "Deny"

  rule {
    name        = "app_rule_collection_rule1"
    description = "app_rule_collection_rule1"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses      = ["10.0.0.1"]
    destination_addresses = ["10.0.0.1"]
    destination_urls      = ["www.google.com/en"]
    terminate_tls         = true
    web_categories        = ["News"]
  }
  rule {
    name        = "app_rule_collection_rule2"
    description = "app_rule_collection_rule2"
    protocols {
      type = "Http"
      port = 80
    }
    source_ip_groups      = [azurerm_ip_group.test_source1.id]
    destination_addresses = ["10.0.0.1"]
    destination_fqdns     = ["pluginsdk.io"]
    terminate_tls         = true
    web_categories        = ["News"]
  }
  rule {
    name        = "app_rule_collection_rule3"
    description = "app_rule_collection_rule3"
    protocols {
      type = "Http"
      port = 80
    }
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses      = ["10.0.0.1", "10.0.0.2"]
    destination_addresses = ["10.0.0.1", "10.0.0.2"]
    destination_urls      = ["www.google.com/en"]
    terminate_tls         = true
    web_categories        = ["News"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FirewallPolicyApplicationRuleCollectionResource) requiresImport(data acceptance.TestData) string {
	template := FirewallPolicyApplicationRuleCollectionResource{}.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_firewall_policy_application_rule_collection" "import" {
  name                     = azurerm_firewall_policy_application_rule_collection.test.name
  rule_collection_group_id = azurerm_firewall_policy_application_rule_collection.test.rule_collection_group_id
  priority                 = azurerm_firewall_policy_application_rule_collection.test.priority
  action                   = azurerm_firewall_policy_application_rule_collection.test.action
  rule {
    name = "app_rule_collection_rule1"
    protocols {
      type = "Https"
      port = 443
    }
    source_addresses  = ["10.0.0.1"]
    destination_fqdns = ["*.microsoft.com"]
  }
}
`, template)
}

func indexFunc[T any](s []T, f func(T) bool) int {
	for i := 0; i < len(s); i++ {
		if f(s[i]) {
			return i
		}
	}
	return -1
}
