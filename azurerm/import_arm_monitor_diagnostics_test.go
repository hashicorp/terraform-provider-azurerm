package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMMonitorDiagnostics_importBasic(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostics.test"
	ri := acctest.RandIntRange(10000, 99999)
	location := testLocation()

	config := testAccAzureRMMonitorDiagnostics_basic(ri, location)

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
