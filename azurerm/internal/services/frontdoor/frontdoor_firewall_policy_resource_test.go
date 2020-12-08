package frontdoor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type FrontDoorFirewallPolicyResource struct {
}

func TestAccFrontDoorFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")
	r := FrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("mode").HasValue("Prevention"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFrontDoorFirewallPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")
	r := FrontDoorFirewallPolicyResource{}

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

func TestAccFrontDoorFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")
	r := FrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.update(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("mode").HasValue("Prevention"),
			),
		},
		{
			Config: r.update(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("mode").HasValue("Prevention"),
				check.That(data.ResourceName).Key("custom_rule.1.name").HasValue("Rule2"),
				check.That(data.ResourceName).Key("custom_rule.2.name").HasValue("Rule3"),
			),
		},
		{
			Config: r.update(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_rule.1.name").DoesNotExist(),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("mode").HasValue("Prevention"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFrontDoorFirewallPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_firewall_policy", "test")
	r := FrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.update(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("testAccFrontDoorWAF%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("mode").HasValue("Prevention"),
				check.That(data.ResourceName).Key("redirect_url").HasValue("https://www.contoso.com"),
				check.That(data.ResourceName).Key("custom_block_response_status_code").HasValue("403"),
				check.That(data.ResourceName).Key("custom_rule.0.name").HasValue("Rule1"),
				check.That(data.ResourceName).Key("custom_rule.1.name").HasValue("Rule2"),
				check.That(data.ResourceName).Key("managed_rule.0.type").HasValue("DefaultRuleSet"),
				check.That(data.ResourceName).Key("managed_rule.0.exclusion.0.match_variable").HasValue("QueryStringArgNames"),
				check.That(data.ResourceName).Key("managed_rule.0.override.1.exclusion.0.selector").HasValue("really_not_suspicious"),
				check.That(data.ResourceName).Key("managed_rule.0.override.1.rule.0.exclusion.0.selector").HasValue("innocent"),
				check.That(data.ResourceName).Key("managed_rule.1.type").HasValue("Microsoft_BotManagerRuleSet"),
			),
		},
		data.ImportStep(),
	})
}

func (FrontDoorFirewallPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.WebApplicationFirewallPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Frontdoor.FrontDoorsPolicyClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Front Door Firewall Policy %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.WebApplicationFirewallPolicyProperties != nil), nil
}

func (FrontDoorFirewallPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                = "testAccFrontDoorWAF%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r FrontDoorFirewallPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor_firewall_policy" "import" {
  name                = azurerm_frontdoor_firewall_policy.test.name
  resource_group_name = azurerm_frontdoor_firewall_policy.test.resource_group_name
}
`, r.basic(data))
}

func (r FrontDoorFirewallPolicyResource) update(data acceptance.TestData, update bool) string {
	if update {
		return r.updated(data)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                              = "testAccFrontDoorWAF%[1]d"
  resource_group_name               = azurerm_resource_group.test.name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  custom_rule {
    name                           = "Rule1"
    enabled                        = true
    priority                       = 1
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
    }
  }

  managed_rule {
    type    = "DefaultRuleSet"
    version = "preview-0.1"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933111"
        enabled = false
        action  = "Block"
      }
    }
  }

  managed_rule {
    type    = "BotProtection"
    version = "preview-0.1"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (FrontDoorFirewallPolicyResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%[2]s"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                              = "testAccFrontDoorWAF%[1]d"
  resource_group_name               = azurerm_resource_group.test.name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  custom_rule {
    name                           = "Rule1"
    enabled                        = true
    priority                       = 1
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
    }
  }

  custom_rule {
    name                           = "Rule2"
    enabled                        = true
    priority                       = 2
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "RemoteAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_condition {
      match_variable     = "RequestHeader"
      selector           = "UserAgent"
      operator           = "Contains"
      negation_condition = false
      match_values       = ["windows"]
      transforms         = ["Lowercase", "Trim"]
    }
  }

  custom_rule {
    name                           = "Rule3"
    enabled                        = true
    priority                       = 3
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Block"

    match_condition {
      match_variable     = "SocketAddr"
      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_condition {
      match_variable     = "RequestHeader"
      selector           = "UserAgent"
      operator           = "Contains"
      negation_condition = false
      match_values       = ["windows"]
      transforms         = ["Lowercase", "Trim"]
    }
  }

  managed_rule {
    type    = "DefaultRuleSet"
    version = "1.0"

    exclusion {
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
      selector       = "not_suspicious"
    }

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Block"
      }
    }

    override {
      rule_group_name = "SQLI"

      exclusion {
        match_variable = "QueryStringArgNames"
        operator       = "Equals"
        selector       = "really_not_suspicious"
      }

      rule {
        rule_id = "942200"
        action  = "Block"

        exclusion {
          match_variable = "QueryStringArgNames"
          operator       = "Equals"
          selector       = "innocent"
        }
      }
    }
  }

  managed_rule {
    type    = "Microsoft_BotManagerRuleSet"
    version = "1.0"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
