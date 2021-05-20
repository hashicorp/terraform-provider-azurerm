package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LinkedServiceAzureSQLDatabaseResource struct {
}

func TestAccDataFactoryLinkedServiceAzureSQLDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_sql_database", "test")
	r := LinkedServiceAzureSQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccDataFactoryLinkedServiceAzureSQLDatabase_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_sql_database", "test")
	r := LinkedServiceAzureSQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.update1(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("2"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("3"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("2"),
				check.That(data.ResourceName).Key("description").HasValue("test description"),
			),
		},
		data.ImportStep("connection_string"),
		{
			Config: r.update2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("3"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("2"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("description").HasValue("test description 2"),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccDataFactoryLinkedServiceAzureSQLDatabase_managed_id(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_sql_database", "test")
	r := LinkedServiceAzureSQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.managed_id(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccDataFactoryLinkedServiceAzureSQLDatabase_KeyVaultReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_sql_database", "test")
	r := LinkedServiceAzureSQLDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.key_vault_reference(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_string").Exists(),
				check.That(data.ResourceName).Key("key_vault_password.0.linked_service_name").HasValue("linkkv"),
				check.That(data.ResourceName).Key("key_vault_password.0.secret_name").HasValue("secret"),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func (t LinkedServiceAzureSQLDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["linkedservices"]

	resp, err := clients.DataFactory.LinkedServiceClient.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Linked Service Azure SQL Database (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (LinkedServiceAzureSQLDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_azure_sql_database" "test" {
  name                = "acctestlssql%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = "data source=serverhostname;initial catalog=master;user id=testUser;Password=test;integrated security=False;encrypt=True;connection timeout=30"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (LinkedServiceAzureSQLDatabaseResource) update1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_azure_sql_database" "test" {
  name                = "acctestlssql%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = "data source=serverhostname;initial catalog=master;user id=testUser;Password=test;integrated security=False;encrypt=True;connection timeout=30"
  annotations         = ["test1", "test2", "test3"]
  description         = "test description"

  parameters = {
    foo = "test1"
    bar = "test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (LinkedServiceAzureSQLDatabaseResource) update2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_azure_sql_database" "test" {
  name                = "acctestlssql%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = "data source=serverhostname;initial catalog=master;user id=testUser;Password=test;integrated security=False;encrypt=True;connection timeout=30;"
  annotations         = ["test1", "test2"]
  description         = "test description 2"

  parameters = {
    foo  = "test1"
    bar  = "test2"
    buzz = "test3"
  }

  additional_properties = {
    foo = "test1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (LinkedServiceAzureSQLDatabaseResource) managed_id(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}
resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}
data "azurerm_client_config" "current" {
}
resource "azurerm_data_factory_linked_service_azure_sql_database" "test" {
  name                 = "acctestlssql%d"
  resource_group_name  = azurerm_resource_group.test.name
  data_factory_name    = azurerm_data_factory.test.name
  connection_string    = "data source=serverhostname;initial catalog=master;user id=testUser;Password=test;integrated security=False;encrypt=True;connection timeout=30"
  use_managed_identity = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (LinkedServiceAzureSQLDatabaseResource) key_vault_reference(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctkv%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name                = "linkkv"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  key_vault_id        = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_azure_sql_database" "test" {
  name                = "acctestlssql%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = "data source=serverhostname;initial catalog=master;user id=testUser;integrated security=False;encrypt=True;connection timeout=30"

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
