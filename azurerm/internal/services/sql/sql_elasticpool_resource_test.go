package sql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SqlElasticPoolResource struct{}

func TestAccSqlElasticPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")
	r := SqlElasticPoolResource{}

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

func TestAccSqlElasticPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")
	r := SqlElasticPoolResource{}

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

func TestAccSqlElasticPool_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")
	r := SqlElasticPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccSqlElasticPool_resizeDtu(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_elasticpool", "test")
	r := SqlElasticPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dtu").HasValue("50"),
				check.That(data.ResourceName).Key("pool_size").HasValue("5000"),
			),
		},
		{
			Config: r.resizedDtu(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dtu").HasValue("100"),
				check.That(data.ResourceName).Key("pool_size").HasValue("10000"),
			),
		},
	})
}

func (r SqlElasticPoolResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ElasticPoolID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Sql.ElasticPoolsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Sql Elastic Pool %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlElasticPoolResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ElasticPoolID(state.ID)
	if err != nil {
		return nil, err
	}
	if _, err := client.Sql.ElasticPoolsClient.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name); err != nil {
		return nil, fmt.Errorf("deleting Sql Elastic Pool %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlElasticPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_elasticpool" "test" {
  name                = "acctest-pool-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  edition             = "Basic"
  dtu                 = 50
  pool_size           = 5000
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SqlElasticPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_elasticpool" "import" {
  name                = azurerm_sql_elasticpool.test.name
  resource_group_name = azurerm_sql_elasticpool.test.resource_group_name
  location            = azurerm_sql_elasticpool.test.location
  server_name         = azurerm_sql_elasticpool.test.server_name
  edition             = azurerm_sql_elasticpool.test.edition
  dtu                 = azurerm_sql_elasticpool.test.dtu
  pool_size           = azurerm_sql_elasticpool.test.pool_size
}
`, r.basic(data))
}

func (r SqlElasticPoolResource) resizedDtu(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_elasticpool" "test" {
  name                = "acctest-pool-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  edition             = "Basic"
  dtu                 = 100
  pool_size           = 10000
}
`, data.RandomInteger, data.Locations.Primary)
}
