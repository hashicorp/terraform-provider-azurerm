package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMEventHubAuthorizationRule_importListen(t *testing.T) {
	resourceName := "azurerm_eventhub_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMEventHubAuthorizationRule_listen(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubAuthorizationRuleDestroy,
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

func TestAccAzureRMEventHubAuthorizationRule_importSend(t *testing.T) {
	resourceName := "azurerm_eventhub_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMEventHubAuthorizationRule_send(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubAuthorizationRuleDestroy,
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

func TestAccAzureRMEventHubAuthorizationRule_importReadWrite(t *testing.T) {
	resourceName := "azurerm_eventhub_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMEventHubAuthorizationRule_readWrite(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubAuthorizationRuleDestroy,
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

func TestAccAzureRMEventHubAuthorizationRule_importManage(t *testing.T) {
	resourceName := "azurerm_eventhub_authorization_rule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMEventHubAuthorizationRule_manage(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubAuthorizationRuleDestroy,
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
