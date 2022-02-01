package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorProfilePolicyResource struct{}

func TestAccFrontdoorProfilePolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_policy", "test")
	r := FrontdoorProfilePolicyResource{}
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

func TestAccFrontdoorProfilePolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_policy", "test")
	r := FrontdoorProfilePolicyResource{}
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

func TestAccFrontdoorProfilePolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_policy", "test")
	r := FrontdoorProfilePolicyResource{}
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

func TestAccFrontdoorProfilePolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_policy", "test")
	r := FrontdoorProfilePolicyResource{}
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

func (r FrontdoorProfilePolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapplicationfirewallpolicies.ParseCdnWebApplicationFirewallPoliciesID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.WebApplicationFirewallPoliciesClient
	resp, err := client.PoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorProfilePolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-cdn-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r FrontdoorProfilePolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_profile_policy" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
  custom_rules {
    rules {
      action        = ""
      enabled_state = ""
      match_conditions {
        match_value      = []
        match_variable   = ""
        negate_condition = false
        operator         = ""
        selector         = ""
        transforms       = ""
      }
      name     = ""
      priority = 0
    }
  }
  etag     = ""
  location = "%s"
  managed_rules {
    managed_rule_sets {
      anomaly_score = 0
      rule_group_overrides {
        rule_group_name = ""
        rules {
          action        = ""
          enabled_state = ""
          rule_id       = ""
        }
      }
      rule_set_type    = ""
      rule_set_version = ""
    }
  }
  policy_settings {
    default_custom_block_response_body        = ""
    default_custom_block_response_status_code = ""
    default_redirect_url                      = ""
    enabled_state                             = ""
    mode                                      = ""
  }
  rate_limit_rules {
    rules {
      action        = ""
      enabled_state = ""
      match_conditions {
        match_value      = []
        match_variable   = ""
        negate_condition = false
        operator         = ""
        selector         = ""
        transforms       = ""
      }
      name                           = ""
      priority                       = 0
      rate_limit_duration_in_minutes = 0
      rate_limit_threshold           = 0
    }
  }
  sku {
    name = ""
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r FrontdoorProfilePolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_policy" "import" {
  name                = azurerm_frontdoor_profile_policy.test.name
  resource_group_name = azurerm_resource_group.test.name
  custom_rules {
    rules {
      action        = ""
      enabled_state = ""
      match_conditions {
        match_value      = []
        match_variable   = ""
        negate_condition = false
        operator         = ""
        selector         = ""
        transforms       = ""
      }
      name     = ""
      priority = 0
    }
  }
  etag     = ""
  location = "%s"
  managed_rules {
    managed_rule_sets {
      anomaly_score = 0
      rule_group_overrides {
        rule_group_name = ""
        rules {
          action        = ""
          enabled_state = ""
          rule_id       = ""
        }
      }
      rule_set_type    = ""
      rule_set_version = ""
    }
  }
  policy_settings {
    default_custom_block_response_body        = ""
    default_custom_block_response_status_code = ""
    default_redirect_url                      = ""
    enabled_state                             = ""
    mode                                      = ""
  }
  rate_limit_rules {
    rules {
      action        = ""
      enabled_state = ""
      match_conditions {
        match_value      = []
        match_variable   = ""
        negate_condition = false
        operator         = ""
        selector         = ""
        transforms       = ""
      }
      name                           = ""
      priority                       = 0
      rate_limit_duration_in_minutes = 0
      rate_limit_threshold           = 0
    }
  }
  sku {
    name = ""
  }

  tags = {
    ENV = "Test"
  }
}
`, config, data.Locations.Primary)
}

func (r FrontdoorProfilePolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_policy" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
  custom_rules {
    rules {
      action        = ""
      enabled_state = ""
      match_conditions {
        match_value      = []
        match_variable   = ""
        negate_condition = false
        operator         = ""
        selector         = ""
        transforms       = ""
      }
      name     = ""
      priority = 0
    }
  }
  etag     = ""
  location = "%s"
  managed_rules {
    managed_rule_sets {
      anomaly_score = 0
      rule_group_overrides {
        rule_group_name = ""
        rules {
          action        = ""
          enabled_state = ""
          rule_id       = ""
        }
      }
      rule_set_type    = ""
      rule_set_version = ""
    }
  }
  policy_settings {
    default_custom_block_response_body        = ""
    default_custom_block_response_status_code = ""
    default_redirect_url                      = ""
    enabled_state                             = ""
    mode                                      = ""
  }
  rate_limit_rules {
    rules {
      action        = ""
      enabled_state = ""
      match_conditions {
        match_value      = []
        match_variable   = ""
        negate_condition = false
        operator         = ""
        selector         = ""
        transforms       = ""
      }
      name                           = ""
      priority                       = 0
      rate_limit_duration_in_minutes = 0
      rate_limit_threshold           = 0
    }
  }
  sku {
    name = ""
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r FrontdoorProfilePolicyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_policy" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
