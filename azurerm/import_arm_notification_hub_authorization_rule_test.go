package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMNotificationHubAuthorizationRule_importListen(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := acctest.RandInt()
	location := testLocation()
	config := testAzureRMNotificationHubAuthorizationRule_listen(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_importManage(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := acctest.RandInt()
	location := testLocation()
	config := testAzureRMNotificationHubAuthorizationRule_manage(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNotificationHubAuthorizationRule_importSend(t *testing.T) {
	resourceName := "azurerm_notification_hub_authorization_rule.test"

	ri := acctest.RandInt()
	location := testLocation()
	config := testAzureRMNotificationHubAuthorizationRule_send(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNotificationHubAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
