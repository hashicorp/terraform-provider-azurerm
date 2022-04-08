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

type FrontdoorSecurityPolicyResource struct{}

func TestAccFrontdoorSecurityPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_security_policy", "test")
	r := FrontdoorSecurityPolicyResource{}
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

func TestAccFrontdoorSecurityPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_security_policy", "test")
	r := FrontdoorSecurityPolicyResource{}
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

func TestAccFrontdoorSecurityPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_security_policy", "test")
	r := FrontdoorSecurityPolicyResource{}
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

func TestAccFrontdoorSecurityPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_security_policy", "test")
	r := FrontdoorSecurityPolicyResource{}
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
	})
}

func (r FrontdoorSecurityPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontdoorSecurityPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontdoorSecurityPoliciesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r FrontdoorSecurityPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-afdx-%[1]d"
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

resource "azurerm_frontdoor_profile_profile" "test" {
  name                = "acctest-c-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r FrontdoorSecurityPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_security_policy" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile_profile.test.id

  web_application_firewall {
    waf_policy_id = azurerm_frontdoor_firewall_policy.test.id

    association {
      domain {
        id = "foo"
      }
      domain {
        id = "bar"
      }
      domain {
        id = "spoon"
      }
      patterns_to_match = ["/foo", "/bar"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r FrontdoorSecurityPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_security_policy" "import" {
  name                 = azurerm_frontdoor_security_policy.test.name
  frontdoor_profile_id = azurerm_frontdoor_profile_profile.test.id

  web_application_firewall {
    waf_policy_id = azurerm_frontdoor_firewall_policy.test.id

    association {
      domain {
        id = "foo"
      }
      domain {
        id = "bar"
      }
      domain {
        id = "spoon"
      }
      patterns_to_match = ["/foo", "/bar"]
    }
  }
}
`, config)
}

func (r FrontdoorSecurityPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_security_policy" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile_profile.test.id

  web_application_firewall {
    waf_policy_id = azurerm_frontdoor_firewall_policy.test.id

    association {
      domain {
        id = "foo"
      }
      domain {
        id = "bar"
      }
      domain {
        id = "spoon"
      }
      patterns_to_match = ["/foo", "/bar"]
    }

    association {
      domain {
        id = "foo"
      }
      domain {
        id      = "badf00d"
        enabled = false
      }

      patterns_to_match = ["foo"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r FrontdoorSecurityPolicyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_security_policy" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile_profile.test.id

  web_application_firewall {
    waf_policy_id = azurerm_frontdoor_firewall_policy.test.id

    association {
      domain {
        id = "foo"
      }
      domain {
        id = "bar"
      }
      domain {
        id      = "spoon"
        enabled = false
      }
      patterns_to_match = ["/spoon", "/bar", "/*"]
    }
  }
}
`, template, data.RandomInteger)
}
