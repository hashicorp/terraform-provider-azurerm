// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/ssh"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ImageResource struct{}

func TestAccImage_standaloneImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupUnmanagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
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
			Config: r.setupUnmanagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
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
			Config: r.setupUnmanagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
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

func TestAccImage_customImageFromVMWithUnmanagedDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupUnmanagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.customImageFromVMWithUnmanagedDisksProvision(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testdestination"),
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
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.customImageFromManagedDiskVMProvision(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testdestination"),
			),
		},
	})
}

func TestAccImage_customImageFromVMSSWithUnmanagedDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image", "test")
	r := ImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setupUnmanagedDisks(data),
			Destroy: false,
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
			),
		},
		{
			Config: r.customImageFromVMSSWithUnmanagedDisksProvision(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineScaleSetExists, "azurerm_virtual_machine_scale_set.testdestination"),
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
			Config: r.setupUnmanagedDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(data), "azurerm_virtual_machine.testsource"),
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

	return utils.Bool(resp.Model != nil), nil
}

func (ImageResource) generalizeVirtualMachine(data acceptance.TestData) func(context.Context, *clients.Client, *pluginsdk.InstanceState) error {
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

		// these are nested in a Set in the Legacy VM resource, simpler to compute them
		userName := fmt.Sprintf("testadmin%d", data.RandomInteger)
		password := fmt.Sprintf("Password1234!%d", data.RandomInteger)

		// first retrieve the Virtual Machine, since we need to find
		nicIdRaw := state.Attributes["network_interface_ids.0"]
		nicId, err := networkParse.NetworkInterfaceID(nicIdRaw)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Retrieving Network Interface..")
		nic, err := client.Network.InterfacesClient.Get(ctx, nicId.ResourceGroup, nicId.Name, "")
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *nicId, err)
		}

		publicIpRaw := ""
		if props := nic.InterfacePropertiesFormat; props != nil {
			if configs := props.IPConfigurations; configs != nil {
				for _, config := range *props.IPConfigurations {
					if config.InterfaceIPConfigurationPropertiesFormat == nil {
						continue
					}

					if config.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress == nil {
						continue
					}

					if config.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress.ID == nil {
						continue
					}

					publicIpRaw = *config.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress.ID
					break
				}
			}
		}
		if publicIpRaw == "" {
			return fmt.Errorf("retrieving %s: could not determine Public IP Address ID", *nicId)
		}

		log.Printf("[DEBUG] Retrieving Public IP Address %q..", publicIpRaw)
		publicIpId, err := networkParse.PublicIpAddressID(publicIpRaw)
		if err != nil {
			return err
		}

		publicIpAddress, err := client.Network.PublicIPsClient.Get(ctx, publicIpId.ResourceGroup, publicIpId.Name, "")
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *publicIpId, err)
		}
		fqdn := ""
		if props := publicIpAddress.PublicIPAddressPropertiesFormat; props != nil {
			if dns := props.DNSSettings; dns != nil {
				if dns.Fqdn != nil {
					fqdn = *dns.Fqdn
				}
			}
		}
		if fqdn == "" {
			return fmt.Errorf("unable to determine FQDN for %q", *publicIpId)
		}

		log.Printf("[DEBUG] Running Generalization Command..")
		sshGeneralizationCommand := ssh.Runner{
			Hostname: fqdn,
			Port:     22,
			Username: userName,
			Password: password,
			CommandsToRun: []string{
				ssh.LinuxAgentDeprovisionCommand,
			},
		}
		if err := sshGeneralizationCommand.Run(ctx); err != nil {
			return fmt.Errorf("Bad: running generalization command: %+v", err)
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
		ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
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
	id, err := commonids.ParseVirtualMachineScaleSetID(state.ID)
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	resp, err := client.Compute.VMScaleSetClient.Get(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s does not exist", *id)
		}

		return fmt.Errorf("Bad: Get on client: %+v", err)
	}

	return nil
}

func (r ImageResource) setupManagedDisks(data acceptance.TestData) string {
	template := r.template(data)
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

resource "azurerm_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  vm_size               = "Standard_D1_v2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = local.admin_username
    admin_password = local.admin_password
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template)
}

func (r ImageResource) setupUnmanagedDisks(data acceptance.TestData) string {
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

resource "azurerm_storage_account" "test" {
  name                     = "accsa${local.random_string}"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

# NOTE: using the legacy vm resource since this test requires an unmanaged disk
resource "azurerm_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  vm_size               = "Standard_D1_v2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = local.admin_username
    admin_password = local.admin_password
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template)
}

func (r ImageResource) standaloneImageProvision(data acceptance.TestData, hyperVGen string) string {
	hyperVGenAtt := ""
	if hyperVGen != "" {
		hyperVGenAtt = fmt.Sprintf(`hyper_v_generation = "%s"`, hyperVGen)
	}

	template := r.setupUnmanagedDisks(data)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  %s

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template, hyperVGenAtt)
}

func (r ImageResource) standaloneImageRequiresImport(data acceptance.TestData) string {
	template := r.standaloneImageProvision(data, "")
	return fmt.Sprintf(`
%s

resource "azurerm_image" "import" {
  name                = azurerm_image.test.name
  location            = azurerm_image.test.location
  resource_group_name = azurerm_image.test.resource_group_name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template)
}

func (r ImageResource) customImageFromVMWithUnmanagedDisksProvision(data acceptance.TestData) string {
	template := r.setupUnmanagedDisks(data)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "testdestination" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

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

resource "azurerm_virtual_machine" "testdestination" {
  name                  = "acctvm"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testdestination.id]
  vm_size               = "Standard_D1_v2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    id = azurerm_image.testdestination.id
  }

  storage_os_disk {
    name          = "myosdisk1"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = local.admin_username
    admin_password = local.admin_password
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template)
}

func (r ImageResource) customImageFromManagedDiskVMProvision(data acceptance.TestData) string {
	template := r.setupManagedDisks(data)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "testdestination" {
  name                      = "acctestdest-${local.number}"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_virtual_machine.testsource.id

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

resource "azurerm_virtual_machine" "testdestination" {
  name                  = "testdestination"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testdestination.id]
  vm_size               = "Standard_D1_v2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    id = azurerm_image.testdestination.id
  }

  storage_os_disk {
    name          = "myosdisk2"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "mdimagetestdest"
    admin_username = local.admin_username
    admin_password = local.admin_password
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template)
}

func (r ImageResource) customImageFromVMSSWithUnmanagedDisksProvision(data acceptance.TestData) string {
	template := r.setupUnmanagedDisks(data)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "testdestination" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_virtual_machine_scale_set" "testdestination" {
  name                = "testdestination"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm${local.number}"
    admin_username       = local.admin_username
    admin_password       = local.admin_password
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = azurerm_subnet.test.id
      primary   = true
    }
  }

  storage_profile_os_disk {
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    id = azurerm_image.testdestination.id
  }
}
`, template)
}

func (r ImageResource) standaloneImageEncrypt(data acceptance.TestData) string {
	template := r.setupUnmanagedDisks(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                        = "acctest%[3]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  purge_protection_enabled    = true
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

  os_disk {
    os_type                = "Linux"
    os_state               = "Generalized"
    blob_uri               = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb                = 30
    caching                = "None"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, template, data.RandomInteger, data.RandomString)
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
  allocation_method   = "Dynamic"
  domain_name_label   = local.domain_name_label
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger)
}
