// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorProfileResource struct{}

func TestAccCdnFrontDoorProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
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

func TestAccCdnFrontDoorProfileWithSystemIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfileWithUserIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfileWithSystemAndUserIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithSystemAndUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
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

func TestAccCdnFrontDoorProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
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

func TestAccCdnFrontDoorProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("240"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("120"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_withSystemIdentity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("240"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateWithSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("120"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("240"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_withUserIdentity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("240"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateWithUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("120"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("240"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_withSystemAndUserIdentity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("240"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateWithSystemAndUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("120"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("response_timeout_seconds").HasValue("240"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_skuDowngrade_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuPremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Test that downgrading from Premium to Standard is not allowed
			Config:      r.skuStandard(data),
			ExpectError: regexp.MustCompile("downgrading `sku_name` from `Premium_AzureFrontDoor` to `Standard_AzureFrontDoor` is not supported"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_standardSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingStandardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_premiumSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingPremiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test valid configuration with QueryStringArgNames and selector
			Config: r.logScrubbingValidQueryStringWithSelector(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Test valid configuration with RequestIPAddress without selector
			Config: r.logScrubbingValidRequestIPAddress(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Test valid configuration with RequestUri without selector
			Config: r.logScrubbingValidRequestUri(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingValidRequestIPAddress(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingMultipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingMultipleRulesDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_disabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_edgeCase_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test empty log_scrubbing block (should fail)
			Config:      r.logScrubbingEmptyBlock(data),
			ExpectError: regexp.MustCompile("empty `log_scrubbing` block: either remove it or specify one or more valid configuration fields"),
		},
		{
			// Test valid configuration with no scrubbing rules but enabled
			Config: r.logScrubbingEnabledNoRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_queryStringArgNames_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test operator being set for QueryStringArgNames (should fail)
			Config:      r.logScrubbingOperatorSet(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.0: `operator` cannot be set when the `match_variable` is `QueryStringArgNames`"),
		},
		{
			// Test empty selector for QueryStringArgNames (should fail)
			Config:      r.logScrubbingEmptySelector(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.0: `selector` is required when the `match_variable` is `QueryStringArgNames`"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_requestIPAddress_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test RequestIPAddress without an operator (should fail)
			Config:      r.logScrubbingOperatorNotSetRequestIPAddress(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.0: `operator` is required when the `match_variable` is `RequestIPAddress`"),
		},
		{
			// Test RequestIPAddress with a selector (should fail)
			Config:      r.logScrubbingSelectorSetRequestIPAddress(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.0: `selector` cannot be set when the `match_variable` is `RequestIPAddress`"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_requestUri_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test RequestUri without an operator (should fail)
			Config:      r.logScrubbingOperatorNotSetRequestUri(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.0: `operator` is required when the `match_variable` is `RequestUri`"),
		},
		{
			// Test RequestUri with a selector (should fail)
			Config:      r.logScrubbingSelectorSetRequestUri(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.0: `selector` cannot be set when the `match_variable` is `RequestUri`"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbingRuleIndex_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test validation error message includes correct rule index (rule 1)
			Config:      r.logScrubbingInvalidRuleAtIndex1(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.1: `selector` cannot be set when the `match_variable` is `RequestIPAddress`"),
		},
		{
			// Test validation error message includes correct rule index (rule 2)
			Config:      r.logScrubbingInvalidRuleAtIndex2(data),
			ExpectError: regexp.MustCompile("log_scrubbing\\.0\\.scrubbing_rule\\.2: `selector` is required when the `match_variable` is `QueryStringArgNames`"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_multipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingMultipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_complexUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingEnabledNoRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingMultipleRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontDoorProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := profiles.ParseProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorProfilesClient
	resp, err := client.Get(ctx, pointer.From(id))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnFrontDoorProfileResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) basicWithSystemIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) basicWithUserIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) basicWithSystemAndUserIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_profile" "import" {
  name                = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
  sku_name            = azurerm_cdn_frontdoor_profile.test.sku_name
}
`, config)
}

func (r CdnFrontDoorProfileResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                     = "acctestprofile-%d"
  resource_group_name      = azurerm_resource_group.test.name
  response_timeout_seconds = 240
  sku_name                 = "Premium_AzureFrontDoor"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  response_timeout_seconds = 120

  tags = {
    ENV = "Production"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) updateWithSystemIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_cdn_frontdoor_profile" "test" {
  name                     = "acctestprofile-%d"
  resource_group_name      = azurerm_resource_group.test.name
  sku_name                 = "Premium_AzureFrontDoor"
  response_timeout_seconds = 120
  identity {
    type = "SystemAssigned"
  }
  tags = {
    ENV = "Production"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) updateWithUserIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_cdn_frontdoor_profile" "test" {
  name                     = "acctestprofile-%d"
  resource_group_name      = azurerm_resource_group.test.name
  sku_name                 = "Premium_AzureFrontDoor"
  response_timeout_seconds = 120
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  tags = {
    ENV = "Production"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) updateWithSystemAndUserIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_cdn_frontdoor_profile" "test" {
  name                     = "acctestprofile-%d"
  resource_group_name      = azurerm_resource_group.test.name
  sku_name                 = "Premium_AzureFrontDoor"
  response_timeout_seconds = 120
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  tags = {
    ENV = "Production"
  }
}
`, template, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestAFD-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingStandardSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = true
    }
  }

  tags = {
    environment = "Test"
    sku         = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingPremiumSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "search_query"
      enabled        = true
    }
  }

  tags = {
    environment = "Test"
    sku         = "Premium"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingValidQueryStringWithSelector(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "custom_param" # This is valid for QueryStringArgNames
      enabled        = true
    }
  }

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingValidRequestIPAddress(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = true
      # No selector specified - this is required for RequestIPAddress
    }
  }

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingValidRequestUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      enabled        = true
      # No selector specified - this is required for RequestUri
    }
  }

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingMultipleRulesDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = false
    }
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      enabled        = false
    }
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "user_id"
      enabled        = true
    }
  }

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = false
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = true
    }
  }

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) skuPremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) skuStandard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingEmptyBlock(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {}
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingEmptySelector(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "" # Empty selector should trigger validation error
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingOperatorSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "param_1"
      operator       = "EqualsAny" # Operator set should trigger validation error
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingOperatorNotSetRequestIPAddress(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true

    # Operator not set - This should trigger a validation error
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingSelectorSetRequestIPAddress(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      selector       = "param_1" # Selector Set - This should trigger validation error
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingOperatorNotSetRequestUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true

    # Operator not set - This should trigger a validation error
    scrubbing_rule {
      match_variable = "RequestUri"
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingSelectorSetRequestUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      selector       = "param_1" # Selector Set - This should trigger validation error
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingEnabledNoRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    # No scrubbing rules defined
  }

  tags = {
    environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingMultipleRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "user_id"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "session_token"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      selector       = "api_key"
      enabled        = false
    }
  }

  tags = {
    environment = "Test"
    rule_count  = "5"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingInvalidRuleAtIndex1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      selector       = "invalid_selector" # This should trigger validation error at index 1
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) logScrubbingInvalidRuleAtIndex2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = true
    }
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      # Missing selector - should trigger validation error at index 2
      enabled = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
