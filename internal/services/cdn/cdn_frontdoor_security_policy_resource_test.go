package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorSecurityPolicyResource struct{}

func TestAccCdnFrontdoorSecurityPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_security_policy", "test")
	r := CdnFrontdoorSecurityPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("azurerm_cdn_frontdoor_firewall_policy.test.cdn_frontdoor_profile_id", "azurerm_cdn_frontdoor_custom_domain.test.cdn_frontdoor_profile_id"),
	})
}

func TestAccCdnFrontdoorSecurityPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_security_policy", "test")
	r := CdnFrontdoorSecurityPolicyResource{}
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

func TestAccCdnFrontdoorSecurityPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_security_policy", "test")
	r := CdnFrontdoorSecurityPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("azurerm_cdn_frontdoor_firewall_policy.test.cdn_frontdoor_profile_id", "azurerm_cdn_frontdoor_custom_domain.test.cdn_frontdoor_profile_id"),
	})
}

func (r CdnFrontdoorSecurityPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontdoorSecurityPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorSecurityPoliciesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnFrontdoorSecurityPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                              = "accTestWAF%[1]d"
  resource_group_name               = azurerm_resource_group.test.name
  cdn_frontdoor_profile_id          = azurerm_cdn_frontdoor_profile.test.id
  sku_name                          = azurerm_cdn_frontdoor_profile.test.sku_name
  enabled                           = true
  mode                              = "Prevention"
  redirect_url                      = "https://www.fabrikam.com"
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

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "accTestCustomDomain-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  dns_zone_id = azurerm_dns_zone.test.id
  host_name   = join(".", ["fabrikam", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontdoorSecurityPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_security_policy" "test" {
  name                     = "accTestSecPol%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = azurerm_cdn_frontdoor_firewall_policy.test.id

      association {
        domain {
          cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.test.id
        }

        patterns_to_match = ["/*"]
      }
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorSecurityPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_security_policy" "import" {
  name                     = "accTestSecPol%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = azurerm_cdn_frontdoor_firewall_policy.test.id

      association {
        domain {
          cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.test.id
        }

        patterns_to_match = ["/*"]
      }
    }
  }
}
`, config, data.RandomInteger)
}

func (r CdnFrontdoorSecurityPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_security_policy" "test" {
  name                     = "accTestSecPol%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = azurerm_cdn_frontdoor_firewall_policy.test.id

      association {
        domain {
          cdn_frontdoor_custom_domain_id = azurerm_cdn_frontdoor_custom_domain.test.id
        }

        patterns_to_match = ["/*"]
      }
    }
  }
}
`, template, data.RandomInteger)
}
