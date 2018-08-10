package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogicAppTriggerHttpRequest_importBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerHttpRequest_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_http_request.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerHttpRequest_importFullSchema(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerHttpRequest_fullSchema(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_http_request.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerHttpRequest_importMethod(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerHttpRequest_method(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_http_request.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppTriggerHttpRequest_importRelativePath(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerHttpRequest_relativePath(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_http_request.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
