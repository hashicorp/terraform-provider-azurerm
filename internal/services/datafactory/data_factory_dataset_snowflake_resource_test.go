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

type DatasetSnowflakeResource struct{}

func TestAccDataFactoryDatasetSnowflake_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_snowflake", "test")
	r := DatasetSnowflakeResource{}

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

func TestAccDataFactoryDatasetSnowflake_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_dataset_snowflake", "test")
	r := DatasetSnowflakeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("2"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("2"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("2"),
				check.That(data.ResourceName).Key("description").HasValue("test description"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("schema_column.#").HasValue("2"),
				check.That(data.ResourceName).Key("parameters.%").HasValue("3"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("3"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("description").HasValue("test description 2"),
			),
		},
		data.ImportStep(),
	})
}

func (t DatasetSnowflakeResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (DatasetSnowflakeResource) basic(data acceptance.TestData) string {
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
  name            = "linkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_snowflake" "test" {
  name            = "linksnowflake"
  data_factory_id = azurerm_data_factory.test.id

  connection_string = "jdbc:snowflake://account.region.snowflakecomputing.com/?user=user&db=db&warehouse=wh"
  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}

resource "azurerm_data_factory_dataset_snowflake" "test" {
  name                = "acctestds%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_snowflake.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (DatasetSnowflakeResource) update1(data acceptance.TestData) string {
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
  name            = "linkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_snowflake" "test" {
  name            = "linksnowflake"
  data_factory_id = azurerm_data_factory.test.id

  connection_string = "jdbc:snowflake://account.region.snowflakecomputing.com/?user=user&db=db&warehouse=wh"
  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}

resource "azurerm_data_factory_dataset_snowflake" "test" {
  name                = "acctestds%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_snowflake.test.name

  description = "test description"
  annotations = ["test1", "test2"]

  table_name  = "testTable"
  schema_name = "testSchema"
  folder      = "testFolder"

  parameters = {
    foo = "test1"
    bar = "test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (DatasetSnowflakeResource) update2(data acceptance.TestData) string {
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
  name            = "linkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_snowflake" "test" {
  name            = "linksnowflake"
  data_factory_id = azurerm_data_factory.test.id

  connection_string = "jdbc:snowflake://account.region.snowflakecomputing.com/?user=user&db=db&warehouse=wh"
  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}

resource "azurerm_data_factory_dataset_snowflake" "test" {
  name                = "acctestds%d"
  data_factory_id     = azurerm_data_factory.test.id
  linked_service_name = azurerm_data_factory_linked_service_snowflake.test.name

  description = "test description 2"
  annotations = ["test1", "test2", "test3"]
  table_name  = "testTable"
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
    name      = "test1"
    type      = "NUMBER"
    precision = 20
    scale     = 1
  }

  schema_column {
    name      = "test2"
    type      = "NUMBER"
    precision = 25
    scale     = 0
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
