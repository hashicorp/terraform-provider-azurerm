package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRecoveryServicesVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesVault_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesVault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesVault_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesVault_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesVault_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMRecoveryServicesVault_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMRecoveryServicesVault_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRecoveryServicesVault_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryServicesVault_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMRecoveryServicesVault_requiresImport),
		},
	})
}

func TestRecoveryServicesVault_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryServicesVault_basicWithIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMRecoveryServicesVaultDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.VaultsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_recovery_services_vault" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Recovery Services Vault still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMRecoveryServicesVaultExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Recovery Services Vault: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Recovery Services Vault %q (resource group: %q) was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on recoveryServicesVaultsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMRecoveryServicesVault_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testRecoveryServicesVault_basicWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  identity {
    type = "SystemAssigned"
  }

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRecoveryServicesVault_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRecoveryServicesVault_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesVault_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_vault" "import" {
  name                = azurerm_recovery_services_vault.test.name
  location            = azurerm_recovery_services_vault.test.location
  resource_group_name = azurerm_recovery_services_vault.test.resource_group_name
  sku                 = azurerm_recovery_services_vault.test.sku
}
`, template)
}
