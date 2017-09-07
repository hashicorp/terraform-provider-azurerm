package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMContainerGroup_basicLinux(t *testing.T) {
	ri := acctest.RandInt()

	name := fmt.Sprintf("acctestcontainergroup-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

	config := testAccAzureRMContainerGroupBasicLinux(name, resourceGroupName, testLocation())

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

	name := fmt.Sprintf("acctestcontainergroup-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

	config := testAccAzureRMContainerGroupBasicWindows(name, resourceGroupName, testLocation())

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

func testAccAzureRMContainerGroupBasicLinux(name string, resourceGroupName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_container_group" "test" {
	
	name = "%s"
	location = "%s"
	resource_group_name = "${azurerm_resource_group.test.name}"
	ip_address_type="public"
	os_type = "linux"
  
	container {
	  name = "%s"
	  image = "nginx:latest"
	  cpu ="1"
	  memory = "1.0"
	  port = "80"
#	  protocol = "TCP"
	}
  
	tags {
	  environment = "testing"
	}
  }
`, resourceGroupName, location, name, location, name)
}

func testAccAzureRMContainerGroupBasicWindows(name string, resourceGroupName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_container_group" "test" {
	
	name = "%s"
	location = "%s"
	resource_group_name = "${azurerm_resource_group.test.name}"
	ip_address_type="public"
	os_type = "windows"
  
	container {
	  name = "%s"
	  image = "winappimage:latest"
	  cpu ="2.0"
	  memory = "3.5"
	  port = "80"
#	  protocol = "TCP"
	}
  
	tags {
	  environment = "testing"
	}
  }
`, resourceGroupName, location, name, location, name)
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
			return fmt.Errorf("Bad: Get on containerGroupsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Container Group %q (resource group: %q) does not exist", name, resourceGroup)
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
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Container Group still exists:\n%#v", resp)
		}
	}

	return nil
}
