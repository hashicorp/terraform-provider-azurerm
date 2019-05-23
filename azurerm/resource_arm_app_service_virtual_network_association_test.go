package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceVirtualNetworkAssociation_basic(t *testing.T) {
	resourceName := "azurerm_app_service_virtual_network_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkAssociationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
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

func TestAccAzureRMAppServiceVirtualNetworkAssociation_update(t *testing.T) {
	resourceName := "azurerm_app_service_virtual_network_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkAssociationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
				),
			},
			{
				Config: testAccAzureRMAppServiceVirtualNetworkAssociation_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkAssociationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceVirtualNetworkAssociation_disappears(t *testing.T) {
	resourceName := "azurerm_app_service_virtual_network_association.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkAssociation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkAssociationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					testCheckAzureRMAppServiceVirtualNetworkAssociationDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMAppServiceVirtualNetworkAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.Attributes["id"]
		parsedID, err := parseAzureResourceID(id)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}
		name := parsedID.Path["sites"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Azure App Service Virtual Network Association: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).appServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Virtual Network Association %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAppServiceVirtualNetworkAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.Attributes["id"]
		parsedID, err := parseAzureResourceID(id)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}
		name := parsedID.Path["sites"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Azure App Service Virtual Network Association: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).appServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.DeleteSwiftVirtualNetwork(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: Delete on appServicesClient: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMAppServiceVirtualNetworkAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).appServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_virtual_network_association" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)

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

func testAccAzureRMAppServiceVirtualNetworkAssociation_base(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  lifecycle {
    ignore_changes = ["ddos_protection_plan"]
  }
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestasp-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestas-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAppServiceVirtualNetworkAssociation_basic(rInt int, location string) string {
	template := testAccAzureRMAppServiceVirtualNetworkAssociation_base(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_association" "test" {
  app_service_id       = "${azurerm_app_service.test.id}"
  subnet_id            = "${azurerm_subnet.test1.id}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
}
`, template)
}

func testAccAzureRMAppServiceVirtualNetworkAssociation_update(rInt int, location string) string {
	template := testAccAzureRMAppServiceVirtualNetworkAssociation_base(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_association" "test" {
  app_service_id       = "${azurerm_app_service.test.id}"
  subnet_id            = "${azurerm_subnet.test2.id}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
}
`, template)
}
