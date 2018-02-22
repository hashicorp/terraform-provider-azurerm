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

func TestAccAzureRMContainerGroup_linuxBasic(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_linuxBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxBasicUpdate(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_linuxBasic(ri, testLocation())
	updatedConfig := testAccAzureRMContainerGroup_linuxBasicUpdated(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxComplete(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_linuxComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.command", "/bin/bash -c ls"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.mount_path", "/aci/logs"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.name", "logs"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "restart_policy", "OnFailure"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_windowsBasic(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_windowsBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_windowsComplete(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_windowsComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.command", "cmd.exe echo hi"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "restart_policy", "Never"),
				),
			},
		},
	})
}

func testAccAzureRMContainerGroup_linuxBasic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
	port   = "80"
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_linuxBasicUpdated(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
	port   = "80"
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_windowsBasic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "windows"

  container {
    name   = "windowsservercore"
    image  = "microsoft/windowsservercore:latest"
    cpu    = "2.0"
    memory = "3.5"
    port   = "80"
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_windowsComplete(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  dns_name_label      = "acctestcontainergroup-%d"
  os_type             = "windows"
  restart_policy      = "Never"

  container {
    name   = "windowsservercore"
    image  = "microsoft/windowsservercore:latest"
    cpu    = "2.0"
    memory = "3.5"
    port   = "80"

	environment_variables {
		"foo"  = "bar"
		"foo1" = "bar1"
	}
	command = "cmd.exe echo hi"
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri, ri)
}

func testAccAzureRMContainerGroup_linuxComplete(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
	name = "acctestss-%d"

	resource_group_name  = "${azurerm_resource_group.test.name}"
	storage_account_name = "${azurerm_storage_account.test.name}"

	quota = 50
}

resource "azurerm_container_group" "test" {
	name                = "acctestcontainergroup-%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	ip_address_type     = "public"
	dns_name_label      = "acctestcontainergroup-%d"
	os_type             = "linux"
	restart_policy      = "OnFailure"

	container {
		name   = "hf"
		image  = "seanmckenna/aci-hellofiles"
		cpu    = "1"
		memory = "1.5"

		port     = "80"
		protocol = "TCP"

		volume {
			name       = "logs"
			mount_path = "/aci/logs"
			read_only  = false
			share_name = "${azurerm_storage_share.test.name}"

			storage_account_name = "${azurerm_storage_account.test.name}"
			storage_account_key = "${azurerm_storage_account.test.primary_access_key}"
		}

		environment_variables {
			"foo" = "bar"
			"foo1" = "bar1"
		}

		command = "/bin/bash -c ls"
	}

	tags {
		environment = "Testing"
	}
}
`, ri, location, ri, ri, ri, ri)
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
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
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
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Container Group still exists:\n%#v", resp)
			}

			return nil
		}

	}

	return nil
}
