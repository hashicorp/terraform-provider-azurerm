// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("captcha_cookie_expiration_in_minutes").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengePolicyBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jSChallengePolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("45"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_captchaPolicyBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.captchaPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("captcha_cookie_expiration_in_minutes").HasValue("45"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengeManagedRuleBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jsChallengeManagedRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengeCustomRuleBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jsChallengeCustomRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_captchaCustomRuleBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.captchaCustomRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingStandardBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingStandardBasic(data),
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

func TestAccCdnFrontDoorFirewallPolicy_drsOnePointOhUpdate(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.drsOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.drsTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.drsOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengePolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("30"),
			),
		},
		data.ImportStep(),
		{
			Config: r.jSChallengePolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("45"),
			),
		},
		data.ImportStep(),
		{
			Config: r.jSChallengePolicyBasicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("1440"),
			),
		},
		data.ImportStep(),
		{
			// NOTE: Since this is an O+C field, when the field is removed
			// from the config it will get the last value from the state
			// file so you need to verify that the last tests value is
			// passed as the value for the field...
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("js_challenge_cookie_expiration_in_minutes").HasValue("1440"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_captchaPolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("captcha_cookie_expiration_in_minutes").HasValue("30"),
			),
		},
		data.ImportStep(),
		{
			Config: r.captchaPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("captcha_cookie_expiration_in_minutes").HasValue("45"),
			),
		},
		data.ImportStep(),
		{
			Config: r.captchaPolicyBasicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("captcha_cookie_expiration_in_minutes").HasValue("1440"),
			),
		},
		data.ImportStep(),
		{
			// NOTE: Since this is an O+C field, when the field is removed
			// from the config it will get the last value from the state
			// file so you need to verify that the last tests value is
			// passed as the value for the field...
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("captcha_cookie_expiration_in_minutes").HasValue("1440"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengeManagedRuleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jsChallengeManagedRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeManagedRuleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeManagedRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeManagedRuleRemove(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeManagedRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengeCustomRuleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jsChallengeCustomRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeCustomRuleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeCustomRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeCustomRuleRemove(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.jsChallengeCustomRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingBasicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingDisabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingDisabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingRuleDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingDisabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingStandardBasicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingStandardBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingStandardDisabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingStandardRuleDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingStandardBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingCompleteUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingDisabledBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_complete(t *testing.T) {
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

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsOnePointOh(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.drsOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsOnePointOhError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.drsOnePointOhError(data),
			ExpectError: regexp.MustCompile(`"AnomalyScoring" is only valid in managed rules where 'type' is DRS`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsOnePointOhTypeError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.drsOnePointOhTypeError(data),
			ExpectError: regexp.MustCompile("If you wish to use the 'Microsoft_DefaultRuleSet' type please update your 'version' field to be '1.1', '2.0' or '2.1'"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsTwoPointOh(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.drsTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsTwoPointOhUpdate(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.drsTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.drsOnePointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.drsTwoPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsTwoPointOhError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.drsTwoPointOhError(data),
			ExpectError: regexp.MustCompile("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsTwoPointOhTypeError(t *testing.T) {
	// NOTE: Regression test case for issue #19088
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.drsTwoPointOhTypeError(data),
			ExpectError: regexp.MustCompile("If you wish to use the 'DefaultRuleSet' type please update your 'version' field to be '1.0' or 'preview-0.1'"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsTwoPointOneAction(t *testing.T) {
	// NOTE: Regression test case for issue #19561
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.drsTwoPointOneActionLog(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_drsTwoPointOneActionError(t *testing.T) {
	// NOTE: Regression test case for issue #19561
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.drsTwoPointOneActionError(data),
			ExpectError: regexp.MustCompile("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above"),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengePolicyStandardSku(t *testing.T) {
	// NOTE: Regression test case for issue #28716
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jsChallengePolicyStandardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_captchaPolicyStandardSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.captchaPolicyStandardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengeDRSError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.jsChallengeDRSError(data),
			ExpectError: regexp.MustCompile(`"JSChallenge" is only valid if the managed rules 'type' is 'Microsoft_BotManagerRuleSet'`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengePolicyStandardSkuError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.jsChallengePolicyStandardSkuError(data),
			ExpectError: regexp.MustCompile(`'js_challenge_cookie_expiration_in_minutes' field is only supported with the "Premium_AzureFrontDoor" sku`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_captchaPolicyStandardSkuError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.captchaPolicyStandardSkuError(data),
			ExpectError: regexp.MustCompile(`'captcha_cookie_expiration_in_minutes' field is only supported with the "Premium_AzureFrontDoor" sku`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_standardSkuManagedRuleError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.standardSkuManagedRuleError(data),
			ExpectError: regexp.MustCompile(`'managed_rule' code block is only supported with the "Premium_AzureFrontDoor" sku`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_jsChallengeStandardSkuCustomRuleActionError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.jsChallengeStandardSkuCustomRuleActionError(data),
			ExpectError: regexp.MustCompile(`'custom_rule' blocks with the 'action' type of 'JSChallenge' are only supported for the "Premium_AzureFrontDoor" sku`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_skuDowngradeError(t *testing.T) {
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
		{
			Config:      r.basicStandard(data),
			ExpectError: regexp.MustCompile(`downgrading from the "Premium_AzureFrontDoor" sku to the "Standard_AzureFrontDoor" sku is not supported`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingRequestIPAddressError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logScrubbingRequestIPAddressError(data),
			ExpectError: regexp.MustCompile(`the "RequestIPAddress" 'match_variable' must use the "EqualsAny" 'operator', got "Equals"`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingRequestUriError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logScrubbingRequestUriError(data),
			ExpectError: regexp.MustCompile(`the "RequestUri" 'match_variable' must use the "EqualsAny" 'operator', got "Equals"`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingSelectorEqualsError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logScrubbingSelectorEqualsError(data),
			ExpectError: regexp.MustCompile(`the 'selector' field must be set when the "Equals" 'operator' is used, got "nil"`),
		},
	})
}

func TestAccCdnFrontDoorFirewallPolicy_logScrubbingSelectorEqualsAnyError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
	r := CdnFrontDoorFirewallPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logScrubbingSelectorEqualsAnyError(data),
			ExpectError: regexp.MustCompile(`the 'selector' field cannot be set when the "EqualsAny" 'operator' is used, got "bar"`),
		},
	})
}

func (CdnFrontDoorFirewallPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := waf.ParseFrontDoorWebApplicationFirewallPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	_, err = clients.Cdn.FrontDoorFirewallPoliciesClient.PoliciesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
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

func (CdnFrontDoorFirewallPolicyResource) templateStandard(data acceptance.TestData) string {
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
  sku_name            = "Standard_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontDoorFirewallPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) basicStandard(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jSChallengePolicyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  js_challenge_cookie_expiration_in_minutes = 45
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) captchaPolicyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  captcha_cookie_expiration_in_minutes = 45
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jSChallengePolicyBasicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  js_challenge_cookie_expiration_in_minutes = 1440
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) captchaPolicyBasicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  captcha_cookie_expiration_in_minutes = 1440
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "foo"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingStandardBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "foo"
    }
  }
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingDisabledBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  log_scrubbing {
    enabled = false

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "foo"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingRuleDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = false
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "foo"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingStandardRuleDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = false
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "foo"
    }
  }
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingStandardDisabledBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
  mode                = "Prevention"

  log_scrubbing {
    enabled = false

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "foo"
    }
  }
}
`, r.templateStandard(data), data.RandomInteger)
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) complete(data acceptance.TestData) string {
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

  js_challenge_cookie_expiration_in_minutes = 30

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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingComplete(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestUri"
      operator       = "EqualsAny"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "keyToBlock"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestCookieNames"
      operator       = "Equals"
      selector       = "Chocolate Chip"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestBodyPostArgNames"
      operator       = "Equals"
      selector       = "var"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestBodyJsonArgNames"
      operator       = "Equals"
      selector       = "JsonValue"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
      selector       = "foo"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingCompleteUpdate(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = false

    scrubbing_rule {
      enabled        = false
      match_variable = "RequestHeaderNames"
      operator       = "Equals"
      selector       = "keyToBlock"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestCookieNames"
      operator       = "Equals"
      selector       = "ChocolateChip"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestBodyJsonArgNames"
      operator       = "Equals"
      selector       = "foo"
    }

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestBodyJsonArgNames"
      operator       = "EqualsAny"
    }

    scrubbing_rule {
      enabled        = false
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
      selector       = "bar"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsOnePointOh(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsOnePointOhError(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsOnePointOhTypeError(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsTwoPointOh(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsTwoPointOhError(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsTwoPointOhTypeError(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsTwoPointOneActionLog(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) drsTwoPointOneActionError(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeDRSError(data acceptance.TestData) string {
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
    version = "preview-0.1"
    action  = "Block"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "JSChallenge"
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengePolicyStandardSku(data acceptance.TestData) string {
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
    name                           = "ShortUserAgents"
    enabled                        = true
    type                           = "MatchRule"
    priority                       = 500
    rate_limit_threshold           = 1
    rate_limit_duration_in_minutes = 5
    action                         = "Allow"

    match_condition {
      match_variable = "RequestHeader"
      selector       = "User-Agent"
      operator       = "LessThanOrEqual"
      match_values   = ["15"]
    }
  }
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) captchaPolicyStandardSku(data acceptance.TestData) string {
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
    name                           = "ShortUserAgents"
    enabled                        = true
    type                           = "MatchRule"
    priority                       = 500
    rate_limit_threshold           = 1
    rate_limit_duration_in_minutes = 5
    action                         = "Allow"

    match_condition {
      match_variable = "RequestHeader"
      selector       = "User-Agent"
      operator       = "LessThanOrEqual"
      match_values   = ["15"]
    }
  }
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengePolicyStandardSkuError(data acceptance.TestData) string {
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

  js_challenge_cookie_expiration_in_minutes = 45
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) captchaPolicyStandardSkuError(data acceptance.TestData) string {
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

  captcha_cookie_expiration_in_minutes = 45
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) standardSkuManagedRuleError(data acceptance.TestData) string {
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
    version = "preview-0.1"
    action  = "Block"

    override {
      rule_group_name = "PHP"

      rule {
        rule_id = "933100"
        enabled = false
        action  = "Allow"
      }
    }
  }
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeStandardSkuCustomRuleActionError(data acceptance.TestData) string {
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
    name                           = "ShortUserAgents"
    enabled                        = true
    type                           = "MatchRule"
    priority                       = 500
    rate_limit_threshold           = 1
    rate_limit_duration_in_minutes = 5
    action                         = "JSChallenge"

    match_condition {
      match_variable = "RequestHeader"
      selector       = "User-Agent"
      operator       = "LessThanOrEqual"
      match_values   = ["15"]
    }
  }
}
`, r.templateStandard(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeManagedRuleBasic(data acceptance.TestData) string {
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
    type    = "Microsoft_BotManagerRuleSet"
    version = "1.0"
    action  = "Log"

    override {
      rule_group_name = "BadBots"

      rule {
        rule_id = "Bot100200"
        enabled = true
        action  = "JSChallenge"
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeCustomRuleBasic(data acceptance.TestData) string {
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
    name                           = "CustomJSChallenge"
    enabled                        = true
    priority                       = 2
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "JSChallenge"

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
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) captchaCustomRuleBasic(data acceptance.TestData) string {
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
    name                           = "CustomCaptcha"
    enabled                        = true
    priority                       = 2
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "CAPTCHA"

    match_condition {
      match_variable     = "RequestUri"
      operator           = "Contains"
      negation_condition = false
      match_values       = ["/path/captcha"]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeManagedRuleUpdate(data acceptance.TestData) string {
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
    type    = "Microsoft_BotManagerRuleSet"
    version = "1.0"
    action  = "Log"

    override {
      rule_group_name = "BadBots"

      rule {
        rule_id = "Bot100200"
        enabled = true
        action  = "Allow"
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeCustomRuleUpdate(data acceptance.TestData) string {
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
    name                           = "CustomJSChallenge"
    enabled                        = true
    priority                       = 2
    rate_limit_duration_in_minutes = 1
    rate_limit_threshold           = 10
    type                           = "MatchRule"
    action                         = "Allow"

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
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeManagedRuleRemove(data acceptance.TestData) string {
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
    type    = "Microsoft_BotManagerRuleSet"
    version = "1.0"
    action  = "Log"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) jsChallengeCustomRuleRemove(data acceptance.TestData) string {
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
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingRequestIPAddressError(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestIPAddress"
      operator       = "Equals"
      selector       = "bar"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingRequestUriError(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestUri"
      operator       = "Equals"
      selector       = "bar"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingSelectorEqualsError(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingSelectorEqualsAnyError(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestBodyJsonArgNames"
      operator       = "EqualsAny"
      selector       = "bar"
    }
  }
}
`, r.template(data), data.RandomInteger)
}
