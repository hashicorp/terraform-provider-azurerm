package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccJitNetworkAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_jit_network_access_policies", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMJitNetworkPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJitNetworkAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckJitNetworkAccessPolicyExists(data),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckJitNetworkAccessPolicyExists(data acceptance.TestData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.JitNetworkAccessPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		jitNetworkAccessPolicyName := fmt.Sprintf("default_%d", data.RandomInteger)
		jitRGName := fmt.Sprintf("acctestJIT_RG_%d", data.RandomInteger)
		location := data.Locations.Primary
		resp, err := client.Get(ctx, jitRGName, jitNetworkAccessPolicyName, location)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("[DEBUG EXISTS] Security Center Subscription Virtual Machine Jit Network Access Policy %q was not found: %+v", jitNetworkAccessPolicyName, err)
			}

			return fmt.Errorf("Bad: Get: %+v", err)
		}

		return nil
	}
}

func testAccJitNetworkAccessPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMJIT_Network_Access_Policy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_jit_network_access_policies" "jit" {
	name = "default_%d"
	asc_location = azurerm_resource_group.rg.location
	resource_group_name = azurerm_resource_group.rg.name
	kind = "Basic"
	virtual_machines {
		name = azurerm_virtual_machine.main.name
		ports {
			port = 22
			protocol = "*"
			allowed_source_address_prefix = "*"
			max_request_access_duration = "PT3H"
		}
		
	}
	
  }
`, template, data.RandomInteger)
}

func testAccAzureRMJIT_Network_Access_Policy_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
	provider "azurerm" {
		features {}
	  }
	  
	  resource "azurerm_resource_group" "rg" {
		name     = "acctestJIT_RG_%[1]d"
		location = "%s"
	  }
	  
	  # create vm
	  resource "azurerm_virtual_network" "vnet" {
		name                = "acctestJIT_network_%[1]d"
		address_space       = ["10.0.0.0/16"]
		location            = azurerm_resource_group.rg.location
		resource_group_name = azurerm_resource_group.rg.name
	  }
	  
	  resource "azurerm_subnet" "subnet" {
		name                 = "subnet__%[1]d"
		resource_group_name  = azurerm_resource_group.rg.name
		virtual_network_name = azurerm_virtual_network.vnet.name
		address_prefixes     = ["10.0.2.0/24"]
	  }
	  
	  resource "azurerm_network_interface" "inet" {
		name                = "nic__%[1]d"
		location            = azurerm_resource_group.rg.location
		resource_group_name = azurerm_resource_group.rg.name
	  
		ip_configuration {
		  name                          = "testconfiguration1"
		  subnet_id                     = azurerm_subnet.subnet.id
		  private_ip_address_allocation = "Dynamic"
		}
	  }

	  resource "azurerm_subnet_network_security_group_association" "nsg_subnet_association" {
		subnet_id                 = azurerm_subnet.subnet.id
		network_security_group_id = azurerm_network_security_group.vmnsg.id
	  }

	  resource "azurerm_network_security_group" "vmnsg"{
		name = "nsg_%[1]d"
		location            = azurerm_resource_group.rg.location
		resource_group_name = azurerm_resource_group.rg.name
		security_rule {
		name						="acctest_secrule"
		priority                   = 1000
		direction                  = "Inbound"
		access                     = "Deny"
		protocol                   = "Tcp"
		source_port_range          = "*"
		destination_port_range     = "*"
		source_address_prefix      = "*"
		destination_address_prefix = "*"
		}
	
	}

	resource "azurerm_public_ip" "public_vm_ip" {
		name                = "vm_ip_%[1]d"
		resource_group_name = azurerm_resource_group.rg.name
		location            = azurerm_resource_group.rg.location
		allocation_method   = "Dynamic"
	  
	  }

	  resource "azurerm_virtual_machine" "main" {
		name                             = "vm_%[1]d"
		location                         = azurerm_resource_group.rg.location
		resource_group_name              = azurerm_resource_group.rg.name
		network_interface_ids            = [azurerm_network_interface.inet.id]
		vm_size                          = "Standard_A1_v2"
		delete_os_disk_on_termination    = true
		delete_data_disks_on_termination = true
	  
		storage_image_reference {
		  publisher = "Canonical"
		  offer     = "UbuntuServer"
		  sku       = "16.04-LTS"
		  version   = "latest"
		}
		storage_os_disk {
		  name              = "disk1"
		  caching           = "ReadWrite"
		  create_option     = "FromImage"
		  managed_disk_type = "Standard_LRS"
		}
		os_profile {
		  computer_name  = "hostname"
		  admin_username = "tester"
		  admin_password = "12345&tester"
		}
		os_profile_linux_config {
		  disable_password_authentication = false
		}
	  }
`, data.RandomInteger, data.Locations.Primary)
}

func testCheckAzureRMJitNetworkPolicyDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.JitNetworkAccessPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_jit_network_access_policies" {
			continue
		}

		jitPolicyName := rs.Primary.Attributes["name"] //"default" // rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		location := rs.Primary.Attributes["asc_location"]

		resp, err := conn.Get(ctx, resourceGroup, jitPolicyName, location)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get Server: %+v", err)
		}

		return fmt.Errorf("JIT Network Access Policy %s still exists", jitPolicyName)
	}

	return nil
}
