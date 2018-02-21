package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: an update test going from External -> Internal -> External

func TestAccAzureRMAppServiceEnvironment_basic(t *testing.T) {
	resourceName := "azurerm_app_service_environment.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceEnvironment_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "number_of_ip_addresses", "1"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_internal(t *testing.T) {
	resourceName := "azurerm_app_service_environment.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceEnvironment_internal(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "number_of_ip_addresses", "0"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_network.0.virtual_network_id"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network.0.subnet_name", fmt.Sprintf("acctestsn-%d", ri)),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_tags(t *testing.T) {
	resourceName := "azurerm_app_service_environment.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAppServiceEnvironment_tags(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Source", "AcceptanceTests"),
				),
			},
		},
	})
}

func testCheckAzureRMAppServiceEnvironmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServiceEnvironmentsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_environment" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAppServiceEnvironmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		appServiceEnvironmentName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Environment: %s", appServiceEnvironmentName)
		}

		client := testAccProvider.Meta().(*ArmClient).appServiceEnvironmentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, appServiceEnvironmentName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Environment %q (resource group: %q) does not exist", appServiceEnvironmentName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServiceEnvironmentsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAppServiceEnvironment_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestrg-%d"
  location = "%s"

  tags {
    "TestName" = "testAccAzureRMAppServiceEnvironment_basic"
  }
}

resource "azurerm_app_service_environment" "test" {
  name = "acctestase-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  number_of_ip_addresses = 1

  frontend_pool {
    vm_size = "Medium"
    number_of_workers = 2
  }

  worker_pool {
    worker_size_id = 0
    worker_size = "Medium"
    worker_count = 2
  }

  worker_pool {
    worker_size_id = 1
    worker_size = "Small"
    worker_count = 3
  }

  worker_pool {
    worker_size_id = 2
    worker_size = "Small"
    worker_count = 4
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServiceEnvironment_internal(rInt int, location string) string {
	// TODO: needs a virtual network
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestrg-%d"
  location = "%s"

  tags {
    "TestName" = "testAccAzureRMAppServiceEnvironment_internal"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_app_service_environment" "test" {
  name = "acctestase-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  number_of_ip_addresses = 0

  virtual_network {
    virtual_network_id = "${azurerm_virtual_network.test.id}"
    subnet_name = "${azurerm_subnet.test.name}"
  }

  frontend_pool {
    vm_size = "Medium"
    number_of_workers = 2
  }

  worker_pool {
    worker_size_id = 0
    worker_size = "Medium"
    worker_count = 2
  }

  worker_pool {
    worker_size_id = 1
    worker_size = "Small"
    worker_count = 3
  }

  worker_pool {
    worker_size_id = 2
    worker_size = "Small"
    worker_count = 4
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceEnvironment_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestrg-%d"
  location = "%s"

  tags {
    "TestName" = "testAccAzureRMAppServiceEnvironment_tags"
  }
}

resource "azurerm_app_service_environment" "test" {
  name = "acctestase-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  number_of_ip_addresses = 1

  frontend_pool {
    vm_size = "Medium"
    number_of_workers = 2
  }

  worker_pool {
    worker_size_id = 0
    worker_size = "Medium"
    worker_count = 2
  }

  worker_pool {
    worker_size_id = 1
    worker_size = "Small"
    worker_count = 3
  }

  worker_pool {
    worker_size_id = 2
    worker_size = "Small"
    worker_count = 4
  }

  tags {
    Source = "AcceptanceTests"
  }
}
`, rInt, location, rInt)
}
