package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMNotificationHubAuthorizationRule_listen(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_listen(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "manage", "false"),
					resource.TestCheckResourceAttr(resourceName, "send", "false"),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_manage(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_manage(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "manage", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "true"),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_send(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_send(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "manage", "false"),
					resource.TestCheckResourceAttr(resourceName, "send", "true"),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_updated(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_listen(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "manage", "false"),
					resource.TestCheckResourceAttr(resourceName, "send", "false"),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
				),
			},
			{
				Config: testAzureRMNotificationHubAuthorizationRule_manage(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "manage", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "true"),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testCheckAzureRMNotificationHubAuthorizationRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).notificationHubsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		notificationHubName := rs.Primary.Attributes["notification_hub_name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		ruleName := rs.Primary.Attributes["name"]

		resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, ruleName)
		if err != nil {
			return fmt.Errorf("Bad: Get on notificationHubsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Notification Hub Authorization Rule does not exist: %s", name)
		}

		return nil
	}
}

func testCheckAzureRMNotificationHubAuthorizationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).notificationHubsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAzureRMNotificationHubAuthorizationRule_listen(ri int, location string) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(ri, location)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = "${azurerm_notification_hub.test.name}"
  namespace_name        = "${azurerm_notification_hub_namespace.test.name}" 
  resource_group_name   = "${azurerm_resource_group.test.name}"
  listen                = true
}
`, template, ri)
}

func testAzureRMNotificationHubAuthorizationRule_send(ri int, location string) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(ri, location)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = "${azurerm_notification_hub.test.name}"
  namespace_name        = "${azurerm_notification_hub_namespace.test.name}" 
  resource_group_name   = "${azurerm_resource_group.test.name}"
  send                  = true
  listen                = true
}
`, template, ri)
}

func testAzureRMNotificationHubAuthorizationRule_manage(ri int, location string) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(ri, location)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "acctestrule-%d"
  notification_hub_name = "${azurerm_notification_hub.test.name}"
  namespace_name        = "${azurerm_notification_hub_namespace.test.name}" 
  resource_group_name   = "${azurerm_resource_group.test.name}"
  manage                = true
  send                  = true
  listen                = true
}
`, template, ri)
}

func testAzureRMNotificationHubAuthorizationRule_template(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"

  sku {
    name = "Free"
  }
}

resource "azurerm_notification_hub" "test" {
  name                = "acctestnh-%d"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}" 
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, ri, location, ri, ri)
}
