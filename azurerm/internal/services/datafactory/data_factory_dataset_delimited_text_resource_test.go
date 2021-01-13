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

type DatasetDelimitedTextResource struct {
}

func TestAccDataFactoryDatasetDelimitedText_http(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_delimited_text", "test")
	r := DatasetDelimitedTextResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.http(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryDatasetDelimitedText_http_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_delimited_text", "test")
	r := DatasetDelimitedTextResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.http_update1(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("2"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("3"),
				check.That(data.ResourceName).Key("schema_column.#").HasValue("1"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("2"),
				check.That(data.ResourceName).Key("description").HasValue("test description"),
			),
		},
		data.ImportStep(),
		{
			Config: r.http_update2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("3"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("2"),
				check.That(data.ResourceName).Key("schema_column.#").HasValue("2"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("description").HasValue("test description 2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryDatasetDelimitedText_blob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_delimited_text", "test")
	r := DatasetDelimitedTextResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.blob(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t DatasetDelimitedTextResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["datasets"]

	resp, err := clients.DataFactory.DatasetClient.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Dataset Delimited Text (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (DatasetDelimitedTextResource) http(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_web" "test" {
  name                = "acctestlsweb%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  authentication_type = "Anonymous"
  url                 = "https://www.bing.com"
}

resource "azurerm_data_factory_dataset_delimited_text" "test" {
  name                = "acctestds%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  linked_service_name = azurerm_data_factory_linked_service_web.test.name

  http_server_location {
    relative_url = "/fizz/buzz/"
    path         = "foo/bar/"
    filename     = "foo.txt"
  }

  column_delimiter    = ","
  row_delimiter       = "NEW"
  encoding            = "UTF-8"
  quote_character     = "x"
  escape_character    = "f"
  first_row_as_header = true
  null_value          = "NULL"

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (DatasetDelimitedTextResource) http_update1(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_web" "test" {
  name                = "acctestlsweb%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  authentication_type = "Anonymous"
  url                 = "http://www.bing.com"
}

resource "azurerm_data_factory_dataset_delimited_text" "test" {
  name                = "acctestds%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  linked_service_name = azurerm_data_factory_linked_service_web.test.name

  http_server_location {
    relative_url = "/fizz/buzz/"
    path         = "foo/bar/"
    filename     = "foo.txt"
  }

  column_delimiter    = ","
  row_delimiter       = "NEW"
  encoding            = "UTF-8"
  quote_character     = "x"
  escape_character    = "f"
  first_row_as_header = true
  null_value          = "NULL"

  description = "test description"
  annotations = ["test1", "test2", "test3"]

  folder = "testFolder"

  parameters = {
    foo = "test1"
    bar = "test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }

  schema_column {
    name        = "test1"
    type        = "Byte"
    description = "description"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (DatasetDelimitedTextResource) http_update2(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_web" "test" {
  name                = "acctestlsweb%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  authentication_type = "Anonymous"
  url                 = "http://www.bing.com"
}

resource "azurerm_data_factory_dataset_delimited_text" "test" {
  name                = "acctestds%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  linked_service_name = azurerm_data_factory_linked_service_web.test.name

  http_server_location {
    relative_url = "/fizz/buzz/"
    path         = "foo/bar/"
    filename     = "foo.txt"
  }

  column_delimiter    = ","
  row_delimiter       = "NEW"
  encoding            = "UTF-8"
  quote_character     = "x"
  escape_character    = "f"
  first_row_as_header = true
  null_value          = "NULL"

  description = "test description 2"
  annotations = ["test1", "test2"]

  folder = "testFolder"

  parameters = {
    foo  = "test1"
    bar  = "test2"
    buzz = "test3"
  }

  additional_properties = {
    foo = "test1"
  }

  schema_column {
    name        = "test1"
    type        = "Byte"
    description = "description"
  }

  schema_column {
    name        = "test2"
    type        = "Byte"
    description = "description"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (DatasetDelimitedTextResource) blob(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestdf%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}


resource "azurerm_data_factory_linked_service_azure_blob_storage" "test" {
  name                = "acctestlsblob%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = azurerm_storage_account.test.primary_connection_string
}

resource "azurerm_data_factory_dataset_delimited_text" "test" {
  name                = "acctestds%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  linked_service_name = azurerm_data_factory_linked_service_azure_blob_storage.test.name

  azure_blob_storage_location {
    container = azurerm_storage_container.test.name
    path      = "foo/bar/"
    filename  = "foo.txt"
  }

  column_delimiter    = ","
  row_delimiter       = "NEW"
  encoding            = "UTF-8"
  quote_character     = "x"
  escape_character    = "f"
  first_row_as_header = true
  null_value          = "NULL"

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
