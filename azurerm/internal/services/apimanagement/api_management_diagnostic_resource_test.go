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

type ApiManagementDiagnosticResource struct {
}

func TestAccApiManagementDiagnostic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

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

func TestAccApiManagementDiagnostic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

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

func TestAccApiManagementDiagnostic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

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

func TestAccApiManagementDiagnostic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
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

func (t ApiManagementDiagnosticResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	diagnosticId, err := parse.DiagnosticID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.DiagnosticClient.Get(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.Name)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Diagnostic (%s): %+v", diagnosticId.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementDiagnosticResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementDiagnosticResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_management_logger_id = azurerm_api_management_logger.test.id
}
`, r.template(data))
}

func (r ApiManagementDiagnosticResource) update(data acceptance.TestData) string {
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

resource "azurerm_api_management_diagnostic" "test" {
  identifier               = "applicationinsights"
  resource_group_name      = azurerm_resource_group.test.name
  api_management_name      = azurerm_api_management.test.name
  api_management_logger_id = azurerm_api_management_logger.test2.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementDiagnosticResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_diagnostic" "import" {
  identifier               = azurerm_api_management_diagnostic.test.identifier
  resource_group_name      = azurerm_api_management_diagnostic.test.resource_group_name
  api_management_name      = azurerm_api_management_diagnostic.test.api_management_name
  api_management_logger_id = azurerm_api_management_diagnostic.test.api_management_logger_id
}
`, r.basic(data))
}

func (r ApiManagementDiagnosticResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_diagnostic" "test" {
  identifier                = "applicationinsights"
  resource_group_name       = azurerm_resource_group.test.name
  api_management_name       = azurerm_api_management.test.name
  api_management_logger_id  = azurerm_api_management_logger.test.id
  sampling_percentage       = 11.1
  always_log_errors         = false
  log_client_ip             = false
  http_correlation_protocol = "Legacy"
  verbosity                 = "error"

  frontend_request {
    body_bytes     = 100
    headers_to_log = ["Accept"]
  }

  frontend_response {
    body_bytes     = 1000
    headers_to_log = ["Content-Length"]
  }

  backend_request {
    body_bytes     = 1
    headers_to_log = ["Host", "Content-Encoding"]
  }

  backend_response {
    body_bytes     = 10
    headers_to_log = ["Content-Type"]
  }
}
`, r.template(data))
}
