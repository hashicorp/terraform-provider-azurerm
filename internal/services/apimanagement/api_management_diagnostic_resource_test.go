// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementDiagnosticResource struct{}

func TestAccApiManagementDiagnostic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

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

func TestAccApiManagementDiagnostic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementDiagnostic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

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

func TestAccApiManagementDiagnostic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

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

func TestAccApiManagementDiagnostic_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_diagnostic", "test")
	r := ApiManagementDiagnosticResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementDiagnosticResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	diagnosticId, err := diagnostic.ParseDiagnosticID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.DiagnosticClient.Get(ctx, *diagnosticId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *diagnosticId, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Properties != nil), nil
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
  sku_name            = "Consumption_0"
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
    data_masking {
      query_params {
        mode  = "Hide"
        value = "backend-Request-Test"
      }
      headers {
        mode  = "Mask"
        value = "backend-Request-Header"
      }
    }
  }

  frontend_response {
    body_bytes     = 1000
    headers_to_log = ["Content-Length"]
    data_masking {
      query_params {
        mode  = "Hide"
        value = "backend-Request-Test"
      }
      headers {
        mode  = "Mask"
        value = "backend-Request-Header"
      }
    }
  }

  backend_request {
    body_bytes     = 1
    headers_to_log = ["Host", "Content-Encoding"]
    data_masking {
      query_params {
        mode  = "Hide"
        value = "backend-Request-Test"
      }
      headers {
        mode  = "Mask"
        value = "backend-Request-Header"
      }
    }
  }

  backend_response {
    body_bytes     = 10
    headers_to_log = ["Content-Type"]
    data_masking {
      query_params {
        mode  = "Hide"
        value = "backend-Request-Test"
      }
      headers {
        mode  = "Mask"
        value = "backend-Request-Header"
      }
    }
  }
  operation_name_format = "Name"
}
`, r.template(data))
}

func (r ApiManagementDiagnosticResource) completeUpdate(data acceptance.TestData) string {
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
  operation_name_format     = "Url"
}
`, r.template(data))
}
