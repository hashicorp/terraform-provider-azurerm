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

type ApiManagementWorkspaceLoggerTestResource struct{}

func TestAccApiManagementWorkspaceLogger_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger", "test")
	r := ApiManagementWorkspaceLoggerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"),
	})
}

func TestAccApiManagementWorkspaceLogger_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger", "test")
	r := ApiManagementWorkspaceLoggerTestResource{}

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

func TestAccApiManagementWorkspaceLogger_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger", "test")
	r := ApiManagementWorkspaceLoggerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"),
	})
}

func TestAccApiManagementWorkspaceLogger_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger", "test")
	r := ApiManagementWorkspaceLoggerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"),
	})
}

func TestAccApiManagementWorkspaceLogger_basicEventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger", "test")
	r := ApiManagementWorkspaceLoggerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.0.connection_string"),
	})
}

func TestAccApiManagementWorkspaceLogger_systemAssignedIdentityEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger", "test")
	r := ApiManagementWorkspaceLoggerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentityEventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.systemAssignedIdentityEventhubUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.systemAssignedIdentityEventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.0.user_assigned_identity_client_id"),
	})
}

func TestAccApiManagementWorkspaceLogger_managedIdentityEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger", "test")
	r := ApiManagementWorkspaceLoggerTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedIdentityEventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.managedIdentityEventhubUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.managedIdentityEventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.0.user_assigned_identity_client_id"),
	})
}

func (ApiManagementWorkspaceLoggerTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r ApiManagementWorkspaceLoggerTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management_workspace_logger" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  application_insights {
    connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
}

%[1]s

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management_workspace_logger" "test" {
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

func (r ApiManagementWorkspaceLoggerTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    application_insights {
      disable_generated_rule = true
    }
  }
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

resource "azurerm_api_management_workspace_logger" "test" {
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

func (r ApiManagementWorkspaceLoggerTestResource) basicEventhub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctestEventHub-%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management_workspace_logger" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  eventhub {
    name              = azurerm_eventhub.test.name
    connection_string = azurerm_eventhub_namespace.test.default_primary_connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerTestResource) managedIdentityEventhub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  eventhub {
    name                             = azurerm_eventhub.test.name
    endpoint_uri                     = "${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
    user_assigned_identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
  depends_on = [azurerm_role_assignment.test]
}
`, r.templateWithManagedIdentity(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerTestResource) managedIdentityEventhubUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  eventhub {
    name                             = azurerm_eventhub.test2.name
    endpoint_uri                     = "${azurerm_eventhub_namespace.test2.name}.servicebus.windows.net"
    user_assigned_identity_client_id = azurerm_user_assigned_identity.test2.client_id
  }
  depends_on = [azurerm_role_assignment.test2]
}
`, r.templateWithManagedIdentity(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerTestResource) systemAssignedIdentityEventhub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  eventhub {
    name         = azurerm_eventhub.test.name
    endpoint_uri = "${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
  }
  depends_on = [azurerm_role_assignment.test]
}
`, r.templateWithSystemAssignedIdentity(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerTestResource) systemAssignedIdentityEventhubUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  eventhub {
    name         = azurerm_eventhub.test2.name
    endpoint_uri = "${azurerm_eventhub_namespace.test2.name}.servicebus.windows.net"
  }
  depends_on = [azurerm_role_assignment.test2]
}
`, r.templateWithSystemAssignedIdentity(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_logger" "import" {
  name                        = azurerm_api_management_workspace_logger.test.name
  api_management_workspace_id = azurerm_api_management_workspace_logger.test.api_management_workspace_id

  application_insights {
    connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceLoggerTestResource) template(data acceptance.TestData) string {
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

func (r ApiManagementWorkspaceLoggerTestResource) templateWithManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "uai-acctestapimnglogger-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_eventhub_namespace" "test2" {
  name                = "acctesteventhubnamespace2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test2" {
  name                = "acctesteventhub2-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test2.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "uai-acctestapimnglogger2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_eventhub.test.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test2" {
  scope                = azurerm_eventhub.test2.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_user_assigned_identity.test2.principal_id
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
      azurerm_user_assigned_identity.test2.id
    ]
  }
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMWS-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r ApiManagementWorkspaceLoggerTestResource) templateWithSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace" "test2" {
  name                = "acctesteventhubnamespace2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test2" {
  name                = "acctesteventhub2-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test2.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMWS-%[1]d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_eventhub.test.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_api_management.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "test2" {
  scope                = azurerm_eventhub.test2.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_api_management.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary)
}
