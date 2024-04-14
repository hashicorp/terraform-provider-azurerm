// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualMachineGalleryApplicationAssignmentResource struct{}

func TestAccVirtualMachineGalleryApplicationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_gallery_application_assignment", "test")
	r := VirtualMachineGalleryApplicationAssignmentResource{}
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

func TestAccVirtualMachineGalleryApplicationAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_gallery_application_assignment", "test")
	r := VirtualMachineGalleryApplicationAssignmentResource{}
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

func TestAccVirtualMachineGalleryApplicationAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_gallery_application_assignment", "test")
	r := VirtualMachineGalleryApplicationAssignmentResource{}
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

func TestAccVirtualMachineGalleryApplicationAssignment_order(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_gallery_application_assignment", "test")
	r := VirtualMachineGalleryApplicationAssignmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.order(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.order(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r VirtualMachineGalleryApplicationAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VirtualMachineGalleryApplicationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	virtualMachine, err := client.Compute.VirtualMachinesClient.Get(ctx, id.VirtualMachineId, virtualmachines.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %s", id, err)
	}

	model := virtualMachine.Model
	if model == nil {
		return nil, fmt.Errorf("retrieving model of %q: %s", id, err)
	}

	if model.Properties == nil || model.Properties.ApplicationProfile == nil || model.Properties.ApplicationProfile.GalleryApplications == nil {
		return pointer.To(false), nil
	}

	for _, application := range pointer.From(model.Properties.ApplicationProfile.GalleryApplications) {
		if strings.EqualFold(id.GalleryApplicationVersionId.ID(), application.PackageReferenceId) {
			return pointer.To(true), nil
		}
	}

	return pointer.To(false), nil
}

func (r VirtualMachineGalleryApplicationAssignmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_gallery_application_assignment" "test" {
  gallery_application_version_id = azurerm_gallery_application_version.test.id
  virtual_machine_id             = azurerm_linux_virtual_machine.test.id
}
`, template)
}

func (r VirtualMachineGalleryApplicationAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_gallery_application_assignment" "import" {
  gallery_application_version_id = azurerm_virtual_machine_gallery_application_assignment.test.gallery_application_version_id
  virtual_machine_id             = azurerm_virtual_machine_gallery_application_assignment.test.virtual_machine_id
}
`, config)
}

func (r VirtualMachineGalleryApplicationAssignmentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_gallery_application_assignment" "test" {
  gallery_application_version_id = azurerm_gallery_application_version.test.id
  virtual_machine_id             = azurerm_linux_virtual_machine.test.id
  configuration_blob_uri         = azurerm_storage_blob.test.id
  order                          = 1
  tag                            = "app"
}
`, template)
}

func (r VirtualMachineGalleryApplicationAssignmentResource) order(data acceptance.TestData, order int) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_gallery_application_assignment" "test" {
  gallery_application_version_id = azurerm_gallery_application_version.test.id
  virtual_machine_id             = azurerm_linux_virtual_machine.test.id
  order                          = %d
}
`, template, order)
}

func (r VirtualMachineGalleryApplicationAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accteststr%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "script"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%[2]d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_shared_image_gallery.test.location
  supported_os_type = "Linux"
}

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  source {
    media_link = azurerm_storage_blob.test.id
  }

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
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
    sku       = "16.04-LTS"
    version   = "latest"
  }

  lifecycle {
    ignore_changes = [
      gallery_application, tags, identity
    ]
  }
}
`, LinuxVirtualMachineResource{}.template(data), data.RandomInteger, data.RandomString)
}
