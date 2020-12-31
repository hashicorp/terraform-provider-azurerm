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

func TestAccAzureRMDataFactoryLinkedServiceAzureFileStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_file_storage", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageExists(data.ResourceName),
				),
			},
			data.ImportStep("connection_string"),
		},
	})
}

func TestAccAzureRMDataFactoryLinkedServiceAzureFileStorage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_file_storage", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_update1(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "annotations.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_properties.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "test description"),
				),
			},
			data.ImportStep("connection_string"),
			{
				Config: testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_update2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "annotations.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_properties.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "test description 2"),
				),
			},
			data.ImportStep("connection_string"),
		},
	})
}

func TestAccAzureRMDataFactoryLinkedServiceAzureFileStorage_file_share_name(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_azure_file_storage", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_file_share_name(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "file_share", "myshare"),
				),
			},
			data.ImportStep("connection_string"),
		},
	})
}

func testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.LinkedServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactoryLinkedServiceClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory Linked Service BlobStorage %q (data factory name: %q / resource group: %q) does not exist", name, dataFactoryName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryLinkedServiceAzureFileStorageDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.LinkedServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_linked_service_azure_file_storage" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory Linked Service BlobStorage still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_basic(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_file_storage" "test" {
  name                = "acctestlsblob%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = "DefaultEndpointsProtocol=https;AccountName=foo;AccountKey=bar"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_update1(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_file_storage" "test" {
  name                = "acctestlsblob%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = "DefaultEndpointsProtocol=https;AccountName=foo2;AccountKey=bar"
  annotations         = ["test1", "test2", "test3"]
  description         = "test description"

  parameters = {
    foO = "test1"
    bar = "test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_update2(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_file_storage" "test" {
  name                = "acctestlsblob%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name

  connection_string = "DefaultEndpointsProtocol=https;AccountName=foo3;AccountKey=bar"
  annotations       = ["Test1", "Test2"]
  description       = "test description 2"

  parameters = {
    foo  = "Test1"
    bar  = "test2"
    buzz = "test3"
  }

  additional_properties = {
    foo = "test1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryLinkedServiceAzureFileStorage_file_share_name(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_azure_file_storage" "test" {
  name                = "acctestlsblob%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  connection_string   = "DefaultEndpointsProtocol=https;AccountName=foo;AccountKey=bar"
  file_share          = "myshare"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
