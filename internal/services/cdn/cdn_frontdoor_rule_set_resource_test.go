// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorRuleSetResource struct{}

func TestAccCdnFrontDoorRuleSet_basic_unattachedRoute(t *testing.T) {
	t.Skip(unattachedFrontDoorRuleSetRegressionSkipMessage)

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRuleSet_basic_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRuleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnFrontDoorRuleSet_complete_unattachedRoute(t *testing.T) {
	t.Skip(unattachedFrontDoorRuleSetRegressionSkipMessage)

	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorRuleSet_complete_attachedRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontDoorRuleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Cdn.FrontDoorRuleSetsClient

	id, err := rulesets.ParseRuleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err = client.Get(ctx, *id); err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (r CdnFrontDoorRuleSetResource) basic(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule_set" "test" {
  name                     = "acctestfdruleset%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, template, data.RandomIntOfLength(8))
}

func (r CdnFrontDoorRuleSetResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data, true)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set" "import" {
  name                     = azurerm_cdn_frontdoor_rule_set.test.name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, config)
}

func (r CdnFrontDoorRuleSetResource) complete(data acceptance.TestData, attachRoute bool) string {
	template := r.templateWithAttachedRoute(data, attachRoute)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_rule_set" "test" {
  name                     = "acctestfdruleset%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, template, data.RandomIntOfLength(8))
}

func (r CdnFrontDoorRuleSetResource) templateWithAttachedRoute(data acceptance.TestData, attachRoute bool) string {
	template := r.template(data)
	if !attachRoute {
		return template
	}

	return fmt.Sprintf(`%s

%s`, template, r.routeTemplate(data))
}

func (CdnFrontDoorRuleSetResource) routeTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "acctest-fdog-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-fdorigin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctest-fdendpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "acctest-fdroute-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  cdn_frontdoor_rule_set_ids    = [azurerm_cdn_frontdoor_rule_set.test.id]
  patterns_to_match             = ["/*"]
  supported_protocols           = ["Http", "Https"]
}
`, data.RandomInteger)
}

func (CdnFrontDoorRuleSetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-cdn-afdx-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-fdprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
