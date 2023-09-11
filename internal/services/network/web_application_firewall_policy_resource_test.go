// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WebApplicationFirewallResource struct{}

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
				check.That(data.ResourceName).Key("custom_rules.#").HasValue("3"),
				check.That(data.ResourceName).Key("custom_rules.0.enabled").HasValue("true"),
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
				check.That(data.ResourceName).Key("custom_rules.2.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.2.name").HasValue("Rule3"),
				check.That(data.ResourceName).Key("custom_rules.2.priority").HasValue("3"),
				check.That(data.ResourceName).Key("custom_rules.2.rule_type").HasValue("MatchRule"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.0.match_variables.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.0.match_variables.0.variable_name").HasValue("RemoteAddr"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.0.operator").HasValue("IPMatch"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.0.negation_condition").HasValue("false"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.0.match_values.#").HasValue("2"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.0.match_values.0").HasValue("192.168.1.0/24"),
				check.That(data.ResourceName).Key("custom_rules.2.match_conditions.0.match_values.1").HasValue("10.0.0.0/24"),
				check.That(data.ResourceName).Key("custom_rules.2.action").HasValue("Block"),
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
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.version").HasValue("3.2"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule_group_name").HasValue("REQUEST-920-PROTOCOL-ENFORCEMENT"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.0.id").HasValue("920300"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.0.action").HasValue("Log"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.1.id").HasValue("920440"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.1.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.1.action").HasValue("Block"),
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
				check.That(data.ResourceName).Key("custom_rules.#").HasValue("3"),
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
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.version").HasValue("3.2"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.#").HasValue("1"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule_group_name").HasValue("REQUEST-920-PROTOCOL-ENFORCEMENT"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.0.id").HasValue("920300"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.0.action").HasValue("Log"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.1.id").HasValue("920440"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.1.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("managed_rules.0.managed_rule_set.0.rule_group_override.0.rule.1.action").HasValue("Block"),
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

func TestAccWebApplicationFirewallPolicy_updateOverrideRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateOverrideRules(data),
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

func TestAccWebApplicationFirewallPolicy_knownCVEs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.knownCVEs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebApplicationFirewallPolicy_OperatorAny(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.operatorAny(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebApplicationFirewallPolicy_excludedRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.excludedRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateExcludedRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebApplicationFirewallPolicy_updateDisabledRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disabledRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateDisabledRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebApplicationFirewallPolicy_LogScrubbing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withLogScrubbing(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebApplicationFirewallPolicy_updateCustomRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_application_firewall_policy", "test")
	r := WebApplicationFirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateCustomRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t WebApplicationFirewallResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.WebApplicationFirewallPolicies.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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
    enabled   = true
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

  custom_rules {
    enabled   = false
    name      = "Rule3"
    priority  = 3
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
      version = "3.2"

      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
        rule {
          id      = "920300"
          enabled = true
          action  = "Log"
        }

        rule {
          id      = "920440"
          enabled = true
          action  = "Block"
        }
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

func (WebApplicationFirewallResource) updateOverrideRules(data acceptance.TestData) string {
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
      version = "3.2"

      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"

        rule {
          id      = "920440"
          enabled = true
          action  = "Block"
        }
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

func (WebApplicationFirewallResource) knownCVEs(data acceptance.TestData) string {
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

  policy_settings {
    enabled                     = true
    file_upload_limit_in_mb     = 100
    max_request_body_size_in_kb = 128
    mode                        = "Prevention"
    request_body_check          = false
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.1"

      rule_group_override {
        disabled_rules = [
          "800112",
          "800111",
          "800110",
          "800100",
          "800113",
        ]
        rule_group_name = "Known-CVEs"
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (WebApplicationFirewallResource) excludedRules(data acceptance.TestData) string {
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

      excluded_rule_set {
        rule_group {
          rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
          excluded_rules = [
            "920100",
            "920120",
          ]
        }
      }
    }

    managed_rule_set {
      type    = "OWASP"
      version = "3.2"

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

func (WebApplicationFirewallResource) operatorAny(data acceptance.TestData) string {
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
        variable_name = "PostArgs"
        selector      = "value"
      }
      operator = "Any"
    }

    action = "Log"
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.2"

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

func (WebApplicationFirewallResource) updateExcludedRules(data acceptance.TestData) string {
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

      excluded_rule_set {
        rule_group {
          rule_group_name = "REQUEST-913-SCANNER-DETECTION"
          excluded_rules = [
            "913100",
            "913101",
          ]
        }

        rule_group {
          rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
          excluded_rules = [
            "920100",
            "920120",
          ]
        }
      }
    }

    managed_rule_set {
      type    = "OWASP"
      version = "3.2"

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

func (WebApplicationFirewallResource) disabledRules(data acceptance.TestData) string {
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

  policy_settings {
    enabled                     = true
    mode                        = "Prevention"
    request_body_check          = true
    file_upload_limit_in_mb     = 100
    max_request_body_size_in_kb = 2000
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.2"

      rule_group_override {
        rule_group_name = "REQUEST-931-APPLICATION-ATTACK-RFI"
        disabled_rules  = ["931130"]
      }

      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
        disabled_rules = [
          "920320", # Missing User Agent Header
          "920230"  # Multiple URL Encoding Detected
        ]
      }

      rule_group_override {
        rule_group_name = "REQUEST-942-APPLICATION-ATTACK-SQLI"
        disabled_rules = [
          "942450",
          "942430",
          "942440",
          "942370",
          "942340",
          "942260",
          "942200",
          "942330",
          "942120",
          "942110",
          "942150",
          "942410",
          "942130",
          "942100"
        ]
      }

      rule_group_override {
        rule_group_name = "REQUEST-941-APPLICATION-ATTACK-XSS"
        disabled_rules = [
          "941340"
        ]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (WebApplicationFirewallResource) updateDisabledRules(data acceptance.TestData) string {
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

  policy_settings {
    enabled                     = true
    mode                        = "Prevention"
    request_body_check          = true
    file_upload_limit_in_mb     = 100
    max_request_body_size_in_kb = 2000
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.2"

      rule_group_override {
        rule_group_name = "REQUEST-931-APPLICATION-ATTACK-RFI"
        disabled_rules  = ["931130"]
      }

      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
        disabled_rules = [
          "920320", # Missing User Agent Header
          "920230"  # Multiple URL Encoding Detected
        ]
      }

      #NEW BLOCK
      rule_group_override {
        rule_group_name = "REQUEST-932-APPLICATION-ATTACK-RCE"
        disabled_rules  = ["932100"]
      }

      rule_group_override {
        rule_group_name = "REQUEST-942-APPLICATION-ATTACK-SQLI"
        disabled_rules = [
          "942450",
          "942430",
          "942440",
          "942370",
          "942340",
          "942260",
          "942200",
          "942330",
          "942120",
          "942110",
          "942150",
          "942410",
          "942130",
          "942100"
        ]
      }

      rule_group_override {
        rule_group_name = "REQUEST-941-APPLICATION-ATTACK-XSS"
        disabled_rules = [
          "941340"
        ]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (WebApplicationFirewallResource) withLogScrubbing(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_web_application_firewall_policy" "test" {
  name                = "acctestwafpolicy-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_settings {
    enabled = true
    mode    = "Detection"
    log_scrubbing {
      enabled = true
      rule {
        enabled                 = true
        match_variable          = "RequestHeaderNames"
        selector_match_operator = "Equals"
        selector                = "User-Agent"
      }
    }
  }

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.2"

      rule_group_override {
        rule_group_name = "REQUEST-920-PROTOCOL-ENFORCEMENT"
        rule {
          id      = "920300"
          enabled = true
          action  = "Log"
        }

        rule {
          id      = "920440"
          enabled = true
          action  = "Block"
        }
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (WebApplicationFirewallResource) customRules(data acceptance.TestData) string {
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
      match_values       = ["Windows"]
    }

    action = "Block"
  }

  custom_rules {
    enabled   = true
    name      = "Rule3"
    priority  = 3
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
    enabled   = false
    name      = "Rule4"
    priority  = 4
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
    enabled   = true
    name      = "Rule5"
    priority  = 5
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
    enabled   = false
    name      = "Rule6"
    priority  = 6
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

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.2"
    }
  }

  policy_settings {
    enabled = true
    mode    = "Detection"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (WebApplicationFirewallResource) updateCustomRules(data acceptance.TestData) string {
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

  custom_rules {
    enabled   = true
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
    enabled   = false
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
      match_values       = ["Windows"]
    }

    action = "Block"
  }

  custom_rules {
    name      = "Rule3"
    priority  = 3
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
    name      = "Rule4"
    priority  = 4
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
    enabled   = false
    name      = "Rule5"
    priority  = 5
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
    enabled   = true
    name      = "Rule6"
    priority  = 6
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

  managed_rules {
    managed_rule_set {
      type    = "OWASP"
      version = "3.2"
    }
  }

  policy_settings {
    enabled = true
    mode    = "Detection"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
