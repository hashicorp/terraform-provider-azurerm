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

type ApiManagementWorkspaceLoggerEventhubTestResource struct{}

func TestAccApiManagementWorkspaceLoggerEventhub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger_eventhub", "test")
	r := ApiManagementWorkspaceLoggerEventhubTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.#", "eventhub.0.%", "eventhub.0.connection_string"),
	})
}

func TestAccApiManagementWorkspaceLoggerEventhub_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger_eventhub", "test")
	r := ApiManagementWorkspaceLoggerEventhubTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.#", "eventhub.0.%", "eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.systemAssignedIdentityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.#", "eventhub.0.%", "eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.#", "eventhub.0.%", "eventhub.0.user_assigned_identity_client_id"),
	})
}

func TestAccApiManagementWorkspaceLoggerEventhub_managedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_logger_eventhub", "test")
	r := ApiManagementWorkspaceLoggerEventhubTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.#", "eventhub.0.%", "eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.managedIdentityUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.#", "eventhub.0.%", "eventhub.0.user_assigned_identity_client_id"),
		{
			Config: r.managedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub.#", "eventhub.0.%", "eventhub.0.user_assigned_identity_client_id"),
	})
}

func (ApiManagementWorkspaceLoggerEventhubTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r ApiManagementWorkspaceLoggerEventhubTestResource) basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_workspace_logger_eventhub" "test" {
  name                        = "acctestapimlogger-%[2]d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  eventhub {
    name              = azurerm_eventhub.test.name
    connection_string = azurerm_eventhub_namespace.test.default_primary_connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceLoggerEventhubTestResource) managedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger_eventhub" "test" {
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

func (r ApiManagementWorkspaceLoggerEventhubTestResource) managedIdentityUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger_eventhub" "test" {
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

func (r ApiManagementWorkspaceLoggerEventhubTestResource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger_eventhub" "test" {
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

func (r ApiManagementWorkspaceLoggerEventhubTestResource) systemAssignedIdentityUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_management_workspace_logger_eventhub" "test" {
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

func (r ApiManagementWorkspaceLoggerEventhubTestResource) template(data acceptance.TestData) string {
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

func (r ApiManagementWorkspaceLoggerEventhubTestResource) templateWithManagedIdentity(data acceptance.TestData) string {
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

func (r ApiManagementWorkspaceLoggerEventhubTestResource) templateWithSystemAssignedIdentity(data acceptance.TestData) string {
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
