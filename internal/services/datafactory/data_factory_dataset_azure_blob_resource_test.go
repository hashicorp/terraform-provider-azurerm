// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatasetAzureBlobResource struct{}

func TestAccDataFactoryDatasetAzureBlob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_azure_blob", "test")
	r := DatasetAzureBlobResource{}

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

func TestAccDataFactoryDatasetAzureBlob_dynamicPathAndFileName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_azure_blob", "test")
	r := DatasetAzureBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dynamicPathAndFileName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryDatasetAzureBlob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_azure_blob", "test")
	r := DatasetAzureBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update1(data),
			Check: acceptance.ComposeTestCheckFunc(
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
			Config: r.update2(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func (t DatasetAzureBlobResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.DatasetClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (DatasetAzureBlobResource) basic(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_blob_storage" "test" {
  name              = "acctestlssql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "DefaultEndpointsProtocol=https;AccountName=foo;AccountKey=bar"
}

resource "azurerm_data_factory_dataset_azure_blob" "test" {
  name                = "acctestds%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_azure_blob_storage.test.name

  path     = "foo"
  filename = "bar.png"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (DatasetAzureBlobResource) update1(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_blob_storage" "test" {
  name              = "acctestlssql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "DefaultEndpointsProtocol=https;AccountName=foo;AccountKey=bar"
}

resource "azurerm_data_factory_dataset_azure_blob" "test" {
  name                = "acctestds%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_azure_blob_storage.test.name

  path     = "foo"
  filename = "bar.png"

  description = "test description"
  annotations = ["test1", "test2", "test3"]
  folder      = "testFolder"

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

func (DatasetAzureBlobResource) update2(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_blob_storage" "test" {
  name              = "acctestlssql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "DefaultEndpointsProtocol=https;AccountName=foo;AccountKey=bar"
}

resource "azurerm_data_factory_dataset_azure_blob" "test" {
  name                = "acctestds%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_azure_blob_storage.test.name

  path     = "foo"
  filename = "bar.png"

  description = "test description 2"
  annotations = ["test1", "test2"]
  folder      = "testFolder"

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

func (DatasetAzureBlobResource) dynamicPathAndFileName(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_blob_storage" "test" {
  name              = "acctestlssql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "DefaultEndpointsProtocol=https;AccountName=foo;AccountKey=bar"
}

resource "azurerm_data_factory_dataset_azure_blob" "test" {
  name                = "acctestds%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_azure_blob_storage.test.name

  path                     = "@concat('run-', pipeline().RunId)"
  filename                 = "@concat('run-', pipeline().RunId, '.txt')"
  dynamic_filename_enabled = true
  dynamic_path_enabled     = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
