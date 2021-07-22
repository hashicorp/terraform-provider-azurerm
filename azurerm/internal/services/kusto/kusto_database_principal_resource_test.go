package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KustoDatabasePrincipalResource struct {
}

func TestAccKustoDatabasePrincipal_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database_principal", "test")
	r := KustoDatabasePrincipalResource{}

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

func (KustoDatabasePrincipalResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Kusto.DatabasesClient
	id, err := parse.DatabasePrincipalID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err = client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName); err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	databasePrincipals, err := client.ListPrincipals(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName)
	if err != nil {
		if !utils.ResponseWasNotFound(databasePrincipals.Response) {
			return nil, fmt.Errorf("retrieving principals for %s: %v", id.String(), err)
		}
	}

	if principals := databasePrincipals.Value; principals != nil {
		for _, currPrincipal := range *principals {
			// kusto database principals are unique when looked at with fqn and role
			if string(currPrincipal.Role) == id.RoleName && currPrincipal.Fqn != nil && *currPrincipal.Fqn == id.FQNName {
				return utils.Bool(true), nil
			}
		}
	}
	return utils.Bool(false), nil
}

func (KustoDatabasePrincipalResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}


resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-kusto-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
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
  cluster_name        = azurerm_kusto_cluster.cluster.name
}

resource "azurerm_kusto_database_principal" "test" {
  resource_group_name = azurerm_resource_group.rg.name
  cluster_name        = azurerm_kusto_cluster.cluster.name
  database_name       = azurerm_kusto_database.test.name

  role      = "Viewer"
  type      = "App"
  client_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
