package policy_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
)

func TestAccAzureRMguestConfigurationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMguestConfigurationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMguestConfigurationAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMguestConfigurationAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMguestConfigurationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMguestConfigurationAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMguestConfigurationAssignment_requiresImport),
		},
	})
}

func TestAccAzureRMguestConfigurationAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMguestConfigurationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMguestConfigurationAssignment_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMguestConfigurationAssignment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMguestConfigurationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMguestConfigurationAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMguestConfigurationAssignment_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMguestConfigurationAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMguestConfigurationAssignment_updateGuestConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMguestConfigurationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMguestConfigurationAssignment_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMguestConfigurationAssignment_updateGuestConfiguration(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMguestConfigurationAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).GuestConfiguration.AssignmentClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("guestConfiguration Assignment not found: %s", resourceName)
		}
		id, err := parse.GuestConfigurationAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, id.VMName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: GuestConfiguration Assignment %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on GuestConfiguration.AssignmentClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMguestConfigurationAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).GuestConfiguration.AssignmentClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_guest_configuration_assignment" {
			continue
		}
		id, err := parse.GuestConfigurationAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, id.VMName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on GuestConfiguration.AssignmentClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMguestConfigurationAssignment_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-guestconfiguration-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
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
  name                     = "acctest-sa-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "accttest-sc-%d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_container" "test" {
  name                  = "accttest-sc-%d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMguestConfigurationAssignment_basic(data acceptance.TestData) string {
	template := testAccAzureRMguestConfigurationAssignment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_guest_configuration_assignment" "test" {
  name = "acctest-gca-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  vm_name = azurerm_storage_container.test.name
}
`, template, data.RandomInteger)
}

func testAccAzureRMguestConfigurationAssignment_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMguestConfigurationAssignment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_guest_configuration_assignment" "import" {
  name = azurerm_guest_configuration_assignment.test.name
  resource_group_name = azurerm_guest_configuration_assignment.test.resource_group_name
  location = azurerm_guest_configuration_assignment.test.location
  vm_name = azurerm_guest_configuration_assignment.test.vm_name
}
`, config)
}

func testAccAzureRMguestConfigurationAssignment_complete(data acceptance.TestData) string {
	template := testAccAzureRMguestConfigurationAssignment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_guest_configuration_assignment" "test" {
  name = "acctest-gca-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  vm_name = azurerm_storage_container.test.name
  context = "Azure policy"
  guest_configuration {
    name = "WhitelistedApplication"
    configuration_parameter {
      name = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }

    configuration_setting {
      action_after_reboot = "ContinueConfiguration"
      allow_module_overwrite = ""
      configuration_mode = "MonitorOnly"
      configuration_mode_frequency_mins = 15
      reboot_if_needed = "False"
      refresh_frequency_mins = 30
    }
    kind = ""
    version = "1.*"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMguestConfigurationAssignment_updateGuestConfiguration(data acceptance.TestData) string {
	template := testAccAzureRMguestConfigurationAssignment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_guest_configuration_assignment" "test" {
  name = "acctest-gca-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  vm_name = azurerm_storage_container.test.name
  context = "Azure policy"
  guest_configuration {
    name = "WhitelistedApplication"
    configuration_parameter {
      name = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }

    configuration_setting {
      action_after_reboot = "ContinueConfiguration"
      allow_module_overwrite = ""
      configuration_mode = "MonitorOnly"
      configuration_mode_frequency_mins = 15
      reboot_if_needed = "False"
      refresh_frequency_mins = 30
    }
    kind = ""
    version = "1.*"
  }
}
`, template, data.RandomInteger)
}
