---
subcategory: "Redis Enterprise"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_enterprise_cluster"
description: |-
  Manages a Redis Enterprise Cluster.
---

# azurerm_redis_enterprise_cluster

Manages a Redis Enterprise Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-redisenterprise"
  location = "West Europe"
}

resource "azurerm_redis_enterprise_cluster" "example" {
  name                = "example-redisenterprise"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "EnterpriseFlash_F300-3"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Redis Enterprise Cluster. Changing this forces a new Redis Enterprise Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Redis Enterprise Cluster should exist. Changing this forces a new Redis Enterprise Cluster to be created.

* `location` - (Required) The Azure Region where the Redis Enterprise Cluster should exist. Changing this forces a new Redis Enterprise Cluster to be created.

* `sku_name` - (Required) The `sku_name` is comprised of two segments separated by a hyphen (e.g. `Enterprise_E10-2`). The first segment of the `sku_name` defines the `name` of the SKU, possible values are `Enterprise_E5`, `Enterprise_E10`, `Enterprise_E20`, `Enterprise_E50`, `Enterprise_E100`, `Enterprise_E200`, `Enterprise_E400`, `EnterpriseFlash_F300`, `EnterpriseFlash_F700` or `EnterpriseFlash_F1500`. The second segment defines the `capacity` of the `sku_name`, possible values for `Enteprise` SKUs are (`2`, `4`, `6`, ...). Possible values for `EnterpriseFlash` SKUs are (`3`, `9`, `15`, ...). Changing this forces a new Redis Enterprise Cluster to be created.

* `minimum_tls_version` - (Optional) The minimum TLS version. Possible values are `1.0`, `1.1` and `1.2`. Defaults to `1.2`. Changing this forces a new Redis Enterprise Cluster to be created.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Redis Enterprise Cluster should be located. Changing this forces a new Redis Enterprise Cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Redis Enterprise Cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Redis Enterprise Cluster.

* `hostname` - DNS name of the cluster endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Redis Enterprise Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Enterprise Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Redis Enterprise Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Redis Enterprise Cluster.

## Import

Redis Enterprise Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_enterprise_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1
```
