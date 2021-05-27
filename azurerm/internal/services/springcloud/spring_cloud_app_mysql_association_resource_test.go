package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SpringCloudAppMysqlAssociationResource struct {
}

func TestAccSpringCloudAppMysqlAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_mysql_association", "test")
	r := SpringCloudAppMysqlAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccSpringCloudAppMysqlAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_mysql_association", "test")
	r := SpringCloudAppMysqlAssociationResource{}

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

func TestAccSpringCloudAppMysqlAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_mysql_association", "test")
	r := SpringCloudAppMysqlAssociationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func (r SpringCloudAppMysqlAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudAppAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.BindingsClient.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r SpringCloudAppMysqlAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_mysql_association" "test" {
  name                = "acctestscamb-%d"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  mysql_server_id     = azurerm_mysql_server.test.id
  database_name       = azurerm_mysql_database.test.name
  username            = azurerm_mysql_server.test.administrator_login
  password            = azurerm_mysql_server.test.administrator_login_password
}
`, r.template(data), data.RandomInteger)
}

func (r SpringCloudAppMysqlAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_mysql_association" "import" {
  name                = azurerm_spring_cloud_app_mysql_association.test.name
  spring_cloud_app_id = azurerm_spring_cloud_app_mysql_association.test.spring_cloud_app_id
  mysql_server_id     = azurerm_spring_cloud_app_mysql_association.test.mysql_server_id
  database_name       = azurerm_spring_cloud_app_mysql_association.test.database_name
  username            = azurerm_spring_cloud_app_mysql_association.test.username
  password            = azurerm_spring_cloud_app_mysql_association.test.password
}
`, r.basic(data))
}

func (r SpringCloudAppMysqlAssociationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_database" "updated" {
  name                = "acctest-db2-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_server.test.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}

resource "azurerm_spring_cloud_app_mysql_association" "test" {
  name                = "acctestscamb-%d"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  mysql_server_id     = azurerm_mysql_server.test.id
  database_name       = azurerm_mysql_database.updated.name
  username            = azurerm_mysql_server.test.administrator_login
  password            = azurerm_mysql_server.test.administrator_login_password
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r SpringCloudAppMysqlAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
  version                          = "5.7"
}

resource "azurerm_mysql_database" "test" {
  name                = "acctest-db-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mysql_server.test.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
