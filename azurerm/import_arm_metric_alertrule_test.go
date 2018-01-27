package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMMetricAlertRule_importVirtualMachineCpu(t *testing.T) {
	resourceName := "azurerm_metric_alertrule.test"

	ri := acctest.RandInt()
	config := testAccAzureRMMetricAlertRule_virtualMachineCpu(ri, testLocation(), true)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMetricAlertRuleDestroy,
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
