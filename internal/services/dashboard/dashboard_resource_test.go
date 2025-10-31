// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dashboard_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/manageddashboards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DashboardResource struct{}

func TestAccDashboard_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard", "test")
	r := DashboardResource{}

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

func TestAccDashboard_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard", "test")
	r := DashboardResource{}

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

func TestAccDashboard_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard", "test")
	r := DashboardResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDashboard_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard", "test")
	r := DashboardResource{}

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

func (r DashboardResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := manageddashboards.ParseDashboardID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Dashboard.ManagedDashboardsClient.DashboardsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r DashboardResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dashboard-%d"
  location = "%s"
}

resource "azurerm_dashboard" "test" {
  name                = "acctest-db-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DashboardResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dashboard" "import" {
  name                = azurerm_dashboard.test.name
  resource_group_name = azurerm_dashboard.test.resource_group_name
  location            = azurerm_dashboard.test.location
}
`, r.basic(data))
}

func (r DashboardResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dashboard-%d"
  location = "%s"
}

resource "azurerm_dashboard" "test" {
  name                = "acctest-db-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "Production"
    cost_center = "IT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DashboardResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dashboard-%d"
  location = "%s"
}

resource "azurerm_dashboard" "test" {
  name                = "acctest-db-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "Production"
    cost_center = "IT"
    owner       = "DevOps Team"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
