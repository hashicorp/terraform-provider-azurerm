// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package portal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LegacyDashboardResource struct{}

func TestAccLegacyDashboard_basic(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("Test no longer valid due to deprecation of the 'azurerm_dashboard' resource in the 4.x version of the provider")
	}

	data := acceptance.BuildTestData(t, "azurerm_dashboard", "test")
	r := LegacyDashboardResource{}
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

func TestAccLegacyDashboard_complete(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("Test no longer valid due to deprecation of the 'azurerm_dashboard' resource in the 4.x version of the provider")
	}

	data := acceptance.BuildTestData(t, "azurerm_dashboard", "test")
	r := LegacyDashboardResource{}
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

func (LegacyDashboardResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dashboard.ParseDashboardID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Portal.DashboardsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (LegacyDashboardResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dashboard" "test" {
  name                 = "my-test-dashboard"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  dashboard_properties = <<DASH
{
   "lenses": {
        "0": {
            "order": 0,
            "parts": {
                "0": {
                    "position": {
                        "x": 0,
                        "y": 0,
                        "rowSpan": 2,
                        "colSpan": 3
                    },
                    "metadata": {
                        "inputs": [],
                        "type": "Extension/HubsExtension/PartType/MarkdownPart",
                        "settings": {
                            "content": {
                                "settings": {
                                    "content": "## This is only a test :)",
                                    "subtitle": "",
                                    "title": "Test MD Tile"
                                }
                            }
                        }
                    }
				}
			}
		}
	}
}
DASH
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LegacyDashboardResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dashboard" "test" {
  name                = "my-test-dashboard"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV = "Test"
  }
  dashboard_properties = <<DASH
{
   "lenses": {
        "0": {
            "order": 0,
            "parts": {
                "0": {
                    "position": {
                        "x": 0,
                        "y": 0,
                        "rowSpan": 2,
                        "colSpan": 3
                    },
                    "metadata": {
                        "inputs": [],
                        "type": "Extension/HubsExtension/PartType/MarkdownPart",
                        "settings": {
                            "content": {
                                "settings": {
                                    "content": "## This is only a test :)",
                                    "subtitle": "",
                                    "title": "Test MD Tile"
                                }
                            }
                        }
                    }
				}
			}
		}
	}
}
DASH
}
`, data.RandomInteger, data.Locations.Primary)
}
