// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorFirewallPolicyResource struct{}

func TestAccCdnFrontDoorFirewallPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
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

func TestAccCdnFrontDoorFirewallPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}

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

func TestAccCdnFrontDoorFirewallPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("request_body_check_enabled").HasValue("false"),
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
		{
			Config: r.update(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_complete(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
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

func TestAccCdnFrontDoorFirewallPolicy_DRSOnePointOh(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.DRSOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSOnePointOhUpdate(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.DRSOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.DRSTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.DRSOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSOnePointOhError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.DRSOnePointOhError(data),
			ExpectError: regexp.MustCompile("'AnomalyScoring' is only valid in managed rules that are DRS 2.0 and above"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSOnePointOhTypeError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.DRSOnePointOhTypeError(data),
			ExpectError: regexp.MustCompile("If you wish to use the 'Microsoft_DefaultRuleSet' type please update your 'version' field to be '1.1', '2.0' or '2.1'"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSTwoPointOh(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.DRSTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSTwoPointOhUpdate(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.DRSTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.DRSOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.DRSTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSTwoPointOhError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.DRSTwoPointOhError(data),
			ExpectError: regexp.MustCompile("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSTwoPointOhTypeError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.DRSTwoPointOhTypeError(data),
			ExpectError: regexp.MustCompile("If you wish to use the 'DefaultRuleSet' type please update your 'version' field to be '1.0' or 'preview-0.1'"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSTwoPointOneAction(t *testing.T) {
	// NOTE: Regression test case for issue #19561
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.DRSTwoPointOneActionLog(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_DRSTwoPointOneActionError(t *testing.T) {
	// NOTE: Regression test case for issue #19561
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.DRSTwoPointOneActionError(data),
			ExpectError: regexp.MustCompile("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above"),
		},
	})
}

func (CdnFrontDoorFirewallPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorFirewallPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cdn.FrontDoorLegacyFirewallPoliciesClient.Get(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (CdnFrontDoorFirewallPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontDoorFirewallPolicyResource) basic(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "import" {
  name                = azurerm_cdn_frontdoor_firewall_policy.test.name
  resource_group_name = azurerm_cdn_frontdoor_firewall_policy.test.resource_group_name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"
}
`, r.basic(data))
}

func (r CdnFrontDoorFirewallPolicyResource) update(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Detection"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="
  request_body_check_enabled        = false

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
    action  = "Block"

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
    action  = "Log"
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) complete(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
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
    type    = "Microsoft_DefaultRuleSet"
    version = "2.0"
    action  = "Block"

    exclusion {
      match_variable = "RequestBodyJsonArgNames"
      operator       = "Equals"
      selector       = "not_suspicious"
    }

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "AnomalyScoring"
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
        action  = "AnomalyScoring"

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
    action  = "Block"
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSOnePointOh(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "DefaultRuleSet"
    version = "1.0"
    action  = "Log"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Block"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSOnePointOhError(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "DefaultRuleSet"
    version = "1.0"
    action  = "Log"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "AnomalyScoring"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSOnePointOhTypeError(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "Microsoft_DefaultRuleSet"
    version = "1.0"
    action  = "Log"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Block"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSTwoPointOh(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "Microsoft_DefaultRuleSet"
    version = "2.0"
    action  = "Log"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "AnomalyScoring"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSTwoPointOhError(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "Microsoft_DefaultRuleSet"
    version = "2.0"
    action  = "Log"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Block"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSTwoPointOhTypeError(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "DefaultRuleSet"
    version = "2.0"
    action  = "Log"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Block"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSTwoPointOneActionLog(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "Microsoft_DefaultRuleSet"
    version = "2.1"
    action  = "Block"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "AnomalyScoring"
      }

      rule {
        rule_id = "933110"
        enabled = false
        action  = "Log"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) DRSTwoPointOneActionError(data acceptance.TestData) string {
	tmp := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%d"
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.contoso.com"
  custom_block_response_status_code = 403
  custom_block_response_body        = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="

  managed_rule {
    type    = "Microsoft_DefaultRuleSet"
    version = "2.1"
    action  = "Block"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "AnomalyScoring"
      }

      rule {
        rule_id = "933110"
        enabled = false
        action  = "Redirect"
      }
    }
  }
}
`, tmp, data.RandomInteger)
}
