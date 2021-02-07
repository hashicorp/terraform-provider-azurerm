package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ImageBuilderTemplateResource struct{}

func TestAccAzureRMImageBuilderTemplate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMImageBuilderTemplate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("customizer.0.shell_sha256_checksum").HasValue("2c6ff6902a4a52deee69e8db26d0036a53388651008aaf31795bb20dabd21fd8"),
				check.That(data.ResourceName).Key("customizer.1.shell_sha256_checksum").HasValue("ade4c5214c3c675e92c66e2d067a870c5b81b9844b3de3cc72c49ff36425fc93"),
				check.That(data.ResourceName).Key("customizer.2.shell_sha256_checksum").HasValue("d9715d72889fb1a0463d06ce9e89d1d2bd33b2c5e5362a736db6f5a25e601a58"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_tags_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags_update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_identity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identity_update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_vnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_windowsPlatformSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.windowsPlatformSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("customizer.0.powershell_sha256_checksum").HasValue("0607c084bdde8ef843cd8b7668e54a37ed07446bb642fe791ba79307a0828ea5"),
				check.That(data.ResourceName).Key("customizer.2.file_sha256_checksum").HasValue("d9715d72889fb1a0463d06ce9e89d1d2bd33b2c5e5362a736db6f5a25e601a58"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_managedImageSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	linuxVMResourceName := "azurerm_linux_virtual_machine.test"
	r := ImageBuilderTemplateResource{}
	rLinuxVMResource := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: rLinuxVMResource.imageFromExistingMachinePrep(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(linuxVMResourceName).ExistsInAzure(rLinuxVMResource),
				generalizeLinuxVirtualMachine("azurerm_linux_virtual_machine.source"),
			),
		},
		{
			Config: r.imageBuilderTemplateFromImage(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// This test case could pass in testing functionality but failed in the end due to conflict happened in resource deletion:
// shared image could not be deleted due to its nested shared image version still exist, while the latter declared deletion succeeded already.
// Issue filed to Azure service team: https://github.com/Azure/azure-rest-api-specs/issues/11559.
func TestAccAzureRMImageBuilderTemplate_sharedImageGallerySource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	linuxVMResourceName := "azurerm_linux_virtual_machine.test"
	r := ImageBuilderTemplateResource{}
	rLinuxVMResource := LinuxVirtualMachineResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: rLinuxVMResource.imageFromExistingMachinePrep(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(linuxVMResourceName).ExistsInAzure(rLinuxVMResource),
				generalizeLinuxVirtualMachine("azurerm_linux_virtual_machine.source"),
			),
		},
		{
			Config: r.imageBuilderTemplateFromSharedImageGallery(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_purchasePlanSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purchasePlanSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_vhdDistribution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vhdDistribution(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_sharedImageDistribution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.sharedImageDistribution(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMImageBuilderTemplate_multipleDistribution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")
	r := ImageBuilderTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleDistribution(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ImageBuilderTemplateResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := client.Compute.VMImageBuilderTemplateClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Image Builder Template: %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r ImageBuilderTemplateResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	if _, err := client.Compute.VMImageBuilderTemplateClient.Delete(ctx, resourceGroup, name); err != nil {
		return nil, fmt.Errorf("deleting Image Builder Template %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r ImageBuilderTemplateResource) basic(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

%[4]s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
%s

resource "azurerm_image_builder_template" "import" {
  name                = azurerm_image_builder_template.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

%s
}
`, template, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) tags_update(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

%s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV         = "Test"
    cost-center = "Ops"
  }

%s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, data.RandomInteger, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) identity_update(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

%s

resource "azurerm_user_assigned_identity" "test1" {
  name                = "acctestUAI-%d-1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_definition" "test1" {
  name  = "acctestRD-%d-1"
  scope = azurerm_resource_group.test.id

  permissions {
    actions = [
      "Microsoft.Compute/galleries/read",
      "Microsoft.Compute/galleries/images/read",
      "Microsoft.Compute/galleries/images/versions/read",
      "Microsoft.Compute/galleries/images/versions/write",
      "Microsoft.Compute/images/write",
      "Microsoft.Compute/images/read",
      "Microsoft.Compute/images/delete"
    ]
    not_actions = []
  }

  assignable_scopes = [
    azurerm_resource_group.test.id,
  ]
}

resource "azurerm_role_assignment" "test1" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = azurerm_role_definition.test1.role_definition_resource_id
  principal_id       = azurerm_user_assigned_identity.test1.principal_id
}

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test1.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

%s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, data.RandomInteger, data.RandomInteger, data.RandomInteger, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) vnet(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-vnet-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestVnet-%[1]d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                                          = "acctestSubnet-%[1]d"
  resource_group_name                           = azurerm_resource_group.test1.name
  virtual_network_name                          = azurerm_virtual_network.test1.name
  address_prefixes                              = ["10.0.1.0/24"]
  enforce_private_link_service_network_policies = true
}

resource "azurerm_network_security_group" "test1" {
  name                = "acctestNSG-%[1]d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "acctestNSR-%[1]d"
  priority                    = 400
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  destination_port_ranges     = ["60000", "60001"]
  source_port_range           = "*"
  source_address_prefix       = "AzureLoadBalancer"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test1.name
  network_security_group_name = azurerm_network_security_group.test1.name
}

resource "azurerm_subnet_network_security_group_association" "test1" {
  subnet_id                 = azurerm_subnet.test1.id
  network_security_group_id = azurerm_network_security_group.test1.id
}

%[3]s

resource "azurerm_role_definition" "test1" {
  name  = "acctestRD-vnet-%[1]d"
  scope = azurerm_resource_group.test1.id

  permissions {
    actions = [
      "Microsoft.Network/virtualNetworks/read",
      "Microsoft.Network/virtualNetworks/subnets/join/action"
    ]
    not_actions = []
  }

  assignable_scopes = [
    azurerm_resource_group.test1.id,
  ]
}

resource "azurerm_role_assignment" "test1" {
  scope              = azurerm_resource_group.test1.id
  role_definition_id = azurerm_role_definition.test1.role_definition_resource_id
  principal_id       = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test1.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

%[4]s

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_role_assignment.test1
  ]

}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) complete(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  build_timeout_minutes = "60"
  size                  = "Standard_D2_v2"
  disk_size_gb          = "10"

  tags = {
    ENV = "Test"
  }

  customizer {
    type             = "Shell"
    name             = "RunScriptFromSource"
    shell_script_uri = "https://raw.githubusercontent.com/danielsollondon/azvmimagebuilder/4afbd7858fb8918edc459a7f09ace43b570d027e/quickquickstarts/customizeScript.sh"
  }

  customizer {
    type             = "Shell"
    name             = "CheckSumCompareShellScript"
    shell_script_uri = "https://raw.githubusercontent.com/danielsollondon/azvmimagebuilder/4afbd7858fb8918edc459a7f09ace43b570d027e/quickquickstarts/customizeScript2.sh"
  }

  customizer {
    type                  = "File"
    name                  = "downloadBuildArtifacts"
    file_source_uri       = "https://raw.githubusercontent.com/danielsollondon/azvmimagebuilder/4afbd7858fb8918edc459a7f09ace43b570d027e/quickquickstarts/exampleArtifacts/buildArtifacts/index.html"
    file_destination_path = "/tmp/index.html"
  }

  customizer {
    type           = "Shell"
    name           = "setupBuildPath"
    shell_commands = ["sudo mkdir -p /buildArtifacts", "sudo cp /tmp/index.html /buildArtifacts/index.html"]
  }

  customizer {
    type           = "Shell"
    name           = "InstallUpgrades"
    shell_commands = ["sudo apt install unattended-upgrades"]
  }

%[4]s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) purchasePlanSource(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "RedHat"
    offer     = "rhel-byos"
    sku       = "rhel-lvm75"
    version   = "latest"
    plan {
      name      = "rhel-lvm75"
      product   = "rhel-byos"
      publisher = "redhat"
    }
  }

  tags = {
    ENV = "Test"
  }

%[4]s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) windowsPlatformSource(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

  customizer {
    type                       = "PowerShell"
    name                       = "CreateBuildPath"
    powershell_run_elevated    = false
    powershell_script_uri      = "https://raw.githubusercontent.com/danielsollondon/azvmimagebuilder/4afbd7858fb8918edc459a7f09ace43b570d027e/testPsScript.ps1"
    powershell_sha256_checksum = "0607c084bdde8ef843cd8b7668e54a37ed07446bb642fe791ba79307a0828ea5"
  }

  customizer {
    type                          = "WindowsRestart"
    name                          = "winRestart"
    windows_restart_command       = "shutdown /r /f /t 0"
    windows_restart_timeout       = "5m"
    windows_restart_check_command = "echo Azure-Image-Builder-Restarted-the-VM > c:\\buildArtifacts\\azureImageBuilderRestart.txt"
  }

  customizer {
    type                  = "File"
    name                  = "downloadBuildArtifacts"
    file_source_uri       = "https://raw.githubusercontent.com/danielsollondon/azvmimagebuilder/4afbd7858fb8918edc459a7f09ace43b570d027e/quickquickstarts/exampleArtifacts/buildArtifacts/index.html"
    file_destination_path = "c:\\buildArtifacts\\index.html"
    file_sha256_checksum  = "d9715d72889fb1a0463d06ce9e89d1d2bd33b2c5e5362a736db6f5a25e601a58"
  }

  customizer {
    type                    = "PowerShell"
    name                    = "settingUpMgmtAgtPath"
    powershell_run_elevated = false
    powershell_commands     = ["mkdir c:\\buildActions", "echo Azure-Image-Builder-Was-Here > c:\\buildActions\\buildActionsOutput.txt"]
  }

  customizer {
    type                           = "WindowsUpdate"
    name                           = "winUpdate"
    windows_update_filters         = ["exclude:$_.Title -like '*Preview*'", "include:$true"]
    windows_update_search_criteria = "IsInstalled=0"
    windows_update_limit           = 20
  }

%[4]s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) vhdDistribution(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionVHDTemplate := r.distributionVHDTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

%[4]s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionVHDTemplate)
}

func (r ImageBuilderTemplateResource) sharedImageDistribution(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionSharedImageTemplate := r.distributionSharedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestSIG%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-IMG-%[1]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  identifier {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
  }
}

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

%[4]s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionSharedImageTemplate)
}

func (r ImageBuilderTemplateResource) multipleDistribution(data acceptance.TestData) string {
	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)
	distributionVHDTemplate := r.distributionVHDTemplate(data)
	distributionSharedImageTemplate := r.distributionSharedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestSIG%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-IMG-%[1]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  identifier {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
  }
}

resource "azurerm_shared_image_gallery" "test1" {
  name                = "acctestSIG%[1]dtest1"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test1" {
  name                = "acctest-IMG-%[1]d-test1"
  gallery_name        = azurerm_shared_image_gallery.test1.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  identifier {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
  }
}

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  tags = {
    ENV = "Test"
  }

%[4]s

  distribution_managed_image {
    name                = "acctestDistManagedImg-%[1]d-1"
    resource_group_name = azurerm_resource_group.test.name
    location            = azurerm_resource_group.test.location
    run_output_name     = "acctest-managedImage-RunOutputName-%[1]d-1"
    tags = {
      ENV = "Test"
    }
  }

%[5]s

%[6]s

  distribution_shared_image {
    id = azurerm_shared_image.test1.id
    replica_regions {
      name = azurerm_resource_group.test.location
    }
    run_output_name = "acctest-sharedImage-RunOutputName-%[1]d-1"
    tags = {
      ENV = "Test"
    }
  }

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, roleTemplate, distributionManagedImageTemplate, distributionVHDTemplate, distributionSharedImageTemplate)
}

func (r ImageBuilderTemplateResource) imageBuilderTemplateFromImage(data acceptance.TestData) string {
	rLinuxVMResource := LinuxVirtualMachineResource{}

	template := rLinuxVMResource.imageFromExistingMachinePrep(data)

	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

%s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_managed_image_id = azurerm_image.test.id

  tags = {
    ENV = "Test"
  }

%s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, template, roleTemplate, data.RandomInteger, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) imageBuilderTemplateFromSharedImageGallery(data acceptance.TestData) string {
	rLinuxVMResource := LinuxVirtualMachineResource{}
	template := rLinuxVMResource.imageFromExistingMachinePrep(data)

	roleTemplate := r.roleTemplate(data)
	distributionManagedImageTemplate := r.distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
%s

resource "azurerm_image" "test" {
  name                      = "capture"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

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
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_resource_group.test.location
    regional_replica_count = "1"
    storage_account_type   = "Standard_LRS"
  }
}

%s

resource "azurerm_image_builder_template" "test" {
  name                = "acctestIBT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  source_shared_image_version_id = azurerm_shared_image_version.test.id

  tags = {
    ENV = "Test"
  }

%s

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, template, roleTemplate, data.RandomInteger, distributionManagedImageTemplate)
}

func (r ImageBuilderTemplateResource) roleTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_definition" "test" {
  name  = "acctestRD-%d"
  scope = azurerm_resource_group.test.id

  permissions {
    actions = [
      "Microsoft.Compute/galleries/read",
      "Microsoft.Compute/galleries/images/read",
      "Microsoft.Compute/galleries/images/versions/read",
      "Microsoft.Compute/galleries/images/versions/write",
      "Microsoft.Compute/images/write",
      "Microsoft.Compute/images/read",
      "Microsoft.Compute/images/delete"
    ]
    not_actions = []
  }

  assignable_scopes = [
    azurerm_resource_group.test.id,
  ]
}

resource "azurerm_role_assignment" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = azurerm_role_definition.test.role_definition_resource_id
  principal_id       = azurerm_user_assigned_identity.test.principal_id
}
`, data.RandomInteger, data.RandomInteger)
}

func (r ImageBuilderTemplateResource) distributionMamagedImageTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
distribution_managed_image {
  name                = "acctestDistManagedImg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  run_output_name     = "acctest-managedImage-RunOutputName-%[1]d"
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger)
}

func (r ImageBuilderTemplateResource) distributionVHDTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
distribution_vhd {
  run_output_name = "acctest-vhd-RunOutputName-%[1]d"
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger)
}

func (r ImageBuilderTemplateResource) distributionSharedImageTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
distribution_shared_image {
  id = azurerm_shared_image.test.id
  replica_regions {
    name = azurerm_resource_group.test.location
  }
  replica_regions {
    name = "westus2"
  }
  run_output_name      = "acctest-sharedImage-RunOutputName-%[1]d"
  storage_account_type = "Standard_ZRS"
  exclude_from_latest  = true
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger)
}
