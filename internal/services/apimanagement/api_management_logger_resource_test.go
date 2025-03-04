// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementLoggerResource struct{}

func TestAccApiManagementLogger_basicEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("1"),
				check.That(data.ResourceName).Key("eventhub.0.name").Exists(),
				check.That(data.ResourceName).Key("eventhub.0.connection_string").Exists(),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"eventhub.0.connection_string"},
		},
	})
}

func TestAccApiManagementLogger_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("1"),
				check.That(data.ResourceName).Key("eventhub.0.name").Exists(),
				check.That(data.ResourceName).Key("eventhub.0.connection_string").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementLogger_systemAssignedIdentityEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentityEventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("1"),
				check.That(data.ResourceName).Key("eventhub.0.name").Exists(),
				check.That(data.ResourceName).Key("eventhub.0.endpoint_uri").Exists(),
			),
		},
		data.ImportStep("eventhub.0.endpoint_uri", "eventhub.0.user_assigned_identity_client_id"),
	})
}

func TestAccApiManagementLogger_managedIdentityEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedIdentityEventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("1"),
				check.That(data.ResourceName).Key("eventhub.0.name").Exists(),
				check.That(data.ResourceName).Key("eventhub.0.endpoint_uri").Exists(),
				check.That(data.ResourceName).Key("eventhub.0.user_assigned_identity_client_id").Exists(),
			),
		},
		data.ImportStep("eventhub.0.endpoint_uri", "eventhub.0.user_assigned_identity_client_id"),
	})
}

func TestAccApiManagementLogger_basicApplicationInsights(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicApplicationInsights(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("0"),
				check.That(data.ResourceName).Key("application_insights.#").HasValue("1"),
				check.That(data.ResourceName).Key("application_insights.0.instrumentation_key").Exists(),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"},
		},
	})
}

func TestAccApiManagementLogger_applicationInsightsConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.applicationInsightsConnectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("0"),
				check.That(data.ResourceName).Key("application_insights.#").HasValue("1"),
				check.That(data.ResourceName).Key("application_insights.0.connection_string").Exists(),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"},
		},
	})
}

func TestAccApiManagementLogger_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "Logger from Terraform test", "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("Logger from Terraform test"),
				check.That(data.ResourceName).Key("buffered").HasValue("false"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("0"),
				check.That(data.ResourceName).Key("application_insights.#").HasValue("1"),
				check.That(data.ResourceName).Key("application_insights.0.instrumentation_key").Exists(),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"application_insights.#", "application_insights.0.connection_string", "application_insights.0.instrumentation_key", "application_insights.0.%"},
		},
	})
}

func TestAccApiManagementLogger_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicApplicationInsights(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("0"),
				check.That(data.ResourceName).Key("application_insights.#").HasValue("1"),
				check.That(data.ResourceName).Key("application_insights.0.instrumentation_key").Exists(),
			),
		},
		{
			Config: r.basicEventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("1"),
				check.That(data.ResourceName).Key("eventhub.0.name").Exists(),
				check.That(data.ResourceName).Key("eventhub.0.connection_string").Exists(),
			),
		},
		{
			Config: r.complete(data, "Logger from Terraform test", "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("false"),
				check.That(data.ResourceName).Key("description").HasValue("Logger from Terraform test"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("0"),
				check.That(data.ResourceName).Key("application_insights.#").HasValue("1"),
				check.That(data.ResourceName).Key("application_insights.0.instrumentation_key").Exists(),
			),
		},
		{
			Config: r.complete(data, "Logger from Terraform update test", "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue("Logger from Terraform update test"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("0"),
				check.That(data.ResourceName).Key("application_insights.#").HasValue("1"),
				check.That(data.ResourceName).Key("application_insights.0.instrumentation_key").Exists(),
			),
		},
		{
			Config: r.complete(data, "Logger from Terraform test", "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("false"),
				check.That(data.ResourceName).Key("description").HasValue("Logger from Terraform test"),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("0"),
				check.That(data.ResourceName).Key("application_insights.#").HasValue("1"),
				check.That(data.ResourceName).Key("application_insights.0.instrumentation_key").Exists(),
			),
		},
		{
			Config: r.basicEventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("buffered").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("eventhub.#").HasValue("1"),
				check.That(data.ResourceName).Key("eventhub.0.name").Exists(),
				check.That(data.ResourceName).Key("eventhub.0.connection_string").Exists(),
			),
		},
	})
}

func (ApiManagementLoggerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := logger.ParseLoggerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.LoggerClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementLoggerResource) basicEventHub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  eventhub {
    name              = azurerm_eventhub.test.name
    connection_string = azurerm_eventhub_namespace.test.default_primary_connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementLoggerResource) systemAssignedIdentityEventHub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  eventhub {
    name         = azurerm_eventhub.test.name
    endpoint_uri = "${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
  }

  depends_on = [azurerm_role_assignment.test]
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_eventhub.test.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_api_management.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementLoggerResource) managedIdentityEventHub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  eventhub {
    name                             = azurerm_eventhub.test.name
    endpoint_uri                     = "${azurerm_eventhub_namespace.test.name}.servicebus.windows.net"
    user_assigned_identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "uai-acctestapimnglogger-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_eventhub.test.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementLoggerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_logger" "import" {
  name                = azurerm_api_management_logger.test.name
  api_management_name = azurerm_api_management_logger.test.api_management_name
  resource_group_name = azurerm_api_management_logger.test.resource_group_name

  eventhub {
    name              = azurerm_eventhub.test.name
    connection_string = azurerm_eventhub_namespace.test.default_primary_connection_string
  }
}
`, r.basicEventHub(data))
}

func (ApiManagementLoggerResource) basicApplicationInsights(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  application_insights {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementLoggerResource) applicationInsightsConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  application_insights {
    connection_string = azurerm_application_insights.test.connection_string
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementLoggerResource) complete(data acceptance.TestData, description, buffered string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "other"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}

resource "azurerm_api_management_logger" "test" {
  name                = "acctestapimnglogger-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  description         = "%s"
  buffered            = %s
  resource_id         = azurerm_application_insights.test.id

  application_insights {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, description, buffered)
}
