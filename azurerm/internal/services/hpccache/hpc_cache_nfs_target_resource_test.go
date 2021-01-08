package hpccache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HPCCacheNFSTargetResource struct {
}

func TestAccHPCCacheNFSTarget_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")
	r := HPCCacheNFSTargetResource{}

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

func TestAccHPCCacheNFSTarget_usageModel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")
	r := HPCCacheNFSTargetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.usageModel(data),
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

func TestAccHPCCacheNFSTarget_namespaceJunction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")
	r := HPCCacheNFSTargetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.namespaceJunction(data),
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

func TestAccHPCCacheNFSTarget_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache_nfs_target", "test")
	r := HPCCacheNFSTargetResource{}

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

func (HPCCacheNFSTargetResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StorageTargetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HPCCache.StorageTargetsClient.Get(ctx, id.ResourceGroup, id.CacheName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving HPC Cache NFS Target (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.BasicStorageTargetProperties != nil), nil
}

func (r HPCCacheNFSTargetResource) basic(data acceptance.TestData) string {
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
`, r.template(data), data.RandomString)
}

func (r HPCCacheNFSTargetResource) usageModel(data acceptance.TestData) string {
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
`, r.template(data), data.RandomString)
}

func (r HPCCacheNFSTargetResource) namespaceJunction(data acceptance.TestData) string {
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
`, r.template(data), data.RandomString)
}

func (r HPCCacheNFSTargetResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (HPCCacheNFSTargetResource) template(data acceptance.TestData) string {
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

`, HPCCacheResource{}.basic(data), data.RandomString)
}
