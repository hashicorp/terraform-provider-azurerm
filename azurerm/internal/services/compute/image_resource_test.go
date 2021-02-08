package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ImageResource struct {
}

func TestAccImage_standaloneImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "testsource")
	r := ImageResource{}

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupUnmanagedDisks(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImageProvision(data, "LRS", ""),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_image.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"delete_data_disks_on_termination",
			"delete_os_disk_on_termination",
		),
	})
}

func TestAccImage_standaloneImage_hyperVGeneration_V2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "testsource")
	r := ImageResource{}
	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupUnmanagedDisks(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImageProvision(data, "LRS", "V2"),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_image.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"delete_data_disks_on_termination",
			"delete_os_disk_on_termination",
		),
	})
}

func TestAccImage_standaloneImageZoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "testsource")
	r := ImageResource{}
	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setupUnmanagedDisks(data, "ZRS"),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImageProvision(data, "ZRS", ""),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_image.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccImage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "testsource")
	r := ImageResource{}

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupUnmanagedDisks(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImageProvision(data, "LRS", ""),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_image.test").ExistsInAzure(r),
			),
		},
		{
			Config:      r.standaloneImageRequiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_image"),
		},
	})
}

func TestAccImage_customImageFromVMWithUnmanagedDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "testsource")
	r := ImageResource{}

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config: r.setupUnmanagedDisks(data, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.customImageFromVMWithUnmanagedDisksProvision(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testdestination"),
			),
		},
	})
}

func TestAccImage_customImageFromVMWithManagedDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "testsource")
	r := ImageResource{}

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setupManagedDisks(data),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.customImageFromManagedDiskVMProvision(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testdestination"),
			),
		},
	})
}

func TestAccImage_customImageFromVMSSWithUnmanagedDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine", "testsource")
	r := ImageResource{}

	resourceGroup := fmt.Sprintf("acctestRG-%d", data.RandomInteger)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", data.RandomInteger)
	sshPort := "22"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// need to create a vm and then reference it in the image creation
			Config:  r.setupUnmanagedDisks(data, "LRS"),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_virtual_machine.testsource"),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.customImageFromVMSSWithUnmanagedDisksProvision(data),
			Check: resource.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineScaleSetExists, "azurerm_virtual_machine_scale_set.testdestination"),
			),
		},
	})
}

func (ImageResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	name := id.Path["images"]

	resp, err := clients.Compute.ImagesClient.Get(ctx, resGroup, name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Image %q", id)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ImageResource) virtualMachineExists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	id, err := parse.VirtualMachineID(state.ID)
	if err != nil {
		return err
	}

	resp, err := client.Compute.VMClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s does not exist", *id)
		}

		return fmt.Errorf("Bad: Get on client: %+v", err)
	}

	return nil
}

func (ImageResource) virtualMachineScaleSetExists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	id, err := parse.VirtualMachineScaleSetID(state.ID)
	if err != nil {
		return err
	}

	resp, err := client.Compute.VMScaleSetClient.Get(ctx, id.ResourceGroup, id.Name)
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

func (r ImageResource) setupUnmanagedDisks(data acceptance.TestData, storageType string) string {
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

resource "azurerm_storage_account" "test" {
  name                     = "accsa${var.random_string}"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "%s"
  allow_blob_public_access = true
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

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
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
}
`, template, storageType)
}

func (r ImageResource) standaloneImageProvision(data acceptance.TestData, storageType string, hyperVGen string) string {
	hyperVGenAtt := ""
	if hyperVGen != "" {
		hyperVGenAtt = fmt.Sprintf(`hyper_v_generation = "%s"`, hyperVGen)
	}

	template := r.setupUnmanagedDisks(data, storageType)
	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                = "accteste"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  zone_resilient      = %t

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
`, template, storageType == "ZRS", hyperVGenAtt)
}

func (r ImageResource) standaloneImageRequiresImport(data acceptance.TestData) string {
	template := r.standaloneImageProvision(data, "LRS", "")
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
	template := r.setupUnmanagedDisks(data, "LRS")
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
  name                = "acctnicdest-${local.random}"
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
	template := r.setupUnmanagedDisks(data, "LRS")
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
    computer_name_prefix = "testvm${local.random}"
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
  address_prefix       = "10.0.2.0/24"
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
