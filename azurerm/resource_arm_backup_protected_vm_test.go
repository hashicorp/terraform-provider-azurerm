package azurerm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBackupProtectedVm_basic(t *testing.T) {
	resourceName := "azurerm_backup_protected_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectedVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectedVm_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBackupProtectedVmExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{ //vault cannot be deleted unless we unregister all backups
				Config: testAccAzureRMBackupProtectedVm_base(ri, acceptance.Location()),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccAzureRMBackupProtectedVm_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_backup_protected_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectedVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectedVm_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBackupProtectedVmExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
				),
			},
			{
				Config:      testAccAzureRMBackupProtectedVm_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_backup_protected_vm"),
			},
			{ //vault cannot be deleted unless we unregister all backups
				Config: testAccAzureRMBackupProtectedVm_base(ri, acceptance.Location()),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccAzureRMBackupProtectedVm_separateResourceGroups(t *testing.T) {
	resourceName := "azurerm_backup_protected_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectedVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBackupProtectedVm_separateResourceGroups(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBackupProtectedVmExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{ //vault cannot be deleted unless we unregister all backups
				Config: testAccAzureRMBackupProtectedVm_additionalVault(ri, acceptance.Location()),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccAzureRMBackupProtectedVm_updateBackupPolicyId(t *testing.T) {
	virtualMachine := "azurerm_virtual_machine.test"
	protectedVmResourceName := "azurerm_backup_protected_vm.test"
	fBackupPolicyResourceName := "azurerm_backup_policy_vm.test"
	sBackupPolicyResourceName := "azurerm_backup_policy_vm.test_change_backup"

	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBackupProtectedVmDestroy,
		Steps: []resource.TestStep{
			{ // Create resources and link first backup policy id
				ResourceName: fBackupPolicyResourceName,
				Config:       testAccAzureRMBackupProtectedVm_linkFirstBackupPolicy(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(protectedVmResourceName, "backup_policy_id", fBackupPolicyResourceName, "id"),
				),
			},
			{ // Modify backup policy id to the second one
				// Set Destroy false to prevent error from cleaning up dangling resource
				ResourceName: sBackupPolicyResourceName,
				Config:       testAccAzureRMBackupProtectedVm_linkSecondBackupPolicy(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(protectedVmResourceName, "backup_policy_id", sBackupPolicyResourceName, "id"),
				),
			},
			{ // Remove backup policy link
				// Backup policy link will need to be removed first so the VM's backup policy subsequently reverts to Default
				// Azure API is quite sensitive, adding the step to control resource cleanup order
				ResourceName: fBackupPolicyResourceName,
				Config:       testAccAzureRMBackupProtectedVm_withVM(ri, acceptance.Location()),
				Check:        resource.ComposeTestCheckFunc(),
			},
			{ // Then VM can be removed
				ResourceName: virtualMachine,
				Config:       testAccAzureRMBackupProtectedVm_withSecondPolicy(ri, acceptance.Location()),
				Check:        resource.ComposeTestCheckFunc(),
			},
			{ // Remove backup policies and vault
				ResourceName: protectedVmResourceName,
				Config:       testAccAzureRMBackupProtectedVm_basePolicyTest(ri, acceptance.Location()),
				Check:        resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testCheckAzureRMBackupProtectedVmDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_backup_protected_vm" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		vmId := rs.Primary.Attributes["source_vm_id"]

		parsedVmId, err := azure.ParseAzureResourceID(vmId)
		if err != nil {
			return fmt.Errorf("[ERROR] Unable to parse source_vm_id '%s': %+v", vmId, err)
		}
		vmName, hasName := parsedVmId.Path["virtualMachines"]
		if !hasName {
			return fmt.Errorf("[ERROR] parsed source_vm_id '%s' doesn't contain 'virtualMachines'", vmId)
		}

		protectedItemName := fmt.Sprintf("VM;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, vmName)
		containerName := fmt.Sprintf("iaasvmcontainer;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, vmName)

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectedItemsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Recovery Services Protected VM still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMBackupProtectedVmExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Recovery Services Protected VM: %q", resourceName)
		}

		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		vmId := rs.Primary.Attributes["source_vm_id"]

		//get VM name from id
		parsedVmId, err := azure.ParseAzureResourceID(vmId)
		if err != nil {
			return fmt.Errorf("[ERROR] Unable to parse source_vm_id '%s': %+v", vmId, err)
		}
		vmName, hasName := parsedVmId.Path["virtualMachines"]
		if !hasName {
			return fmt.Errorf("[ERROR] parsed source_vm_id '%s' doesn't contain 'virtualMachines'", vmId)
		}

		protectedItemName := fmt.Sprintf("VM;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, vmName)
		containerName := fmt.Sprintf("iaasvmcontainer;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, vmName)

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ProtectedItemsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Recovery Services Protected VM %q (resource group: %q) was not found: %+v", protectedItemName, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on recoveryServicesVaultsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMBackupProtectedVm_base(rInt int, location string) string {
	rstr := strconv.Itoa(rInt)
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctest_subnet"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  address_prefix       = "10.0.10.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctest_nic"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "acctestipconfig"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-ip"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"
  domain_name_label   = "acctestip%[1]d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]s"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-datadisk"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1023"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  vm_size               = "Standard_A0"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "acctest-osdisk"
    managed_disk_type = "Standard_LRS"
    caching           = "ReadWrite"
    create_option     = "FromImage"
  }

  storage_data_disk {
    name              = "acctest-datadisk"
    managed_disk_id   = "${azurerm_managed_disk.test.id}"
    managed_disk_type = "Standard_LRS"
    disk_size_gb      = "${azurerm_managed_disk.test.disk_size_gb}"
    create_option     = "Attach"
    lun               = 0
  }

  os_profile {
    computer_name  = "acctest"
    admin_username = "vmadmin"
    admin_password = "Password123!@#"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = "${azurerm_storage_account.test.primary_blob_endpoint}"
  }
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, rInt, location, rstr[len(rstr)-5:])
}

func testAccAzureRMBackupProtectedVm_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  source_vm_id        = "${azurerm_virtual_machine.test.id}"
  backup_policy_id    = "${azurerm_backup_policy_vm.test.id}"
}
`, testAccAzureRMBackupProtectedVm_base(rInt, location))
}

// For update backup policy id test
func testAccAzureRMBackupProtectedVm_basePolicyTest(rInt int, location string) string {
	rstr := strconv.Itoa(rInt)
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery1-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctest_subnet"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  address_prefix       = "10.0.10.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctest_nic"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "acctestipconfig"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-ip"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"
  domain_name_label   = "acctestip%[1]d"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]s"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-datadisk"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1023"
}
`, rInt, location, rstr[len(rstr)-5:])
}

// For update backup policy id test
func testAccAzureRMBackupProtectedVm_withVault(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%[2]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}
`, testAccAzureRMBackupProtectedVm_basePolicyTest(rInt, location), rInt)
}

// For update backup policy id test
func testAccAzureRMBackupProtectedVm_withFirstPolicy(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, testAccAzureRMBackupProtectedVm_withVault(rInt, location), rInt)
}

// For update backup policy id test
func testAccAzureRMBackupProtectedVm_withSecondPolicy(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_backup_policy_vm" "test_change_backup" {
  name                = "acctest2-%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 15
  }
}
`, testAccAzureRMBackupProtectedVm_withFirstPolicy(rInt, location), rInt)
}

// For update backup policy id test
func testAccAzureRMBackupProtectedVm_withVM(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_machine" "test" {
  name                          = "acctestvm-%[2]d"
  location                      = "${azurerm_resource_group.test.location}"
  resource_group_name           = "${azurerm_resource_group.test.name}"
  vm_size                       = "Standard_A0"
  network_interface_ids         = ["${azurerm_network_interface.test.id}"]
  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "acctest-osdisk"
    managed_disk_type = "Standard_LRS"
    caching           = "ReadWrite"
    create_option     = "FromImage"
  }

  storage_data_disk {
    name              = "acctest-datadisk"
    managed_disk_id   = "${azurerm_managed_disk.test.id}"
    managed_disk_type = "Standard_LRS"
    disk_size_gb      = "${azurerm_managed_disk.test.disk_size_gb}"
    create_option     = "Attach"
    lun               = 0
  }

  os_profile {
    computer_name  = "acctest"
    admin_username = "vmadmin"
    admin_password = "Password123!@#"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = "${azurerm_storage_account.test.primary_blob_endpoint}"
  }
}
`, testAccAzureRMBackupProtectedVm_withSecondPolicy(rInt, location), rInt)
}

// For update backup policy id test
func testAccAzureRMBackupProtectedVm_linkFirstBackupPolicy(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  source_vm_id        = "${azurerm_virtual_machine.test.id}"
  backup_policy_id    = "${azurerm_backup_policy_vm.test.id}"
}
`, testAccAzureRMBackupProtectedVm_withVM(rInt, location))
}

// For update backup policy id test
func testAccAzureRMBackupProtectedVm_linkSecondBackupPolicy(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  source_vm_id        = "${azurerm_virtual_machine.test.id}"
  backup_policy_id    = "${azurerm_backup_policy_vm.test_change_backup.id}"
}
`, testAccAzureRMBackupProtectedVm_withVM(rInt, location))
}

func testAccAzureRMBackupProtectedVm_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "import" {
  resource_group_name = "${azurerm_backup_protected_vm.test.resource_group_name}"
  recovery_vault_name = "${azurerm_backup_protected_vm.test.recovery_vault_name}"
  source_vm_id        = "${azurerm_backup_protected_vm.test.source_vm_id}"
  backup_policy_id    = "${azurerm_backup_protected_vm.test.backup_policy_id}"
}
`, testAccAzureRMBackupProtectedVm_basic(rInt, location))
}

func testAccAzureRMBackupProtectedVm_additionalVault(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-recovery2-%[2]d"
  location = "%[3]s"
}

resource "azurerm_recovery_services_vault" "test2" {
  name                = "acctest2-%[2]d"
  location            = "${azurerm_resource_group.test2.location}"
  resource_group_name = "${azurerm_resource_group.test2.name}"
  sku                 = "Standard"
}

resource "azurerm_backup_policy_vm" "test2" {
  name                = "acctest2-%[2]d"
  resource_group_name = "${azurerm_resource_group.test2.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test2.name}"

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, testAccAzureRMBackupProtectedVm_base(rInt, location), rInt, location)
}

func testAccAzureRMBackupProtectedVm_separateResourceGroups(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_protected_vm" "test" {
  resource_group_name = "${azurerm_resource_group.test2.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test2.name}"
  backup_policy_id    = "${azurerm_backup_policy_vm.test2.id}"
  source_vm_id        = "${azurerm_virtual_machine.test.id}"
}
`, testAccAzureRMBackupProtectedVm_additionalVault(rInt, location))
}
