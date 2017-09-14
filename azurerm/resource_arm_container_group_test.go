package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerGroup_basicLinux(t *testing.T) {
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroupBasicLinux(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists("azurerm_container_group.test"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_basicWindows(t *testing.T) {
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroupBasicWindows(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists("azurerm_container_group.test"),
				),
			},
		},
	})
}

func testAccAzureRMContainerGroupBasicLinux(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
	
	name = "acctestcontainergroup-%d"
	location = "%s"
	resource_group_name = "${azurerm_resource_group.test.name}"
	ip_address_type="public"
	os_type = "linux"
  
	container {
        name = "hw"
        image = "microsoft/aci-helloworld:latest"
        cpu ="0.5"
        memory =  "1.5"
        port = "80"
    }
    container {
        name = "sidecar"
        image = "microsoft/aci-tutorial-sidecar"
        cpu="0.5"
        memory="1.5"
    }
  
	tags {
	  environment = "testing"
	}
  }
`, ri, location, ri, location)
}

func testAccAzureRMContainerGroupBasicWindows(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
	
	name = "acctestcontainergroup-%d"
	location = "%s"
	resource_group_name = "${azurerm_resource_group.test.name}"
	ip_address_type="public"
	os_type = "windows"
  
	container {
	  name = "winapp"
	  image = "winappimage:latest"
	  cpu ="2.0"
	  memory = "3.5"
	  port = "80"
	}
  
	tags {
	  environment = "testing"
	}
  }
`, ri, location, ri, location)
}

func testCheckAzureRMContainerGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Registry: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).containerGroupsClient

		resp, err := conn.Get(resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Container Group %q (resource group: %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on containerGroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMContainerGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).containerGroupsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Container Group still exists:\n%#v", resp)
			}

			return nil
		}

	}

	return nil
}
