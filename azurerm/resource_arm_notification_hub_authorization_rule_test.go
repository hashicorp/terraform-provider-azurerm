package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMNotificationHubAuthorizationRule_listen(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_listen(ri, testLocation()),
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

func TestAccAzureRMNotificationHubAuthorizationRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_notification_hub_authorization_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_listen(ri, testLocation()),
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
				Config:      testAzureRMNotificationHubAuthorizationRule_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_notification_hub_authorization_rule"),
			},
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_manage(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := tf.AccRandTimeInt()
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

	ri := tf.AccRandTimeInt()
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

func TestAccAzureRMNotificationHubAuthorizationRule_multi(t *testing.T) {
	resourceOneName := "azurerm_notification_hub_authorization_rule.test1"
	resourceTwoName := "azurerm_notification_hub_authorization_rule.test2"
	resourceThreeName := "azurerm_notification_hub_authorization_rule.test3"

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNotificationHubAuthorizationRule_multi(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceOneName),
					resource.TestCheckResourceAttr(resourceOneName, "manage", "false"),
					resource.TestCheckResourceAttr(resourceOneName, "send", "true"),
					resource.TestCheckResourceAttr(resourceOneName, "listen", "true"),
					resource.TestCheckResourceAttrSet(resourceOneName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceOneName, "secondary_access_key"),
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
			{
				ResourceName:      resourceOneName,
				ImportState:       true,
				ImportStateVerify: true,
			},
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
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := tf.AccRandTimeInt()
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

func testCheckAzureRMNotificationHubAuthorizationRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).notificationHubs.HubsClient
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
			return fmt.Errorf("Notification Hub Authorization Rule does not exist: %s", ruleName)
		}

		return nil
	}
}

func testCheckAzureRMNotificationHubAuthorizationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).notificationHubs.HubsClient
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

func testAzureRMNotificationHubAuthorizationRule_requiresImport(ri int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "import" {
  name                  = "${azurerm_notification_hub_authorization_rule.test.name}"
  notification_hub_name = "${azurerm_notification_hub_authorization_rule.test.notification_hub_name}"
  namespace_name        = "${azurerm_notification_hub_authorization_rule.test.namespace_name}"
  resource_group_name   = "${azurerm_notification_hub_authorization_rule.test.resource_group_name}"
  listen                = "${azurerm_notification_hub_authorization_rule.test.listen}"
}
`, testAzureRMNotificationHubAuthorizationRule_listen(ri, location))
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

func testAzureRMNotificationHubAuthorizationRule_multi(ri int, location string) string {
	template := testAzureRMNotificationHubAuthorizationRule_template(ri, location)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_authorization_rule" "test1" {
  name                  = "acctestruleone-%d"
  notification_hub_name = "${azurerm_notification_hub.test.name}"
  namespace_name        = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  send                  = true
  listen                = true
}

resource "azurerm_notification_hub_authorization_rule" "test2" {
	name                  = "acctestruletwo-%d"
	notification_hub_name = "${azurerm_notification_hub.test.name}"
	namespace_name        = "${azurerm_notification_hub_namespace.test.name}"
	resource_group_name   = "${azurerm_resource_group.test.name}"
	send                  = true
	listen                = true
}

resource "azurerm_notification_hub_authorization_rule" "test3" {
	name                  = "acctestrulethree-%d"
	notification_hub_name = "${azurerm_notification_hub.test.name}"
	namespace_name        = "${azurerm_notification_hub_namespace.test.name}"
	resource_group_name   = "${azurerm_resource_group.test.name}"
	send                  = true
	listen                = true
}

`, template, ri, ri, ri)
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
  name     = "acctestRG-%d"
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
