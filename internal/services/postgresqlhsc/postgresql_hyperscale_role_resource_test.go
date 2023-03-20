package postgresqlhsc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLHyperScaleRoleResource struct{}

func TestPostgreSQLHyperScaleRole_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_role", "test")
	r := PostgreSQLHyperScaleRoleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "H@Sh1CoR3!"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestPostgreSQLHyperScaleRole_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_role", "test")
	r := PostgreSQLHyperScaleRoleResource{}
	password := "H@Sh1CoR3!"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, password),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, password)
		}),
	})
}

func TestPostgreSQLHyperScaleRole_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_hyperscale_role", "test")
	r := PostgreSQLHyperScaleRoleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "H@Sh1CoR3!"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "H@Sh1CoR4!"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PostgreSQLHyperScaleRoleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := roles.ParseRoleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PostgreSQLHSC.RolesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r PostgreSQLHyperScaleRoleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresqlhsc-%d"
  location = "%s"
}

resource "azurerm_postgresql_hyperscale_server_group" "test" {
  name                = "acctest-postgresqlhscsg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PostgreSQLHyperScaleRoleResource) basic(data acceptance.TestData, password string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_role" "test" {
  name            = "acctest-postgresqlhsc-%d"
  server_group_id = azurerm_postgre_sql_hsc_cluster.test.id
  password        = "%s"
}
`, r.template(data), data.RandomInteger, password)
}

func (r PostgreSQLHyperScaleRoleResource) requiresImport(data acceptance.TestData, password string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_hyperscale_role" "import" {
  name            = azurerm_postgresql_hyperscale_role.test.name
  server_group_id = azurerm_postgresql_hyperscale_role.test.server_group_id
  password        = azurerm_postgresql_hyperscale_role.test.password
}
`, r.basic(data, password))
}
