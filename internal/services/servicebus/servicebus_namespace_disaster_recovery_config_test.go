package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceBusNamespaceDisasterRecoveryConfigResource struct {
}

func TestAccAzureRMServiceBusNamespacePairing_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_disaster_recovery_config", "pairing_test")
	r := ServiceBusNamespaceDisasterRecoveryConfigResource{}
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

func (t ServiceBusNamespaceDisasterRecoveryConfigResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NamespaceDisasterRecoveryConfigID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.DisasterRecoveryConfigsClient.Get(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus NameSpace (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ServiceBusNamespaceDisasterRecoveryConfigResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "primary" {
  name     = "acctest1RG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "secondary" {
  name     = "acctest2RG-%[1]d"
  location = "%[3]s"
}

resource "azurerm_servicebus_namespace" "primary_namespace_test" {
  name                = "acctest1-%[1]d"
  location            = azurerm_resource_group.primary.location
  resource_group_name = azurerm_resource_group.primary.name
  sku                 = "Premium"
  capacity            = "1"
}

resource "azurerm_servicebus_topic" "example" {
  name                = "topic-test"
  resource_group_name = azurerm_resource_group.primary.name
  namespace_name      = azurerm_servicebus_namespace.primary_namespace_test.name
}

resource "azurerm_servicebus_queue" "example" {
  name                = "queue-test"
  resource_group_name = azurerm_resource_group.primary.name
  namespace_name      = azurerm_servicebus_namespace.primary_namespace_test.name
}

resource "azurerm_servicebus_namespace_authorization_rule" "example" {
  name                = "example_namespace_rule"
  namespace_name      = azurerm_servicebus_namespace.primary_namespace_test.name
  resource_group_name = azurerm_resource_group.primary.name
  manage              = true
  listen              = true
  send                = true
}

resource "azurerm_servicebus_queue_authorization_rule" "example" {
  name                = "example_queue_rule"
  namespace_name      = azurerm_servicebus_namespace.primary_namespace_test.name
  queue_name          = azurerm_servicebus_queue.example.name
  resource_group_name = azurerm_resource_group.primary.name
  manage              = true
  listen              = true
  send                = true
}

resource "azurerm_servicebus_topic_authorization_rule" "example" {
  name                = "example_topic_rule"
  namespace_name      = azurerm_servicebus_namespace.primary_namespace_test.name
  topic_name          = azurerm_servicebus_topic.example.name
  resource_group_name = azurerm_resource_group.primary.name
  manage              = true
  listen              = true
  send                = true
}

resource "azurerm_servicebus_namespace" "secondary_namespace_test" {
  name                = "acctest2-%[1]d"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  sku                 = "Premium"
  capacity            = "1"

  depends_on = [
    azurerm_servicebus_topic_authorization_rule.example,
    azurerm_servicebus_queue_authorization_rule.example,
    azurerm_servicebus_namespace_authorization_rule.example
  ]
}

resource "azurerm_servicebus_namespace_disaster_recovery_config" "pairing_test" {
  name                 = "acctest-alias-%[1]d"
  primary_namespace_id = azurerm_servicebus_namespace.primary_namespace_test.id
  partner_namespace_id = azurerm_servicebus_namespace.secondary_namespace_test.id
}

`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
