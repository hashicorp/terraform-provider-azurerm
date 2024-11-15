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

type StackHCIWindowsVirtualMachineResource struct{}

func TestAccStackHCIWindowsVirtualMachine_basic(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCIWindowsVirtualMachine_complete(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCIWindowsVirtualMachine_update(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTag(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCIWindowsVirtualMachine_requiresImport(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_windows_virtual_machine", "test")
	r := StackHCIWindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
	clusterClient := client.AzureStackHCI.VirtualMachineInstances
	id, err := parse.StackHCIVirtualMachineID(state.ID)
	if err != nil {
		return nil, err
	}

	arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroup, id.MachineName)
	scopeId := commonids.NewScopeID(arcMachineId.ID())
	resp, err := clusterClient.Get(ctx, scopeId)
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
  arc_machine_id = azurerm_arc_machine.test.id
  custom_location_id = %[3]q

  hardware_profile {
    vm_size = "Custom"
    processor_number = 2
    memory_mb = 8192
  }

  network_profile {
    network_interface_ids = [azurerm_stack_hci_network_interface.test.id]
  }

  os_profile {
    admin_username = "adminuser"
    admin_password = "!password!@#$"
    computer_name  = "1a"
  }

  storage_profile {
    data_disk_ids = [azurerm_stack_hci_virtual_hard_disk.test.id]
    image_id = azurerm_stack_hci_marketplace_gallery_image.test.id
  }

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_virtual_hard_disk" "test" {
  name                = "acctest-vhd-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  disk_size_in_gb     = 2

  tags = {
    foo = "bar"
  }

  lifecycle {
    ignore_changes = [storage_path_id]
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) updateTag(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_virtual_hard_disk" "test" {
  name                = "acctest-vhd-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  disk_size_in_gb     = 2

  tags = {
    env = "test"
    foo = "bar"
  }

  lifecycle {
    ignore_changes = [storage_path_id]
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_storage_path" "test" {
  name                = "acctest-sp-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-%[2]s"
}

resource "azurerm_stack_hci_virtual_hard_disk" "test" {
  name                     = "acctest-vhd-%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  custom_location_id       = %[3]q
  disk_size_in_gb          = 2
  dynamic_enabled          = false
  hyperv_generation        = "V2"
  physical_sector_in_bytes = 4096
  logical_sector_in_bytes  = 512
  block_size_in_bytes      = 1024
  disk_file_format         = "vhdx"
  storage_path_id          = azurerm_stack_hci_storage_path.test.id

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv))
}

func (r StackHCIWindowsVirtualMachineResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_virtual_hard_disk" "import" {
  name                = azurerm_stack_hci_virtual_hard_disk.test.name
  resource_group_name = azurerm_stack_hci_virtual_hard_disk.test.resource_group_name
  location            = azurerm_stack_hci_virtual_hard_disk.test.location
  custom_location_id  = azurerm_stack_hci_virtual_hard_disk.test.custom_location_id
  disk_size_in_gb     = azurerm_stack_hci_virtual_hard_disk.test.disk_size_in_gb
}
`, config)
}

func (r StackHCIWindowsVirtualMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-vm-%[2]s"
  location = %[1]q
}

resource "azurerm_stack_hci_logical_network" "test" {
  name                = "acctest-ln-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  virtual_switch_name = "ConvergedSwitch(managementcompute)"

  subnet {
    ip_allocation_method = "Dynamic"
  }
}

resource "azurerm_stack_hci_network_interface" "test" {
  name                = "acctest-ni-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q

  ip_configuration {
    subnet_id = azurerm_stack_hci_logical_network.test.id
  }

  lifecycle {
    ignore_changes = [mac_address]
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

//resource "azurerm_stack_hci_marketplace_gallery_image" "test" {
//  name                = "acctest-mgi-%[2]s"
//  resource_group_name = azurerm_resource_group.test.name
//  location            = azurerm_resource_group.test.location
//  custom_location_id  = %[3]q
//  hyperv_generation   = "V2"
//  os_type             = "Windows"
//  version             = "20348.2655.240905"
//  identifier {
//    publisher = "MicrosoftWindowsServer"
//    offer     = "WindowsServer"
//    sku       = "2022-datacenter-azure-edition-core"
//  }
//
//  depends_on = [azurerm_role_assignment.test]
//}

resource "azurerm_stack_hci_virtual_hard_disk" "test" {
  name                = "acctest-vhd-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  disk_size_in_gb     = 2

  tags = {
    env = "test"
    foo = "bar"
  }

  lifecycle {
    ignore_changes = [storage_path_id]
  }
}

resource "azurerm_arc_machine" "test" {
  name                = "acctest-hcivm-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "HCI"

  //identity {
  //  type = "SystemAssigned"
  //}
}
`, data.Locations.Primary, data.RandomString, os.Getenv(customLocationIdEnv))
}
