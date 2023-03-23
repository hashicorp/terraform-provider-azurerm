package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2022-08-01/compute"
)

func TestAccWindowsVirtualMachine_imageFromImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// create the original VM
			Config: r.imageFromExistingMachinePrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.generalizeVirtualMachine, "azurerm_windows_virtual_machine.source"),
			),
		},
		{
			// then create an image from that VM, and then create a VM from that image
			Config: r.imageFromImage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_imageFromPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}
	publisher := "plesk"
	offer := "plesk-onyx-windows"
	sku := "plsk-win-hst-azr-m"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientWithoutResource(r.cancelExistingAgreement(publisher, offer, sku)),
			),
		},
		{
			Config: r.imageFromPlan(data, publisher, offer, sku),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_imageFromSharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// create the original VM
			Config: r.imageFromExistingMachinePrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.generalizeVirtualMachine, "azurerm_windows_virtual_machine.source"),
			),
		},
		{
			// then create an image from that VM, and then create a VM from that image
			Config: r.imageFromSharedImageGallery(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachine_imageFromSourceImageReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func (WindowsVirtualMachineResource) imageFromExistingMachineDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  vm_name = "acctvm-%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = local.vm_name
}

resource "azurerm_network_interface" "public" {
  name                = "acctnicsource-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, data.RandomString, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r WindowsVirtualMachineResource) imageFromExistingMachinePrep(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "source" {
  name                = "${local.vm_name}1"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.public.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, r.imageFromExistingMachineDependencies(data))
}

func (r WindowsVirtualMachineResource) imageFromImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_windows_virtual_machine.source.id
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = "${local.vm_name}2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  source_image_id     = azurerm_image.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
}
`, r.imageFromExistingMachinePrep(data))
}

func (r WindowsVirtualMachineResource) imageFromPlan(data acceptance.TestData, publisher string, offer string, sku string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_marketplace_agreement" "test" {
  publisher = "%[2]s"
  offer     = "%[3]s"
  plan      = "%[4]s"
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  plan {
    publisher = "%[2]s"
    product   = "%[3]s"
    name      = "%[4]s"
  }

  source_image_reference {
    publisher = "%[2]s"
    offer     = "%[3]s"
    sku       = "%[4]s"
    version   = "latest"
  }

  depends_on = ["azurerm_marketplace_agreement.test"]
}
`, r.template(data), publisher, offer, sku)
}

func (r WindowsVirtualMachineResource) imageFromSharedImageGallery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_windows_virtual_machine.source.id
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Windows"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = "${local.vm_name}2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  source_image_id     = azurerm_shared_image_version.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
}
`, r.imageFromExistingMachinePrep(data))
}

func (r WindowsVirtualMachineResource) imageFromSourceImageReference(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, r.template(data))
}

func (WindowsVirtualMachineResource) empty() string {
	return `
provider "azurerm" {
  features {}
}
`
}

func (WindowsVirtualMachineResource) generalizeVirtualMachine(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.VirtualMachineID(state.ID)
	if err != nil {
		return err
	}

	command := []string{
		"$cmd = \"$Env:SystemRoot\\system32\\sysprep\\sysprep.exe\"",
		"$args = \"/generalize /oobe /mode:vm /quit\"",
		"Start-Process powershell -Argument \"$cmd $args\" -Wait",
	}
	runCommand := compute.RunCommandInput{
		CommandID: utils.String("RunPowerShellScript"),
		Script:    &command,
	}

	future, err := client.Compute.VMClient.RunCommand(ctx, id.ResourceGroup, id.Name, runCommand)
	if err != nil {
		return fmt.Errorf("Bad: Error in running sysprep: %+v", err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Compute.VMClient.Client); err != nil {
		return fmt.Errorf("Bad: Error waiting for Windows VM to sysprep: %+v", err)
	}

	// Upgrading to the 2021-07-01 exposed a new hibernate parameter in the GET method
	daFuture, err := client.Compute.VMClient.Deallocate(ctx, id.ResourceGroup, id.Name, utils.Bool(false))
	if err != nil {
		return fmt.Errorf("Bad: Deallocation error: %+v", err)
	}

	if err := daFuture.WaitForCompletionRef(ctx, client.Compute.VMClient.Client); err != nil {
		return fmt.Errorf("Bad: Deallocation error: %+v", err)
	}

	if _, err = client.Compute.VMClient.Generalize(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("Bad: Generalizing error: %+v", err)
	}

	return nil
}

func (r WindowsVirtualMachineResource) cancelExistingAgreement(publisher string, offer string, sku string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		client := clients.Compute.MarketplaceAgreementsClient
		id := parse.NewPlanID(client.SubscriptionID, publisher, offer, sku)

		existing, err := client.Get(ctx, id.AgreementName, id.OfferName, id.Name)
		if err != nil {
			return err
		}

		if props := existing.AgreementProperties; props != nil {
			if accepted := props.Accepted; accepted != nil && *accepted {
				resp, err := client.Cancel(ctx, id.AgreementName, id.OfferName, id.Name)
				if err != nil {
					if utils.ResponseWasNotFound(resp.Response) {
						return fmt.Errorf("marketplace agreement %q does not exist", id)
					}
					return fmt.Errorf("canceling Marketplace Agreement : %+v", err)
				}
			}
		}

		return nil
	}
}
