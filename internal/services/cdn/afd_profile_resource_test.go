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

type AfdProfileResource struct{}

func TestAccCdnAfdProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := AfdProfileResource{}

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

func TestAccCdnAfdProfile_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := AfdProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnAfdProfile_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := AfdProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nonStandardCasing(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.nonStandardCasing(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	})
}

func TestAccCdnAfdProfile_standardToPremiumFrontDoor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := AfdProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardFrontDoor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.premiumFrontDoor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnAfdProfile_standardFrontDoor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := AfdProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardFrontDoor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku", "Standard_AzureFrontDoor"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnAfdProfile_premiumFrontDoor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := AfdProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premiumFrontDoor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "sku", "Premium_AzureFrontDoor"),
			),
		},
		data.ImportStep(),
	})
}

func (r AfdProfileResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cdn.ProfilesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Cdn Profile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r AfdProfileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AfdProfileResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_AzureFrontDoor"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AfdProfileResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_AzureFrontDoor"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AfdProfileResource) nonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "standard_azurefrontdoor"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AfdProfileResource) standardFrontDoor(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AfdProfileResource) premiumFrontDoor(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnprof%[1]d"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary)
}
