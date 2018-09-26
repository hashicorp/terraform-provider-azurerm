package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"strconv"
)

//basic daily
//basic weekly
//complete daily
//complete weekly

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateDailyToWeekly(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeeklyToDaily(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMRecoveryServicesProtectionPolicyVm_updateWeeklyToPartial(t *testing.T) {
	resourceName := "azurerm_recovery_services_protection_policy_vm.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName, ri),
			},
			{
				Config: testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(resourceName, ri),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMRecoveryServicesProtectionPolicyVmDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).recoveryServicesProtectionPoliciesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_recovery_services_protection_policy_vm" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		policyName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Recovery Services Vault Policy still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).recoveryServicesProtectionPoliciesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Recovery Services Vault Policy: %q", resourceName)
		}

		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		policyName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Recovery Services Vault Policy %q (resource group: %q) was not found: %+v", policyName, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on recoveryServicesVaultsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt int, location string) string {
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
  name                         = "acctest-ip"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "acctestip%[1]d"
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
    disk_size_gb      = "1023"
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

`, rInt, location, strconv.Itoa(rInt)[0:5])
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  
  backup = {
    frequency = "Daily"
    time      = "23:00"
  } 

  retention_daily = {
    count = 10
  }
}

`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
}

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "10"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  
  backup = {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday"]
  } 

  retention_weekly = {
    count    = 42
    weekdays = ["Sunday", "Wednesday"]
  }
}

`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
}

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_basicWeekly(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Weekly"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "2"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  
  backup = {
    frequency = "Daily"
    time      = "23:00"
  } 

  retention_daily = {
    count = 10
  }

  retention_weekly = {
    count    = 42
    weekdays = ["Sunday", "Wednesday"]
  }

  retention_monthly = {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly = {
    count    = 77
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }
}
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
}

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeDaily(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Daily"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "retention_daily.0.count", "10"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.count", "7"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.count", "77"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.months.#", "2"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  
  backup = {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday", "Saturday"]
  } 

  retention_weekly = {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_monthly = {
    count    = 7
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly = {
    count    = 77
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }
}
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
}

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeekly(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Weekly"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.count", "7"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.count", "77"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.months.#", "2"),
	)
}

func testAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  
  backup = {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday", "Saturday"]
  } 

  retention_weekly = {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday"]
  }

  retention_monthly = {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly = {
    count    = 77
    weekdays = ["Sunday"]
    weeks    = ["Last"]
    months   = ["January"]
  }
}
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_base(rInt, location), rInt)
}

func checkAccAzureRMRecoveryServicesProtectionPolicyVm_completeWeeklyPartial(resourceName string, ri int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "resource_group_name", fmt.Sprintf("acctestRG-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "recovery_vault_name", fmt.Sprintf("acctest-%d", ri)),
		resource.TestCheckResourceAttr(resourceName, "backup.0.frequency", "Weekly"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.time", "23:00"),
		resource.TestCheckResourceAttr(resourceName, "backup.0.weekdays.#", "4"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.count", "42"),
		resource.TestCheckResourceAttr(resourceName, "retention_weekly.0.weekdays.#", "3"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.count", "7"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weekdays.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_monthly.0.weeks.#", "2"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.count", "77"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weekdays.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.weeks.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "retention_yearly.0.months.#", "1"),
	)
}
