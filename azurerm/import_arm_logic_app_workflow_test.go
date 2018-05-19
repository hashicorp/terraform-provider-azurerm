package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogicAppWorkflow_importEmpty(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppWorkflow_empty(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorklowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppWorkflow_importTags(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppWorkflow_tags(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorklowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogicAppWorkflow_importActionHTTP(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppWorkflow_actionHTTP(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorklowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
