package sql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SqlFailoverGroupResource struct{}

func TestAccSqlFailoverGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_failover_group", "test")
	r := SqlFailoverGroupResource{}

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

func TestAccSqlFailoverGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_failover_group", "test")
	r := SqlFailoverGroupResource{}

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

func TestAccSqlFailoverGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_failover_group", "test")
	r := SqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccSqlFailoverGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_failover_group", "test")
	r := SqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (r SqlFailoverGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FailoverGroupID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Sql.FailoverGroupsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Sql Failover Group %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlFailoverGroupResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FailoverGroupID(state.ID)
	if err != nil {
		return nil, err
	}
	if _, err := client.Sql.FailoverGroupsClient.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name); err != nil {
		return nil, fmt.Errorf("deleting Sql Failover Group %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlFailoverGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test_primary" {
  name                         = "acctestmssql%[1]d-primary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_server" "test_secondary" {
  name                         = "acctestmssql%[1]d-secondary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = "%[3]s"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test_primary.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_sql_failover_group" "test" {
  name                = "acctestsfg%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test_primary.name
  databases           = [azurerm_sql_database.test.id]

  partner_servers {
    id = azurerm_sql_server.test_secondary.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SqlFailoverGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_failover_group" "import" {
  name                = azurerm_sql_failover_group.test.name
  resource_group_name = azurerm_sql_failover_group.test.resource_group_name
  server_name         = azurerm_sql_failover_group.test.server_name
  databases           = azurerm_sql_failover_group.test.databases

  partner_servers {
    id = azurerm_sql_failover_group.test.partner_servers[0].id
  }

  read_write_endpoint_failover_policy {
    mode          = azurerm_sql_failover_group.test.read_write_endpoint_failover_policy[0].mode
    grace_minutes = azurerm_sql_failover_group.test.read_write_endpoint_failover_policy[0].grace_minutes
  }
}
`, r.basic(data))
}

func (r SqlFailoverGroupResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test_primary" {
  name                         = "acctestmssql%[1]d-primary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_server" "test_secondary" {
  name                         = "acctestmssql%[1]d-secondary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = "%[3]s"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test_primary.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_sql_failover_group" "test" {
  name                = "acctestsfg%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test_primary.name
  databases           = [azurerm_sql_database.test.id]

  partner_servers {
    id = azurerm_sql_server.test_secondary.id
  }
  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SqlFailoverGroupResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test_primary" {
  name                         = "acctestmssql%[1]d-primary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_server" "test_secondary" {
  name                         = "acctestmssql%[1]d-secondary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = "%[3]s"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test_primary.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_sql_failover_group" "test" {
  name                = "acctestsfg%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test_primary.name
  databases           = [azurerm_sql_database.test.id]

  partner_servers {
    id = azurerm_sql_server.test_secondary.id
  }
  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
