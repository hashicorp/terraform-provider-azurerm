package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualharddisks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCIVirtualHardDiskResource struct{}

func TestAccStackHCIVirtualHardDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_virtual_hard_disk", "test")
	r := StackHCIVirtualHardDiskResource{}

	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

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

func TestAccStackHCIVirtualHardDisk_complete(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_virtual_hard_disk", "test")
	r := StackHCIVirtualHardDiskResource{}

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

func TestAccStackHCIVirtualHardDisk_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_virtual_hard_disk", "test")
	r := StackHCIVirtualHardDiskResource{}

	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func TestAccStackHCIVirtualHardDisk_requiresImport(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_virtual_hard_disk", "test")
	r := StackHCIVirtualHardDiskResource{}

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

func (r StackHCIVirtualHardDiskResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.VirtualHardDisks
	id, err := virtualharddisks.ParseVirtualHardDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCIVirtualHardDiskResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_virtual_hard_disk" "import" {
  name                = azurerm_stack_hci_virtual_hard_disk.test.name
  resource_group_name = azurerm_stack_hci_virtual_hard_disk.test.resource_group_name
  location            = azurerm_stack_hci_virtual_hard_disk.test.location
  custom_location_id  = azurerm_stack_hci_virtual_hard_disk.test.custom_location_id
  disk_size_in_gb = azurerm_stack_hci_virtual_hard_disk.test.disk_size_in_gb
}
`, config)
}

func (r StackHCIVirtualHardDiskResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_virtual_hard_disk" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  disk_size_in_gb     = 2
  lifecycle {
    ignore_changes = [storage_path_id]
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCIVirtualHardDiskResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_virtual_hard_disk" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  disk_size_in_gb     = 2


  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCIVirtualHardDiskResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}

variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-ln-${var.random_string}"
  location = var.primary_location
}
`, data.Locations.Primary, data.RandomString)
}
