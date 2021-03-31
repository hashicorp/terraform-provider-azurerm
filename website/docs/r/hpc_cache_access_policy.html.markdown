---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache_access_policy"
description: |-
  Manages a HPC Cache Access Policy.
---

# azurerm_hpc_cache_access_policy

Manages a HPC Cache Access Policy.

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

resource "azurerm_subnet" "example" {
  name                 = "examplesubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_hpc_cache" "example" {
  name                = "examplehpccache"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.example.id
  sku_name            = "Standard_2G"
}

resource "azurerm_hpc_cache_access_policy" "example" {
  name         = "example"
  hpc_cache_id = azurerm_hpc_cache.example.id

  access_rule {
    scope  = "default"
    access = "rw"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this HPC Cache Access Policy. Changing this forces a new HPC Cache Access Policy to be created.

* `hpc_cache_id` - (Required) The ID of the HPC Cache that this HPC Cache Access Policy resides in. Changing this forces a new HPC Cache Access Policy to be created.

* `access_rule` - (Required) Up to three `access_rule` blocks as defined below.

---

An `access_rule` block supports the following:

* `scope` - (Required) The scope of this rule. The `scope` and (potentially) the `filter` determine which clients match the rule. Possible values are: `default`, `network`, `host`.

~> **NOTE**: Each `access_rule` should set a unique `scope`.

* `access` - (Required) The access level for this rule. Possible values are: `rw`, `ro`, `no`.

* `filter` - (Optional) The filter applied to the `scope` for this rule. The filter's format depends on its scope: `default` scope matches all clients and has no filter value; `network` scope takes a CIDR format; `host` takes an IP address or fully qualified domain name. If a client does not match any filter rule and there is no default rule, access is denied.

* `suid_enabled` - (Optional) Whether [SUID](https://docs.microsoft.com/en-us/azure/hpc-cache/access-policies#suid) is allowed? Defaults to `false`.

* `submount_access_enabled` - (Optional) Whether allow access to subdirectories under the root export? Defaults to `false`.

* `root_squash_enabled` - (Optional) Whether to enable [root squash](https://docs.microsoft.com/en-us/azure/hpc-cache/access-policies#root-squash)? Defaults to `false`.

* `anonymous_uid` - (Optional) The anonymous UID used when `root_squash_enabled` is `true`.
 
* `anonymous_gid` - (Optional) The anonymous GID used when `root_squash_enabled` is `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the HPC Cache Access Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HPC Cache Access Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the HPC Cache Access Policy.
* `update` - (Defaults to 30 minutes) Used when updating the HPC Cache Access Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the HPC Cache Access Policy.

## Import

HPC Cache Access Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hpc_cache_access_policy.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.StorageCache/caches/cache1/cacheAccessPolicies/policy1
```
