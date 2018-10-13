package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementApi_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApiManagementApi_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_complete(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApiManagementApi_complete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMApiManagementApiDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).apiManagementApiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["service_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		apiId := fmt.Sprintf("%s;rev=%d", name, 1)

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMApiManagementApiExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["service_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for API Management API: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).apiManagementApiClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		apiId := fmt.Sprintf("%s;rev=%d", name, 1)
		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management API %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApi_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
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

resource "azurerm_api_management_api" "test" {
  name            = "acctestAMA-%d"
	service_name    = "${azurerm_api_management.test.name}"
	path 						= "api1"

	import {
    content_value  = "${file("testdata/api_management_api_swagger.json")}"
    content_format = "swagger-json"
  }

  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMApiManagementApi_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
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

resource "azurerm_api_management_api" "test" {
  name            = "acctestAMA-%d"
	service_name    = "${azurerm_api_management.test.name}"
	path 						= "api1"

	import {
		content_value = "${file("testdata/api_management_api_wsdl.xml")}"
    content_format = "wsdl"

    wsdl_selector {
      service_name  = "Calculator"
      endpoint_name = "CalculatorHttpsSoap11Endpoint"
    }
  }

	soap_pass_through = true

	subscription_key_parameter_names {
		header = "test1"
		query = "test2"
	}

	protocols = ["http", "https"]

  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt)
}
