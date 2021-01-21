package compute_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"golang.org/x/crypto/ssh"
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
			Config: r.standaloneImage_setup(data, userName, password, hostName, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImage_provision(data, userName, password, hostName, "LRS", ""),
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
			Config: r.standaloneImage_setup(data, userName, password, hostName, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImage_provision(data, userName, password, hostName, "LRS", "V2"),
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
			Config:  r.standaloneImage_setup(data, userName, password, hostName, "ZRS"),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImage_provision(data, userName, password, hostName, "ZRS", ""),
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
			Config: r.standaloneImage_setup(data, userName, password, hostName, "LRS"),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.standaloneImage_provision(data, userName, password, hostName, "LRS", ""),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_image.test").ExistsInAzure(r),
			),
		},
		{
			Config:      r.standaloneImage_requiresImport(data, userName, password, hostName),
			ExpectError: acceptance.RequiresImportError("azurerm_image"),
		},
	})
}

func TestAccImage_customImageVMFromVHD(t *testing.T) {
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
			Config: r.customImage_fromVHD_setup(data, userName, password, hostName),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.customImage_fromVHD_provision(data, userName, password, hostName),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testdestination", true),
			),
		},
	})
}

func TestAccImage_customImageVMFromVM(t *testing.T) {
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
			Config:  r.customImage_fromVM_sourceVM(data, userName, password, hostName),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.customImage_fromVM_destinationVM(data, userName, password, hostName),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testdestination", true),
			),
		},
	})
}

func TestAccImageVMSS_customImageVMSSFromVHD(t *testing.T) {
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
			Config:  r.vmss_customImage_fromVHD_setup(data, userName, password, hostName),
			Destroy: false,
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
				testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, data.Locations.Primary),
			),
		},
		{
			Config: r.vmss_customImage_fromVHD_provision(data, userName, password, hostName),
			Check: resource.ComposeTestCheckFunc(
				testCheckAzureVMSSExists("azurerm_virtual_machine_scale_set.testdestination", true),
			),
		},
	})
}

func testGeneralizeWindowsVMImage(resourceGroup string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		command := []string{
			"$cmd = \"$Env:SystemRoot\\system32\\sysprep\\sysprep.exe\"",
			"$args = \"/generalize /oobe /mode:vm /quit\"",
			"Start-Process powershell -Argument \"$cmd $args\" -Wait",
		}
		runCommand := compute.RunCommandInput{
			CommandID: utils.String("RunPowerShellScript"),
			Script:    &command,
		}

		future, err := client.RunCommand(ctx, resourceGroup, name, runCommand)
		if err != nil {
			return fmt.Errorf("Bad: Error in running sysprep: %+v", err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: Error waiting for Windows VM to sysprep: %+v", err)
		}

		daFuture, err := client.Deallocate(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Deallocation error: %+v", err)
		}

		if err := daFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: Deallocation error: %+v", err)
		}

		_, err = client.Generalize(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Generalizing error: %+v", err)
		}

		return nil
	}
}

// nolint unparam
func testGeneralizeVMImage(resourceGroup string, vmName string, userName string, password string, hostName string, port string, location string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		armClient := acceptance.AzureProvider.Meta().(*clients.Client)
		vmClient := armClient.Compute.VMClient
		ctx := armClient.StopContext

		normalizedLocation := azure.NormalizeLocation(location)
		suffix := armClient.Account.Environment.ResourceManagerVMDNSSuffix
		dnsName := fmt.Sprintf("%s.%s.%s", hostName, normalizedLocation, suffix)

		if err := deprovisionVM(userName, password, dnsName, port); err != nil {
			return fmt.Errorf("Bad: Deprovisioning error %+v", err)
		}

		future, err := vmClient.Deallocate(ctx, resourceGroup, vmName)
		if err != nil {
			return fmt.Errorf("Bad: Deallocating error %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, vmClient.Client); err != nil {
			return fmt.Errorf("Bad: Deallocating error %+v", err)
		}

		if _, err = vmClient.Generalize(ctx, resourceGroup, vmName); err != nil {
			return fmt.Errorf("Bad: Generalizing error %+v", err)
		}

		return nil
	}
}

func deprovisionVM(userName string, password string, hostName string, port string) error {
	// SSH into the machine and execute a waagent deprovisioning command
	var b bytes.Buffer
	cmd := "sudo waagent -verbose -deprovision+user -force"

	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	log.Printf("[INFO] Connecting to %s:%v remote server...", hostName, port)

	hostAddress := strings.Join([]string{hostName, port}, ":")
	client, err := ssh.Dial("tcp", hostAddress, config)
	if err != nil {
		return fmt.Errorf("Bad: deprovisioning error %+v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Bad: deprovisioning error, failure creating session %+v", err)
	}
	defer session.Close()

	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("Bad: deprovisioning error, failure running command %+v", err)
	}

	return nil
}

// nolint unparam
func testCheckAzureVMExists(sourceVM string, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		log.Printf("[INFO] testing MANAGED IMAGE VM EXISTS - BEGIN.")

		vmRs, vmOk := s.RootModule().Resources[sourceVM]
		if !vmOk {
			return fmt.Errorf("VM Not found: %s", sourceVM)
		}
		vmName := vmRs.Primary.Attributes["name"]

		resourceGroup, hasResourceGroup := vmRs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for VM: %s", vmName)
		}

		resp, err := client.Get(ctx, resourceGroup, vmName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Bad: VM %q (resource group %q) does not exist", vmName, resourceGroup)
		}
		if resp.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Bad: VM %q (resource group %q) still exists", vmName, resourceGroup)
		}

		log.Printf("[INFO] testing MANAGED IMAGE VM EXISTS - END.")

		return nil
	}
}

func (t ImageResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func testCheckAzureVMSSExists(sourceVMSS string, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vmssClient := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMScaleSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		log.Printf("[INFO] testing MANAGED IMAGE VMSS EXISTS - BEGIN.")

		vmRs, vmOk := s.RootModule().Resources[sourceVMSS]
		if !vmOk {
			return fmt.Errorf("VMSS Not found: %s", sourceVMSS)
		}
		vmssName := vmRs.Primary.Attributes["name"]

		resourceGroup, hasResourceGroup := vmRs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for VMSS: %s", vmssName)
		}

		resp, err := vmssClient.Get(ctx, resourceGroup, vmssName)
		if err != nil {
			return fmt.Errorf("Bad: Get on vmssClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Bad: VMSS %q (resource group %q) does not exist", vmssName, resourceGroup)
		}
		if resp.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Bad: VMSS %q (resource group %q) still exists", vmssName, resourceGroup)
		}

		log.Printf("[INFO] testing MANAGED IMAGE VMSS EXISTS - END.")

		return nil
	}
}

func (ImageResource) standaloneImage_setup(data acceptance.TestData, userName string, password string, hostName string, storageType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
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
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "%s"
  allow_blob_public_access = true

  tags = {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
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
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, hostName, data.RandomInteger, data.RandomInteger, storageType, userName, password)
}

func (r ImageResource) standaloneImage_provision(data acceptance.TestData, userName string, password string, hostName string, storageType string, hyperVGen string) string {
	hyperVGenAtt := ""
	if hyperVGen != "" {
		hyperVGenAtt = fmt.Sprintf(`hyper_v_generation = "%s"`, hyperVGen)
	}

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
`, r.standaloneImage_setup(data, userName, password, hostName, storageType), storageType == "ZRS", hyperVGenAtt)
}

func (r ImageResource) standaloneImage_requiresImport(data acceptance.TestData, userName string, password string, hostName string) string {
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
`, r.standaloneImage_provision(data, userName, password, hostName, "LRS", ""))
}

func (ImageResource) customImage_fromVHD_setup(data acceptance.TestData, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
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
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
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
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, hostName, data.RandomInteger, data.RandomInteger, userName, password)
}

func (ImageResource) customImage_fromVHD_provision(data acceptance.TestData, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
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
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
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
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "%s"
    admin_password = "%s"
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
  name                = "acctnicdest-%d"
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
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, hostName, data.RandomInteger, data.RandomInteger, userName, password, data.RandomInteger, userName, password)
}

func (ImageResource) customImage_fromVM_sourceVM(data acceptance.TestData, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
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
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, hostName, data.RandomInteger, userName, password)
}

func (ImageResource) customImage_fromVM_destinationVM(data acceptance.TestData, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
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
    admin_username = "%s"
    admin_password = "%s"
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
  name                      = "acctestdest-%d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_virtual_machine.testsource.id

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

resource "azurerm_network_interface" "testdestination" {
  name                = "acctnicdest-%d"
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
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, hostName, data.RandomInteger, userName, password, data.RandomInteger, data.RandomInteger, userName, password)
}

func (ImageResource) vmss_customImage_fromVHD_setup(data acceptance.TestData, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  domain_name_label   = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
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
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
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
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, hostName, data.RandomInteger, data.RandomInteger, userName, password)
}

func (r ImageResource) vmss_customImage_fromVHD_provision(data acceptance.TestData, userName string, password string, hostName string) string {
	setup := r.vmss_customImage_fromVHD_setup(data, userName, password, hostName)
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
    computer_name_prefix = "testvm%d"
    admin_username       = "%s"
    admin_password       = "%s"
  }

  network_profile {
    name    = "TestNetworkProfile%d"
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
`, setup, data.RandomInteger, userName, password, data.RandomInteger)
}
