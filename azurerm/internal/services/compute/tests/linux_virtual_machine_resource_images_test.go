package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccLinuxVirtualMachine_imageFromImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				// create the original VM
				Config: testLinuxVirtualMachine_imageFromExistingMachinePrep(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists("azurerm_linux_virtual_machine.source"),
					generalizeLinuxVirtualMachine("azurerm_linux_virtual_machine.source"),
				),
			},
			{
				// then create an image from that VM, and then create a VM from that image
				Config: testLinuxVirtualMachine_imageFromImage(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccLinuxVirtualMachine_imageFromPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_imageFromPlan(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccLinuxVirtualMachine_imageFromSharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				// create the original VM
				Config: testLinuxVirtualMachine_imageFromExistingMachinePrep(data),
				Check: resource.ComposeTestCheckFunc(

					checkLinuxVirtualMachineExists("azurerm_linux_virtual_machine.source"),
					generalizeLinuxVirtualMachine("azurerm_linux_virtual_machine.source"),
				),
			},
			{
				// then create an image from that VM, and then create a VM from that image
				Config: testLinuxVirtualMachine_imageFromSharedImageGallery(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccLinuxVirtualMachine_imageFromSourceImageReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkLinuxVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testLinuxVirtualMachine_imageFromSourceImageReference(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func testLinuxVirtualMachine_imageFromExistingMachineDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
# note: whilst these aren't used in all tests, it saves us redefining these everywhere
locals {
  first_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9ddAwoR0XBoT9kixLX6atgX9dovt9fpR1HO/R9jYwYnuB+SZ845KSqat+U0m6oagZhpsfcEEwjGGjQz6Z1rB6mvffsKq6i74cmm0jO564nBnZQeh31q3sFNs+XdrDtFmnYRqdPHhhr1sw0C/rxbiaE6nYZWRfHW//81nEePKMpjiN8JsrYQNbzEpz8QOBSquwBmXO+LVx//zAbY4jGTa4hjGeNzIgMJZ8Jk/11XbcxSK1PK43BrejHg6kctmEkYvMH/o12RfAeB8okGCRW3scwOozxVrHwxaPgEf03jig+Ag9V+GXNBabL5AWtxcuPN63rUfaAXEIXTHmndwVOxlpLrUf5ox1+ddGyWbLMXzd7akPioof5MNJMq/yuFGC5dY0Z6/+yGRNtShQesVo/czhKEPGIcsIi5gnKdfDB4i9ay2yz8ystnW6jbabcyqejk1Qc61wapaFdhUHL0iD/GW/5ZujDs5C3BT7EIgKLIfAaAx5TBEJyE1KQ/GEOifB8ztDl/gp99o+i2HKABtmYv12y4JVlEUkRckeLrw6luEb3ColHshsQcQGfudGFFgdEdcgBrV4Ch7IkLxVYQl3pegzZiirMPnRKh10r/Hrg6uYxn7sLeTJoD5VOKmqmeK4kFXsZMVtA6/SnxQtUKkKlfLBwBSDrrdgLjBV+KOndiwC7Q=="
  second_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
  vm_name           = "acctestsourcevm-%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = local.vm_name
}

resource "azurerm_network_interface" "public" {
  name                = "acctnicsource-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testLinuxVirtualMachine_imageFromExistingMachinePrep(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_imageFromExistingMachineDependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "source" {
  name                            = "acctestsourceVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "Eung6ahthane2ied"

  network_interface_ids = [
    azurerm_network_interface.public.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_imageFromImage(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_imageFromExistingMachinePrep(data)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "Eung6ahthane2ied"
  source_image_id                 = azurerm_image.test.id

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
}

`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_imageFromPlan(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_marketplace_agreement" "test" {
  publisher = "cloudbees"
  offer     = "jenkins-operations-center"
  plan      = "jenkins-operations-center-solo"
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "Eung6ahthane2ied"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  plan {
    name      = "jenkins-operations-center-solo"
    product   = "jenkins-operations-center"
    publisher = "cloudbees"
  }

  source_image_reference {
    publisher = "cloudbees"
    offer     = "jenkins-operations-center"
    sku       = "jenkins-operations-center-solo"
    version   = "latest"
  }

  depends_on = ["azurerm_marketplace_agreement.test"]
}
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_imageFromSharedImageGallery(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_imageFromExistingMachinePrep(data)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "Eung6ahthane2ied"
  source_image_id                 = azurerm_shared_image_version.test.id

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
}
`, template, data.RandomInteger)
}

func testLinuxVirtualMachine_imageFromSourceImageReference(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  disable_password_authentication = false
  admin_password                  = "Eung6ahthane2ied"

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func generalizeLinuxVirtualMachine(resourceName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]
		username := rs.Primary.Attributes["admin_username"]
		password := rs.Primary.Attributes["admin_password"]
		port := "22"
		location := rs.Primary.Attributes["location"]

		return testGeneralizeVMImage(resourceGroup, name, username, password, name, port, location)(s)
	}
}
