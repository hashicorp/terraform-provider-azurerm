package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	nw "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
)

func TestAccAzureRMLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	var lb network.LoadBalancer
	addressPoolName := fmt.Sprintf("%d-address-pool", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(data, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	var lb network.LoadBalancer
	addressPoolName := fmt.Sprintf("%d-address-pool", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(data, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerBackEndAddressPool_requiresImport(data, addressPoolName),
				ExpectError: acceptance.RequiresImportError(data.ResourceType),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	var lb network.LoadBalancer
	addressPoolName := fmt.Sprintf("%d-address-pool", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_removal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolNotExists(addressPoolName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	var lb network.LoadBalancer
	addressPoolName := fmt.Sprintf("%d-address-pool", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_basic(data, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolDisappears(addressPoolName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMLoadBalancerBackEndAddressPool_withNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")

	var lb network.LoadBalancer
	addressPoolName := fmt.Sprintf("%d-address-pool", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerBackEndAddressPool_withNIC(data, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName, &lb),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolExists(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := nw.FindLoadBalancerBackEndAddressPoolByName(lb, addressPoolName)
		if !exists {
			return fmt.Errorf("A BackEnd Address Pool with name %q cannot be found.", addressPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolNotExists(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := nw.FindLoadBalancerBackEndAddressPoolByName(lb, addressPoolName)
		if exists {
			return fmt.Errorf("A BackEnd Address Pool with name %q has been found.", addressPoolName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerBackEndAddressPoolDisappears(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		_, i, exists := nw.FindLoadBalancerBackEndAddressPoolByName(lb, addressPoolName)
		if !exists {
			return fmt.Errorf("A BackEnd Address Pool with name %q cannot be found.", addressPoolName)
		}

		currentPools := *lb.LoadBalancerPropertiesFormat.BackendAddressPools
		pools := append(currentPools[:i], currentPools[i+1:]...)
		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &pools

		id, err := azure.ParseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %+v", err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
	}
}

func testAccAzureRMLoadBalancerBackEndAddressPool_basic(data acceptance.TestData, addressPoolName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, addressPoolName)
}

func testAccAzureRMLoadBalancerBackEndAddressPool_requiresImport(data acceptance.TestData, name string) string {
	template := testAccAzureRMLoadBalancerBackEndAddressPool_basic(data, name)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "import" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, template)
}

func testAccAzureRMLoadBalancerBackEndAddressPool_removal(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLoadBalancerBackEndAddressPool_withNIC(data acceptance.TestData, addressPoolName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVnet-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestPip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "acctestLB-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "%[3]s"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_lb_outbound_rule" "test" {
  name                    = "acctestOR-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "PublicIPAddress"
  }
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                   = "acctestVMSS-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  upgrade_policy_mode    = "Automatic"
  overprovision          = false
  single_placement_group = true

  sku {
    name     = "Standard_B2ms"
    tier     = "Standard"
    capacity = 3
  }

  storage_profile_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-datacenter-core-smalldisk-g2"
    version   = "latest"
  }

  storage_profile_os_disk {
    create_option     = "FromImage"
    caching           = "ReadWrite"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "azureuser"
    admin_password       = "Qwerty12345!"
  }

  os_profile_windows_config {
    provision_vm_agent        = true
    enable_automatic_upgrades = false
  }

  network_profile {
    name    = "primary"
    primary = true

    ip_configuration {
      name      = "primary"
      primary   = true
      subnet_id = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [
        azurerm_lb_backend_address_pool.test.id
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, addressPoolName)
}
