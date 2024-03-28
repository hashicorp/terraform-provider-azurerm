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

type LinkedServiceAzurePostgreSQLResource struct{}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update1(data),
			Check: acceptance.ComposeTestCheckFunc(
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
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_PasswordKeyVaultReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.key_vault_reference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_string").Exists(),
				check.That(data.ResourceName).Key("key_vault_password.0.linked_service_name").HasValue("linkkv"),
				check.That(data.ResourceName).Key("key_vault_password.0.secret_name").HasValue("secret"),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_ConnectionStringKeyVaultReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connection_string_key_vault_reference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_connection_string.0.linked_service_name").HasValue("linkkv"),
				check.That(data.ResourceName).Key("key_vault_connection_string.0.secret_name").HasValue("connection_string"),
				check.That(data.ResourceName).Key("key_vault_password.0.linked_service_name").HasValue("linkkv"),
				check.That(data.ResourceName).Key("key_vault_password.0.secret_name").HasValue("password"),
			),
		},
		data.ImportStep(),
	})
}

func (t LinkedServiceAzurePostgreSQLResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.LinkedServiceClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Azure PostgreSQL Database (%s): %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (LinkedServiceAzurePostgreSQLResource) basic(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name              = "acctestlsazpsql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Host=serverHostname;Port=5432;Database=postgres;UID=psqladmin@serverHostname;EncryptionMethod=0;validateservercertificate=1;Password=password123"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (LinkedServiceAzurePostgreSQLResource) update1(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name              = "acctestlsazpsql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Host=serverHostname;Port=5432;Database=postgres;UID=psqladmin@serverHostname;EncryptionMethod=0;validateservercertificate=1;Password=password123"
  annotations       = ["test1", "test2", "test3"]
  description       = "test description"

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

func (LinkedServiceAzurePostgreSQLResource) update2(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name              = "acctestlsazpsql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Host=serverHostname;Port=5432;Database=postgres;UID=psqladmin@serverHostname;EncryptionMethod=0;validateservercertificate=1;Password=password123"
  annotations       = ["test1", "test2"]
  description       = "test description 2"

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

func (LinkedServiceAzurePostgreSQLResource) key_vault_reference(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name              = "acctestlsazpsql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Host=serverHostname;Port=5432;Database=postgres;UID=psqladmin@serverHostname;EncryptionMethod=0;validateservercertificate=1"

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (LinkedServiceAzurePostgreSQLResource) connection_string_key_vault_reference(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name            = "acctestlsazpsql%d"
  data_factory_id = azurerm_data_factory.test.id

  key_vault_connection_string {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "connection_string"
  }

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "password"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
