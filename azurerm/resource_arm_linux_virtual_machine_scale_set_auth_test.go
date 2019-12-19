package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMLinuxVirtualMachineScaleSet_authPassword(t *testing.T) {
	resourceName := "azurerm_linux_virtual_machine_scale_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authPassword(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"admin_password",
					"terraform_should_roll_instances_when_required",
				},
			},
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_authSSHKey(t *testing.T) {
	resourceName := "azurerm_linux_virtual_machine_scale_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKey(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"terraform_should_roll_instances_when_required",
				},
			},
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_authSSHKeyAndPassword(t *testing.T) {
	resourceName := "azurerm_linux_virtual_machine_scale_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKeyAndPassword(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"admin_password",
					"terraform_should_roll_instances_when_required",
				},
			},
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_authMultipleSSHKeys(t *testing.T) {
	resourceName := "azurerm_linux_virtual_machine_scale_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authMultipleSSHKeys(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"terraform_should_roll_instances_when_required",
				},
			},
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_authUpdatingSSHKeys(t *testing.T) {
	resourceName := "azurerm_linux_virtual_machine_scale_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKey(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"terraform_should_roll_instances_when_required",
				},
			},
			{
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKeyUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"terraform_should_roll_instances_when_required",
				},
			},
		},
	})
}

func TestAccAzureRMLinuxVirtualMachineScaleSet_authDisablePasswordAuthUpdate(t *testing.T) {
	resourceName := "azurerm_linux_virtual_machine_scale_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLinuxVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				// disable it
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKey(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_password",
					"terraform_should_roll_instances_when_required",
				},
			},
			{
				// enable it
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authPassword(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_password",
					"terraform_should_roll_instances_when_required",
				},
			},
			{
				// disable it
				Config: testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKey(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLinuxVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_password",
					"terraform_should_roll_instances_when_required",
				},
			},
		},
	})
}

func testAccAzureRMLinuxVirtualMachineScaleSet_authPassword(rInt int, location string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, rInt)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKey(rInt int, location string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, rInt)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKeyUpdated(rInt int, location string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.second_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, rInt)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_authSSHKeyAndPassword(rInt int, location string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssw0rd1234!"

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, rInt)
}

func testAccAzureRMLinuxVirtualMachineScaleSet_authMultipleSSHKeys(rInt int, location string) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.second_public_key
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template, rInt)
}
