package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAzureRMDataFactoryV2_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryV2Exists("azurerm_data_factory_v2.test"),
				),
			},
		},
	})
}

func testCheckAzureRMDataFactoryV2Exists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).dataFactoryClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactoryClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Data Factory %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryV2Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).dataFactoryClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_v2" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory still exists:\n%#v", resp.FactoryProperties)
		}
	}

	return nil
}

func testAccAzureRMDataFactoryV2_basic(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }
  resource "azurerm_data_factory_v2" "test" {
    name                = "acctestdfv2%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
  }
`, rInt, location, rInt)
}
