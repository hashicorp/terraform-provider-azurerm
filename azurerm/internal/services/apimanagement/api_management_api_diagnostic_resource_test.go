package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementApiDiagnosticResource struct {
}

func TestAccApiManagementApiDiagnostic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_diagnostic", "test")
	r := ApiManagementApiDiagnosticResource{}

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

func TestAccApiManagementApiDiagnostic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_diagnostic", "test")
	r := ApiManagementApiDiagnosticResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiDiagnostic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_diagnostic", "test")
	r := ApiManagementApiDiagnosticResource{}

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

func TestAccApiManagementApiDiagnostic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_diagnostic", "test")
	r := ApiManagementApiDiagnosticResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ApiManagementApiDiagnosticResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ApiDiagnosticID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiDiagnosticClient.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.DiagnosticName)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagementApiDiagnostic (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementApiDiagnosticResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%[1]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  application_insights {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestAMA-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  revision            = "1"
  display_name        = "Test API"
  path                = "test"
  protocols           = ["https"]

  import {
    content_format = "swagger-link-json"
    content_value  = "http://conferenceapi.azurewebsites.net/?format=json"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementApiDiagnosticResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_name                 = azurerm_api_management_api.test.name
  api_management_logger_id = azurerm_api_management_logger.test.id
}
`, r.template(data))
}

func (r ApiManagementApiDiagnosticResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsightsUpdate-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_api_management_logger" "test2" {
  name                = "acctestapimngloggerUpdate-%[2]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  application_insights {
    instrumentation_key = azurerm_application_insights.test2.instrumentation_key
  }
}

resource "azurerm_api_management_api_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_name                 = azurerm_api_management_api.test.name
  api_management_logger_id = azurerm_api_management_logger.test2.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiDiagnosticResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_diagnostic" "import" {
  identifier               = azurerm_api_management_api_diagnostic.test.identifier
  resource_group_name      = azurerm_api_management_api_diagnostic.test.resource_group_name
  api_management_name      = azurerm_api_management_api_diagnostic.test.api_management_name
  api_name                 = azurerm_api_management_api.test.name
  api_management_logger_id = azurerm_api_management_api_diagnostic.test.api_management_logger_id
}
`, r.basic(data))
}

func (r ApiManagementApiDiagnosticResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_diagnostic" "test" {
  identifier                = "applicationinsights"
  resource_group_name       = azurerm_resource_group.test.name
  api_management_name       = azurerm_api_management.test.name
  api_name                  = azurerm_api_management_api.test.name
  api_management_logger_id  = azurerm_api_management_logger.test.id
  always_log_errors         = true
  log_client_ip             = true
  http_correlation_protocol = "W3C"
  verbosity                 = "verbose"

  backend_request {
    body_bytes     = 1
    headers_to_log = ["Host"]
  }

  backend_response {
    body_bytes     = 2
    headers_to_log = ["Content-Type"]
  }

  frontend_request {
    body_bytes     = 3
    headers_to_log = ["Accept"]
  }

  frontend_response {
    body_bytes     = 4
    headers_to_log = ["Content-Length"]
  }
}
`, r.template(data))
}
