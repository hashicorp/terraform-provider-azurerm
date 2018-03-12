package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMServiceBusTopicAuthorizationRule_importListen(t *testing.T) {
	resourceName := "azurerm_servicebus_topic_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_listen(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
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

func TestAccAzureRMServiceBusTopicAuthorizationRule_importSend(t *testing.T) {
	resourceName := "azurerm_servicebus_topic_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_send(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
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

func TestAccAzureRMServiceBusTopicAuthorizationRule_importReadWrite(t *testing.T) {
	resourceName := "azurerm_servicebus_topic_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_readWrite(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
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

func TestAccAzureRMServiceBusTopicAuthorizationRule_importManage(t *testing.T) {
	resourceName := "azurerm_servicebus_topic_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusTopicAuthorizationRule_manage(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusTopicAuthorizationRuleDestroy,
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
