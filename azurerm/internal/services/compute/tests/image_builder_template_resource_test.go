package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMImageBuilderTemplate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMImageBuilderTemplate_requiresImport),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "customizer.0.shell_sha256_checksum", "2c6ff6902a4a52deee69e8db26d0036a53388651008aaf31795bb20dabd21fd8"),
					resource.TestCheckResourceAttr(data.ResourceName, "customizer.1.shell_sha256_checksum", "ade4c5214c3c675e92c66e2d067a870c5b81b9844b3de3cc72c49ff36425fc93"),
					resource.TestCheckResourceAttr(data.ResourceName, "customizer.2.file_sha256_checksum", "d9715d72889fb1a0463d06ce9e89d1d2bd33b2c5e5362a736db6f5a25e601a58"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_tags_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMImageBuilderTemplate_tags_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMImageBuilderTemplate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_identity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMImageBuilderTemplate_identity_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMImageBuilderTemplate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_vnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_vnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_windowsPlatformSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_windowsPlatformSource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "customizer.0.powershell_sha256_checksum", "0607c084bdde8ef843cd8b7668e54a37ed07446bb642fe791ba79307a0828ea5"),
					resource.TestCheckResourceAttr(data.ResourceName, "customizer.2.file_sha256_checksum", "d9715d72889fb1a0463d06ce9e89d1d2bd33b2c5e5362a736db6f5a25e601a58"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_managedImageSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				// create the original VM
				Config: testLinuxVirtualMachine_imageFromExistingMachinePrep(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists("azurerm_linux_virtual_machine.source"),
					generalizeLinuxVirtualMachine("azurerm_linux_virtual_machine.source"),
				),
			},
			{
				// then create an image builder template by consuming the image id as source after generalizing the original VM
				Config: testLinuxVirtualMachine_imageBuilderTemplateFromImage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_sharedImageGallerySource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				// create the original VM
				Config: testLinuxVirtualMachine_imageFromExistingMachinePrep(data),
				Check: resource.ComposeTestCheckFunc(
					checkLinuxVirtualMachineExists("azurerm_linux_virtual_machine.source"),
					generalizeLinuxVirtualMachine("azurerm_linux_virtual_machine.source"),
				),
			},
			{
				// then create an image builder template by consuming the image version in SIG as source after generalizing the original VM
				Config: testLinuxVirtualMachine_imageBuilderTemplateFromSharedImageGallery(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_purchasePlanSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_purchasePlanSource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_vhdDistribution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_vhdDistribution(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_sharedImageDistribution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_sharedImageDistribution(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMImageBuilderTemplate_multipleDistribution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_image_builder_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMImageBuilderTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMImageBuilderTemplate_multipleDistribution(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageBuilderTemplateExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMImageBuilderTemplateDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMImageBuilderTemplateClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_image_builder_template" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Image Builder Template still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMImageBuilderTemplate_basic(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMImageBuilderTemplate_basic(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

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

func testAccAzureRMImageBuilderTemplate_tags_update(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
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

func testAccAzureRMImageBuilderTemplate_identity_update(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_vnet(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-vnet-%[1]d"
  location = "%[2]s"
  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_complete(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_purchasePlanSource(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_windowsPlatformSource(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_vhdDistribution(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionVHDTemplate := testAccAzureRMImageBuilderTemplate_distributionVHDTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_sharedImageDistribution(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionSharedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionSharedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testAccAzureRMImageBuilderTemplate_multipleDistribution(data acceptance.TestData) string {
	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)
	distributionVHDTemplate := testAccAzureRMImageBuilderTemplate_distributionVHDTemplate(data)
	distributionSharedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionSharedImageTemplate(data)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

  tags = {
    "StorageType" = "Standard_LRS"
    "type"        = "test"
  }
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

func testCheckAzureRMImageBuilderTemplateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMImageBuilderTemplateClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Image Builder Template Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("No resource group found in state for Image Builder Template: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Get Image Builder Template: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("image builder template %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testLinuxVirtualMachine_imageBuilderTemplateFromImage(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_imageFromExistingMachinePrep(data)

	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

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

func testLinuxVirtualMachine_imageBuilderTemplateFromSharedImageGallery(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_imageFromExistingMachinePrep(data)

	roleTemplate := testAccAzureRMImageBuilderTemplate_roleTemplate(data)
	distributionManagedImageTemplate := testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data)

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

func testAccAzureRMImageBuilderTemplate_roleTemplate(data acceptance.TestData) string {
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

func testAccAzureRMImageBuilderTemplate_distributionMamagedImageTemplate(data acceptance.TestData) string {
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

func testAccAzureRMImageBuilderTemplate_distributionVHDTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
distribution_vhd {
  run_output_name = "acctest-vhd-RunOutputName-%[1]d"
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger)
}

func testAccAzureRMImageBuilderTemplate_distributionSharedImageTemplate(data acceptance.TestData) string {
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
