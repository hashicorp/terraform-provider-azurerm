// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorProfileDataSource struct{}

func TestAccCdnFrontDoorProfileDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_profile", "test")
	d := CdnFrontDoorProfileDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_AzureFrontDoor"),
			),
		},
	})
}

func TestAccCdnFrontDoorProfileDataSource_basicWithSystemIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_profile", "test")
	d := CdnFrontDoorProfileDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basicWithSystemIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard_AzureFrontDoor"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
	})
}

func TestAccCdnFrontDoorProfileDataSource_basicWithUserIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_profile", "test")
	d := CdnFrontDoorProfileDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basicWithUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard_AzureFrontDoor"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
	})
}

func TestAccCdnFrontDoorProfileDataSource_basicWithSystemAndUserIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_profile", "test")
	d := CdnFrontDoorProfileDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basicWithSystemAndUserIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard_AzureFrontDoor"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
			),
		},
	})
}

func (CdnFrontDoorProfileDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_profile" "test" {
  name                = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontDoorProfileResource{}.complete(data))
}

func (CdnFrontDoorProfileDataSource) basicWithSystemIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_profile" "test" {
  name                = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontDoorProfileResource{}.basicWithSystemIdentity(data))
}

func (CdnFrontDoorProfileDataSource) basicWithUserIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_profile" "test" {
  name                = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontDoorProfileResource{}.basicWithUserIdentity(data))
}

func (CdnFrontDoorProfileDataSource) basicWithSystemAndUserIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_profile" "test" {
  name                = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontDoorProfileResource{}.basicWithSystemAndUserIdentity(data))
}
