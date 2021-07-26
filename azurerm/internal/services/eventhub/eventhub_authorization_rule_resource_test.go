package eventhub_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type EventHubAuthorizationRuleResource struct {
}

func TestAccEventHubAuthorizationRule_listen(t *testing.T) {
	testAccEventHubAuthorizationRule(t, true, false, false)
}

func TestAccEventHubAuthorizationRule_send(t *testing.T) {
	testAccEventHubAuthorizationRule(t, false, true, false)
}

func TestAccEventHubAuthorizationRule_listensend(t *testing.T) {
	testAccEventHubAuthorizationRule(t, true, true, false)
}

func TestAccEventHubAuthorizationRule_manage(t *testing.T) {
	testAccEventHubAuthorizationRule(t, true, true, true)
}

func testAccEventHubAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_authorization_rule", "test")
	r := EventHubAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base(data, listen, send, manage),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
				check.That(data.ResourceName).Key("eventhub_name").Exists(),
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

func TestAccEventHubAuthorizationRule_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_authorization_rule", "test1")
	r := EventHubAuthorizationRuleResource{}
	resourceTwoName := "azurerm_eventhub_authorization_rule.test2"
	resourceThreeName := "azurerm_eventhub_authorization_rule.test3"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multi(data, true, true, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("manage").HasValue("false"),
				check.That(data.ResourceName).Key("send").HasValue("true"),
				check.That(data.ResourceName).Key("listen").HasValue("true"),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(resourceTwoName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(resourceTwoName, "manage", "false"),
				resource.TestCheckResourceAttr(resourceTwoName, "send", "true"),
				resource.TestCheckResourceAttr(resourceTwoName, "listen", "true"),
				resource.TestCheckResourceAttrSet(resourceTwoName, "primary_connection_string"),
				resource.TestCheckResourceAttrSet(resourceTwoName, "secondary_connection_string"),
				check.That(resourceThreeName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(resourceThreeName, "manage", "false"),
				resource.TestCheckResourceAttr(resourceThreeName, "send", "true"),
				resource.TestCheckResourceAttr(resourceThreeName, "listen", "true"),
				resource.TestCheckResourceAttrSet(resourceThreeName, "primary_connection_string"),
				resource.TestCheckResourceAttrSet(resourceThreeName, "secondary_connection_string"),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      resourceTwoName,
			ImportState:       true,
			ImportStateVerify: true,
		},
		{
			ResourceName:      resourceThreeName,
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccEventHubAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_authorization_rule", "test")
	r := EventHubAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base(data, true, true, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, true, true, true),
			ExpectError: acceptance.RequiresImportError("azurerm_eventhub_authorization_rule"),
		},
	})
}

func TestAccEventHubAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_authorization_rule", "test")
	r := EventHubAuthorizationRuleResource{}

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

func TestAccEventHubAuthorizationRule_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_authorization_rule", "test")
	r := EventHubAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withAliasConnectionString(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (EventHubAuthorizationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	name := id.Path["authorizationRules"]
	namespaceName := id.Path["namespaces"]
	eventHubName := id.Path["eventhubs"]

	resp, err := clients.Eventhub.EventHubsClient.GetAuthorizationRule(ctx, id.ResourceGroup, namespaceName, eventHubName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Event Hub Authorization Rule %q (eventhub %s / namespace %s / resource group: %s) does not exist", name, eventHubName, namespaceName, id.ResourceGroup)
	}

	return utils.Bool(resp.AuthorizationRuleProperties != nil), nil
}

func (EventHubAuthorizationRuleResource) base(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func (EventHubAuthorizationRuleResource) multi(data acceptance.TestData, listen, send, manage bool) string {
	template := EventHubAuthorizationRuleResource{}.base(data, listen, send, manage)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_authorization_rule" "test1" {
  name                = "acctestruleone-%d"
  eventhub_name       = azurerm_eventhub.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  send                = true
  listen              = true
}

resource "azurerm_eventhub_authorization_rule" "test2" {
  name                = "acctestruletwo-%d"
  eventhub_name       = azurerm_eventhub.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  send                = true
  listen              = true
}

resource "azurerm_eventhub_authorization_rule" "test3" {
  name                = "acctestrulethree-%d"
  eventhub_name       = azurerm_eventhub.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  send                = true
  listen              = true
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (EventHubAuthorizationRuleResource) requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	template := EventHubAuthorizationRuleResource{}.base(data, listen, send, manage)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_authorization_rule" "import" {
  name                = azurerm_eventhub_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_authorization_rule.test.namespace_name
  eventhub_name       = azurerm_eventhub_authorization_rule.test.eventhub_name
  resource_group_name = azurerm_eventhub_authorization_rule.test.resource_group_name
  listen              = azurerm_eventhub_authorization_rule.test.listen
  send                = azurerm_eventhub_authorization_rule.test.send
  manage              = azurerm_eventhub_authorization_rule.test.manage
}
`, template)
}

func (EventHubAuthorizationRuleResource) withAliasConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ehar-%[1]d-1"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-ehar-%[1]d-2"
  location = "%[3]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_eventhub_namespace" "test2" {
  name                = "acctesteventhubnamespace2-%[1]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name

  sku = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "acctest-EHN-DRC-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  namespace_name       = azurerm_eventhub_namespace.test.name
  partner_namespace_id = azurerm_eventhub_namespace.test2.id
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = true
  manage = true

  depends_on = [azurerm_eventhub_namespace_disaster_recovery_config.test]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
