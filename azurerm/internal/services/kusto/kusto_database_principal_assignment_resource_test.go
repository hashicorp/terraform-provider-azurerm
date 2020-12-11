package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KustoDatabasePrincipalAssignmentResource struct {
}

func TestAccKustoDatabasePrincipalAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database_principal_assignment", "test")
	r := KustoDatabasePrincipalAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoDatabasePrincipalAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database_principal_assignment", "test")
	r := KustoDatabasePrincipalAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (KustoDatabasePrincipalAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DatabasePrincipalAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.DatabasePrincipalAssignmentsClient.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.PrincipalAssignmentName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.DatabasePrincipalProperties != nil), nil
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
