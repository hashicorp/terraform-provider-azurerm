// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mixedreality_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mixedreality/2021-01-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpatialAnchorsAccountResource struct{}

func TestAccSpatialAnchorsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spatial_anchors_account", "test")
	r := SpatialAnchorsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("account_id").Exists(),
				check.That(data.ResourceName).Key("account_domain").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpatialAnchorsAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spatial_anchors_account", "test")
	r := SpatialAnchorsAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
	})
}

func (SpatialAnchorsAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := resource.ParseSpatialAnchorsAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MixedReality.SpatialAnchorsAccountClient.SpatialAnchorsAccountsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (SpatialAnchorsAccountResource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (SpatialAnchorsAccountResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mr-%d"
  location = "%s"
}

resource "azurerm_spatial_anchors_account" "test" {
  name                = "acCTestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    Environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
