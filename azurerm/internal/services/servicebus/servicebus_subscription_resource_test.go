package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceBusSubscriptionResource struct {
}

func TestAccServiceBusSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

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

func TestAccServiceBusSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

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

func TestAccServiceBusSubscription_defaultTtl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDefaultTtl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr("azurerm_servicebus_subscription.test", "default_message_ttl", "PT1H"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusSubscription_updateEnableBatched(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateEnableBatched(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enable_batched_operations").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusSubscription_updateRequiresSession(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateRequiresSession(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("requires_session").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusSubscription_updateForwardTo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

	expectedValue := fmt.Sprintf("acctestservicebustopic-forward_to-%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateForwardTo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("forward_to").HasValue(expectedValue),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusSubscription_updateForwardDeadLetteredMessagesTo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

	expectedValue := fmt.Sprintf("acctestservicebustopic-forward_dl_messages_to-%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateForwardDeadLetteredMessagesTo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("forward_dead_lettered_messages_to").HasValue(expectedValue),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusSubscription_updateDeadLetteringOnFilterEvaluationExceptions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateDeadLetteringOnFilterEvaluationExceptions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusSubscription_status(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription", "test")
	r := ServiceBusSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Active"),
			),
		},
		{
			Config: r.status(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Disabled"),
			),
		},
		{
			Config: r.status(data, "ReceiveDisabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("ReceiveDisabled"),
			),
		},
		{
			Config: r.status(data, "Active"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("status").HasValue("Active"),
			),
		},
	})
}

func (t ServiceBusSubscriptionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.SubscriptionsClient.Get(ctx, id.ResourceGroup, id.NamespaceName, id.TopicName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus NameSpace Subscription (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

const testAccServiceBusSubscription_tfTemplate = `
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctestservicebustopic-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_servicebus_subscription" "test" {
  name                = "acctestservicebussubscription-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  topic_name          = "${azurerm_servicebus_topic.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  max_delivery_count  = 10
	%s
}
`

func (ServiceBusSubscriptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(testAccServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, "")
}

func (r ServiceBusSubscriptionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_subscription" "import" {
  name                = azurerm_servicebus_subscription.test.name
  namespace_name      = azurerm_servicebus_subscription.test.namespace_name
  topic_name          = azurerm_servicebus_subscription.test.topic_name
  resource_group_name = azurerm_servicebus_subscription.test.resource_group_name
  max_delivery_count  = azurerm_servicebus_subscription.test.max_delivery_count
}
`, r.basic(data))
}

func (ServiceBusSubscriptionResource) withDefaultTtl(data acceptance.TestData) string {
	return fmt.Sprintf(testAccServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"default_message_ttl = \"PT1H\"\n")
}

func (ServiceBusSubscriptionResource) updateEnableBatched(data acceptance.TestData) string {
	return fmt.Sprintf(testAccServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"enable_batched_operations = true\n")
}

func (ServiceBusSubscriptionResource) updateRequiresSession(data acceptance.TestData) string {
	return fmt.Sprintf(testAccServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"requires_session = true\n")
}

func (ServiceBusSubscriptionResource) updateForwardTo(data acceptance.TestData) string {
	forwardToTf := testAccServiceBusSubscription_tfTemplate + `


resource "azurerm_servicebus_topic" "forward_to" {
  name                = "acctestservicebustopic-forward_to-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}


`
	return fmt.Sprintf(forwardToTf, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"forward_to = \"${azurerm_servicebus_topic.forward_to.name}\"\n", data.RandomInteger)
}

func (ServiceBusSubscriptionResource) updateForwardDeadLetteredMessagesTo(data acceptance.TestData) string {
	forwardToTf := testAccServiceBusSubscription_tfTemplate + `


resource "azurerm_servicebus_topic" "forward_dl_messages_to" {
  name                = "acctestservicebustopic-forward_dl_messages_to-%d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}


`
	return fmt.Sprintf(forwardToTf, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"forward_dead_lettered_messages_to = \"${azurerm_servicebus_topic.forward_dl_messages_to.name}\"\n", data.RandomInteger)
}

func (ServiceBusSubscriptionResource) status(data acceptance.TestData, status string) string {
	return fmt.Sprintf(testAccServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		fmt.Sprintf("status = \"%s\"", status))
}

func (ServiceBusSubscriptionResource) updateDeadLetteringOnFilterEvaluationExceptions(data acceptance.TestData) string {
	return fmt.Sprintf(testAccServiceBusSubscription_tfTemplate, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		"dead_lettering_on_filter_evaluation_error = false\n")
}
