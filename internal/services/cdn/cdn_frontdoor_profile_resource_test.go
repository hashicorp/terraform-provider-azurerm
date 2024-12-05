// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
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

func TestAccCdnFrontDoorProfileWithSystemIdentity_update(t *testing.T) {
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
	})
}

func TestAccCdnFrontDoorProfileWithUserIdentity_update(t *testing.T) {
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
	})
}

func TestAccCdnFrontDoorProfileWithSystemAndUserIdentity_update(t *testing.T) {
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

func (r CdnFrontDoorProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	profileId := profiles.ProfileId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroup,
		ProfileName:       id.ProfileName,
	}

	client := clients.Cdn.FrontDoorProfilesClient
	resp, err := client.Get(ctx, profileId)
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
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

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
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

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
  name                = "acctestprofile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium_AzureFrontDoor"

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

  name = "acctestAFD-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
