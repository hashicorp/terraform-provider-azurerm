// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ImageResource struct{}

func TestAccImage_standaloneImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupManagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.testsource"),
			),
		},
		{
			Config: r.standaloneImageProvision(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccImage_standaloneImage_hyperVGeneration_V2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupManagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.testsource"),
			),
		},
		{
			Config: r.standaloneImageProvision(data, "V2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccImage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupManagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.testsource"),
			),
		},
		{
			Config: r.standaloneImageProvision(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.standaloneImageRequiresImport),
	})
}

func TestAccImage_customImageFromVMWithExplicitOsDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupManagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.testsource"),
			),
		},
		{
			Config: r.customImageFromVMWithManagedDisksProvision(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testdestination"),
			),
		},
	})
}

func TestAccImage_customImageFromVMWithManagedDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setupManagedDisks(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.testsource"),
			),
		},
		{
			Config: r.customImageFromManagedDiskVMProvision(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testdestination"),
			),
		},
	})
}

func TestAccImage_customImageFromVMSSWithExplicitOsDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setupManagedDisks(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.testsource"),
			),
		},
		{
			Config: r.customImageFromVMSSWithUnmanagedDisksProvision(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineScaleSetExists, "azurerm_linux_virtual_machine_scale_set.testdestination"),
			),
		},
	})
}

func TestAccImage_standaloneImageEncrypt(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupManagedDisksWithKV(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_linux_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.testsource"),
			),
		},
		{
			Config: r.standaloneImageEncrypt(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ImageResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := images.ParseImageID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.ImagesClient.Get(ctx, *id, images.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving Compute Image %q", id)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ImageResource) generalizeVirtualMachine() func(context.Context, *clients.Client, *pluginsdk.InstanceState) error {
	return func(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := virtualmachines.ParseVirtualMachineID(state.ID)
		if err != nil {
			return err
		}

		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()
		}

		log.Printf("[DEBUG] Deallocating VM..")
		if err := client.Compute.VirtualMachinesClient.DeallocateThenPoll(ctx, *id, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
			return fmt.Errorf("Bad: deallocating %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] Generalizing VM..")
		if _, err = client.Compute.VirtualMachinesClient.Generalize(ctx, *id); err != nil {
			return fmt.Errorf("Bad: Generalizing %s: %+v", *id, err)
		}

		return nil
	}
}

func (ImageResource) virtualMachineExists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := virtualmachines.ParseVirtualMachineID(state.ID)
	if err != nil {
		return err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}
	resp, err := client.Compute.VirtualMachinesClient.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return nil
}

func (ImageResource) virtualMachineScaleSetExists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(state.ID)
	if err != nil {
		return err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}
	resp, err := client.Compute.VirtualMachineScaleSetsClient.Get(ctx, *id, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return nil
}

func (r ImageResource) setupManagedDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_linux_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  size                  = "Standard_D1_v2"

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    name                 = "myosdisk1${local.number}"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  computer_name                   = "mdimagetestsource"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

data "azurerm_managed_disk" "testsource" {
  name                = azurerm_linux_virtual_machine.testsource.os_disk.0.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (r ImageResource) setupManagedDisksWithKV(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults       = false
      purge_soft_delete_on_destroy          = false
      purge_soft_deleted_keys_on_destroy    = false
      purge_soft_deleted_secrets_on_destroy = false
    }
  }
}

%s

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_linux_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  size                  = "Standard_D1_v2"

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    name                 = "myosdisk1${local.number}"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  computer_name                   = "mdimagetestsource"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

data "azurerm_managed_disk" "testsource" {
  name                = azurerm_linux_virtual_machine.testsource.os_disk.0.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctdatadisk-${local.number}"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}
`, template)
}

func (r ImageResource) standaloneImageProvision(data acceptance.TestData, hyperVGen string) string {
	hyperVGenAtt := ""
	if hyperVGen != "" {
		hyperVGenAtt = fmt.Sprintf(`hyper_v_generation = "%s"`, hyperVGen)
	}

	osDisk := `
  os_disk {
    os_type         = "Linux"
    os_state        = "Generalized"
    managed_disk_id = data.azurerm_managed_disk.testsource.id
    size_gb         = 30
    caching         = "None"
    storage_type    = "Standard_LRS"
  }`

	template := r.setupManagedDisks(data)

	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

%[2]s

%[3]s

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template, hyperVGenAtt, osDisk)
}

func (r ImageResource) standaloneImageRequiresImport(data acceptance.TestData) string {
	template := r.standaloneImageProvision(data, "")

	osDisk := `
  os_disk {
	os_type         = "Linux"
	os_state        = "Generalized"
	managed_disk_id = data.azurerm_managed_disk.testsource.id
	size_gb         = 30
	caching         = "None"
	storage_type    = "Standard_LRS"
  }`

	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "import" {
  name                = azurerm_image.test.name
  location            = azurerm_image.test.location
  resource_group_name = azurerm_image.test.resource_group_name

%[2]s

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template, osDisk)
}

func (r ImageResource) customImageFromVMWithManagedDisksProvision(data acceptance.TestData) string {
	osDisk := `
  os_disk {
	os_type         = "Linux"
	os_state        = "Generalized"
	managed_disk_id = data.azurerm_managed_disk.testsource.id
	size_gb         = 30
	caching         = "None"
	storage_type    = "Standard_LRS"
  }`

	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "testdestination" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

%[2]s

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_network_interface" "testdestination" {
  name                = "acctnicdest-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration2"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "testdestination" {
  name                  = "acctvm"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testdestination.id]
  size                  = "Standard_D1_v2"

  source_image_id = azurerm_image.testdestination.id

  os_disk {
    name                 = "myosdisk2"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  computer_name                   = "mdimagetestsource"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, r.setupManagedDisks(data), osDisk)
}

func (r ImageResource) customImageFromManagedDiskVMProvision(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_image" "testdestination" {
  name                      = "acctestdest-${local.number}"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.testsource.id

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

resource "azurerm_network_interface" "testdestination" {
  name                = "acctnicdest-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration2"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "testdestination" {
  name                  = "testdestination"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testdestination.id]
  size                  = "Standard_D1_v2"

  source_image_id = azurerm_image.testdestination.id

  os_disk {
    name                 = "myosdisk2${local.number}"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  computer_name                   = "mdimagetestdest"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, r.setupManagedDisks(data))
}

func (r ImageResource) customImageFromVMSSWithUnmanagedDisksProvision(data acceptance.TestData) string {
	osDisk := `
  os_disk {
	os_type         = "Linux"
	os_state        = "Generalized"
	managed_disk_id = data.azurerm_managed_disk.testsource.id
	size_gb         = 30
	caching         = "None"
	storage_type    = "Standard_LRS"
  }`

	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "testdestination" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

%[2]s

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_linux_virtual_machine_scale_set" "testdestination" {
  name                            = "testdestination"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  sku                             = "Standard_D1_v2"
  instances                       = 2
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false

  source_image_id = azurerm_image.testdestination.id

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = azurerm_subnet.test.id
      primary   = true
    }
  }
}
`, r.setupManagedDisks(data), osDisk)
}

func (r ImageResource) standaloneImageEncrypt(data acceptance.TestData) string {
	osDisk := `
  os_disk {
    os_type                = "Linux"
    os_state               = "Generalized"
    managed_disk_id        = data.azurerm_managed_disk.testsource.id
    size_gb                = 30
    caching                = "None"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
    storage_type           = "Standard_LRS"
    }`

	dataDisk := `
  data_disk {
    managed_disk_id        = azurerm_managed_disk.test.id
    size_gb                = 10
    caching                = "None"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
    storage_type           = "Standard_LRS"
    }`

	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                        = "acctest%[3]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  rbac_authorization_enabled  = false
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  purge_protection_enabled    = true
  soft_delete_retention_days  = 7
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_image" "test" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  %[4]s

  %[5]s

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, r.setupManagedDisksWithKV(data), data.RandomInteger, data.RandomString, osDisk, dataDisk)
}

func (ImageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  number            = "%d"
  location          = %q
  domain_name_label = "acctestvm-%s"
  random_string     = %q
  admin_username    = "testadmin%d"
  admin_password    = "Password1234!%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.number}"
  location = local.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-${local.number}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = local.domain_name_label
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}
