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

func TestAccAzureRMApiManagementApiOperation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiOperation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementApiOperation_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementApiOperation_customMethod(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_customMethod(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "method", "HAMMERTIME"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiOperation_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_headers(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiOperation_requestRepresentations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_requestRepresentation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagementApiOperation_requestRepresentationUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiOperation_representations(t *testing.T) {
	// TODO: once `azurerm_api_management_schema` is supported add `request.0.representation.0.schema_id`
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_representation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagementApiOperation_representationUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApiManagementApiOperationDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiOperationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiOperationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		operationId := rs.Primary.Attributes["operation_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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

func testAccAzureRMApiManagementApiOperation_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  display_name        = "DELETE Resource"
  method              = "DELETE"
  url_template        = "/resource"
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_customMethod(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  display_name        = "HAMMERTIME Resource"
  method              = "HAMMERTIME"
  url_template        = "/resource"
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "import" {
  operation_id        = azurerm_api_management_api_operation.test.operation_id
  api_name            = azurerm_api_management_api_operation.test.api_name
  api_management_name = azurerm_api_management_api_operation.test.api_management_name
  resource_group_name = azurerm_api_management_api_operation.test.resource_group_name
  display_name        = azurerm_api_management_api_operation.test.display_name
  method              = azurerm_api_management_api_operation.test.method
  url_template        = azurerm_api_management_api_operation.test.url_template
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_requestRepresentation(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
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

func testAccAzureRMApiManagementApiOperation_requestRepresentationUpdated(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
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

func testAccAzureRMApiManagementApiOperation_headers(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
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

func testAccAzureRMApiManagementApiOperation_representation(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
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

func testAccAzureRMApiManagementApiOperation_representationUpdated(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiOperation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
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

func testAccAzureRMApiManagementApiOperation_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
