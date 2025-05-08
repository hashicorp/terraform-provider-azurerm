// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorRuleSetResource struct{}

func TestAccCdnFrontDoorRuleSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
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

func TestAccCdnFrontDoorRuleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
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

func TestAccCdnFrontDoorRuleSet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set", "test")
	r := CdnFrontDoorRuleSetResource{}
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

func (r CdnFrontDoorRuleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Cdn.FrontDoorRuleSetsClient

	id, err := rulesets.ParseRuleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	_, err = client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnFrontDoorRuleSetResource) basic(data acceptance.TestData) string {
	template := r.template(data)
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
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set" "import" {
  name                     = azurerm_cdn_frontdoor_rule_set.test.name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, config)
}

func (r CdnFrontDoorRuleSetResource) complete(data acceptance.TestData) string {
	template := r.template(data)
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
