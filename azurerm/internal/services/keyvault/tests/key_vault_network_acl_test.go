package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultNetworkAcl_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_network_acl", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultNetworkAcl_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultNetworkAclExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultNetworkAcl_updateRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_network_acl", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultNetworkAcl_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultNetworkAclExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.bypass", "None"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKeyVaultNetworkAcl_updateTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultNetworkAclExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.bypass", "AzureServices"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultNetworkAcl_addIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_network_acl", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultNetworkAcl_addIp(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultNetworkAclExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.ip_rules.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVaultNetworkAcl_addVnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_network_acl", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultNetworkAcl_addVnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultNetworkAclExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_acls.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKeyVault_common(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
	features {}
}  

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%[1]s"
	location = "%[2]s"
}  

resource "azurerm_key_vault" "test" {
	name                = "acctestkv%[1]s"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	sku_name            = "standard"
	tenant_id           = data.azurerm_client_config.current.tenant_id
}  

resource "azurerm_virtual_network" "test" {
	name                = "acctestRG-%[1]s-network"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	address_space       = ["10.0.0.0/16"]
}  

resource "azurerm_subnet" "test" {
	name                 = "acctestRG-%[1]s-subnet"
	virtual_network_name = azurerm_virtual_network.test.name
	resource_group_name  = azurerm_resource_group.test.name
	address_prefixes     = ["10.0.1.0/24"]
	service_endpoints    = ["Microsoft.KeyVault"]
}
`, fmt.Sprintf("%d", data.RandomInteger)[8:], data.Locations.Primary)
}

func testAccAzureRMKeyVaultNetworkAcl_basicTemplate(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_network_acl" "test" {
	key_vault_name      = azurerm_key_vault.test.name
	resource_group_name = azurerm_resource_group.test.name
	network_acls {
	  default_action = "Allow"
	  bypass         = "None"
	}
}
`, template)
}

func testAccAzureRMKeyVaultNetworkAcl_updateTemplate(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_network_acl" "test" {
	key_vault_name      = azurerm_key_vault.test.name
	resource_group_name = azurerm_resource_group.test.name
	network_acls {
	  default_action = "Deny"
	  bypass         = "AzureServices"
	}
}
`, template)
}

func testAccAzureRMKeyVaultNetworkAcl_addIp(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_network_acl" "test" {
	key_vault_name      = azurerm_key_vault.test.name
	resource_group_name = azurerm_resource_group.test.name
	network_acls {
	  default_action = "Deny"
	  bypass         = "AzureServices"
	  ip_rules = ["43.0.0.0/24"]
	}
}
`, template)
}

func testAccAzureRMKeyVaultNetworkAcl_addVnet(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_network_acl" "test" {
	key_vault_name      = azurerm_key_vault.test.name
	resource_group_name = azurerm_resource_group.test.name
	network_acls {
	  default_action = "Deny"
	  bypass         = "AzureServices"
	  virtual_network_subnet_ids = [azurerm_subnet.test.id]
	}
}
`, template)
}

func testCheckAzureRMKeyVaultNetworkAclExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}
		resourceGroup := id.ResourceGroup
		keyVaultName := id.Path["vaults"]

		keyVault, err := client.Get(ctx, resourceGroup, keyVaultName)
		if err != nil {
			if utils.ResponseWasNotFound(keyVault.Response) {
				return fmt.Errorf("Key Vault %q (Resource Group %q) was not found", keyVaultName, resourceGroup)
			}

			return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", keyVaultName, resourceGroup, err)
		}

		if rules := keyVault.Properties.NetworkAcls; rules == nil {
			return fmt.Errorf("Network Acl for Azure Key Vault %q (Resource Group %q): %+v does not exist", keyVaultName, resourceGroup, err)
		}

		return nil
	}
}
