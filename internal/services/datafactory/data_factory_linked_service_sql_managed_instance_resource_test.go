// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinkedServiceSQLManagedInstanceResource struct{}

func TestAccDataFactoryLinkedServiceSQLManagedInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sql_managed_instance", "test")
	r := LinkedServiceSQLManagedInstanceResource{}

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

func TestAccDataFactoryLinkedServiceSQLManagedInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sql_managed_instance", "test")
	r := LinkedServiceSQLManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal_key", "key_vault_password"),
	})
}

func TestAccDataFactoryLinkedServiceSQLManagedInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sql_managed_instance", "test")
	r := LinkedServiceSQLManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataFactoryLinkedServiceSQLManagedInstance_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sql_managed_instance", "test")
	r := LinkedServiceSQLManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("service_principal_key"),
	})
}

func TestAccDataFactoryLinkedServiceSQLManagedInstance_keyVaultPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sql_managed_instance", "test")
	r := LinkedServiceSQLManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryLinkedServiceSQLManagedInstance_keyVaultConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sql_managed_instance", "test")
	r := LinkedServiceSQLManagedInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultConnectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("key_vault_password.0.linked_service_name", "key_vault_connection_string.0.secret_name"),
	})
}

func (t LinkedServiceSQLManagedInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := linkedservices.ParseLinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.LinkedServiceClient.Get(ctx, id.ResourceGroupName, id.FactoryName, id.LinkedServiceName, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory SQL Managed Instance(%s): %+v", *id, err)
	}

	return pointer.To(resp.ID != nil), nil
}

func (LinkedServiceSQLManagedInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_sql_managed_instance" "test" {
  name              = "acctestlssqlmi%[1]d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Server=myserver.database.windows.net;Database=mydatabase;User ID=myuser;Password=mypassword"
  annotations       = ["test1", "test2", "test3"]
  description       = "test description"

  parameters = {
    param1 = "value1"
    param2 = "value2"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LinkedServiceSQLManagedInstanceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_sql_managed_instance" "test" {
  name              = "acctestlssqlmi%[1]d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Server=myserver.database.windows.net;Database=mydatabase;User ID=myuser;Password=mypassword"
  annotations       = ["test1", "test2"]
  description       = "test description 2"

  parameters = {
    param1 = "value1"
    param2 = "value2"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LinkedServiceSQLManagedInstanceResource) servicePrincipal(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_sql_managed_instance" "test" {
  name                  = "acctestlssqlmi%[1]d"
  data_factory_id       = azurerm_data_factory.test.id
  connection_string     = "Server=myserver.database.windows.net;Database=mydatabase"
  service_principal_id  = "00000000-0000-0000-0000-000000000000"
  service_principal_key = "testkey"
  tenant                = "11111111-1111-1111-1111-111111111111"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LinkedServiceSQLManagedInstanceResource) keyVaultPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = "%[3]s"
  sku_name            = "standard"
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name            = "acctestlinkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_sql_managed_instance" "test" {
  name              = "acctestlssqlmi%[1]d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Server=myserver.database.windows.net;Database=mydatabase;User ID=myuser"

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Client().TenantID)
}

func (LinkedServiceSQLManagedInstanceResource) keyVaultConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                = "acckv%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = "%[3]s"
  sku_name            = "standard"
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name            = "acctestlinkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_sql_managed_instance" "test" {
  name            = "acctestlssqlmi%[1]d"
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
`, data.RandomInteger, data.Locations.Primary, data.Client().TenantID)
}

func (r LinkedServiceSQLManagedInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_integration_runtime_azure" "test" {
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  name            = "acctestlssqlmi%[2]d"
}

resource "azurerm_data_factory_linked_service_sql_managed_instance" "test" {
  name                     = "acctestlssqlmi%[2]d"
  data_factory_id          = azurerm_data_factory.test.id
  connection_string        = "Server=myserver.database.windows.net;Database=mydatabase;"
  service_principal_id     = "00000000-0000-0000-0000-000000000000"
  service_principal_key    = "testkey"
  tenant                   = "11111111-1111-1111-1111-111111111111"
  integration_runtime_name = azurerm_data_factory_integration_runtime_azure.test.name
  annotations              = ["test1", "test2", "test3"]
  description              = "complete test description"

  parameters = {
    param1 = "value1"
    param2 = "value2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (LinkedServiceSQLManagedInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
