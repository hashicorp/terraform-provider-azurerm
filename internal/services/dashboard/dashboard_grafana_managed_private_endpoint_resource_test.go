// // Copyright (c) HashiCorp, Inc.
// // SPDX-License-Identifier: MPL-2.0

package dashboard_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/managedprivateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedPrivateEndpointResource struct{}

func TestAccDashboardGrafanaManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDashboardGrafanaManagedPrivateEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDashboardGrafanaManagedPrivateEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDashboardGrafanaManagedPrivateEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDashboardGrafanaManagedPrivateEndpoint_withSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana_managed_private_endpoint", "test")
	r := ManagedPrivateEndpointResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.essential(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagedPrivateEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := grafanaresource.ParseGrafanaID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Dashboard.GrafanaResourceClient
	resp, err := client.GrafanaGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagedPrivateEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_dashboard_grafana" "test" {
  name                  = "a-dg-%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  grafana_major_version = "10"
}

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-mw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_dashboard_grafana_managed_private_endpoint" "test" {
  managed_grafana_id           = azurerm_dashboard_grafana.test.id
  name                         = "acctest-mpe-%d"
  location                     = azurerm_dashboard_grafana.test.location
  private_link_resource_id     = azurerm_monitor_workspace.test.id
  group_ids                    = ["prometheusMetrics"]
  private_link_resource_region = azurerm_dashboard_grafana.test.location
}
`, template, data.RandomInteger)
}

func (r ManagedPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_dashboard_grafana_managed_private_endpoint" "import" {
  managed_grafana_id           = azurerm_dashboard_grafana_managed_private_endpoint.test.managed_grafana_id
  name                         = azurerm_dashboard_grafana_managed_private_endpoint.test.name
  location                     = azurerm_dashboard_grafana_managed_private_endpoint.test.location
  private_link_resource_id     = azurerm_dashboard_grafana_managed_private_endpoint.test.private_link_resource_id
}
`, config)
// may be needed in the config above
//   group_ids                    = ["prometheusMetrics"]
//  private_link_resource_region = azurerm_dashboard_grafana.test.location
}

func (r ManagedPrivateEndpointResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_dashboard_grafana_managed_private_endpoint" "test" {
  managed_grafana_id           = azurerm_dashboard_grafana.test.id
  name                         = "acctest-mpe-%d"
  location                     = azurerm_dashboard_grafana.test.location
  private_link_resource_id     = azurerm_monitor_workspace.test.id
  group_ids                    = ["prometheusMetrics"]
  private_link_resource_region = azurerm_dashboard_grafana.test.location

  tags = {
    key = "value"
  }
  
  request_message = "please approve"
}
`, template, data.RandomInteger)
}

func (r ManagedPrivateEndpointResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_dashboard_grafana_managed_private_endpoint" "test" {
  managed_grafana_id           = azurerm_dashboard_grafana.test.id
  name                         = "acctest-mpe-%d"
  location                     = azurerm_dashboard_grafana.test.location
  private_link_resource_id     = azurerm_monitor_workspace.test.id
  group_ids                    = ["prometheusMetrics"]
  private_link_resource_region = azurerm_dashboard_grafana.test.location

  tags = {
    key2 = "value2"
  }
}
`, template, data.RandomInteger)
}
