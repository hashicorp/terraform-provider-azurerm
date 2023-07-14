// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IotHubEndpointServiceBusTopicResource struct{}

func TestAccIotHubEndpointServiceBusTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")
	r := IotHubEndpointServiceBusTopicResource{}

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

func TestAccIotHubEndpointServiceBusTopic_IotHubIdAndTwoResourceGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")
	r := IotHubEndpointServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIotHubIdAndTwoResourceGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubEndpointServiceBusTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")
	r := IotHubEndpointServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iothub_endpoint_servicebus_topic"),
		},
	})
}

func TestAccIotHubEndpointServiceBusTopic_AuthenticationTypeSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")
	r := IotHubEndpointServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTypeSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubEndpointServiceBusTopic_AuthenticationTypeUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")
	r := IotHubEndpointServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubEndpointServiceBusTopic_AuthenticationTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_servicebus_topic", "test")
	r := IotHubEndpointServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.authenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.authenticationTypeSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.authenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (IotHubEndpointServiceBusTopicResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name         = "acctestservicebustopic-%[1]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name     = "acctest-%[1]d"
  topic_id = azurerm_servicebus_topic.test.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_servicebus_topic" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  connection_string = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r IotHubEndpointServiceBusTopicResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_servicebus_topic" "import" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  connection_string = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
}
`, r.basic(data))
}

func (IotHubEndpointServiceBusTopicResource) withIotHubIdAndTwoResourceGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name         = "acctestservicebustopic-%[1]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name     = "acctest-%[1]d"
  topic_id = azurerm_servicebus_topic.test.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_servicebus_topic" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctest"
  iothub_id           = azurerm_iothub.test.id

  connection_string = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
}
`, data.RandomInteger, data.Locations.Primary)
}

func (t IotHubEndpointServiceBusTopicResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.EndpointServiceBusTopicID(state.ID)
	if err != nil {
		return nil, err
	}

	iothub, err := clients.IoTHub.ResourceClient.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil || iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	if endpoints := iothub.Properties.Routing.Endpoints.ServiceBusTopics; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, id.EndpointName) {
					return utils.Bool(true), nil
				}
			}
		}
	}

	return utils.Bool(false), nil
}

func (r IotHubEndpointServiceBusTopicResource) authenticationTypeDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_servicebus_topic" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  connection_string = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointServiceBusTopicResource) authenticationTypeSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_servicebus_topic" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  authentication_type = "identityBased"
  endpoint_uri        = "sb://${azurerm_servicebus_namespace.test.name}.servicebus.windows.net"
  entity_path         = azurerm_servicebus_topic.test.name

  depends_on = [
    azurerm_role_assignment.test_azure_service_bus_data_sender_system,
  ]
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointServiceBusTopicResource) authenticationTypeUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_servicebus_topic" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"


  authentication_type = "identityBased"
  identity_id         = azurerm_user_assigned_identity.test.id
  endpoint_uri        = "sb://${azurerm_servicebus_namespace.test.name}.servicebus.windows.net"
  entity_path         = azurerm_servicebus_topic.test.name
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointServiceBusTopicResource) authenticationTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name         = "acctestservicebustopic-%[1]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name     = "acctest-%[1]d"
  topic_id = azurerm_servicebus_topic.test.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test_azure_service_bus_data_sender_user" {
  role_definition_name = "Azure Service Bus Data Sender"
  scope                = azurerm_servicebus_topic.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_azure_service_bus_data_sender_user,
  ]
}

resource "azurerm_role_assignment" "test_azure_service_bus_data_sender_system" {
  role_definition_name = "Azure Service Bus Data Sender"
  scope                = azurerm_servicebus_topic.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary)
}
