package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAutomationRunbook_importRunbookPSWorkflow(t *testing.T) {
	resourceName := "azurerm_automation_runbook.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAutomationRunbook_PSWorkflow(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				// publish content link is not returned after the runbook is created
				ExpectNonEmptyPlan: true,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
