package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMMonitorDiagnostics_importBasic(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostics.test"
	ri := acctest.RandIntRange(10000, 99999)
	config := testAccAzureRMMonitorDiagnostics_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}
