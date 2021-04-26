package eventhub_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type EventHubNamespaceAuthorizationRuleResource struct {
}

func TestAccEventHubNamespaceAuthorizationRule_listen(t *testing.T) {
	testAccEventHubNamespaceAuthorizationRule(t, true, false, false)
}

func TestAccEventHubNamespaceAuthorizationRule_send(t *testing.T) {
	testAccEventHubNamespaceAuthorizationRule(t, false, true, false)
}

func TestAccEventHubNamespaceAuthorizationRule_listensend(t *testing.T) {
	testAccEventHubNamespaceAuthorizationRule(t, true, true, false)
}

func TestAccEventHubNamespaceAuthorizationRule_manage(t *testing.T) {
	testAccEventHubNamespaceAuthorizationRule(t, true, true, true)
}

func testAccEventHubNamespaceAuthorizationRule(t *testing.T, listen, send, manage bool) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test")
	r := EventHubNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base(data, listen, send, manage),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("namespace_name").Exists(),
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

func TestAccEventHubNamespaceAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test")
	r := EventHubNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base(data, true, true, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, true, true, true),
			ExpectError: acceptance.RequiresImportError("azurerm_eventhub_namespace_authorization_rule"),
		},
	})
}

func TestAccEventHubNamespaceAuthorizationRule_rightsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test")
	r := EventHubNamespaceAuthorizationRuleResource{}

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

func TestAccEventHubNamespaceAuthorizationRule_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test")
	r := EventHubNamespaceAuthorizationRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// `primary_connection_string_alias` and `secondary_connection_string_alias` are still `nil` in `azurerm_eventhub_namespace_authorization_rule` after created `azurerm_eventhub_namespace` successfully since `azurerm_eventhub_namespace_disaster_recovery_config` hasn't been created.
			// So these two properties should be checked in the second run.
			// And `depends_on` cannot be applied to `azurerm_eventhub_namespace_authorization_rule`.
			// Because it would throw error message `BreakPairing operation is only allowed on primary namespace with valid secondary namespace.` while destroying `azurerm_eventhub_namespace_disaster_recovery_config` if `depends_on` is applied.
			Config: r.withAliasConnectionString(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withAliasConnectionString(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string_alias").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespaceAuthorizationRule_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_authorization_rule", "test1")
	r := EventHubNamespaceAuthorizationRuleResource{}
	resourceTwoName := "azurerm_eventhub_namespace_authorization_rule.test2"
	resourceThreeName := "azurerm_eventhub_namespace_authorization_rule.test3"

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
		data.ImportStepFor(resourceTwoName),
		data.ImportStepFor(resourceThreeName),
	})
}

func (EventHubNamespaceAuthorizationRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NamespaceAuthorizationRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Eventhub.NamespacesClient.GetAuthorizationRule(ctx, id.ResourceGroup, id.NamespaceName, id.AuthorizationRuleName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.AuthorizationRuleProperties != nil), nil
}

func (EventHubNamespaceAuthorizationRuleResource) base(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "acctest-EHN-AR%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func (EventHubNamespaceAuthorizationRuleResource) withAliasConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ehnar-%[1]d-1"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-ehnar-%[1]d-2"
  location = "%[3]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "test2" {
  name                = "acctesteventhubnamespace2-%[1]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "acctest-EHN-DRC-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  namespace_name       = azurerm_eventhub_namespace.test.name
  partner_namespace_id = azurerm_eventhub_namespace.test2.id
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "acctest-EHN-AR%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = true
  manage = true
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (EventHubNamespaceAuthorizationRuleResource) requiresImport(data acceptance.TestData, listen, send, manage bool) string {
	template := EventHubNamespaceAuthorizationRuleResource{}.base(data, listen, send, manage)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_authorization_rule" "import" {
  name                = azurerm_eventhub_namespace_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace_authorization_rule.test.namespace_name
  resource_group_name = azurerm_eventhub_namespace_authorization_rule.test.resource_group_name
  listen              = azurerm_eventhub_namespace_authorization_rule.test.listen
  send                = azurerm_eventhub_namespace_authorization_rule.test.send
  manage              = azurerm_eventhub_namespace_authorization_rule.test.manage
}
`, template)
}

func (EventHubNamespaceAuthorizationRuleResource) multi(data acceptance.TestData, listen, send, manage bool) string {
	template := EventHubNamespaceAuthorizationRuleResource{}.base(data, listen, send, manage)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_authorization_rule" "test1" {
  name                = "acctestruleone-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  send   = true
  listen = true
  manage = false
}

resource "azurerm_eventhub_namespace_authorization_rule" "test2" {
  name                = "acctestruletwo-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  send   = true
  listen = true
  manage = false
}

resource "azurerm_eventhub_namespace_authorization_rule" "test3" {
  name                = "acctestrulethree-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  send   = true
  listen = true
  manage = false
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
