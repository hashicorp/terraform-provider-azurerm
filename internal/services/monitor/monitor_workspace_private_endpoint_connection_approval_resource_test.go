// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkspacePrivateEndpointConnectionTestResource struct{}

func TestAccMonitorWorkspacePrivateEndpointConnectionApproval_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace_private_endpoint_connection_approval", "test")
	r := WorkspacePrivateEndpointConnectionTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").IsNotEmpty(),
				check.That(data.ResourceName).Key("workspace_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("private_endpoint_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("approval_message").HasValue("Approved via Terraform"),
			),
		},
		data.ImportStep("private_endpoint_connection_name"),
	})
}

func TestAccMonitorWorkspacePrivateEndpointConnectionApproval_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace_private_endpoint_connection_approval", "test")
	r := WorkspacePrivateEndpointConnectionTestResource{}
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

func TestAccMonitorWorkspacePrivateEndpointConnectionApproval_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace_private_endpoint_connection_approval", "test")
	r := WorkspacePrivateEndpointConnectionTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("private_endpoint_connection_name"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("approval_message").HasValue("Updated approval message"),
			),
		},
		data.ImportStep("private_endpoint_connection_name"),
	})
}

func TestAccMonitorWorkspacePrivateEndpointConnectionApproval_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_workspace_private_endpoint_connection_approval", "test")
	r := WorkspacePrivateEndpointConnectionTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").IsNotEmpty(),
				check.That(data.ResourceName).Key("workspace_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("private_endpoint_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("approval_message").HasValue("Approved by automation"),
			),
		},
		data.ImportStep("private_endpoint_connection_name"),
	})
}

func (r WorkspacePrivateEndpointConnectionTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azuremonitorworkspaces.ParseAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Monitor.WorkspacesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return utils.Bool(false), nil
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.PrivateEndpointConnections == nil {
		return utils.Bool(false), nil
	}

	// Check if there's an approved connection
	for _, conn := range *resp.Model.Properties.PrivateEndpointConnections {
		if conn.Properties != nil &&
			conn.Properties.PrivateLinkServiceConnectionState.Status != nil &&
			*conn.Properties.PrivateLinkServiceConnectionState.Status == azuremonitorworkspaces.PrivateEndpointServiceConnectionStatusApproved {
			return utils.Bool(true), nil
		}
	}

	return utils.Bool(false), nil
}

func (r WorkspacePrivateEndpointConnectionTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-mamw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_dashboard_grafana" "test" {
  name                  = "a-dg-%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  grafana_major_version = "10"
}

resource "azurerm_dashboard_grafana_managed_private_endpoint" "test" {
  grafana_id                   = azurerm_dashboard_grafana.test.id
  name                         = "acctest-mpe-%d"
  location                     = azurerm_dashboard_grafana.test.location
  private_link_resource_id     = azurerm_monitor_workspace.test.id
  group_ids                    = ["prometheusMetrics"]
  private_link_resource_region = azurerm_dashboard_grafana.test.location
}

data "azurerm_monitor_workspace" "test" {
  name                = azurerm_monitor_workspace.test.name
  resource_group_name = azurerm_monitor_workspace.test.resource_group_name

  depends_on = [azurerm_dashboard_grafana_managed_private_endpoint.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomIntOfLength(8))
}

func (r WorkspacePrivateEndpointConnectionTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_workspace_private_endpoint_connection_approval" "test" {
  workspace_id                     = azurerm_monitor_workspace.test.id
  private_endpoint_connection_name = data.azurerm_monitor_workspace.test.private_endpoint_connections[0].name
}
`, template)
}

func (r WorkspacePrivateEndpointConnectionTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_workspace_private_endpoint_connection_approval" "import" {
  workspace_id                     = azurerm_monitor_workspace_private_endpoint_connection_approval.test.workspace_id
  private_endpoint_connection_name = azurerm_monitor_workspace_private_endpoint_connection_approval.test.private_endpoint_connection_name
}
`, config)
}

func (r WorkspacePrivateEndpointConnectionTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_workspace_private_endpoint_connection_approval" "test" {
  workspace_id                     = azurerm_monitor_workspace.test.id
  private_endpoint_connection_name = data.azurerm_monitor_workspace.test.private_endpoint_connections[0].name
  approval_message                 = "Updated approval message"
}
`, template)
}

func (r WorkspacePrivateEndpointConnectionTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_workspace_private_endpoint_connection_approval" "test" {
  workspace_id                     = azurerm_monitor_workspace.test.id
  private_endpoint_connection_name = data.azurerm_monitor_workspace.test.private_endpoint_connections[0].name
  approval_message                 = "Approved by automation"
}
`, template)
}
