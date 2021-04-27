package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementApiOperationResource struct {
}

func TestAccApiManagementApiOperation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementApiOperation_customMethod(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.customMethod(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("method").HasValue("HAMMERTIME"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.headers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_requestRepresentations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.requestRepresentation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.requestRepresentationUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_representations(t *testing.T) {
	// TODO: once `azurerm_api_management_schema` is supported add `request.0.representation.0.schema_id`
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.representation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.representationUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiOperationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiId := id.Path["apis"]
	operationId := id.Path["operations"]

	resp, err := clients.ApiManagement.ApiOperationsClient.Get(ctx, resourceGroup, serviceName, apiId, operationId)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagementApi Operation (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementApiOperationResource) basic(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r ApiManagementApiOperationResource) customMethod(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r ApiManagementApiOperationResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (r ApiManagementApiOperationResource) requestRepresentation(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r ApiManagementApiOperationResource) requestRepresentationUpdated(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r ApiManagementApiOperationResource) headers(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r ApiManagementApiOperationResource) representation(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r ApiManagementApiOperationResource) representationUpdated(data acceptance.TestData) string {
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
`, r.template(data))
}

func (ApiManagementApiOperationResource) template(data acceptance.TestData) string {
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
