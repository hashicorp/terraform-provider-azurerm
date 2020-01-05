package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkInterface_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					testCheckAzureRMNetworkInterfaceDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkInterface_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_interface"),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterface_setNetworkSecurityGroupId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	config := testAccAzureRMNetworkInterface_basic(data)
	updatedConfig := testAccAzureRMNetworkInterface_basicWithNetworkSecurityGroup(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_removeNetworkSecurityGroupId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	config := testAccAzureRMNetworkInterface_basicWithNetworkSecurityGroup(data)
	updatedConfig := testAccAzureRMNetworkInterface_basic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_security_group_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_multipleSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	config := testAccAzureRMNetworkInterface_multipleSubnets(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_multipleSubnetsPrimary(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	config := testAccAzureRMNetworkInterface_multipleSubnets(data)
	updatedConfig := testAccAzureRMNetworkInterface_multipleSubnetsUpdatedPrimary(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.primary", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.name", "testconfiguration1"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.1.primary", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.1.name", "testconfiguration2"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.primary", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.name", "testconfiguration2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.1.primary", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.1.name", "testconfiguration1"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_enableIPForwarding(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_ipForwarding(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_ip_forwarding", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_enableAcceleratedNetworking(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_acceleratedNetworking(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_accelerated_networking", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_multipleLoadBalancers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test1")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_multipleLoadBalancers(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists("azurerm_network_interface.test1"),
					testCheckAzureRMNetworkInterfaceExists("azurerm_network_interface.test2"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_applicationGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_applicationGatewayBackendPool(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists("azurerm_network_interface.test"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.application_gateway_backend_address_pools_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMNetworkInterface_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_IPAddressesBug1286(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_withIPAddresses(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.0.private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.0.public_ip_address_id"),
				),
			},
			{
				Config: testAccAzureRMNetworkInterface_withIPAddressesUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.0.private_ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.public_ip_address_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_IPAddressesFeature2543(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_withIPv6Addresses(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.1.private_ip_address_version", "IPv6"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_address"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_bug7986(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test1")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_bug7986(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists("azurerm_network_interface.test1"),
					testCheckAzureRMNetworkInterfaceExists("azurerm_network_interface.test2"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_applicationSecurityGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_applicationSecurityGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.application_security_group_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkInterface_importPublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_publicIP(data),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkInterfaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for availability set: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Interface %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ifaceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkInterfaceDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for availability set: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Error deleting Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for the deletion of Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkInterfaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_interface" {
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

		return fmt.Errorf("Network Interface still exists:\n%#v", resp.InterfacePropertiesFormat)
	}

	return nil
}

func testAccAzureRMNetworkInterface_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "import" {
  name                = "${azurerm_network_interface.test.name}"
  location            = "${azurerm_network_interface.test.location}"
  resource_group_name = "${azurerm_network_interface.test.resource_group_name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, template)
}

func testAccAzureRMNetworkInterface_basicWithNetworkSecurityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_interface" "test" {
  name                      = "acctestni-%d"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  network_security_group_id = "${azurerm_network_security_group.test.id}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_multipleSubnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "testconfiguration2"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_multipleSubnetsUpdatedPrimary(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration2"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_ipForwarding(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                 = "acctestni-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  enable_ip_forwarding = true

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_acceleratedNetworking(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                          = "acctestni-%d"
  location                      = "${azurerm_resource_group.test.location}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  enable_ip_forwarding          = false
  enable_accelerated_networking = true

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_withIPAddresses(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "test-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.2.9"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_withIPAddressesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "test-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_withIPv6Addresses(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "test-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "testconfiguration2"
    private_ip_address_version    = "IPv6"
    private_ip_address_allocation = "dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_multipleLoadBalancers(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "testext" {
  name                = "acctestpip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "testext" {
  name                = "acctestlb-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "publicipext"
    public_ip_address_id = "${azurerm_public_ip.testext.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "testext" {
  name                = "testbackendpoolext"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.testext.id}"
}

resource "azurerm_lb_nat_rule" "testext" {
  name                           = "testnatruleext"
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.testext.id}"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3390
  frontend_ip_configuration_name = "publicipext"
}

resource "azurerm_public_ip" "testint" {
  name                = "testpublicipint"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "testint" {
  name                = "testlbint"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                          = "publicipint"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "testint" {
  name                = "testbackendpoolint"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.testint.id}"
}

resource "azurerm_lb_nat_rule" "testint" {
  name                           = "testnatruleint"
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.testint.id}"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3391
  frontend_ip_configuration_name = "publicipint"
}

resource "azurerm_network_interface" "test1" {
  name                 = "acctestnic1-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  enable_ip_forwarding = true

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"

    load_balancer_backend_address_pools_ids = [
      "${azurerm_lb_backend_address_pool.testext.id}",
      "${azurerm_lb_backend_address_pool.testint.id}",
    ]
  }
}

resource "azurerm_network_interface" "test2" {
  name                 = "acctestnic2-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  enable_ip_forwarding = true

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"

    load_balancer_inbound_nat_rules_ids = [
      "${azurerm_lb_nat_rule.testext.id}",
      "${azurerm_lb_nat_rule.testint.id}",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_applicationGatewayBackendPool(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "gateway" {
  name                 = "subnet-gateway-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pubip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"
}

resource "azurerm_application_gateway" "test" {
  name                = "acctestgw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "gw-ip-config1"
    subnet_id = "${azurerm_subnet.gateway.id}"
  }

  frontend_port {
    name = "port-8080"
    port = 8080
  }

  frontend_ip_configuration {
    name                 = "ip-config-public"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "pool-1"
  }

  backend_http_settings {
    name                  = "backend-http-1"
    port                  = 8080
    protocol              = "Http"
    cookie_based_affinity = "Enabled"
    request_timeout       = 30
  }

  http_listener {
    name                           = "listener-1"
    frontend_ip_configuration_name = "ip-config-public"
    frontend_port_name             = "port-8080"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule-basic-1"
    rule_type                  = "Basic"
    http_listener_name         = "listener-1"
    backend_address_pool_name  = "pool-1"
    backend_http_settings_name = "backend-http-1"
  }

  tags = {
    environment = "tf01"
  }
}

resource "azurerm_network_interface" "test" {
  name                 = "acctestnic-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  enable_ip_forwarding = true

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"

    application_gateway_backend_address_pools_ids = [
      "${azurerm_application_gateway.test.backend_address_pool.0.id}",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_bug7986(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-%d-nsg"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "test1"
  priority                    = 101
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  network_security_group_name = "${azurerm_network_security_group.test.name}"
}

resource "azurerm_network_security_rule" "test2" {
  name                        = "test2"
  priority                    = 102
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  network_security_group_name = "${azurerm_network_security_group.test.name}"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-%d-pip"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-%d-vn"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "first"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test1" {
  name                = "acctest-%d-nic1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_network_interface" "test2" {
  name                = "acctest-%d-nic2"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_publicIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "testext" {
  name                = "acctestip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.testext.id}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkInterface_applicationSecurityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                           = "testconfiguration1"
    subnet_id                      = "${azurerm_subnet.test.id}"
    private_ip_address_allocation  = "Dynamic"
    application_security_group_ids = ["${azurerm_application_security_group.test.id}"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
