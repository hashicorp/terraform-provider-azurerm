// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mixedreality_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type SpatialAnchorsAccountDataSource struct{}

func TestAccSpatialAnchorsAccountDataSource_basic(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("Skipping since `azurerm_spatial_anchors_account` is deprecated and will be removed in 5.0")
	}

	data := acceptance.BuildTestData(t, "data.azurerm_spatial_anchors_account", "test")
	r := SpatialAnchorsAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("account_id").Exists(),
				check.That(data.ResourceName).Key("account_domain").Exists(),
			),
		},
	})
}

func (SpatialAnchorsAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mr-%d"
  location = "%s"
}

resource "azurerm_spatial_anchors_account" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_spatial_anchors_account" "test" {
  name                = azurerm_spatial_anchors_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
