package serviceconnector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/sdk/2022-05-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceConnectorAppServiceResource struct{}

func (r ServiceConnectorAppServiceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servicelinker.ParseScopedLinkerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ServiceConnector.ServiceLinkerClient.LinkerGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccServiceConnectorAppServiceCosmosdb_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorAppServiceCosmosdb_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.cosmosdbUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorAppService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_connection", "test")
	r := ServiceConnectorAppServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ServiceConnectorAppServiceResource) cosmosdbBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  auth_info {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) cosmosdbUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "update" {
  name                = "cosmos-sql-db-update"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "update" {
  name                = "test-containerupdate%[2]s"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.update.name
  partition_key_path  = "/definitionupdate"
}

resource "azurerm_service_plan" "update" {
  location            = azurerm_resource_group.test.location
  name                = "testserviceplanupdate%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "P1v3"
  os_type             = "Linux"
}

resource "azurerm_linux_web_app" "update" {
  location            = azurerm_resource_group.test.location
  name                = "linuxwebappupdate%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.update.id
  site_config {}
}

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  app_service_id     = azurerm_linux_web_app.update.id
  target_resource_id = azurerm_cosmosdb_sql_database.update.id
  auth_info {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_app_service_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  app_service_id     = azurerm_linux_web_app.test.id
  target_resource_id = azurerm_cosmosdb_sql_database.test.id
  client_type        = "java"
  vnet_solution      = "privateLink"
  auth_info {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomInteger)
}

func (r ServiceConnectorAppServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestacc%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "test-container%[3]s"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_path  = "/definition"
}

resource "azurerm_service_plan" "test" {
  location            = azurerm_resource_group.test.location
  name                = "testserviceplan%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "P1v2"
  os_type             = "Linux"
}

resource "azurerm_linux_web_app" "test" {
  location            = azurerm_resource_group.test.location
  name                = "linuxwebapp%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id
  site_config {}
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
