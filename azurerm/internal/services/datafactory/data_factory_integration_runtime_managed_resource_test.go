package datafactory_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataFactoryIntegrationRuntimeManaged_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryIntegrationRuntimeManagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryIntegrationRuntimeManaged_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryIntegrationRuntimeManagedExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryIntegrationRuntimeManaged_vnetIntegration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryIntegrationRuntimeManagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryIntegrationRuntimeManaged_vnetIntegration(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryIntegrationRuntimeManagedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "vnet_integration.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vnet_integration.0.vnet_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vnet_integration.0.subnet_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryIntegrationRuntimeManaged_catalogInfo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryIntegrationRuntimeManagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryIntegrationRuntimeManaged_catalogInfo(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryIntegrationRuntimeManagedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "catalog_info.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "catalog_info.0.server_endpoint"),
					resource.TestCheckResourceAttr(data.ResourceName, "catalog_info.0.administrator_login", "ssis_catalog_admin"),
					resource.TestCheckResourceAttr(data.ResourceName, "catalog_info.0.administrator_password", "my-s3cret-p4ssword!"),
					resource.TestCheckResourceAttr(data.ResourceName, "catalog_info.0.pricing_tier", "Basic"),
				),
			},
			data.ImportStep("catalog_info.0.administrator_password"),
		},
	})
}

func TestAccAzureRMDataFactoryIntegrationRuntimeManaged_customSetupScript(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryIntegrationRuntimeManagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryIntegrationRuntimeManaged_customSetupScript(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryIntegrationRuntimeManagedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_setup_script.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "custom_setup_script.0.blob_container_uri"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "custom_setup_script.0.sas_token"),
				),
			},
			data.ImportStep("catalog_info.0.administrator_password", "custom_setup_script.0.sas_token"),
		},
	})
}

func testAccAzureRMDataFactoryIntegrationRuntimeManaged_basic(data acceptance.TestData) string {
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
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  node_size                        = "Standard_D8_v3"
  number_of_nodes                  = 2
  max_parallel_executions_per_node = 8
  edition                          = "Standard"
  license_type                     = "LicenseIncluded"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDataFactoryIntegrationRuntimeManaged_vnetIntegration(data acceptance.TestData) string {
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
  name               = "acctestsubnet%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  node_size = "Standard_D8_v3"

  vnet_integration {
    vnet_id     = "${azurerm_virtual_network.test.id}"
    subnet_name = "${azurerm_subnet.test.name}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryIntegrationRuntimeManaged_catalogInfo(data acceptance.TestData) string {
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
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

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

func testAccAzureRMDataFactoryIntegrationRuntimeManaged_customSetupScript(data acceptance.TestData) string {
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
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  node_size = "Standard_D8_v3"

  custom_setup_script {
    blob_container_uri = "${azurerm_storage_account.test.primary_blob_endpoint}/${azurerm_storage_container.test.name}"
    sas_token          = "${data.azurerm_storage_account_blob_container_sas.test.sas}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func testCheckAzureRMDataFactoryIntegrationRuntimeManagedExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.IntegrationRuntimesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		factoryName := rs.Primary.Attributes["data_factory_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory Managed Integration Runtime: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on IntegrationRuntimesClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory Managed Integration Runtime %q (Resource Group: %q, Data Factory %q) does not exist", name, factoryName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryIntegrationRuntimeManagedDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_integration_managed" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		factoryName := rs.Primary.Attributes["data_factory_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory Managed Integration Runtime still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}
