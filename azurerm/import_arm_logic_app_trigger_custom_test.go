package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogicAppTriggerCustom_importBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppTriggerCustom_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_trigger_custom.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
