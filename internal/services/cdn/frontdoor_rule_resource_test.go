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

type FrontdoorRuleResource struct{}

func TestAccFrontdoorRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_rule", "test")
	r := FrontdoorRuleResource{}
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

func TestAccFrontdoorRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_rule", "test")
	r := FrontdoorRuleResource{}
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

func TestAccFrontdoorRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_rule", "test")
	r := FrontdoorRuleResource{}
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

func TestAccFrontdoorRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_rule", "test")
	r := FrontdoorRuleResource{}
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

func (r FrontdoorRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := rules.ParseRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontdoorRulesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-afdx-%d"
  location = "%s"
}

resource "azurerm_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_frontdoor_rule_set" "test" {
  name           = "acctest-c-%d"
  cdn_profile_id = azurerm_frontdoor_profile_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r FrontdoorRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_rule" "test" {
  name                  = "acctest-c-%d"
  frontdoor_rule_set_id = azurerm_frontdoor_rule_set.test.id

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

func (r FrontdoorRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_rule" "import" {
  name                  = azurerm_frontdoor_rule.test.name
  frontdoor_rule_set_id = azurerm_frontdoor_rule_set.test.id

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

func (r FrontdoorRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_rule" "test" {
  name                  = "acctest-c-%d"
  frontdoor_rule_set_id = azurerm_frontdoor_rule_set.test.id

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

func (r FrontdoorRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_rule" "test" {
  name                  = "acctest-c-%d"
  frontdoor_rule_set_id = azurerm_frontdoor_rule_set.test.id

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
