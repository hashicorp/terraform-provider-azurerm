// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceLoggerApplicationInsightsResource struct{}

func TestAccApiManagementWorkspaceLoggerApplicationInsights_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger_application_insights", "test")
	r := ApiManagementWorkspaceLoggerApplicationInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.%", "application_insights.0.connection_string", "application_insights.0.instrumentation_key"),
	})
}

func TestAccApiManagementWorkspaceLoggerApplicationInsights_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger_application_insights", "test")
	r := ApiManagementWorkspaceLoggerApplicationInsightsResource{}

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

func TestAccApiManagementWorkspaceLoggerApplicationInsights_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger_application_insights", "test")
	r := ApiManagementWorkspaceLoggerApplicationInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.%", "application_insights.0.connection_string", "application_insights.0.instrumentation_key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.%", "application_insights.0.connection_string", "application_insights.0.instrumentation_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.%", "application_insights.0.connection_string", "application_insights.0.instrumentation_key"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.%", "application_insights.0.connection_string", "application_insights.0.instrumentation_key"),
	})
}

func TestAccApiManagementWorkspaceLoggerApplicationInsights_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger_application_insights", "test")
	r := ApiManagementWorkspaceLoggerApplicationInsightsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.%", "application_insights.0.connection_string", "application_insights.0.instrumentation_key"),
	})
}

func (ApiManagementWorkspaceLoggerApplicationInsightsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := logger.ParseWorkspaceLoggerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.LoggerClient_v2024_05_01.WorkspaceLoggerGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management_workspace_logger_application_insights" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  application_insights {
    connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management_workspace_logger_application_insights" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  description                 = "Logger from Terraform test"
  buffering_enabled           = false
  resource_id                 = azurerm_application_insights.test.id

  application_insights {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_application_insights" "test2" {
  name                = "acctestappinsights2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management_workspace_logger_application_insights" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  buffering_enabled           = true
  description                 = "Logger from Terraform test update"
  resource_id                 = azurerm_application_insights.test2.id

  application_insights {
    instrumentation_key = azurerm_application_insights.test2.instrumentation_key
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_logger_application_insights" "import" {
  name                        = azurerm_api_management_workspace_logger_application_insights.test.name
  api_management_workspace_id = azurerm_api_management_workspace_logger_application_insights.test.api_management_workspace_id

  application_insights {
    connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMWS-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
}
`, data.RandomInteger, data.Locations.Primary)
}
