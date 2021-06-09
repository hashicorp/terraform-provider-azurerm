package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type WebApplicationFirewallResource struct {
}

func TestAccWebApplicationFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

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

func TestAccWebApplicationFirewallPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.0.name").HasValue("Rule1"),
				check.That(data.ResourceName).Key("custom_rules.0.priority").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.0.rule_type").HasValue("MatchRule"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_variables.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_variables.0.variable_name").HasValue("RemoteAddr"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.operator").HasValue("IPMatch"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.negation_condition").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_values.#").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_values.0").HasValue("192.168.1.0/24"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_values.1").HasValue("10.0.0.0/24"),
				check.That(data.ResourceName).Key("custom_rules.0.action").HasValue("Block"),
				check.That(data.ResourceName).Key("custom_rules.1.name").HasValue("Rule2"),
				check.That(data.ResourceName).Key("custom_rules.1.priority").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.1.rule_type").HasValue("MatchRule"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.#").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_variables.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_variables.0.variable_name").HasValue("RemoteAddr"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.operator").HasValue("IPMatch"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.negation_condition").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_values.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_values.0").HasValue("192.168.1.0/24"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_variables.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_variables.0.variable_name").HasValue("RequestHeaders"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_variables.0.selector").HasValue("UserAgent"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.operator").HasValue("Contains"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.negation_condition").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_values.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_values.0").HasValue("windows"),
				check.That(data.ResourceName).Key("custom_rules.1.action").HasValue("Block"),
				check.That(data.ResourceName).Key("managed_rules.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.#").HasValue("2"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.0.match_variable").HasValue("RequestHeaderNames"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.0.selector").HasValue("x-shared-secret"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.0.selector_match_operator").HasValue("Equals"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.1.match_variable").HasValue("RequestCookieNames"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.1.selector").HasValue("too-much-fun"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.1.selector_match_operator").HasValue("EndsWith"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.type").HasValue("OWASP"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.version").HasValue("3.1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule_group_name").HasValue("REQUEST-920-PROTOCOL-ENFORCEMENT"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.disabled_rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.disabled_rules.0").HasValue("920300"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.disabled_rules.1").HasValue("920440"),
				check.That(data.ResourceName).Key("policy_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("policy_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("policy_settings.0.mode").HasValue("Prevention"),
				check.That(data.ResourceName).Key("policy_settings.0.request_body_check").HasValue("true"),
				check.That(data.ResourceName).Key("policy_settings.0.file_upload_limit_in_mb").HasValue("100"),
				check.That(data.ResourceName).Key("policy_settings.0.max_request_body_size_in_kb").HasValue("128"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebApplicationFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.0.name").HasValue("Rule1"),
				check.That(data.ResourceName).Key("custom_rules.0.priority").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.0.rule_type").HasValue("MatchRule"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_variables.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_variables.0.variable_name").HasValue("RemoteAddr"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.operator").HasValue("IPMatch"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.negation_condition").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_values.#").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_values.0").HasValue("192.168.1.0/24"),
				check.That(data.ResourceName).Key("custom_rules.0.match_conditions.0.match_values.1").HasValue("10.0.0.0/24"),
				check.That(data.ResourceName).Key("custom_rules.0.action").HasValue("Block"),
				check.That(data.ResourceName).Key("custom_rules.1.name").HasValue("Rule2"),
				check.That(data.ResourceName).Key("custom_rules.1.priority").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.1.rule_type").HasValue("MatchRule"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.#").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_variables.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_variables.0.variable_name").HasValue("RemoteAddr"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.operator").HasValue("IPMatch"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.negation_condition").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_values.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.0.match_values.0").HasValue("192.168.1.0/24"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_variables.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_variables.0.variable_name").HasValue("RequestHeaders"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_variables.0.selector").HasValue("UserAgent"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.operator").HasValue("Contains"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.negation_condition").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_values.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.1.match_conditions.1.match_values.0").HasValue("windows"),
				check.That(data.ResourceName).Key("custom_rules.1.action").HasValue("Block"),
				check.That(data.ResourceName).Key("managed_rules.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.#").HasValue("2"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.0.match_variable").HasValue("RequestHeaderNames"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.0.selector").HasValue("x-shared-secret"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.0.selector_match_operator").HasValue("Equals"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.1.match_variable").HasValue("RequestCookieNames"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.1.selector").HasValue("too-much-fun"),
				check.That(data.ResourceName).Key("managed_rules.0.exclusion.1.selector_match_operator").HasValue("EndsWith"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.type").HasValue("OWASP"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.version").HasValue("3.1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule_group_name").HasValue("REQUEST-920-PROTOCOL-ENFORCEMENT"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.disabled_rules.#").HasValue("2"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.disabled_rules.0").HasValue("920300"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.disabled_rules.1").HasValue("920440"),
				check.That(data.ResourceName).Key("policy_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("policy_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("policy_settings.0.mode").HasValue("Prevention"),
				check.That(data.ResourceName).Key("policy_settings.0.request_body_check").HasValue("true"),
				check.That(data.ResourceName).Key("policy_settings.0.file_upload_limit_in_mb").HasValue("100"),
				check.That(data.ResourceName).Key("policy_settings.0.max_request_body_size_in_kb").HasValue("128"),
			),
		},
		data.ImportStep(),
	})
}

func (t WebApplicationFirewallResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	name := id.Path["ApplicationGatewayWebApplicationFirewallPolicies"]

	resp, err := clients.Network.WebApplicationFirewallPoliciesClient.Get(ctx, resGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Web Application Firewall (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (WebApplicationFirewallResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_web_application_firewall_policy" "test" {
  name                = "acctestwafpolicy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location


  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"
    }
  }

  policy_settings {
    enabled = true
    mode    = "Detection"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (WebApplicationFirewallResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_web_application_firewall_policy" "test" {
  name                = "acctestwafpolicy-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    env = "test"
  }

  custom_rules {
    name      = "Rule1"
    priority  = 1
    rule_type = "MatchRule"

    match_conditions {
      match_variables {
        variable_name = "RemoteAddr"
      }

      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24", "10.0.0.0/24"]
    }

    action = "Block"
  }

  custom_rules {
    name      = "Rule2"
    priority  = 2
    rule_type = "MatchRule"

    match_conditions {
      match_variables {
        variable_name = "RemoteAddr"
      }

      operator           = "IPMatch"
      negation_condition = false
      match_values       = ["192.168.1.0/24"]
    }

    match_conditions {
      match_variables {
        variable_name = "RequestHeaders"
        selector      = "UserAgent"
      }

      operator           = "Contains"
      negation_condition = false
      match_values       = ["windows"]
      transforms         = ["Lowercase"]
    }

    action = "Block"
  }

  managed_rules {
    exclusion {
      match_variable          = "RequestHeaderNames"
      selector                = "x-shared-secret"
      selector_match_operator = "Equals"
    }

    exclusion {
      match_variable          = "RequestCookieNames"
      selector                = "too-much-fun"
      selector_match_operator = "EndsWith"
    }

    managed_rule_set {
      type    = "OWASP"
      version = "3.1"

      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
        disabled_rules = [
          "920300",
          "920440",
        ]
      }
    }
  }

  policy_settings {
    enabled = true
    mode    = "Prevention"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
