// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IntegrationRuntimeManagedSsisResource struct{}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("managed-integration-runtime"),
				check.That(data.ResourceName).Key("location").HasValue(azure.NormalizeLocation(data.Locations.Primary)),
				check.That(data.ResourceName).Key("node_size").HasValue("Standard_D8_v3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "Basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"catalog_info.0.administrator_password",
			"custom_setup_script.0.sas_token",
			"express_custom_setup.0.component.0.license",
			"express_custom_setup.0.command_key.0.password",
		),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_SSISDBpricingtier_GP_S_Gen5_1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "GP_S_Gen5_1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"catalog_info.0.administrator_password",
			"custom_setup_script.0.sas_token",
			"express_custom_setup.0.component.0.license",
			"express_custom_setup.0.command_key.0.password",
		),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_SSISDBpricingtier_S0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "S0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"catalog_info.0.administrator_password",
			"custom_setup_script.0.sas_token",
			"express_custom_setup.0.component.0.license",
			"express_custom_setup.0.command_key.0.password",
		),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_vnetIntegration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vnetIntegration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"catalog_info.0.administrator_password",
			"custom_setup_script.0.sas_token",
			"express_custom_setup.0.component.0.license",
			"express_custom_setup.0.command_key.0.password",
		),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_keyVaultSecretReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultSecretReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"catalog_info.0.administrator_password",
			"custom_setup_script.0.sas_token",
		),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_aadAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aadAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_userAssignedManagedCredentials(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedManagedCredentials(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_expressVnetInjection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expressVnetInjection(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_withElasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withElasticPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("catalog_info.0.administrator_password"),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_computeScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.computeScale(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_computeScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.computeScale(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (IntegrationRuntimeManagedSsisResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  node_size       = "Standard_D8_v3"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeManagedSsisResource) complete(data acceptance.TestData, pricingTier string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test1" {
  name                = "acctpip1%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
  domain_name_label   = "acctpip1%[1]d"
}

resource "azurerm_public_ip" "test2" {
  name                = "acctpip2%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
  domain_name_label   = "acctpip2%[1]d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "setup-files"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_share" "test" {
  name                 = "sharename"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 30
}

data "azurerm_storage_account_blob_container_sas" "test" {
  connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  container_name    = "${azurerm_storage_container.test.name}"
  https_only        = true

  start  = "2017-03-21"
  expiry = "2022-03-21"

  permissions {
    read   = true
    add    = false
    create = false
    write  = true
    delete = false
    list   = true
  }
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsql%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "ssis_catalog_admin"
  administrator_login_password = "my-s3cret-p4ssword!"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_custom_service" "test" {
  name                 = "acctestls%[1]d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.test.primary_connection_string}"
}
JSON
}

resource "azurerm_data_factory_linked_custom_service" "file_share_linked_service" {
  name                 = "acctestls1%[1]d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "AzureFileStorage"
  type_properties_json = <<JSON
{
  "host": "${azurerm_storage_share.test.url}",
  "password": {
    "type": "SecureString",
    "value": "${azurerm_storage_account.test.primary_access_key}"
  }
}
JSON
}

resource "azurerm_data_factory_integration_runtime_self_hosted" "test" {
  name            = "acctestSIRsh%[1]d"
  data_factory_id = azurerm_data_factory.test.id
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "acctestiras%[1]d"
  description     = "acctest"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size                        = "Standard_D8_v3"
  number_of_nodes                  = 2
  max_parallel_executions_per_node = 8
  edition                          = "Standard"
  license_type                     = "LicenseIncluded"

  vnet_integration {
    vnet_id     = "${azurerm_virtual_network.test.id}"
    subnet_name = "${azurerm_subnet.test.name}"
    public_ips  = [azurerm_public_ip.test1.id, azurerm_public_ip.test2.id]
  }

  catalog_info {
    server_endpoint        = "${azurerm_mssql_server.test.fully_qualified_domain_name}"
    administrator_login    = "ssis_catalog_admin"
    administrator_password = "my-s3cret-p4ssword!"
    pricing_tier           = "%[4]s"
    dual_standby_pair_name = "dual_name"
  }

  custom_setup_script {
    blob_container_uri = "${azurerm_storage_account.test.primary_blob_endpoint}/${azurerm_storage_container.test.name}"
    sas_token          = "${data.azurerm_storage_account_blob_container_sas.test.sas}"
  }

  express_custom_setup {
    powershell_version = "6.2.0"

    environment = {
      Env = "test"
      Foo = "Bar"
    }

    component {
      name    = "SentryOne.TaskFactory"
      license = "license"
    }

    component {
      name = "oh22is.HEDDA.IO"
    }

    command_key {
      target_name = "name1"
      user_name   = "username1"
      password    = "password1"
    }
  }

  package_store {
    name                = "store1"
    linked_service_name = azurerm_data_factory_linked_custom_service.file_share_linked_service.name
  }

  proxy {
    self_hosted_integration_runtime_name = azurerm_data_factory_integration_runtime_self_hosted.test.name
    staging_storage_linked_service_name  = azurerm_data_factory_linked_custom_service.test.name
    path                                 = "containerpath"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, pricingTier)
}

func (IntegrationRuntimeManagedSsisResource) vnetIntegration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "acctestiras%[1]d"
  description     = "acctest"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size = "Standard_D8_v3"

  vnet_integration {
    subnet_id = "${azurerm_subnet.test.id}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (IntegrationRuntimeManagedSsisResource) keyVaultSecretReference(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctkv%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = ["Get", "Delete", "List", "Purge", "Recover", "Set"]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "password"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test1" {
  name                = "acctpip1%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
  domain_name_label   = "acctpip1%[1]d"
}

resource "azurerm_public_ip" "test2" {
  name                = "acctpip2%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
  domain_name_label   = "acctpip2%[1]d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "setup-files"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_share" "test" {
  name                 = "sharename"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 30
}

data "azurerm_storage_account_blob_container_sas" "test" {
  connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  container_name    = "${azurerm_storage_container.test.name}"
  https_only        = true

  start  = "2017-03-21"
  expiry = "2022-03-21"

  permissions {
    read   = true
    add    = false
    create = false
    write  = true
    delete = false
    list   = true
  }
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsql%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "ssis_catalog_admin"
  administrator_login_password = "my-s3cret-p4ssword!"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_custom_service" "test" {
  name                 = "acctestls%[1]d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.test.primary_connection_string}"
}
JSON
}

resource "azurerm_data_factory_linked_custom_service" "file_share_linked_service" {
  name                 = "acctestls1%[1]d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "AzureFileStorage"
  type_properties_json = <<JSON
{
  "host": "${azurerm_storage_share.test.url}",
  "password": {
    "type": "SecureString",
    "value": "${azurerm_storage_account.test.primary_access_key}"
  }
}
JSON
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name            = "acctestlinkkv%[1]d"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_integration_runtime_self_hosted" "test" {
  name            = "acctestSIRsh%[1]d"
  data_factory_id = azurerm_data_factory.test.id
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "acctestiras%[1]d"
  description     = "acctest"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size                        = "Standard_D8_v3"
  number_of_nodes                  = 2
  max_parallel_executions_per_node = 8
  edition                          = "Standard"
  license_type                     = "LicenseIncluded"

  vnet_integration {
    vnet_id     = "${azurerm_virtual_network.test.id}"
    subnet_name = "${azurerm_subnet.test.name}"
    public_ips  = [azurerm_public_ip.test1.id, azurerm_public_ip.test2.id]
  }

  catalog_info {
    server_endpoint        = "${azurerm_mssql_server.test.fully_qualified_domain_name}"
    administrator_login    = "ssis_catalog_admin"
    administrator_password = "my-s3cret-p4ssword!"
    pricing_tier           = "Basic"
    dual_standby_pair_name = "dual_name"
  }

  custom_setup_script {
    blob_container_uri = "${azurerm_storage_account.test.primary_blob_endpoint}/${azurerm_storage_container.test.name}"
    sas_token          = "${data.azurerm_storage_account_blob_container_sas.test.sas}"
  }

  express_custom_setup {
    powershell_version = "6.2.0"

    environment = {
      Env = "test"
      Foo = "Bar"
    }

    component {
      name = "SentryOne.TaskFactory"
      key_vault_license {
        linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
        secret_name         = azurerm_key_vault_secret.test.name
        secret_version      = azurerm_key_vault_secret.test.version
        parameters = {
          "Env" : "Pro",
        }
      }
    }

    component {
      name = "oh22is.HEDDA.IO"
    }

    command_key {
      target_name = "name1"
      user_name   = "username1"
      key_vault_password {
        linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
        secret_name         = azurerm_key_vault_secret.test.name
        secret_version      = azurerm_key_vault_secret.test.version
        parameters = {
          "Env" : "Pro",
        }
      }
    }
  }

  package_store {
    name                = "store1"
    linked_service_name = azurerm_data_factory_linked_custom_service.file_share_linked_service.name
  }

  proxy {
    self_hosted_integration_runtime_name = azurerm_data_factory_integration_runtime_self_hosted.test.name
    staging_storage_linked_service_name  = azurerm_data_factory_linked_custom_service.test.name
    path                                 = "containerpath"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (IntegrationRuntimeManagedSsisResource) aadAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsql%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "ssis_catalog_admin"
  administrator_login_password = "my-s3cret-p4ssword!"

  azuread_administrator {
    login_username = "AzureAD Admin"
    object_id      = data.azurerm_client_config.test.object_id
  }
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  node_size       = "Standard_D8_v3"

  catalog_info {
    server_endpoint = azurerm_mssql_server.test.fully_qualified_domain_name
    pricing_tier    = "Basic"
  }

  depends_on = [azurerm_mssql_server.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (IntegrationRuntimeManagedSsisResource) expressVnetInjection(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  node_size       = "Standard_D8_v3"
  express_vnet_integration {
    subnet_id = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IntegrationRuntimeManagedSsisResource) withElasticPool(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-ssis-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 4.8828125

  sku {
    name     = "BasicPool"
    tier     = "Basic"
    capacity = 50
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 5
  }
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  node_size       = "Standard_D8_v3"

  catalog_info {
    server_endpoint        = "${azurerm_mssql_server.test.fully_qualified_domain_name}"
    administrator_login    = "ssis_catalog_admin"
    administrator_password = "my-s3cret-p4ssword!"
    elastic_pool_name      = "${azurerm_mssql_elasticpool.test.name}"
    dual_standby_pair_name = "dual_name"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (IntegrationRuntimeManagedSsisResource) userAssignedManagedCredentials(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  location            = azurerm_resource_group.test.location
  name                = "testuser%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_data_factory_credential_user_managed_identity" "test" {
  name            = "credential%d"
  data_factory_id = azurerm_data_factory.test.id
  identity_id     = azurerm_user_assigned_identity.test.id
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsql%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "ssis_catalog_admin"
  administrator_login_password = "my-s3cret-p4ssword!"

  azuread_administrator {
    login_username = "AzureAD Admin"
    object_id      = data.azurerm_client_config.test.object_id
  }
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "managed-integration-runtime-%d"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  node_size       = "Standard_D8_v3"
  credential_name = azurerm_data_factory_credential_user_managed_identity.test.name

  catalog_info {
    server_endpoint = azurerm_mssql_server.test.fully_qualified_domain_name
    pricing_tier    = "Basic"
  }

  depends_on = [azurerm_mssql_server.test]
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (IntegrationRuntimeManagedSsisResource) computeScale(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[2]d"
  location = "%[1]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  node_size       = "Standard_D8_v3"

  copy_compute_scale {
    data_integration_unit = 8
    time_to_live          = 6
  }

  pipeline_external_compute_scale {
    number_of_external_nodes = 6
    number_of_pipeline_nodes = 6
    time_to_live             = 8
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (t IntegrationRuntimeManagedSsisResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IntegrationRuntimeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.IntegrationRuntimesClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
