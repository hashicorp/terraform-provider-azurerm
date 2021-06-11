package servicebus_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceBusQueueAuthorizationRuleResource struct {
}

func TestAccServiceBusQueueAuthorizationRule_listen(t *testing.T) {
	testAccServiceBusQueueAuthorizationRule(t, true, false, false)
}

func TestAccServiceBusQueueAuthorizationRule_send(t *testing.T) {
	testAccServiceBusQueueAuthorizationRule(t, false, true, false)
}

func TestAccServiceBusQueueAuthorizationRule_listensend(t *testing.T) {
	testAccServiceBusQueueAuthorizationRule(t, true, true, false)
}

func TestAccServiceBusQueueAuthorizationRule_manage(t *testing.T) {
	testAccServiceBusQueueAuthorizationRule(t, true, true, true)
}

func testAccServiceBusQueueAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data, listen, send, manage),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("queue_name").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("listen").HasValue(strconv.FormatBool(listen)),
				check.That(data.ResourceName).Key("send").HasValue(strconv.FormatBool(send)),
				check.That(data.ResourceName).Key("manage").HasValue(strconv.FormatBool(manage)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueueAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data, true, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
			),
		},
		{
			Config: r.base(data, true, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("true"),
				check.That(data.ResourceName).Key("manage").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusQueueAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue_authorization_rule", "test")
	r := ServiceBusQueueAuthorizationRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data, true, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
			),
		},
		{
			Config:      r.requiresImport(data, true, false, false),
			ExpectError: acceptance.RequiresImportError("azurerm_servicebus_queue_authorization_rule"),
		},
	})
}

func (t ServiceBusQueueAuthorizationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.QueueAuthorizationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.QueuesClient.GetAuthorizationRule(ctx, id.ResourceGroup, id.NamespaceName, id.QueueName, id.AuthorizationRuleName)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus NameSpace Queue Authorization Rule (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ServiceBusQueueAuthorizationRuleResource) base(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_servicebus_namespace.test.name

  enable_partitioning = true
}

resource "azurerm_servicebus_queue_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  queue_name          = azurerm_servicebus_queue.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func (r ServiceBusQueueAuthorizationRuleResource) requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_queue_authorization_rule" "import" {
  name                = azurerm_servicebus_queue_authorization_rule.test.name
  namespace_name      = azurerm_servicebus_queue_authorization_rule.test.namespace_name
  queue_name          = azurerm_servicebus_queue_authorization_rule.test.queue_name
  resource_group_name = azurerm_servicebus_queue_authorization_rule.test.resource_group_name

  listen = azurerm_servicebus_queue_authorization_rule.test.listen
  send   = azurerm_servicebus_queue_authorization_rule.test.send
  manage = azurerm_servicebus_queue_authorization_rule.test.manage
}
`, r.base(data, listen, send, manage))
}
