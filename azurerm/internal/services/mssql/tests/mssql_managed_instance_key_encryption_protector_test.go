package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMsSqlManagedInstanceEncryption_keyEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceEncryptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstanceEncryption_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceEncryptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstanceEncryption_ServiceManagedEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceEncryptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstanceEncryption_ServiceManaged(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceEncryptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstanceEncryption_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceEncryptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstanceEncryption_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceEncryptionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMssqlManagedInstanceEncryption_requiresImport),
		},
	})
}

func TestAccAzureRMMsSqlManagedInstanceEncryption_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceEncryptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlManagedInstanceEncryption_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceEncryptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstanceEncryption_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceEncryptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlManagedInstanceEncryption_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceEncryptionExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAzureRMMsSqlManagedInstanceEncryptionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ManagedInstanceEncryptionProtectorsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		managedInstanceId := rs.Primary.Attributes["managed_instance_id"]
		id, _ := azure.ParseAzureResourceID(managedInstanceId)
		managedInstanceName := id.Path["managedInstances"]
		resourceGroup := id.ResourceGroup

		resp, err := client.Get(ctx, resourceGroup, managedInstanceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on ManagedInstanceEncryption: %+v", err)
		}

		if string(resp.ManagedInstanceEncryptionProtectorProperties.ServerKeyType) == "ServiceManaged" {
			return fmt.Errorf("Bad: Managed Instance Encryption (Managed Sql Instance %q, resource group: %q) does not exist", managedInstanceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlManagedInstanceEncryptionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ManagedInstanceEncryptionProtectorsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_managed_instance_key" {
			continue
		}

		managedInstanceId := rs.Primary.Attributes["managed_instance_id"]
		id, _ := azure.ParseAzureResourceID(managedInstanceId)
		managedInstanceName := id.Path["managedInstances"]
		resourceGroup := id.ResourceGroup

		if resp, err := client.Get(ctx, resourceGroup, managedInstanceName); err != nil {
			if string(resp.ManagedInstanceEncryptionProtectorProperties.ServerKeyType) != "ServiceManaged" {
				return fmt.Errorf("Get on managed instance key Client: %+v", err)
			}
		}
		return nil
	}

	return nil
}

func testAccAzureRMMsSqlManagedInstanceEncryption_basic(data acceptance.TestData) string {
	keyvaultName := "acctst-kv-" + data.RandomString
	template := testAccAzureRMMsSqlManagedInstanceEncryption_prepareDependencies(data, keyvaultName)
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_instance_key" "test" {
	key_name                          = "${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.version}"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	uri 					 = azurerm_key_vault_key.test.id
  }

  resource "azurerm_mssql_managed_instance_encryption_protector" "test" {
	server_key_name                          = "${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.version}"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	server_key_type = "AzureKeyVault"
  }
`, template)
}

func testAccAzureRMMsSqlManagedInstanceEncryption_ServiceManaged(data acceptance.TestData) string {
	keyvaultName := "acctst-kv-" + data.RandomString
	template := testAccAzureRMMsSqlManagedInstanceEncryption_prepareDependencies(data, keyvaultName)
	return fmt.Sprintf(`%s

  resource "azurerm_mssql_managed_instance_encryption_protector" "test" {
	server_key_name                          = "ServiceManaged"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	server_key_type = "ServiceManaged"
  }
`, template)
}

func testAccAzureRMMsSqlManagedInstanceEncryption_update(data acceptance.TestData) string {
	keyvaultName := "acctst-kv-" + data.RandomString
	template := testAccAzureRMMsSqlManagedInstanceEncryption_prepareDependencies(data, keyvaultName)
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_instance_key" "test" {
	key_name                          = "${azurerm_key_vault_key.test1.name}_${azurerm_key_vault_key.test1.name}_${azurerm_key_vault_key.test1.version}"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	uri 					 = azurerm_key_vault_key.test1.id
  }

  resource "azurerm_mssql_managed_instance_encryption_protector" "test" {
	server_key_name                          = "${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.version}"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	server_key_type = "AzureKeyVault"
  }

`, template)
}

func testAccAzureRMMssqlManagedInstanceEncryption_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlManagedInstanceEncryption_basic(data)
	return fmt.Sprintf(`%s

resource "azurerm_mssql_managed_instance_encryption_protector" "import" {
	server_key_name                    = "${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.name}_${azurerm_key_vault_key.test.version}"
	managed_instance_name          = azurerm_mssql_managed_instance.test.name
	resource_group_name				= azurerm_mssql_managed_instance.test.resource_group_name
	server_key_type = "AzureKeyVault"
  }
`, template)
}

func testAccAzureRMMsSqlManagedInstanceEncryption_prepareDependencies(data acceptance.TestData, keyvaultName string) string {
	return fmt.Sprintf(`provider "azurerm" {
	  features {}
	}

	data "azurerm_client_config" "current" {}
	
	resource "azurerm_resource_group" "test" {
	  name     = "acctestRG-%[1]d"
	  location = "%[2]s"
	}
	
	resource "azurerm_network_security_group" "test" {
	  name                = "accTestNetworkSecurityGroup-%[1]d"
	  location            = "%[2]s"
	  resource_group_name = azurerm_resource_group.test.name
	}
	
	resource "azurerm_virtual_network" "test" {
	  name                = "acctest-%[1]d-network"
	  resource_group_name = azurerm_resource_group.test.name
	  location            = "%[2]s"
	  address_space       = ["10.0.0.0/16"]
	}
	
	resource "azurerm_subnet" "test" {
	  name                 = "internal"
	  virtual_network_name = azurerm_virtual_network.test.name
	  resource_group_name  = azurerm_resource_group.test.name
	  address_prefixes     = ["10.0.1.0/24"]
	  delegation {
		name = "miDelegation"
		service_delegation {
		  name = "Microsoft.Sql/managedInstances"
		}
	  }
	}
	
	resource "azurerm_subnet_network_security_group_association" "test" {
	  subnet_id                 = azurerm_subnet.test.id
	  network_security_group_id = azurerm_network_security_group.test.id
	}
	
	resource "azurerm_route_table" "test" {
	  name                = "test-routetable-%[1]d"
	  location            = azurerm_resource_group.test.location
	  resource_group_name = azurerm_resource_group.test.name
	  route {
		name                   = "test"
		address_prefix         = "10.100.0.0/14"
		next_hop_type          = "VirtualAppliance"
		next_hop_in_ip_address = "10.10.1.1"
	  }
	}
	
	resource "azurerm_subnet_route_table_association" "test" {
	  subnet_id      = azurerm_subnet.test.id
	  route_table_id = azurerm_route_table.test.id
	}
	
	resource "azurerm_mssql_managed_instance" "test" {
	  name                         = "acctest-mi-%[1]d"
	  resource_group_name          = azurerm_resource_group.test.name
	  location                     = azurerm_resource_group.test.location
	  administrator_login          = "AcceptanceTestUser"
	  administrator_login_password = "LengthyPassword@1234"
	  subnet_id                    = azurerm_subnet.test.id
	  identity {
		type = "SystemAssigned"
	  }
	  sku {
		capacity = 8
		family   = "Gen5"
		name     = "GP_Gen5"
		tier     = "GeneralPurpose"
	  }
	  depends_on = [
		azurerm_subnet_network_security_group_association.test,
		azurerm_subnet_route_table_association.test,
	  ]
	},

	resource "azurerm_key_vault" "test" {
		name                        = "%[3]s"
		location                    = azurerm_resource_group.test.location
		resource_group_name         = azurerm_resource_group.test.name
		enabled_for_disk_encryption = true
		tenant_id                   = data.azurerm_client_config.current.tenant_id
		soft_delete_enabled         = true
		soft_delete_retention_days  = 7
		purge_protection_enabled    = false
	  
		sku_name = "standard"
	  
		access_policy {
		  tenant_id = data.azurerm_client_config.current.tenant_id
		  object_id = data.azurerm_client_config.current.object_id
	  
		  key_permissions = [
			"get",
			"create",
			"list"
		  ]
	  
		  secret_permissions = [
			"get",
		  ]
	  
		  storage_permissions = [
			"get",
		  ]
		}
	  },

	  resource "azurerm_key_vault_key" "test" {
		name         = "acc-test1"
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
	  }

	  resource "azurerm_key_vault_key" "test1" {
		name         = "acc-test2"
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
	  }

	`, data.RandomInteger, data.Locations.Primary, keyvaultName)
}
