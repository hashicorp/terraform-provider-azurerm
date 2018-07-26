package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLogicAppActionCustom_importBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMLogicAppActionCustom_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogicAppWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_logic_app_action_custom.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
