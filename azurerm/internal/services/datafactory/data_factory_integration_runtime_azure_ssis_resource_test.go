package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IntegrationRuntimeManagedSsisResource struct {
}

func TestAccDataFactoryIntegrationRuntimeManagedSsis_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_azure_ssis", "test")
	r := IntegrationRuntimeManagedSsisResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("managed-integration-runtime"),
				check.That(data.ResourceName).Key("data_factory_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
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
			Config: r.complete(data),
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
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  node_size           = "Standard_D8_v3"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeManagedSsisResource) complete(data acceptance.TestData) string {
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
  address_prefix       = "10.0.2.0/24"
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

resource "azurerm_sql_server" "test" {
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
  name                = "acctestSIRsh%[1]d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name                = "acctestiras%[1]d"
  description         = "acctest"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

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
    server_endpoint        = "${azurerm_sql_server.test.fully_qualified_domain_name}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (IntegrationRuntimeManagedSsisResource) aadAuth(data acceptance.TestData) string {
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

data "azuread_service_principal" "test" {
  display_name = azurerm_data_factory.test.name
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = azurerm_data_factory.test.name
  tenant_id           = azurerm_data_factory.test.identity.0.tenant_id
  object_id           = data.azuread_service_principal.test.application_id
}

resource "azurerm_data_factory_integration_runtime_azure_ssis" "test" {
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  node_size           = "Standard_D8_v3"

  catalog_info {
    server_endpoint = azurerm_sql_server.test.fully_qualified_domain_name
    pricing_tier    = "Basic"
  }

  depends_on = [azurerm_sql_active_directory_administrator.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (t IntegrationRuntimeManagedSsisResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["integrationruntimes"]

	resp, err := clients.DataFactory.IntegrationRuntimesClient.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Integration Runtime Managed SSIS (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
