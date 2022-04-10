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

type CdnFrontdoorRuleResource struct{}

func TestAccCdnFrontdoorRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontdoorRuleResource{}
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

func TestAccCdnFrontdoorRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontdoorRuleResource{}
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

func TestAccCdnFrontdoorRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontdoorRuleResource{}
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

func TestAccCdnFrontdoorRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule", "test")
	r := CdnFrontdoorRuleResource{}
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

func (r CdnFrontdoorRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontdoorRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontdoorRulesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnFrontdoorRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_rule_set" "test" {
  name                     = "acctest-c-%d"
  frontdoor_cdn_profile_id = azurerm_cdn_frontdoor_profile_profile.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r CdnFrontdoorRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  name                      = "acctest-c-%d"
  frontdoor_cdn_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_rule" "import" {
  name                      = azurerm_cdn_frontdoor_rule.test.name
  frontdoor_cdn_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  order = 0
}
`, config)
}

func (r CdnFrontdoorRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  name                      = "acctest-c-%d"
  frontdoor_cdn_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  actions    = ["CacheExpiration", "UrlRedirect", "OriginGroupOverride"]
  conditions = ["HostName", "IsDevice", "PostArgs", "RequestMethod"]

  match_processing_behavior = "Continue"
  order                     = 1
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorRuleResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_cdn_frontdoor_rule" "test" {
  name                      = "acctest-c-%d"
  frontdoor_cdn_rule_set_id = azurerm_cdn_frontdoor_rule_set.test.id

  actions    = ["CacheExpiration", "UrlRedirect"]
  conditions = ["HostName", "IsDevice", "RequestMethod"]

  match_processing_behavior = "Stop"
  order                     = 2
}
`, template, data.RandomInteger)
}
