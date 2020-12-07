package hpccache_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMHPCCacheNFSTarget_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHPCCacheNFSTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHPCCacheNFSTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHPCCacheNFSTarget_usageModel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHPCCacheNFSTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHPCCacheNFSTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMHPCCacheNFSTarget_usageModel(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMHPCCacheNFSTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHPCCacheNFSTarget_namespaceJunction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHPCCacheNFSTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHPCCacheNFSTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMHPCCacheNFSTarget_namespaceJunction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMHPCCacheNFSTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHPCCacheNFSTarget_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHPCCacheNFSTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHPCCacheNFSTarget_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHPCCacheNFSTargetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMHPCCacheNFSTarget_requiresImport),
		},
	})
}

func testCheckAzureRMHPCCacheNFSTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).HPCCache.StorageTargetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("HPC Cache NFS Target not found: %s", resourceName)
		}

		id, err := parse.StorageTargetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.CacheName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: HPC Cache NFS Target %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on Storage.StorageTargetsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMHPCCacheNFSTargetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).HPCCache.StorageTargetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hpc_cache_nfs_target" {
			continue
		}

		id, err := parse.StorageTargetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.CacheName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Storage.StorageTargetsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMHPCCacheNFSTarget_basic(data acceptance.TestData) string {
	template := testAccAzureRMHPCCacheNFSTarget_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_nfs_target" "test" {
  name                = "acctest-HPCCTGT-%s"
  resource_group_name = azurerm_resource_group.test.name
  cache_name          = azurerm_hpc_cache.test.name
  target_host_name    = azurerm_linux_virtual_machine.test.private_ip_address
  usage_model         = "READ_HEAVY_INFREQ"
  namespace_junction {
    namespace_path = "/nfs/a1"
    nfs_export     = "/export/a"
    target_path    = "1"
  }
  namespace_junction {
    namespace_path = "/nfs/b"
    nfs_export     = "/export/b"
  }
}
`, template, data.RandomString)
}

func testAccAzureRMHPCCacheNFSTarget_usageModel(data acceptance.TestData) string {
	template := testAccAzureRMHPCCacheNFSTarget_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_nfs_target" "test" {
  name                = "acctest-HPCCTGT-%s"
  resource_group_name = azurerm_resource_group.test.name
  cache_name          = azurerm_hpc_cache.test.name
  target_host_name    = azurerm_linux_virtual_machine.test.private_ip_address
  usage_model         = "WRITE_WORKLOAD_15"
  namespace_junction {
    namespace_path = "/nfs/a1"
    nfs_export     = "/export/a"
    target_path    = "1"
  }
  namespace_junction {
    namespace_path = "/nfs/b"
    nfs_export     = "/export/b"
  }
}
`, template, data.RandomString)
}

func testAccAzureRMHPCCacheNFSTarget_namespaceJunction(data acceptance.TestData) string {
	template := testAccAzureRMHPCCacheNFSTarget_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_nfs_target" "test" {
  name                = "acctest-HPCCTGT-%s"
  resource_group_name = azurerm_resource_group.test.name
  cache_name          = azurerm_hpc_cache.test.name
  target_host_name    = azurerm_linux_virtual_machine.test.private_ip_address
  usage_model         = "WRITE_WORKLOAD_15"
  namespace_junction {
    namespace_path = "/nfs/a"
    nfs_export     = "/export/a"
    target_path    = ""
  }
}
`, template, data.RandomString)
}

func testAccAzureRMHPCCacheNFSTarget_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMHPCCacheNFSTarget_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache_nfs_target" "import" {
  name                = azurerm_hpc_cache_nfs_target.test.name
  resource_group_name = azurerm_hpc_cache_nfs_target.test.resource_group_name
  cache_name          = azurerm_hpc_cache_nfs_target.test.cache_name
  target_host_name    = azurerm_hpc_cache_nfs_target.test.target_host_name
  usage_model         = azurerm_hpc_cache_nfs_target.test.usage_model
  namespace_junction {
    namespace_path = "/nfs/a1"
    nfs_export     = "/export/a"
    target_path    = "1"
  }
  namespace_junction {
    namespace_path = "/nfs/b"
    nfs_export     = "/export/b"
  }
}
`, template)
}

func testAccAzureRMHPCCacheNFSTarget_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "testvm" {
  name                 = "acctest-sub-vm-%[2]s"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.3.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctest-nic-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.testvm.id
    private_ip_address_allocation = "Dynamic"
  }
}

locals {
  custom_data = <<CUSTOM_DATA
#!/bin/bash
sudo -i 
apt-get install -y nfs-kernel-server
mkdir -p /export/a/1
mkdir -p /export/a/2
mkdir -p /export/b
cat << EOF > /etc/exports
/export/a *(rw,fsid=0,insecure,no_subtree_check,async)
/export/b *(rw,fsid=0,insecure,no_subtree_check,async)
EOF
systemctl start nfs-server
exportfs -arv
CUSTOM_DATA
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctest-vm-%[2]s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  custom_data = base64encode(local.custom_data)
}

`, testAccAzureRMHPCCache_basic(data), data.RandomString)
}
