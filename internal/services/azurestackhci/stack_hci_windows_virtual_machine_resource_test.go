package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// https://learn.microsoft.com/en-us/azure/azure-local/manage/virtual-machine-image-azure-marketplace?tabs=azurecli
// Provisioning a VM requires an image. We support downloading an image using the `azurerm_stack_hci_marketplace_gallery_image` resource;
// however, downloading an image takes more than one hour.
// To speed up testing, we reuse a pre-downloaded image.
const (
	vmImageIdEnv = "ARM_TEST_STACK_HCI_VM_IMAGE_ID"
)

type StackHCIWindowsVirtualMachineResource struct{}

func TestAccStackHCIWindowsVirtualMachine(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// the test environment network limitation
	// (which our test suite can't easily work around)

	testCases := map[string]func(t *testing.T){
		"basic":          testAccStackHCIWindowsVirtualMachine_basic,
		"complete":       testAccStackHCIWindowsVirtualMachine_complete,
		"update":         testAccStackHCIWindowsVirtualMachine_update,
		"requiresImport": testAccStackHCIWindowsVirtualMachine_requiresImport,
	}
	for name, m := range testCases {
		t.Run(name, func(t *testing.T) {
			m(t)
		})
	}
}

func testAccStackHCIWindowsVirtualMachine_basic(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	if os.Getenv(vmImageIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", vmImageIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.admin_password"),
	})
}

func testAccStackHCIWindowsVirtualMachine_complete(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	if os.Getenv(vmImageIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", vmImageIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.admin_password", "http_proxy_configuration.0.http_proxy", "http_proxy_configuration.0.https_proxy"),
	})
}

func testAccStackHCIWindowsVirtualMachine_update(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	if os.Getenv(vmImageIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", vmImageIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.admin_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.admin_password", "http_proxy_configuration.0.http_proxy", "http_proxy_configuration.0.https_proxy"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.admin_password"),
	})
}

func testAccStackHCIWindowsVirtualMachine_requiresImport(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	if os.Getenv(vmImageIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", vmImageIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r StackHCIWindowsVirtualMachineResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StackHCIVirtualMachineID(state.ID)
	if err != nil {
		return nil, err
	}

	arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
	scopeId := commonids.NewScopeID(arcMachineId.ID())
	resp, err := client.AzureStackHCI.VirtualMachineInstances.Get(ctx, scopeId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCIWindowsVirtualMachineResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_windows_virtual_machine" "test" {
  arc_machine_id     = azurerm_arc_machine.test.id
  custom_location_id = %[3]q

  hardware_profile {
    vm_size          = "Custom"
    processor_number = 2
    memory_in_mb     = 8192
  }

  network_profile {
    network_interface_ids = [azurerm_stack_hci_network_interface.test.id]
  }

  os_profile {
    admin_username = "adminuser"
    admin_password = "!password!@#$"
    computer_name  = "testvm"
  }

  storage_profile {
    data_disk_ids = [azurerm_stack_hci_virtual_hard_disk.test.id]
    image_id      = %[4]q
  }

  depends_on = [azurerm_role_assignment.test]

  lifecycle {
    ignore_changes = [storage_profile.0.vm_config_storage_path_id]
  }
}
`, template, data.RandomInteger, os.Getenv(customLocationIdEnv), os.Getenv(vmImageIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_windows_virtual_machine" "import" {
  arc_machine_id     = azurerm_stack_hci_windows_virtual_machine.test.arc_machine_id
  custom_location_id = %[3]q

  hardware_profile {
    vm_size          = "Custom"
    processor_number = 2
    memory_in_mb     = 8192
  }

  network_profile {
    network_interface_ids = [azurerm_stack_hci_network_interface.test.id]
  }

  os_profile {
    admin_username = "adminuser"
    computer_name  = "testvm"
  }

  storage_profile {
    data_disk_ids = [azurerm_stack_hci_virtual_hard_disk.test.id]
    image_id      = %[4]q
  }

  depends_on = [azurerm_role_assignment.test]

  lifecycle {
    ignore_changes = [storage_profile.0.vm_config_storage_path_id]
  }
}
`, config, data.RandomInteger, os.Getenv(customLocationIdEnv), os.Getenv(vmImageIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "tls_private_key" "rsa-4096-example" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "azurerm_stack_hci_virtual_hard_disk" "test2" {
  name                = "acctest-vhd2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  disk_size_in_gb     = 2

  lifecycle {
    ignore_changes        = [storage_path_id]
    create_before_destroy = true
  }
}

resource "azurerm_stack_hci_network_interface" "test2" {
  name                = "acctest-ni2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q

  ip_configuration {
    subnet_id = azurerm_stack_hci_logical_network.test.id
  }

  lifecycle {
    ignore_changes        = [mac_address, ip_configuration.0.private_ip_address]
    create_before_destroy = true
  }
}

resource "azurerm_stack_hci_windows_virtual_machine" "test" {
  arc_machine_id     = azurerm_arc_machine.test.id
  custom_location_id = %[3]q

  hardware_profile {
    vm_size          = "Custom"
    processor_number = 2
    memory_in_mb     = 8192
  }

  network_profile {
    network_interface_ids = [azurerm_stack_hci_network_interface.test.id, azurerm_stack_hci_network_interface.test2.id]
  }

  os_profile {
    admin_username = "adminuser"
    admin_password = "!password!@#$"
    computer_name  = "testvm"
  }

  storage_profile {
    data_disk_ids = [azurerm_stack_hci_virtual_hard_disk.test.id, azurerm_stack_hci_virtual_hard_disk.test2.id]
    image_id      = %[4]q
  }

  depends_on = [azurerm_role_assignment.test]

  lifecycle {
    ignore_changes = [storage_profile.0.vm_config_storage_path_id]
  }
}
`, template, data.RandomInteger, os.Getenv(customLocationIdEnv), os.Getenv(vmImageIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "tls_private_key" "rsa-4096-example" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "azurerm_stack_hci_storage_path" "test" {
  name                = "acctest-sp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp2-%[4]s"
}

resource "azurerm_stack_hci_virtual_hard_disk" "test2" {
  name                = "acctest-vhd2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  disk_size_in_gb     = 2

  lifecycle {
    ignore_changes = [storage_path_id]
  }
}

resource "azurerm_stack_hci_network_interface" "test2" {
  name                = "acctest-ni2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  dns_servers         = ["192.168.1.254"]

  ip_configuration {
    private_ip_address = "192.168.1.18"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }

  lifecycle {
    ignore_changes = [mac_address]
  }
}

resource "azurerm_stack_hci_windows_virtual_machine" "test" {
  arc_machine_id      = azurerm_arc_machine.test.id
  custom_location_id  = %[3]q
  tpm_enabled         = true
  secure_boot_enabled = true

  hardware_profile {
    vm_size          = "Custom"
    processor_number = 2
    memory_in_mb     = 8192
    dynamic_memory {
      maximum_memory_in_mb            = 8192
      minimum_memory_in_mb            = 512
      target_memory_buffer_percentage = 20
    }
  }

  network_profile {
    network_interface_ids = [azurerm_stack_hci_network_interface.test.id, azurerm_stack_hci_network_interface.test2.id]
  }

  os_profile {
    admin_username                    = "adminuser"
    admin_password                    = "!password!@#$"
    computer_name                     = "testvm2"
    automatic_update_enabled          = true
    time_zone                         = "UTC"
    provision_vm_agent_enabled        = true
    provision_vm_config_agent_enabled = true
    ssh_public_key {
      path     = "C:\\Users\\adminuser\\.ssh\\rsa.pub"
      key_data = tls_private_key.rsa-4096-example.public_key_openssh
    }
  }

  storage_profile {
    data_disk_ids             = [azurerm_stack_hci_virtual_hard_disk.test.id]
    image_id                  = %[5]q
    vm_config_storage_path_id = azurerm_stack_hci_storage_path.test.id
  }

  http_proxy_configuration {
    http_proxy  = "http://proxy.example.com:3128"
    https_proxy = "http://proxy.example.com:3128"
    no_proxy = [
      "localhost",
      "127.0.0.1",
      "mcr.microsoft.com"
    ]
    trusted_ca = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ6RENDQVMyZ0F3SUJBZ0lCQVRBS0JnZ3Foa2pPUFFRREJEQVFNUTR3REFZRFZRUUtFd1ZGVGtOUFRUQWUKRncwd09URXhNVEF5TXpBd01EQmFGdzB4TURBMU1Ea3lNekF3TURCYU1CQXhEakFNQmdOVkJBb1RCVVZPUTA5TgpNSUdiTUJBR0J5cUdTTTQ5QWdFR0JTdUJCQUFqQTRHR0FBUUJBcUN1Um94NU4zTVRVOHdUdUllSUJYRjdpTW5oCm50cW1HVktRMGhmUUZEUUd2K0x5ZHVvN0pQcUZwL1kyamxYU2ROckFkejVXeGJyWStrRHhJcGtCUXRJQWtJREQKWlZtVHVlcTNaREFmY0dkRU5uek5KVkNhUGxIWEpMdkVFSU5jb0prVU8rK2NWeXl3ZHJlVkpjNjd2aE54MVRkWApWM3BwN2YrUmJPbU5LYm5WUkJ5ak5UQXpNQTRHQTFVZER3RUIvd1FFQXdJSGdEQVRCZ05WSFNVRUREQUtCZ2dyCkJnRUZCUWNEQVRBTUJnTlZIUk1CQWY4RUFqQUFNQW9HQ0NxR1NNNDlCQU1FQTRHTUFEQ0JpQUpDQWJiYjdzdkkKNXR1aEN5QTNqUVRTZ0E4enB2azBZV05Ya1owN3h6ZFY4amRNTXVtQ2FXOXljRUlxSjVLU3F1dVBoVXc5b2VregpCNTFkYXliVjFWUVhWVmRWQWtJQStrTU1TSnp3dHpIcU5BVVRtaVpQY2c3SDh2MUFTbDR0UjZscEtUcFVQWTJYCmxYT0N0MllmNGRzRnNpanV2emJKQmR4NzVkNEVmNVRSSFBjZytQSE5aZ2c9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
  }

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomInteger, os.Getenv(customLocationIdEnv), data.RandomString, os.Getenv(vmImageIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-vm-%[2]d"
  location = %[1]q
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["192.168.1.254"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "192.168.1.0/24"
    ip_pool {
      start = "192.168.1.0"
      end   = "192.168.1.255"
    }

    route {
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "192.168.1.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  dns_servers         = ["192.168.1.254"]
  mac_address         = "02:ec:01:0c:00:08"

  ip_configuration {
    private_ip_address = "192.168.1.15"
    subnet_id          = azurerm_stack_hci_logical_network.test.id
  }
}

// service principal of 'Microsoft.AzureStackHCI Resource Provider'
data "azuread_service_principal" "hciRp" {
  client_id = "1412d89f-b8a8-4111-b4fd-e82905cbd85d"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Azure Connected Machine Resource Manager"
  principal_id         = data.azuread_service_principal.hciRp.object_id
}

resource "azurerm_stack_hci_virtual_hard_disk" "test" {
  name                = "acctest-vhd-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  disk_size_in_gb     = 2

  lifecycle {
    ignore_changes = [storage_path_id]
  }
}

resource "azurerm_arc_machine" "test" {
  name                = "acctest-hcivm-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "HCI"

  identity {
    type = "SystemAssigned"
  }
}
`, data.Locations.Primary, data.RandomInteger, os.Getenv(customLocationIdEnv))
}
