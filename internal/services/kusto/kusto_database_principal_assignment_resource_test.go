// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databaseprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoDatabasePrincipalAssignmentResource struct{}

func TestAccKustoDatabasePrincipalAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database_principal_assignment", "test")
	r := KustoDatabasePrincipalAssignmentResource{}

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

func TestAccKustoDatabasePrincipalAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database_principal_assignment", "test")
	r := KustoDatabasePrincipalAssignmentResource{}

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

func (KustoDatabasePrincipalAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := databaseprincipalassignments.ParseDatabasePrincipalAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.DatabasePrincipalAssignmentsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	exists := resp.Model != nil

	return &exists, nil
}

func (KustoDatabasePrincipalAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-kusto-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_database_principal_assignment" "test" {
  name                = "acctestkdpa%d"
  resource_group_name = azurerm_resource_group.rg.name
  cluster_name        = azurerm_kusto_cluster.test.name
  database_name       = azurerm_kusto_database.test.name

  tenant_id      = data.azurerm_client_config.current.tenant_id
  principal_id   = data.azurerm_client_config.current.client_id
  principal_type = "App"
  role           = "Viewer"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (KustoDatabasePrincipalAssignmentResource) requiresImport(data acceptance.TestData) string {
	template := KustoDatabasePrincipalAssignmentResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_database_principal_assignment" "import" {
  name                = azurerm_kusto_database_principal_assignment.test.name
  resource_group_name = azurerm_kusto_database_principal_assignment.test.resource_group_name
  cluster_name        = azurerm_kusto_database_principal_assignment.test.cluster_name
  database_name       = azurerm_kusto_database_principal_assignment.test.database_name

  tenant_id      = azurerm_kusto_database_principal_assignment.test.tenant_id
  principal_id   = azurerm_kusto_database_principal_assignment.test.principal_id
  principal_type = azurerm_kusto_database_principal_assignment.test.principal_type
  role           = azurerm_kusto_database_principal_assignment.test.role
}
`, template)
}
