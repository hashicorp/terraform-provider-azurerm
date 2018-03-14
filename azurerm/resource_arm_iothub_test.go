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
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_basicStandard(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists("azurerm_iothub.test"),
				),
			},
		},
	})

}

func testCheckAzureRMIotHubDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).iothubResourceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("IotHub %s still exists in resource group %s", name, resourceGroup)
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

		client := testAccProvider.Meta().(*ArmClient).iothubResourceClient
		resp, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: IotHub %q (resource group: %q) does not exist", iothubName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}

		return nil

	}
}

func testAccAzureRMIotHub_basicStandard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "foo" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.foo.name}"
  location            = "${azurerm_resource_group.foo.location}"
  sku {
    name = "S1"
    tier = "Standard"
    capacity = "1"
  }

  tags {
    "purpose" = "testing"
  }
}
`, rInt, location, rInt)
}
