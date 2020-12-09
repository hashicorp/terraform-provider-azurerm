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

func TestAccAzureRMDataFactoryLinkedServiceWeb_anon_auth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_web", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceWebDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryLinkedServiceWeb_anon_auth(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceWebExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryLinkedServiceWeb_basic_auth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_web", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceWebDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryLinkedServiceWeb_basic_auth(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceWebExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "http://www.bing.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "authentication_type", "Basic"),
					resource.TestCheckResourceAttr(data.ResourceName, "username", "foo"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "password"),
				),
			},
			data.ImportStep("password"),
		},
	})
}

func TestAccAzureRMDataFactoryLinkedServiceWeb_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_web", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceWebDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryLinkedServiceWeb_update1(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceWebExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "annotations.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_properties.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "http://www.google.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "authentication_type", "Anonymous"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataFactoryLinkedServiceWeb_update2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceWebExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "annotations.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "additional_properties.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Test Description 2"),
					resource.TestCheckResourceAttr(data.ResourceName, "url", "http://www.yahoo.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "authentication_type", "Anonymous"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataFactoryLinkedServiceWebExists(name string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: Data Factory Linked Service Web %q (data factory name: %q / resource group: %q) does not exist", name, dataFactoryName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryLinkedServiceWebDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.LinkedServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_linked_service_web" {
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
			return fmt.Errorf("Data Factory Linked Service Web still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMDataFactoryLinkedServiceWeb_anon_auth(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_web" "test" {
  name                = "acctestlsweb%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  authentication_type = "Anonymous"
  url                 = "http://www.bing.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryLinkedServiceWeb_basic_auth(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_web" "test" {
  name                = "acctestlsweb%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  authentication_type = "Basic"
  url                 = "http://www.bing.com"
  username            = "foo"
  password            = "bar"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryLinkedServiceWeb_update1(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_web" "test" {
  name                = "acctestlsweb%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  authentication_type = "Anonymous"
  url                 = "http://www.google.com"
  annotations         = ["test1", "test2", "test3"]
  description         = "test description"

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

func testAccAzureRMDataFactoryLinkedServiceWeb_update2(data acceptance.TestData) string {
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

resource "azurerm_data_factory_linked_service_web" "test" {
  name                = "acctestlsweb%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  authentication_type = "Anonymous"
  url                 = "http://www.yahoo.com"
  annotations         = ["test1", "test2"]
  description         = "Test Description 2"

  parameters = {
    foo  = "Test1"
    bar  = "Test2"
    buzz = "Test3"
  }

  additional_properties = {
    foo = "Test1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
