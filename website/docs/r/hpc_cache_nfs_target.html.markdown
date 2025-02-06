---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache_nfs_target"
description: |-
  Manages a NFS Target within a HPC Cache.
---

# azurerm_hpc_cache_nfs_target

Manages a NFS Target within a HPC Cache.

~> **NOTE:**: By request of the service team the provider no longer automatically registering the `Microsoft.StorageCache` Resource Provider for this resource. To register it you can run `az provider register --namespace 'Microsoft.StorageCache'`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevn"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example_hpc" {
  name                 = "examplesubnethpc"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_hpc_cache" "example" {
  name                = "examplehpccache"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.example_hpc.id
  sku_name            = "Standard_2G"
}

resource "azurerm_subnet" "example_vm" {
  name                 = "examplesubnetvm"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "examplenic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example_vm.id
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

resource "azurerm_linux_virtual_machine" "example" {
  name                = "examplevm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/id_rsa.pub")
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  custom_data = base64encode(local.custom_data)
}

resource "azurerm_hpc_cache_nfs_target" "example" {
  name                = "examplehpcnfstarget"
  resource_group_name = azurerm_resource_group.example.name
  cache_name          = azurerm_hpc_cache.example.name
  target_host_name    = azurerm_linux_virtual_machine.example.private_ip_address
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HPC Cache NFS Target. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the HPC Cache NFS Target. Changing this forces a new resource to be created.

* `cache_name` - (Required) The name HPC Cache, which the HPC Cache NFS Target will be added to. Changing this forces a new resource to be created.

* `target_host_name` - (Required) The IP address or fully qualified domain name (FQDN) of the HPC Cache NFS target. Changing this forces a new resource to be created.

* `usage_model` - (Required) The type of usage of the HPC Cache NFS Target. Possible values are: `READ_HEAVY_INFREQ`, `READ_HEAVY_CHECK_180`, `READ_ONLY`, `READ_WRITE`, `WRITE_WORKLOAD_15`, `WRITE_AROUND`, `WRITE_WORKLOAD_CHECK_30`, `WRITE_WORKLOAD_CHECK_60` and `WRITE_WORKLOAD_CLOUDWS`.

* `namespace_junction` - (Required) Can be specified multiple times to define multiple `namespace_junction`. Each `namespace_junction` block supports fields documented below.

* `verification_timer_in_seconds` - (Optional) The amount of time the cache waits before it checks the back-end storage for file updates. Possible values are between `1` and `31536000`.

* `write_back_timer_in_seconds` - (Optional) The amount of time the cache waits after the last file change before it copies the changed file to back-end storage. Possible values are between `1` and `31536000`.

---

A `namespace_junction` block supports the following:

* `namespace_path` - (Required) The client-facing file path of this NFS target within the HPC Cache NFS Target.

* `nfs_export` - (Required) The NFS export of this NFS target within the HPC Cache NFS Target.

* `target_path` - (Optional) The relative subdirectory path from the `nfs_export` to map to the `namespace_path`. Defaults to `""`, in which case the whole `nfs_export` is exported.

* `access_policy_name` - (Optional) The name of the access policy applied to this target. Defaults to `default`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the HPC Cache NFS Target.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HPC Cache NFS Target.
* `read` - (Defaults to 5 minutes) Used when retrieving the HPC Cache NFS Target.
* `update` - (Defaults to 30 minutes) Used when updating the HPC Cache NFS Target.
* `delete` - (Defaults to 30 minutes) Used when deleting the HPC Cache NFS Target.

## Import

NFS Target within a HPC Cache can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hpc_cache_nfs_target.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageCache/caches/cache1/storageTargets/target1
```
