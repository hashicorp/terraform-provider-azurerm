package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorProfileRuleResource struct{}

func TestAccFrontdoorProfileRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_rule", "test")
	r := FrontdoorProfileRuleResource{}
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

func TestAccFrontdoorProfileRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_rule", "test")
	r := FrontdoorProfileRuleResource{}
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

func TestAccFrontdoorProfileRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_rule", "test")
	r := FrontdoorProfileRuleResource{}
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

func TestAccFrontdoorProfileRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_profile_rule", "test")
	r := FrontdoorProfileRuleResource{}
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

func (r FrontdoorProfileRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := rules.ParseRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontdoorProfileRulesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorProfileRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-cdn-%d"
  location = "%s"
}
resource "azurerm_frontdoor_profile_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_frontdoor_profile_rule_set" "test" {
  name           = "acctest-c-%d"
  cdn_profile_id = azurerm_frontdoor_profile_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r FrontdoorProfileRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_profile_rule" "test" {
  name            = "acctest-c-%d"
  cdn_rule_set_id = azurerm_frontdoor_profile_rule_set.test.id
  actions {
    name = ""
  }
  conditions {
    name = ""
  }
  match_processing_behavior = ""
  order                     = 0
}
`, template, data.RandomInteger)
}

func (r FrontdoorProfileRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_rule" "import" {
  name            = azurerm_frontdoor_profile_rule.test.name
  cdn_rule_set_id = azurerm_frontdoor_profile_rule_set.test.id
  actions {
    name = ""
  }
  conditions {
    name = ""
  }
  match_processing_behavior = ""
  order                     = 0
}
`, config)
}

func (r FrontdoorProfileRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_rule" "test" {
  name            = "acctest-c-%d"
  cdn_rule_set_id = azurerm_frontdoor_profile_rule_set.test.id
  actions {
    name = ""
  }
  conditions {
    name = ""
  }
  match_processing_behavior = ""
  order                     = 0
}
`, template, data.RandomInteger)
}

func (r FrontdoorProfileRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_profile_rule" "test" {
  name            = "acctest-c-%d"
  cdn_rule_set_id = azurerm_frontdoor_profile_rule_set.test.id
  actions {
    name = ""
  }
  conditions {
    name = ""
  }
  match_processing_behavior = ""
  order                     = 0
}
`, template, data.RandomInteger)
}
