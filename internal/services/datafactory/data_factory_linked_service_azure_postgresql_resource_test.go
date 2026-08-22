// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinkedServiceAzurePostgreSQLResource struct{}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedManagedIdentityAuthType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAuthType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedManagedIdentityAuthType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_basicAuthType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAuthType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_basicAuthTypeRequiresUsername(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basicAuthTypeWithoutUsername(data),
			ExpectError: regexp.MustCompile("`username` must be specified when `authentication_type` is `Basic`"),
		},
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_userAssignedManagedIdentityAuthType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedManagedIdentityAuthType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryLinkedServiceAzurePostgreSQL_userAssignedManagedIdentityRequiresCredentialName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_postgresql", "test")
	r := LinkedServiceAzurePostgreSQLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.userAssignedManagedIdentityAuthTypeWithoutCredentialName(data),
			ExpectError: regexp.MustCompile("`credential_name` must be specified when `authentication_type` is `UserAssignedManagedIdentity`"),
		},
	})
}

func (t LinkedServiceAzurePostgreSQLResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := linkedservices.ParseLinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.DataFactory.LinkedServicesClient.Get(ctx, *id, linkedservices.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Linked Service PostgreSQL (%s): %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (LinkedServiceAzurePostgreSQLResource) template(data acceptance.TestData) string {
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

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LinkedServiceAzurePostgreSQLResource) systemAssignedManagedIdentityAuthType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name                = "acctestadfpostgresql%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "SystemAssignedManagedIdentity"
  server              = "acctest-server.postgres.database.azure.com"
  port                = 5432
  database_name       = "acctestdb"
}
`, r.template(data), data.RandomInteger)
}

func (r LinkedServiceAzurePostgreSQLResource) basicAuthType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctkv%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name            = "linkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name                = "acctestadfpostgresql%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "Basic"
  server              = "acctest-server.postgres.database.azure.com"
  port                = 5432
  database_name       = "acctestdb"
  username            = "testuser"

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r LinkedServiceAzurePostgreSQLResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_linked_service_azure_postgresql" "import" {
  name                = azurerm_data_factory_linked_service_azure_postgresql.test.name
  data_factory_id     = azurerm_data_factory_linked_service_azure_postgresql.test.data_factory_id
  authentication_type = azurerm_data_factory_linked_service_azure_postgresql.test.authentication_type
  server              = azurerm_data_factory_linked_service_azure_postgresql.test.server
  port                = azurerm_data_factory_linked_service_azure_postgresql.test.port
  database_name       = azurerm_data_factory_linked_service_azure_postgresql.test.database_name
  username            = azurerm_data_factory_linked_service_azure_postgresql.test.username

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}
`, r.basicAuthType(data))
}

func (r LinkedServiceAzurePostgreSQLResource) basicAuthTypeWithoutUsername(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctkv%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name            = "linkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name                = "acctestadfpostgresql%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "Basic"
  server              = "acctest-server.postgres.database.azure.com"
  port                = 5432
  database_name       = "acctestdb"

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r LinkedServiceAzurePostgreSQLResource) userAssignedManagedIdentityAuthType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "testadf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_data_factory_credential_user_managed_identity" "test" {
  name            = azurerm_user_assigned_identity.test.name
  description     = "Test ADF PostgreSQL DB UMI"
  data_factory_id = azurerm_data_factory.test.id
  identity_id     = azurerm_user_assigned_identity.test.id
}

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name                = "acctestadfpostgresql%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "UserAssignedManagedIdentity"
  server              = "acctest-server.postgres.database.azure.com"
  port                = 5432
  database_name       = "acctestdb"

  credential_name = azurerm_data_factory_credential_user_managed_identity.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LinkedServiceAzurePostgreSQLResource) userAssignedManagedIdentityAuthTypeWithoutCredentialName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name                = "acctestadfpostgresql%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "UserAssignedManagedIdentity"
  server              = "acctest-server.postgres.database.azure.com"
  port                = 5432
  database_name       = "acctestdb"
}
`, r.template(data), data.RandomInteger)
}

func (r LinkedServiceAzurePostgreSQLResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_integration_runtime_self_hosted" "test" {
  name            = "acctest-adf-%d"
  data_factory_id = azurerm_data_factory.test.id
}

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name                     = "acctestadfpostgresql%d"
  data_factory_id          = azurerm_data_factory.test.id
  authentication_type      = "SystemAssignedManagedIdentity"
  server                   = "acctest-server-complete.postgres.database.azure.com"
  port                     = 5431
  database_name            = "acctestdbcomplete"
  annotations              = ["test1", "test2"]
  description              = "test description"
  ssl_mode                 = "2"
  integration_runtime_name = azurerm_data_factory_integration_runtime_self_hosted.test.name
  parameters = {
    foo  = "test1"
    bar  = "test2"
    buzz = "test3"
  }

}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}
