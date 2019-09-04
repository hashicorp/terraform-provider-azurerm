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

func TestAccAzureRMApiManagementApi_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_online", "false"),
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

func TestAccAzureRMApiManagementApi_wordRevision(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_wordRevision(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "revision", "one-point-oh"),
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

func TestAccAzureRMApiManagementApi_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementApi_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_api"),
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_soapPassthrough(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_soapPassthrough(ri, location),
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

func TestAccAzureRMApiManagementApi_importSwagger(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_importSwagger(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_importWsdl(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_importWsdl(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_importUpdate(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_importWsdl(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
			{
				Config: testAccAzureRMApiManagementApi_importSwagger(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"import",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementApi_complete(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApi_complete(ri, location),
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
	conn := testAccProvider.Meta().(*ArmClient).apiManagement.ApiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		revision := rs.Primary.Attributes["revision"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		apiId := fmt.Sprintf("%s;rev=%s", name, revision)

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
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		revision := rs.Primary.Attributes["revision"]

		conn := testAccProvider.Meta().(*ArmClient).apiManagement.ApiClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		apiId := fmt.Sprintf("%s;rev=%s", name, revision)
		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API %q Revision %q (API Management Service %q / Resource Group: %q) does not exist", name, revision, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApi_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}
`, template, rInt)
}

func testAccAzureRMApiManagementApi_wordRevision(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "one-point-oh"
}
`, template, rInt)
}

func testAccAzureRMApiManagementApi_soapPassthrough(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
  soap_pass_through   = true
}
`, template, rInt)
}

func testAccAzureRMApiManagementApi_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "import" {
  name                = "${azurerm_api_management_api.test.name}"
  resource_group_name = "${azurerm_api_management_api.test.resource_group_name}"
  api_management_name = "${azurerm_api_management_api.test.api_management_name}"
  display_name        = "${azurerm_api_management_api.test.display_name}"
  path                = "${azurerm_api_management_api.test.path}"
  protocols           = "${azurerm_api_management_api.test.protocols}"
  revision            = "${azurerm_api_management_api.test.revision}"
}
`, template)
}

func testAccAzureRMApiManagementApi_importSwagger(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"

  import {
    content_value  = "${file("testdata/api_management_api_swagger.json")}"
    content_format = "swagger-json"
  }
}
`, template, rInt)
}

func testAccAzureRMApiManagementApi_importWsdl(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"

  import {
    content_value  = "${file("testdata/api_management_api_wsdl.xml")}"
    content_format = "wsdl"

    wsdl_selector {
      service_name  = "Calculator"
      endpoint_name = "CalculatorHttpsSoap11Endpoint"
    }
  }
}
`, template, rInt)
}

func testAccAzureRMApiManagementApi_complete(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "Butter Parser"
  path                = "butter-parser"
  protocols           = ["https", "http"]
  revision            = "3"
  description         = "What is my purpose? You parse butter."
  service_url         = "https://example.com/foo/bar"

  subscription_key_parameter_names {
    header = "X-Butter-Robot-API-Key"
    query  = "location"
  }
}
`, template, rInt)
}

func testAccAzureRMApiManagementApi_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
`, rInt, location, rInt)
}
