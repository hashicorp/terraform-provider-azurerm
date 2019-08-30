package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkProfile_basic(t *testing.T) {
	resourceName := "azurerm_network_profile.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkProfile_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkProfileExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "container_network_interface_ids.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNetworkProfile_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_network_profile.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkProfile_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkProfileExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkProfile_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_network_profile"),
			},
		},
	})
}

func TestAccAzureRMNetworkProfile_withTags(t *testing.T) {
	resourceName := "azurerm_network_profile.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkProfile_withTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMNetworkProfile_withUpdatedTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Staging"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMNetworkProfile_disappears(t *testing.T) {
	resourceName := "azurerm_network_profile.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkProfile_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkProfileExists(resourceName),
					testCheckAzureRMNetworkProfileDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNetworkProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Network Profile: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).network.ProfileClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Profile %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on netProfileClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkProfileDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Network Profile: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).network.ProfileClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		if _, err := client.Delete(ctx, resourceGroup, name); err != nil {
			return fmt.Errorf("Bad: Delete on netProfileClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.ProfileClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_profile" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Network Profile still exists:\n%#v", resp.ProfilePropertiesFormat)
	}

	return nil
}

func testAccAzureRMNetworkProfile_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.1.0.0/24"

  delegation {
    name = "acctestdelegation-%d"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_network_profile" "test" {
  name                = "acctestnetprofile-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  container_network_interface {
    name = "acctesteth-%d"

    ip_configuration {
      name      = "acctestipconfig-%d"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetworkProfile_requiresImport(rInt int, location string) string {
	basicConfig := testAccAzureRMNetworkProfile_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_profile" "import" {
  name                = "${azurerm_network_profile.test.name}"
  location            = "${azurerm_network_profile.test.location}"
  resource_group_name = "${azurerm_network_profile.test.resource_group_name}"

  container_network_interface {
    name = "${azurerm_network_profile.test.container_network_interface.0.name}"

    ip_configuration {
      name      = "${azurerm_network_profile.test.container_network_interface.0.ip_configuration.0.name}"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }
}
`, basicConfig)
}

func testAccAzureRMNetworkProfile_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.1.0.0/24"

  delegation {
    name = "acctestdelegation-%d"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_network_profile" "test" {
  name                = "acctestnetprofile-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  container_network_interface {
    name = "acctesteth-%d"

    ip_configuration {
      name      = "acctestipconfig-%d"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetworkProfile_withUpdatedTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.1.0.0/24"

  delegation {
    name = "acctestdelegation-%d"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_network_profile" "test" {
  name                = "acctestnetprofile-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  container_network_interface {
    name = "acctesteth-%d"

    ip_configuration {
      name      = "acctestipconfig-%d"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  tags = {
    environment = "Staging"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}
