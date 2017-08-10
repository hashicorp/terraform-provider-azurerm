package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAutomationSchedule_importScheduleOneTime(t *testing.T) {
	resourceName := "azurerm_automation_schedule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAutomationSchedule_oneTime(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
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
