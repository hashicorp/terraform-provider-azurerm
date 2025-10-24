// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnProfileResource struct{}

func TestAccCdnProfile_basic(t *testing.T) {
	if cdn.IsCdnDeprecatedForCreation() {
		t.Skip(cdn.CreateDeprecationMessage)
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")
	r := CdnProfileResource{}

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

func TestAccCdnProfile_requiresImport(t *testing.T) {
	if cdn.IsCdnDeprecatedForCreation() {
		t.Skip(cdn.CreateDeprecationMessage)
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")
	r := CdnProfileResource{}

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

func TestAccCdnProfile_withTags(t *testing.T) {
	if cdn.IsCdnDeprecatedForCreation() {
		t.Skip(cdn.CreateDeprecationMessage)
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")
	r := CdnProfileResource{}

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

func TestAccCdnProfile_skuDeprecation(t *testing.T) {
	if cdn.IsCdnDeprecatedForCreation() {
		t.Skip(cdn.CreateDeprecationMessage)
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")
	r := CdnProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.standardAkamai(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(cdn.AkamaiDeprecationMessage),
		},
		{
			Config:      r.standardVerizon(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(cdn.VerizonDeprecationMessage),
		},
		{
			Config:      r.premiumVerizon(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(cdn.VerizonDeprecationMessage),
		},
	})
}

func TestAccCdnProfile_createShouldFail(t *testing.T) {
	if !cdn.IsCdnDeprecatedForCreation() {
		t.Skip("CDN is not deprecated for creation until October 1, 2025")
	}

	expectedError := cdn.CreateDeprecationMessage
	if cdn.IsCdnFullyRetired() {
		expectedError = cdn.FullyRetiredMessage
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")
	r := CdnProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basic(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(expectedError),
		},
	})
}

func TestAccCdnProfile_locationNormalization(t *testing.T) {
	if !cdn.IsCdnDeprecatedForCreation() {
		t.Skip("This test validates the fix for issue #30768, which only affects behavior after CDN deprecation date (October 1, 2025)")
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")
	r := CdnProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This test verifies the fix for issue #30768
			// When location is stored in state as normalized (e.g., "eastus") but config uses
			// human-readable format (e.g., "East US"), the plan should not trigger the
			// deprecation error because there's no actual change requiring recreation
			Config:   r.locationHumanReadable(data),
			PlanOnly: true,
			// ExpectNonEmptyPlan is false by default - we expect no changes
			// If the bug exists, this would fail with the deprecation error
		},
	})
}

func (r CdnProfileResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r CdnProfileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnProfileResource) standardAkamai(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Akamai"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnProfileResource) standardVerizon(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnProfileResource) premiumVerizon(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium_Verizon"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnProfileResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_profile" "import" {
  name                = azurerm_cdn_profile.test.name
  location            = azurerm_cdn_profile.test.location
  resource_group_name = azurerm_cdn_profile.test.resource_group_name
  sku                 = azurerm_cdn_profile.test.sku
}
`, template)
}

func (r CdnProfileResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnProfileResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnProfileResource) locationHumanReadable(data acceptance.TestData) string {
	// This config uses human-readable location format to test that normalization works correctly
	// The location will be stored in state as normalized (e.g., "eastus") but config uses
	// the human-readable format (e.g., "East US") to ensure the deprecation check doesn't
	// trigger false positives
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}
