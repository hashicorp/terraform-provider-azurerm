package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMAutomationCredential_importCredential(t *testing.T) {
	resourceName := "azurerm_automation_credential.test"

	ri := acctest.RandInt()
	config := testAccAzureRMAutomationCredential_testCredential(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationCredentialDestroy,
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
