package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccLinuxVirtualMachineScaleSet_imagesAutomaticUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imagesAutomaticUpdate(data, "16.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.imagesAutomaticUpdate(data, "18.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_imagesFromCapturedVirtualMachineImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// provision a standard Virtual Machine with an Unmanaged Disk
			Config: r.imagesFromVirtualMachinePrerequisitesWithVM(data),
		},
		{
			// then delete the Virtual Machine
			Config: r.imagesFromVirtualMachinePrerequisites(data),
		},
		{
			// then capture two images of the Virtual Machine
			Config: r.imagesFromVirtualMachinePrerequisitesWithImage(data),
		},
		{
			// then provision a Virtual Machine Scale Set using this image
			Config: r.imagesFromVirtualMachine(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			// then update the image on this Virtual Machine Scale Set
			Config: r.imagesFromVirtualMachine(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_imagesManualUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imagesManualUpdate(data, "16.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.imagesManualUpdate(data, "18.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_imagesManualUpdateExternalRoll(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imagesManualUpdateExternalRoll(data, "16.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.imagesManualUpdateExternalRoll(data, "18.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_imagesRollingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imagesRollingUpdate(data, "16.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
		{
			Config: r.imagesRollingUpdate(data, "18.04-LTS"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func TestAccLinuxVirtualMachineScaleSet_imagesPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.imagesPlan(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"admin_password",
		),
	})
}

func (r LinuxVirtualMachineScaleSetResource) imagesAutomaticUpdate(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  name                           = "test"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "test" {
  name                           = "AccTestLBRule"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Automatic"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "%s"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
    }
  }

  automatic_os_upgrade_policy {
    disable_automatic_rollback  = true
    enable_automatic_os_upgrade = true
  }

  rolling_upgrade_policy {
    max_batch_instance_percent              = 100
    max_unhealthy_instance_percent          = 100
    max_unhealthy_upgraded_instance_percent = 100
    pause_time_between_batches              = "PT30S"
  }

  depends_on = ["azurerm_lb_rule.test"]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, version)
}

func (r LinuxVirtualMachineScaleSetResource) imagesFromVirtualMachinePrerequisites(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "source" {
  name                = "source-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_interface" "source" {
  name                = "sourcenic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "source"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.source.id
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r LinuxVirtualMachineScaleSetResource) imagesFromVirtualMachinePrerequisitesWithVM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine" "source" {
  name                  = "source"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.source.id]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "osdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/osdisk.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = 30
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "mradministrator"
    admin_password = "P@ssword1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
`, r.imagesFromVirtualMachinePrerequisites(data))
}

func (r LinuxVirtualMachineScaleSetResource) imagesFromVirtualMachinePrerequisitesWithImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_image" "first" {
  name                = "first"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/osdisk.vhd"
    size_gb  = 30
    caching  = "None"
  }
}

resource "azurerm_image" "second" {
  name                = "second"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/osdisk.vhd"
    size_gb  = 30
    caching  = "None"
  }

  depends_on = ["azurerm_image.first"]
}
`, r.imagesFromVirtualMachinePrerequisites(data))
}

func (r LinuxVirtualMachineScaleSetResource) imagesFromVirtualMachine(data acceptance.TestData, image string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "mradministrator"
  admin_password      = "P@ssword1234!"

  disable_password_authentication = false
  source_image_id                 = azurerm_image.%s.id

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "None"
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
`, r.imagesFromVirtualMachinePrerequisitesWithImage(data), data.RandomInteger, image)
}

func (r LinuxVirtualMachineScaleSetResource) imagesManualUpdate(data acceptance.TestData, version string) string {
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
    sku       = "%s"
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
`, r.template(data), data.RandomInteger, version)
}

func (r LinuxVirtualMachineScaleSetResource) imagesManualUpdateExternalRoll(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine_scale_set {
      roll_instances_when_required = false
    }
  }
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "%s"
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
`, r.template(data), data.RandomInteger, version)
}

func (r LinuxVirtualMachineScaleSetResource) imagesRollingUpdate(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_nat_pool" "test" {
  name                           = "test"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "acctest-lb-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "test" {
  name                           = "AccTestLBRule"
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  probe_id                       = azurerm_lb_probe.test.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.test.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  health_probe_id     = azurerm_lb_probe.test.id
  upgrade_mode        = "Rolling"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "%s"
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
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.test.id]
    }
  }

  rolling_upgrade_policy {
    max_batch_instance_percent              = 21
    max_unhealthy_instance_percent          = 22
    max_unhealthy_upgraded_instance_percent = 23
    pause_time_between_batches              = "PT30S"
  }

  depends_on = ["azurerm_lb_rule.test"]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, version)
}

func (r LinuxVirtualMachineScaleSetResource) imagesPlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_marketplace_agreement" "test" {
  publisher = "cloudbees"
  offer     = "jenkins-operations-center"
  plan      = "jenkins-operations-center-solo"
}

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
    publisher = "cloudbees"
    offer     = "jenkins-operations-center"
    sku       = "jenkins-operations-center-solo"
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

  plan {
    name      = "jenkins-operations-center-solo"
    product   = "jenkins-operations-center"
    publisher = "cloudbees"
  }

  depends_on = ["azurerm_marketplace_agreement.test"]
}
`, r.template(data), data.RandomInteger)
}
