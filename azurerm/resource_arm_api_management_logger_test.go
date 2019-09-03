package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementLogger_basicEventHub(t *testing.T) {
	resourceName := "azurerm_api_management_logger.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.connection_string"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"eventhub.0.connection_string"},
			},
		},
	})
}

func TestAccAzureRMApiManagementLogger_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_logger.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.connection_string"),
				),
			},
			{
				Config:      testAccAzureRMApiManagementLogger_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_logger"),
			},
		},
	})
}

func TestAccAzureRMApiManagementLogger_basicApplicationInsights(t *testing.T) {
	resourceName := "azurerm_api_management_logger.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicApplicationInsights(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application_insights.#", "application_insights.0.instrumentation_key"},
			},
		},
	})
}

func TestAccAzureRMApiManagementLogger_complete(t *testing.T) {
	resourceName := "azurerm_api_management_logger.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_complete(ri, location, "Logger from Terraform test", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Logger from Terraform test"),
					resource.TestCheckResourceAttr(resourceName, "buffered", "false"),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application_insights.#", "application_insights.0.instrumentation_key"},
			},
		},
	})
}

func TestAccAzureRMApiManagementLogger_update(t *testing.T) {
	resourceName := "azurerm_api_management_logger.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementLoggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementLogger_basicApplicationInsights(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.connection_string"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_complete(ri, location, "Logger from Terraform test", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Logger from Terraform test"),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_complete(ri, location, "Logger from Terraform update test", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Logger from Terraform update test"),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_complete(ri, location, "Logger from Terraform test", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Logger from Terraform test"),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "application_insights.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "application_insights.0.instrumentation_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementLogger_basicEventHub(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementLoggerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "buffered", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "eventhub.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub.0.connection_string"),
				),
			},
		},
	})
}

func testCheckAzureRMApiManagementLoggerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("API Management Logger not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagement.LoggerClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).apiManagement.LoggerClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMApiManagementLogger_basicEventHub(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  eventhub {
    name              = "${azurerm_eventhub.test.name}"
    connection_string = "${azurerm_eventhub_namespace.test.default_primary_connection_string}"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMApiManagementLogger_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementLogger_basicEventHub(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_logger" "import" {
  name                = "${azurerm_api_management_logger.test.name}"
  api_management_name = "${azurerm_api_management_logger.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_logger.test.resource_group_name}"

  eventhub {
    name              = "${azurerm_eventhub.test.name}"
    connection_string = "${azurerm_eventhub_namespace.test.default_primary_connection_string}"
  }
}
`, template)
}

func testAccAzureRMApiManagementLogger_basicApplicationInsights(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "Other"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  application_insights {
    instrumentation_key = "${azurerm_application_insights.test.instrumentation_key}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMApiManagementLogger_complete(rInt int, location string, description, buffered string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "Other"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  description         = "%s"
  buffered            = %s

  application_insights {
    instrumentation_key = "${azurerm_application_insights.test.instrumentation_key}"
  }
}
`, rInt, location, rInt, rInt, rInt, description, buffered)
}
