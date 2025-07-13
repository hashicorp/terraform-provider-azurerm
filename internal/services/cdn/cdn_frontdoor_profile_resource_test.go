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

func TestAccCdnFrontDoorProfile_standardSku_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicStandardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_premiumSku_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPremiumSku(data),
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
			Config: r.basicStandardSku(data),
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
			Config: r.basicPremiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.basicStandardSku(data),
			ExpectError: regexp.MustCompile("downgrading `sku_name` from `Premium_AzureFrontDoor` to `Standard_AzureFrontDoor` is not supported"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_maxRules_standardSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingStandardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.match_variable").HasValue("QueryStringArgNames"),
				check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.1.match_variable").HasValue("RequestIPAddress"),
				check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.2.match_variable").HasValue("RequestUri"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_maxRules_premiumSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logScrubbingPremiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.match_variable").HasValue("QueryStringArgNames"),
				check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.1.match_variable").HasValue("RequestIPAddress"),
				check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.2.match_variable").HasValue("RequestUri"),
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
			Config: r.basicPremiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_scrubbing.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.logScrubbingPremiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.match_variable").HasValue("QueryStringArgNames"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicPremiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_scrubbing.#").HasValue("0"),
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
			Config: r.basicPremiumSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_scrubbing.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_emptyBlock(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logScrubbingEmpty(data),
			ExpectError: regexp.MustCompile("when the `log_scrubbing` block is defined, at least one `scrubbing_rule` must be specified"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_duplicateScrubbingRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.logScrubbingDuplicateScrubbingRules(data),
			ExpectError: regexp.MustCompile("duplicate `QueryStringArgNames` rule found in `log_scrubbing.0.scrubbing_rule.2.match_variable`"),
		},
	})
}

func TestAccCdnFrontDoorProfile_logScrubbing_ruleLimit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This test works because the schema validation of `MaxItems: 3` executes before the provider validation
			// since that validation is implemented within the terraform core runtime itself
			Config:      r.logScrubbingRuleLimit(data),
			ExpectError: regexp.MustCompile(`No more than 3 "scrubbing_rule" blocks are allowed`),
		},
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

func (r CdnFrontDoorProfileResource) template(data acceptance.TestData) string {
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

func (r CdnFrontDoorProfileResource) basicStandardSku(data acceptance.TestData) string {
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

  tags = {
    environment = "Production"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) basicPremiumSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  tags = {
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
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

  tags = {
    environment = "Production"
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

  tags = {
    environment = "Production"
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

  tags = {
    environment = "Production"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) requiresImport(data acceptance.TestData) string {
	config := r.basicStandardSku(data)
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
    environment = "Production"
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
    environment = "Production"
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
    environment = "Production"
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
    environment = "Production"
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) updateWithSystemAndUserIdentity(data acceptance.TestData) string {
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
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingPremiumSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }

	scrubbing_rule {
      match_variable = "RequestIPAddress"
    }

	scrubbing_rule {
      match_variable = "RequestUri"
    }
  }

  tags = {
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingStandardSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"

  log_scrubbing {
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }

	scrubbing_rule {
      match_variable = "RequestIPAddress"
    }

	scrubbing_rule {
      match_variable = "RequestUri"
    }
  }

  tags = {
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingEmpty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

  log_scrubbing {}

  tags = {
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingDuplicateScrubbingRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"

  log_scrubbing {
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }

	scrubbing_rule {
      match_variable = "RequestIPAddress"
    }

	scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }
  }

  tags = {
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingRuleLimit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"

  log_scrubbing {
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }

	scrubbing_rule {
      match_variable = "RequestIPAddress"
    }

	scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }

	scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }
  }

  tags = {
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}
