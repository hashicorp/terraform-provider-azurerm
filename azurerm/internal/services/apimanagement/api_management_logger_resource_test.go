package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementLogger_basicEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.connection_string"),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"eventhub.0.connection_string"},
			},
		},
	})
}

func TestAccAzureRMApiManagementLogger_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.connection_string"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementLogger_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementLogger_basicApplicationInsights(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicApplicationInsights(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application_insights.#", "application_insights.0.instrumentation_key"},
			},
		},
	})
}

func TestAccAzureRMApiManagementLogger_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_complete(data, "Logger from Terraform test", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Logger from Terraform test"),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application_insights.#", "application_insights.0.instrumentation_key"},
			},
		},
	})
}

func TestAccAzureRMApiManagementLogger_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicApplicationInsights(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.connection_string"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_complete(data, "Logger from Terraform test", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Logger from Terraform test"),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_complete(data, "Logger from Terraform update test", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Logger from Terraform update test"),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_complete(data, "Logger from Terraform test", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Logger from Terraform test"),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub.0.connection_string"),
				),
			},
		},
	})
}

func testCheckAzureRMApiManagementLoggerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.LoggerClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("API Management Logger not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Logger %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagement.LoggerClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementLoggerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.LoggerClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_logger" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on apiManagement.LoggerClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMApiManagementLogger_basicEventHub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  eventhub {
    name              = azurerm_eventhub.test.name
    connection_string = azurerm_eventhub_namespace.test.default_primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementLogger_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementLogger_basicEventHub(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_logger" "import" {
  name                = azurerm_api_management_logger.test.name
  api_management_name = azurerm_api_management_logger.test.api_management_name
  resource_group_name = azurerm_api_management_logger.test.resource_group_name

  eventhub {
    name              = azurerm_eventhub.test.name
    connection_string = azurerm_eventhub_namespace.test.default_primary_connection_string
  }
}
`, template)
}

func testAccAzureRMApiManagementLogger_basicApplicationInsights(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  application_insights {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementLogger_complete(data acceptance.TestData, description, buffered string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  description         = "%s"
  buffered            = %s

  application_insights {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, description, buffered)
}
