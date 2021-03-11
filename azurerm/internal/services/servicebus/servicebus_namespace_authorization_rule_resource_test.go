package servicebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceBusNamespaceAuthorizationRuleResource struct {
}

func TestAccServiceBusNamespaceAuthorizationRule_listen(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, true, false, false)
}

func TestAccServiceBusNamespaceAuthorizationRule_send(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, false, true, false)
}

func TestAccServiceBusNamespaceAuthorizationRule_listensend(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, true, true, false)
}

func TestAccServiceBusNamespaceAuthorizationRule_manage(t *testing.T) {
	testAccServiceBusNamespaceAuthorizationRule(t, true, true, true)
}

func testAccServiceBusNamespaceAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base(data, listen, send, manage),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceBusNamespaceAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base(data, true, false, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
			),
		},
		{
			Config: r.base(data, true, true, true),
			Check: resource.ComposeTestCheckFunc(
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

func TestAccServiceBusNamespaceAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace_authorization_rule", "test")
	r := ServiceBusNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base(data, true, false, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("send").HasValue("false"),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
			),
		},
		{
			Config:      r.requiresImport(data, true, false, false),
			ExpectError: acceptance.RequiresImportError("azurerm_servicebus_namespace_authorization_rule"),
		},
	})
}

func (t ServiceBusNamespaceAuthorizationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NamespaceAuthorizationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.NamespacesClient.GetAuthorizationRule(ctx, id.ResourceGroup, id.NamespaceName, id.AuthorizationRuleName)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus Name Space Authorization Rule (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ServiceBusNamespaceAuthorizationRuleResource) base(data acceptance.TestData, listen, send, manage bool) string {
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

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func (r ServiceBusNamespaceAuthorizationRuleResource) requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace_authorization_rule" "import" {
  name                = azurerm_servicebus_namespace_authorization_rule.test.name
  namespace_name      = azurerm_servicebus_namespace_authorization_rule.test.namespace_name
  resource_group_name = azurerm_servicebus_namespace_authorization_rule.test.resource_group_name

  listen = azurerm_servicebus_namespace_authorization_rule.test.listen
  send   = azurerm_servicebus_namespace_authorization_rule.test.send
  manage = azurerm_servicebus_namespace_authorization_rule.test.manage
}
`, r.base(data, listen, send, manage))
}
