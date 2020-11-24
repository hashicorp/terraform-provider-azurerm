package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMNotificationHubAuthorizationRule_listen(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_listen(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_listen(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.RequiresImportErrorStep(testAzureRMNotificationHubAuthorizationRule_requiresImport),
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_manage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_manage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_send(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_send(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_multi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test1")
	resourceTwoName := "azurerm_notification_hub_authorization_rule.test2"
	resourceThreeName := "azurerm_notification_hub_authorization_rule.test3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_multi(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceTwoName),
					resource.TestCheckResourceAttr(resourceTwoName, "manage", "false"),
					resource.TestCheckResourceAttr(resourceTwoName, "send", "true"),
					resource.TestCheckResourceAttr(resourceTwoName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceTwoName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceTwoName, "secondary_access_key"),
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceThreeName),
					resource.TestCheckResourceAttr(resourceThreeName, "manage", "false"),
					resource.TestCheckResourceAttr(resourceThreeName, "send", "true"),
					resource.TestCheckResourceAttr(resourceThreeName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceThreeName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceThreeName, "secondary_access_key"),
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
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_updated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_authorization_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_listen(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			{
				Config: testAzureRMNotificationHubAuthorizationRule_manage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "manage", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "send", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).NotificationHubs.HubsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		notificationHubName := rs.Primary.Attributes["notification_hub_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		ruleName := rs.Primary.Attributes["name"]

		resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, ruleName)
		if err != nil {
			return fmt.Errorf("Bad: Get on notificationHubsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Notification Hub Authorization Rule does not exist: %s", ruleName)
		}

		return nil
	}
}

func testCheckAzureRMNotificationHubAuthorizationRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).NotificationHubs.HubsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_notification_hub_authorization_rule" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		notificationHubName := rs.Primary.Attributes["notification_hub_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		ruleName := rs.Primary.Attributes["name"]
		resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, ruleName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Notification Hub Authorization Rule still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMNotificationHubAuthorizationRule_listen(data acceptance.TestData) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  listen                = true
}
`, template, data.RandomInteger)
}

func testAzureRMNotificationHubAuthorizationRule_requiresImport(data acceptance.TestData) string {
	template := testAzureRMNotificationHubAuthorizationRule_listen(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "import" {
  name                  = azurerm_notification_hub_authorization_rule.test.name
  notification_hub_name = azurerm_notification_hub_authorization_rule.test.notification_hub_name
  namespace_name        = azurerm_notification_hub_authorization_rule.test.namespace_name
  resource_group_name   = azurerm_notification_hub_authorization_rule.test.resource_group_name
  listen                = azurerm_notification_hub_authorization_rule.test.listen
}
`, template)
}

func testAzureRMNotificationHubAuthorizationRule_send(data acceptance.TestData) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}
`, template, data.RandomInteger)
}

func testAzureRMNotificationHubAuthorizationRule_multi(data acceptance.TestData) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test1" {
  name                  = "acctestruleone-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}

resource "azurerm_notification_hub_authorization_rule" "test2" {
  name                  = "acctestruletwo-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}

resource "azurerm_notification_hub_authorization_rule" "test3" {
  name                  = "acctestrulethree-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  send                  = true
  listen                = true
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAzureRMNotificationHubAuthorizationRule_manage(data acceptance.TestData) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = azurerm_notification_hub.test.name
  namespace_name        = azurerm_notification_hub_namespace.test.name
  resource_group_name   = azurerm_resource_group.test.name
  manage                = true
  send                  = true
  listen                = true
}
`, template, data.RandomInteger)
}

func testAzureRMNotificationHubAuthorizationRule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  namespace_type      = "NotificationHub"
  sku_name            = "Free"
}

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = azurerm_notification_hub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
