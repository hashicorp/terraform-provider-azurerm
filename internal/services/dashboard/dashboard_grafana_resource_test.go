// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dashboard_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/grafanaresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DashboardGrafanaResource struct{}

func TestAccDashboardGrafana_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana", "test")
	r := DashboardGrafanaResource{}
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

func TestAccDashboardGrafana_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana", "test")
	r := DashboardGrafanaResource{}
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

func TestAccDashboardGrafana_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana", "test")
	r := DashboardGrafanaResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("smtp.0.password"),
	})
}

func TestAccDashboardGrafana_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana", "test")
	r := DashboardGrafanaResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("smtp.0.password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("smtp.0.password"),
	})
}

func TestAccDashboardGrafana_withSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dashboard_grafana", "test")
	r := DashboardGrafanaResource{}
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

func (r DashboardGrafanaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r DashboardGrafanaResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-mw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DashboardGrafanaResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_dashboard_grafana" "test" {
  name                  = "a-dg-%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  grafana_major_version = "11"
}
`, template, data.RandomInteger)
}

func (r DashboardGrafanaResource) essential(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_dashboard_grafana" "test" {
  name                  = "a-dg-%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  grafana_major_version = "10"

  sku = "Essential"
}
`, template, data.RandomInteger)
}

func (r DashboardGrafanaResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_dashboard_grafana" "import" {
  name                  = azurerm_dashboard_grafana.test.name
  resource_group_name   = azurerm_dashboard_grafana.test.resource_group_name
  location              = azurerm_dashboard_grafana.test.location
  grafana_major_version = "11"
}
`, config)
}

func (r DashboardGrafanaResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%[1]s
resource "azurerm_user_assigned_identity" "test" {
  name                = "a-uid-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_dashboard_grafana" "test" {
  name                              = "a-dg-%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  api_key_enabled                   = true
  deterministic_outbound_ip_enabled = true
  public_network_access_enabled     = false
  grafana_major_version             = "10"
  smtp {
    enabled          = true
    host             = "localhost:25"
    user             = "user"
    password         = "password"
    from_address     = "admin@grafana.localhost"
    from_name        = "Grafana"
    start_tls_policy = "OpportunisticStartTLS"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  azure_monitor_workspace_integrations {
    resource_id = azurerm_monitor_workspace.test.id
  }

  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger)
}

func (r DashboardGrafanaResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_monitor_workspace" "test2" {
  name                = "acctest-mw2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_dashboard_grafana" "test" {
  name                  = "a-dg-%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  grafana_major_version = "11"

  identity {
    type = "SystemAssigned"
  }

  azure_monitor_workspace_integrations {
    resource_id = azurerm_monitor_workspace.test.id
  }

  azure_monitor_workspace_integrations {
    resource_id = azurerm_monitor_workspace.test2.id
  }

  smtp {
    enabled          = true
    host             = "localhost:26"
    user             = "user"
    password         = "password"
    from_address     = "admin@grafana.localhost"
    from_name        = "Grafana"
    start_tls_policy = "OpportunisticStartTLS"
  }

  tags = {
    key2 = "value2"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}
