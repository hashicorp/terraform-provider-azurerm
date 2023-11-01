// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiOperationResource struct{}

func TestAccApiManagementApiOperation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementApiOperation_customMethod(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customMethod(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.headers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_requestRepresentations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.requestRepresentation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.requestRepresentationUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.representation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.representationUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_templateParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.templateParameter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiOperation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation", "test")
	r := ApiManagementApiOperationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiOperationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apioperation.ParseOperationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiOperationsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
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

func (r ApiManagementApiOperationResource) complete(data acceptance.TestData) string {
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

  response {
    status_code = 200

    header {
      name          = "test"
      required      = true
      type          = "string"
      default_value = "default"
      description   = "This is a test description"
      values        = ["multipart/form-data"]
    }

    representation {
      content_type = "multipart/form-data"

      form_parameter {
        default_value = "multipart/form-data"
        description   = "This is a test description"
        name          = "test"
        required      = true
        type          = "string"
        values        = ["multipart/form-data"]
      }

      example {
        name           = "test"
        description    = "This is a test description"
        external_value = "https://example.com/foo/bar"
        summary        = "This is a test summary"
      }
    }
  }

  request {
    description = "Created user object"

    query_parameter {
      default_value = "multipart/form-data"
      description   = "This is a test description"
      name          = "test"
      required      = true
      type          = "string"
      values        = ["multipart/form-data"]
    }

    header {
      name          = "test"
      required      = true
      type          = "string"
      default_value = "default"
      description   = "This is a test description"
    }

    representation {
      content_type = "multipart/form-data"

      example {
        description    = "This is a test description"
        external_value = "https://example.com/foo/bar"
        name           = "test"
        summary        = "This is a test summary"
        value          = "backend-Request-Test"
      }

      form_parameter {
        default_value = "multipart/form-data"
        description   = "This is a test description"
        name          = "test"
        required      = true
        type          = "string"
        values        = ["multipart/form-data"]
      }
    }
  }
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
resource "azurerm_api_management_api_schema" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  schema_id           = "acctestSchema%d"
  content_type        = "application/json"
  value               = file("testdata/api_management_api_schema_swagger.json")
}

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
      values   = ["application/x-www-form-urlencoded"]
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

      type_name = "User"
      schema_id = azurerm_api_management_api_schema.test.schema_id
      example {
        description    = "This is a test description"
        external_value = "https://example.com/foo/bar"
        name           = "test"
        summary        = "This is a test summary"
        value          = "backend-Request-Test"
      }
    }

    representation {
      content_type = "application/xml"

      example {
        name  = "sample"
        value = <<SAMPLE
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
}
`, r.template(data), data.RandomInteger)
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
      type_name    = "User"

      example {
        name  = "sample"
        value = <<SAMPLE
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

      example {
        name  = "sample"
        value = <<SAMPLE
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

    representation {
      content_type = "application/json"

      example {
        name  = "sample"
        value = <<SAMPLE
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
}
`, r.template(data))
}

func (r ApiManagementApiOperationResource) templateParameter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  operation_id        = "acctest-operation"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  display_name        = "Acceptance Test Operation"
  method              = "DELETE"
  url_template        = "/users/{id}/delete"
  description         = "This can only be done by the logged in user."

  template_parameter {
    name     = "id"
    type     = "number"
    required = true
  }

  response {
    status_code = 200
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
  sku_name            = "Consumption_0"
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
