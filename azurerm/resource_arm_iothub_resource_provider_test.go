package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMIotHub_basicStandard(t *testing.T) {
	name := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_basicStandard(name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists("azurerm_iothub.test"),
				),
			},
		},
	})

}

func testCheckAzureRMIotHubDestroy(s *terraform.State) error {
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	conn := testAccProvider.Meta().(*ArmClient).expressRouteCircuitClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("IotHub still exists:\n%#v", resp.ExpressRouteCircuitPropertiesFormat)
		}
	}
	return nil
}
func testCheckAzureRMIotHubExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		iothubName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IotHub: %s", iothubName)
		}

		conn := testAccProvider.Meta().(*ArmClient).iothubResourceClient

		resp, err := conn.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: IotHub %q (resource group: %q) does not exist", iothubName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}

		return nil

	}
}

func testAccAzureRMIotHub_basicStandard(name string) string {
	return fmt.Sprintf(`
resource "azurerm_iothub" "test" {
	name                             = "%s"
	resource_group_name              = "shelley.bess"
	location                         = "eastus"
		sku {
			name = "S1"
			tier = "Standard"
			capacity = "1"
		}

		tags {
			"purpose" = "testing"
		}
}
`, name)
}
