package compute_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

// NOTE: Test `TestAccVirtualMachine_enableAnWithVM` requires a machine of size `D8_v3` which is large/expensive - you may wish to ignore this test"

func TestAccVirtualMachine_basicLinuxMachine_managedDisk_standardSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_standardSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_os_disk.0.managed_disk_type").HasValue("StandardSSD_LRS"),
			),
		},
	})
}

func TestAccVirtualMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_standardSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_machine"),
		},
	})
}

func TestAccVirtualMachine_basicLinuxMachine_managedDisk_explicit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_explicit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachine_basicLinuxMachine_managedDisk_implicit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_implicit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachine_basicLinuxMachine_managedDisk_attach(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_attach(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachine_withDataDisk_managedDisk_explicit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDataDisk_managedDisk_explicit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachine_withDataDisk_managedDisk_implicit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDataDisk_managedDisk_implicit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachine_deleteManagedDiskOptOut(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	var vm compute.VirtualMachine
	var osd string
	var dtd string

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Destroy: false,
			Config:  r.withDataDisk_managedDisk_implicit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testLookupVirtualMachineManagedDiskID(&vm, "myosdisk1", &osd),
				testLookupVirtualMachineManagedDiskID(&vm, "mydatadisk1", &dtd),
			),
		},
		{
			Config: r.basicLinuxMachineDeleteVM_managedDisk(data),
			Check: resource.ComposeTestCheckFunc(
				testCheckVirtualMachineManagedDiskExists(&osd, true),
				testCheckVirtualMachineManagedDiskExists(&dtd, true),
			),
		},
	})
}

func TestAccVirtualMachine_deleteManagedDiskOptIn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	var vm compute.VirtualMachine
	var osd string
	var dtd string

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Destroy: false,
			Config:  r.basicLinuxMachine_managedDisk_DestroyDisksBefore(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testLookupVirtualMachineManagedDiskID(&vm, "myosdisk1", &osd),
				testLookupVirtualMachineManagedDiskID(&vm, "mydatadisk1", &dtd),
			),
		},
		{
			Config: r.basicLinuxMachine_managedDisk_DestroyDisksAfter(data),
			Check: resource.ComposeTestCheckFunc(
				testCheckVirtualMachineManagedDiskExists(&osd, false),
				testCheckVirtualMachineManagedDiskExists(&dtd, false),
			),
		},
	})
}

func TestAccVirtualMachine_osDiskTypeConflict(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.osDiskTypeConflict(data),
			ExpectError: regexp.MustCompile("conflicts with storage_os_disk.0.managed_disk_type"),
		},
	})
}

func TestAccVirtualMachine_dataDiskTypeConflict(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.dataDiskTypeConflict(data),
			ExpectError: regexp.MustCompile("Conflict between `vhd_uri`"),
		},
	})
}

func TestAccVirtualMachine_bug33(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.bug33(data),
		},
	})
}

func TestAccVirtualMachine_attachSecondDataDiskWithAttachOption(t *testing.T) {
	var afterCreate, afterUpdate compute.VirtualMachine
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_data_disk.0.create_option").HasValue("Empty"),
			),
		},
		{
			Config: r.basicLinuxMachine_managedDisk_attach(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testAccCheckVirtualMachineRecreated(t, &afterCreate, &afterUpdate),
				check.That(data.ResourceName).Key("storage_data_disk.0.create_option").HasValue("Empty"),
				check.That(data.ResourceName).Key("storage_data_disk.1.create_option").HasValue("Attach"),
			),
		},
	})
}

func TestAccVirtualMachine_linuxNoConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.linuxNoConfig(data),
			ExpectError: regexp.MustCompile("Error: either a `os_profile_linux_config` or a `os_profile_windows_config` must be specified."),
		},
	})
}

func TestAccVirtualMachine_windowsNoConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.windowsNoConfig(data),
			ExpectError: regexp.MustCompile("Error: either a `os_profile_linux_config` or a `os_profile_windows_config` must be specified."),
		},
	})
}

func TestAccVirtualMachine_multipleNICs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	prefix := fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/networkInterfaces", subscriptionId, data.RandomInteger)
	firstNicName := fmt.Sprintf("%s/acctni1-%d", prefix, data.RandomInteger)
	secondNicName := fmt.Sprintf("%s/acctni2-%d", prefix, data.RandomInteger)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleNICs(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("network_interface_ids.0").HasValue(firstNicName),
				check.That(data.ResourceName).Key("network_interface_ids.1").HasValue(secondNicName),
			),
		},
	})
}

func TestAccVirtualMachine_managedServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withManagedServiceIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
				resource.TestMatchOutput("principal_id", validate.UUIDRegExp),
			),
		},
	})
}

func TestAccVirtualMachine_enableAnWithVM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.anWithVM(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachine_basicLinuxMachine_managedDisk_changeOsWriteAcceleratorEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_withOsWriteAcceleratorEnabled(data, "true"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_os_disk.0.write_accelerator_enabled").HasValue("true"),
			),
		},
		{
			Config: r.basicLinuxMachine_managedDisk_withOsWriteAcceleratorEnabled(data, "false"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_os_disk.0.write_accelerator_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccVirtualMachine_basicLinuxMachine_managedDisk_withWriteAcceleratorEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_withWriteAcceleratorEnabled(data, "true"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_data_disk.0.write_accelerator_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccVirtualMachine_basicLinuxMachine_managedDisk_changeWriteAcceleratorEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_withWriteAcceleratorEnabled(data, "false"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_data_disk.0.write_accelerator_enabled").HasValue("false"),
			),
		},
		{
			Config: r.basicLinuxMachine_managedDisk_withWriteAcceleratorEnabled(data, "true"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_data_disk.0.write_accelerator_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccVirtualMachine_winRMCerts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.winRMCerts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachine_hasDiskInfoWhenStopped(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.hasDiskInfoWhenStopped(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_os_disk.0.managed_disk_type").HasValue("Standard_LRS"),
				check.That(data.ResourceName).Key("storage_data_disk.0.disk_size_gb").HasValue("64"),
			),
		},
		{
			Config: r.hasDiskInfoWhenStopped(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClient(r.deallocate),
				check.That(data.ResourceName).Key("storage_os_disk.0.managed_disk_type").HasValue("Standard_LRS"),
				check.That(data.ResourceName).Key("storage_data_disk.0.disk_size_gb").HasValue("64"),
			),
		},
	})
}

func TestAccVirtualMachine_importBasic_withZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicLinuxMachine_managedDisk_implicit_withZone(data),
		},

		data.ImportStep(
			"delete_data_disks_on_termination",
			"delete_os_disk_on_termination",
		),
	})
}

// this test requires eastus2 location and availability zone 1
// might be worth creating a `azurerm_compute_zones` resource to get the zones to use?
func TestAccVirtualMachine_ultraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "test")
	r := VirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ultraSSD(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_withOsWriteAcceleratorEnabled(data acceptance.TestData, enabled string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_M64s"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name                      = "osd-%d"
    create_option             = "FromImage"
    caching                   = "None"
    disk_size_gb              = "50"
    managed_disk_type         = "Premium_LRS"
    write_accelerator_enabled = %s
  }

  storage_data_disk {
    name              = "acctmd-%d"
    create_option     = "Empty"
    disk_size_gb      = "1"
    managed_disk_type = "Standard_LRS"
    lun               = 0
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabled, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_withWriteAcceleratorEnabled(data acceptance.TestData, enabled string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_M64s"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name                      = "acctmd-%d"
    create_option             = "Empty"
    disk_size_gb              = "1"
    managed_disk_type         = "Premium_LRS"
    lun                       = 0
    write_accelerator_enabled = %s
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabled, data.RandomInteger)
}

func (VirtualMachineResource) winRMCerts(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "prefix" {
  default = "acctest%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-resources"
  location = "%s"

  tags = {
    source = "TestAccVirtualMachine_winRMCerts"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "${var.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "${var.prefix}-nic"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_key_vault" "test" {
  name                = "${var.prefix}-keyvault"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku_name = "standard"

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

    key_permissions = [
      "backup",
      "create",
      "decrypt",
      "delete",
      "encrypt",
      "get",
      "import",
      "list",
      "purge",
      "recover",
      "restore",
      "sign",
      "unwrapKey",
      "update",
      "verify",
      "wrapKey",
    ]

    secret_permissions = [
      "backup",
      "delete",
      "get",
      "list",
      "purge",
      "recover",
      "restore",
      "set",
    ]

    certificate_permissions = [
      "create",
      "delete",
      "deleteissuers",
      "get",
      "getissuers",
      "import",
      "list",
      "listissuers",
      "managecontacts",
      "manageissuers",
      "purge",
      "setissuers",
      "update",
    ]
  }

  enabled_for_deployment          = true
  enabled_for_template_deployment = true
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "${var.prefix}-cert"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=${azurerm_network_interface.test.private_ip_address}"
      validity_in_months = 12
    }
  }
}

resource "azurerm_virtual_machine" "test" {
  name                          = "${var.prefix}-vm"
  location                      = "${azurerm_resource_group.test.location}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  network_interface_ids         = ["${azurerm_network_interface.test.id}"]
  vm_size                       = "Standard_F2"
  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "${var.prefix}-osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "${var.prefix}-vm"
    admin_username = "mradministrator"
    admin_password = "Th15IsD0g1234!"
  }

  os_profile_windows_config {
    provision_vm_agent = true

    winrm {
      protocol        = "https"
      certificate_url = "${azurerm_key_vault_certificate.test.secret_id}"
    }
  }

  os_profile_secrets {
    source_vault_id = "${azurerm_key_vault.test.id}"

    vault_certificates {
      certificate_url   = "${azurerm_key_vault_certificate.test.secret_id}"
      certificate_store = "My"
    }
  }
}
`, data.RandomString, data.Locations.Primary)
}

func (VirtualMachineResource) withManagedServiceIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  identity {
    type = "SystemAssigned"
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}

output "principal_id" {
  value = azurerm_virtual_machine.test.identity[0]["principal_id"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_explicit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_standardSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "StandardSSD_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine" "import" {
  name                  = azurerm_virtual_machine.test.name
  location              = azurerm_virtual_machine.test.location
  resource_group_name   = azurerm_virtual_machine.test.resource_group_name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "StandardSSD_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, r.basicLinuxMachine_managedDisk_standardSSD(data), data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_implicit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "osd-%d"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "50"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_implicit_withZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"
  zones                 = ["1"]

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "osd-%d"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "50"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_attach(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctmd2-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name              = "acctmd-%d"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 0
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name            = azurerm_managed_disk.test.name
    create_option   = "Attach"
    disk_size_gb    = "1"
    lun             = 1
    managed_disk_id = azurerm_managed_disk.test.id
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name              = "acctmd-%d"
    create_option     = "Empty"
    disk_size_gb      = "1"
    managed_disk_type = "Standard_LRS"
    lun               = 0
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_DestroyDisksBefore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  delete_os_disk_on_termination = true

  storage_data_disk {
    name          = "mydatadisk1"
    disk_size_gb  = "1"
    create_option = "Empty"
    lun           = 0
  }

  delete_data_disks_on_termination = true

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachine_managedDisk_DestroyDisksAfter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) basicLinuxMachineDeleteVM_managedDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) withDataDisk_managedDisk_explicit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name              = "dtd-%d"
    disk_size_gb      = "1"
    create_option     = "Empty"
    caching           = "ReadWrite"
    lun               = 0
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) withDataDisk_managedDisk_implicit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  storage_data_disk {
    name          = "mydatadisk1"
    disk_size_gb  = "1"
    create_option = "Empty"
    caching       = "ReadWrite"
    lun           = 0
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) osDiskTypeConflict(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "10"
    managed_disk_type = "Standard_LRS"
    vhd_uri           = "should_cause_conflict"
  }

  storage_data_disk {
    name              = "mydatadisk1"
    caching           = "ReadWrite"
    create_option     = "Empty"
    disk_size_gb      = "45"
    managed_disk_type = "Standard_LRS"
    lun               = "0"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) dataDiskTypeConflict(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "10"
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name              = "mydatadisk1"
    caching           = "ReadWrite"
    create_option     = "Empty"
    disk_size_gb      = "45"
    managed_disk_type = "Standard_LRS"
    lun               = "0"
  }

  storage_data_disk {
    name              = "mydatadisk1"
    vhd_uri           = "should_cause_conflict"
    caching           = "ReadWrite"
    create_option     = "Empty"
    disk_size_gb      = "45"
    managed_disk_type = "Standard_LRS"
    lun               = "1"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) bug33(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm%s"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F1"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctvm%s"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}

func (VirtualMachineResource) linuxNoConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F1"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctvm%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) windowsNoConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F1"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctvm%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) multipleNICs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "first" {
  name                = "acctni1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface" "second" {
  name                = "acctni2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                         = "acctvm%s"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  network_interface_ids        = [azurerm_network_interface.first.id, azurerm_network_interface.second.id]
  primary_network_interface_id = azurerm_network_interface.first.id
  vm_size                      = "Standard_F1"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctvm%s"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}

func (VirtualMachineResource) anWithVM(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                          = "acctestni-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  enable_ip_forwarding          = false
  enable_accelerated_networking = true

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                         = "acctestvm-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  primary_network_interface_id = azurerm_network_interface.test.id
  network_interface_ids        = [azurerm_network_interface.test.id]

  // Only large VMs allow AN
  vm_size                       = "Standard_D8_v3"
  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "antest-%d-OSDisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
    disk_size_gb      = 32
  }

  os_profile {
    computer_name  = "antestMachine-%d"
    admin_username = "antestuser"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) hasDiskInfoWhenStopped(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
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

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_DS1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "acctest-osdisk-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name          = "acctest-datadisk-%d"
    caching       = "ReadWrite"
    create_option = "Empty"
    lun           = 0
    disk_size_gb  = 64
  }

  os_profile {
    computer_name  = "acctest-machine-%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineResource) ultraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/28"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/29"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D2S_V3"
  zones                 = ["1"]

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    os_type           = "linux"
    managed_disk_type = "Premium_LRS"
    disk_size_gb      = "64"
  }

  storage_data_disk {
    name              = "mydatadisk1"
    caching           = "None"
    create_option     = "Empty"
    managed_disk_type = "UltraSSD_LRS"
    disk_size_gb      = "64"
    lun               = 1
  }

  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
