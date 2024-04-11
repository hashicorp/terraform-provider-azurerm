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

type IntegrationRuntimeManagedResource struct{}

func TestAccDataFactoryIntegrationRuntimeManaged_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

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

func TestAccDataFactoryIntegrationRuntimeManaged_vnetIntegration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vnetIntegration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vnet_integration.#").HasValue("1"),
				check.That(data.ResourceName).Key("vnet_integration.0.vnet_id").Exists(),
				check.That(data.ResourceName).Key("vnet_integration.0.subnet_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManaged_catalogInfo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.catalogInfo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("catalog_info.#").HasValue("1"),
				check.That(data.ResourceName).Key("catalog_info.0.server_endpoint").Exists(),
				check.That(data.ResourceName).Key("catalog_info.0.administrator_login").HasValue("ssis_catalog_admin"),
				check.That(data.ResourceName).Key("catalog_info.0.administrator_password").HasValue("my-s3cret-p4ssword!"),
				check.That(data.ResourceName).Key("catalog_info.0.pricing_tier").HasValue("Basic"),
			),
		},
		data.ImportStep("catalog_info.0.administrator_password"),
	})
}

func TestAccDataFactoryIntegrationRuntimeManaged_customSetupScript(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customSetupScript(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_setup_script.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_setup_script.0.blob_container_uri").Exists(),
				check.That(data.ResourceName).Key("custom_setup_script.0.sas_token").Exists(),
			),
		},
		data.ImportStep("catalog_info.0.administrator_password", "custom_setup_script.0.sas_token"),
	})
}

func TestAccDataFactoryIntegrationRuntimeManaged_credential(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.credential(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManaged_aadAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

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

func (IntegrationRuntimeManagedResource) basic(data acceptance.TestData) string {
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

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size                        = "Standard_D8_v3"
  number_of_nodes                  = 2
  max_parallel_executions_per_node = 8
  edition                          = "Standard"
  license_type                     = "LicenseIncluded"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeManagedResource) vnetIntegration(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size = "Standard_D8_v3"

  vnet_integration {
    vnet_id     = "${azurerm_virtual_network.test.id}"
    subnet_name = "${azurerm_subnet.test.name}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (IntegrationRuntimeManagedResource) catalogInfo(data acceptance.TestData) string {
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

resource "azurerm_sql_server" "test" {
  name                         = "acctestsql%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "ssis_catalog_admin"
  administrator_login_password = "my-s3cret-p4ssword!"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size = "Standard_D8_v3"

  catalog_info {
    server_endpoint        = "${azurerm_sql_server.test.fully_qualified_domain_name}"
    administrator_login    = "ssis_catalog_admin"
    administrator_password = "my-s3cret-p4ssword!"
    pricing_tier           = "Basic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (IntegrationRuntimeManagedResource) customSetupScript(data acceptance.TestData) string {
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

resource "azurerm_storage_account" "test" {
  name                      = "acctestsa%s"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  location                  = "${azurerm_resource_group.test.location}"
  account_kind              = "BlobStorage"
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  access_tier               = "Hot"
  enable_https_traffic_only = true
}

resource "azurerm_storage_container" "test" {
  name                  = "setup-files"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
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

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size = "Standard_D8_v3"

  custom_setup_script {
    blob_container_uri = "${azurerm_storage_account.test.primary_blob_endpoint}/${azurerm_storage_container.test.name}"
    sas_token          = "${data.azurerm_storage_account_blob_container_sas.test.sas}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (IntegrationRuntimeManagedResource) aadAuth(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsql%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "ssis_catalog_admin"
  administrator_login_password = "my-s3cret-p4ssword!"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = azurerm_data_factory.test.name
  tenant_id           = azurerm_data_factory.test.identity.0.tenant_id
  object_id           = azurerm_data_factory.test.identity.0.principal_id
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  node_size       = "Standard_D8_v3"

  catalog_info {
    server_endpoint = azurerm_sql_server.test.fully_qualified_domain_name
    pricing_tier    = "Basic"
  }

  depends_on = [azurerm_sql_active_directory_administrator.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (IntegrationRuntimeManagedResource) credential(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[2]d"
  location = "%[1]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestdf%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_data_factory_credential_user_managed_identity" "test" {
  name            = "credential%[2]d"
  description     = "ORIGINAL DESCRIPTION"
  data_factory_id = azurerm_data_factory.test.id
  identity_id     = azurerm_user_assigned_identity.test.id
  annotations     = ["1"]
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name            = "managed-integration-runtime"
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location

  node_size = "Standard_D8_v3"

  credential_name = azurerm_data_factory_credential_user_managed_identity.test.name
}
`, data.Locations.Primary, data.RandomInteger)
}

func (t IntegrationRuntimeManagedResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IntegrationRuntimeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.IntegrationRuntimesClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Integration Runtime Managed (%s): %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
