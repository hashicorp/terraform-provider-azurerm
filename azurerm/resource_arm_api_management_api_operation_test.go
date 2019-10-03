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

func TestAccAzureRMApiManagementApiOperation_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api_operation.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
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

func TestAccAzureRMApiManagementApiOperation_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_api_operation.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementApiOperation_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_api_operation"),
			},
		},
	})
}

func TestAccAzureRMApiManagementApiOperation_customMethod(t *testing.T) {
	resourceName := "azurerm_api_management_api_operation.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_customMethod(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "method", "HAMMERTIME"),
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

func TestAccAzureRMApiManagementApiOperation_headers(t *testing.T) {
	resourceName := "azurerm_api_management_api_operation.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_headers(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
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

func TestAccAzureRMApiManagementApiOperation_requestRepresentations(t *testing.T) {
	resourceName := "azurerm_api_management_api_operation.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_requestRepresentation(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMApiManagementApiOperation_requestRepresentationUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
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

func TestAccAzureRMApiManagementApiOperation_representations(t *testing.T) {
	// TODO: once `azurerm_api_management_schema` is supported add `request.0.representation.0.schema_id`
	resourceName := "azurerm_api_management_api_operation.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_representation(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMApiManagementApiOperation_representationUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
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

func testCheckAzureRMApiManagementApiOperationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).apiManagement.ApiOperationsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_operation" {
			continue
		}

		operationId := rs.Primary.Attributes["operation_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationId)
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

func testCheckAzureRMApiManagementApiOperationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		operationId := rs.Primary.Attributes["operation_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).apiManagement.ApiOperationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Operation %q (API %q / API Management Service %q / Resource Group: %q) does not exist", operationId, apiName, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementApiOperationsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApiOperation_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  display_name        = "DELETE Resource"
  method              = "DELETE"
  url_template        = "/resource"
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_customMethod(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  display_name        = "HAMMERTIME Resource"
  method              = "HAMMERTIME"
  url_template        = "/resource"
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "import" {
  operation_id        = "${azurerm_api_management_api_operation.test.operation_id}"
  api_name            = "${azurerm_api_management_api_operation.test.api_name}"
  api_management_name = "${azurerm_api_management_api_operation.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_api_operation.test.resource_group_name}"
  display_name        = "${azurerm_api_management_api_operation.test.display_name}"
  method              = "${azurerm_api_management_api_operation.test.method}"
  url_template        = "${azurerm_api_management_api_operation.test.url_template}"
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_requestRepresentation(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  display_name        = "Acceptance Test Operation"
  method              = "DELETE"
  url_template        = "/user1"
  description         = "This can only be done by the logged in user."

  request {
    description = "Created user object"

    representation {
      content_type = "application/json"
      type_name    = "User"
    }
  }
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_requestRepresentationUpdated(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  display_name        = "Acceptance Test Operation"
  method              = "DELETE"
  url_template        = "/user1"
  description         = "This can only be done by the logged in user."

  request {
    description = "Created user object"

    representation {
      content_type = "application/json"
      type_name    = "User"
    }
  }
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_headers(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  display_name        = "Acceptance Test Operation"
  method              = "DELETE"
  url_template        = "/user1"
  description         = "This can only be done by the logged in user."

  request {
    description = "Created user object"

    header {
      name     = "X-Test-Operation"
      required = true
      type     = "string"
    }

    representation {
      content_type = "application/json"
      type_name    = "User"
    }
  }

  response {
    status_code = 200
    description = "successful operation"

    header {
      name     = "X-Test-Operation"
      required = true
      type     = "string"
    }

    representation {
      content_type = "application/xml"

      sample = <<SAMPLE
<response>
  <user name="bravo24">
    <groups>
      <group id="abc123" name="First Group" />
      <group id="bcd234" name="Second Group" />
    </groups>
  </user> 
</response>
SAMPLE
    }
  }
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_representation(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  display_name        = "Acceptance Test Operation"
  method              = "DELETE"
  url_template        = "/user1"
  description         = "This can only be done by the logged in user."

  request {
    description = "Created user object"

    representation {
      content_type = "application/json"
      type_name    = "User"
    }
  }

  response {
    status_code = 200
    description = "successful operation"

    representation {
      content_type = "application/xml"

      sample = <<SAMPLE
<response>
  <user name="bravo24">
    <groups>
      <group id="abc123" name="First Group" />
      <group id="bcd234" name="Second Group" />
    </groups>
  </user> 
</response>
SAMPLE
    }
  }
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_representationUpdated(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  display_name        = "Acceptance Test Operation"
  method              = "DELETE"
  url_template        = "/user1"
  description         = "This can only be done by the logged in user."

  request {
    description = "Created user object"

    representation {
      content_type = "application/json"
      type_name    = "User"
    }
  }

  response {
    status_code = 200
    description = "successful operation"

    representation {
      content_type = "application/xml"

      sample = <<SAMPLE
<response>
  <user name="bravo24">
    <groups>
      <group id="abc123" name="First Group" />
      <group id="bcd234" name="Second Group" />
    </groups>
  </user> 
</response>
SAMPLE
    }

    representation {
      content_type = "application/json"

      sample = <<SAMPLE
{
  "user": {
    "groups": [
      {
        "id": "abc123",
        "name": "First Group"
      },
      {
        "id": "bcd234",
        "name": "Second Group"
      }
    ]
  }
}
SAMPLE
    }
  }
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_template(rInt int, location string) string {
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
`, rInt, location, rInt, rInt)
}
