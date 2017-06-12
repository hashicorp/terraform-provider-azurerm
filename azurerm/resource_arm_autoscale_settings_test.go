package azurerm

import (
	"testing"
)

func TestAccAzureRMAutoscaleSettings_basic(t *testing.T) {

	// ri := acctest.RandInt()
	// config := fmt.Sprintf(testAccAzureRMVAvailabilitySet_basic, ri, ri)

	// resource.Test(t, resource.TestCase{
	// 	PreCheck:     func() { testAccPreCheck(t) },
	// 	Providers:    testAccProviders,
	// 	CheckDestroy: testCheckAzureRMAvailabilitySetDestroy,
	// 	Steps: []resource.TestStep{
	// 		resource.TestStep{
	// 			Config: config,
	// 			Check: resource.ComposeTestCheckFunc(
	// 				testCheckAzureRMAvailabilitySetExists("azurerm_availability_set.test"),
	// 				resource.TestCheckResourceAttr(
	// 					"azurerm_availability_set.test", "platform_update_domain_count", "5"),
	// 				resource.TestCheckResourceAttr(
	// 					"azurerm_availability_set.test", "platform_fault_domain_count", "3"),
	// 			),
	// 		},
	// 	},
	// })
}
