---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_cluster"
description: |-
  Manages Kusto (also known as Azure Data Explorer) Cluster
---

# azurerm_kusto_cluster

Manages a Kusto (also known as Azure Data Explorer) Cluster

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "my-kusto-cluster-rg"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "kustocluster"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto Cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Cluster should exist. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `double_encryption_enabled` - (Optional) Is the cluster's double encryption enabled? Defaults to `false`. Changing this forces a new resource to be created.

* `identity` - (Optional) An identity block.

* `enable_disk_encryption` - (Optional) Specifies if the cluster's disks are encrypted.

* `enable_streaming_ingest` - (Optional) Specifies if the streaming ingest is enabled.

* `enable_purge` - (Optional) Specifies if the purge operations are enabled.

* `virtual_network_configuration`- (Optional) A `virtual_network_configuration` block as defined below. Changing this forces a new resource to be created.

* `language_extensions` - (Optional) An list of `language_extensions` to enable. Valid values are: `PYTHON` and `R`.

* `optimized_auto_scale` - (Optional) An `optimized_auto_scale` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `trusted_external_tenants` - (Optional) Specifies a list of tenant IDs that are trusted by the cluster.

* `zones` - (Optional) A list of Availability Zones in which the cluster instances should be created in. Changing this forces a new resource to be created.

* `engine` - (Optional). The engine type that should be used. Possible values are `V2` and `V3`. Defaults to `V2`.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU. Valid values are: `Dev(No SLA)_Standard_D11_v2`, `Dev(No SLA)_Standard_E2a_v4`, `Standard_D11_v2`, `Standard_D12_v2`, `Standard_D13_v2`, `Standard_D14_v2`, `Standard_DS13_v2+1TB_PS`, `Standard_DS13_v2+2TB_PS`, `Standard_DS14_v2+3TB_PS`, `Standard_DS14_v2+4TB_PS`, `Standard_E16as_v4+3TB_PS`, `Standard_E16as_v4+4TB_PS`, `Standard_E16a_v4`, `Standard_E2a_v4`, `Standard_E4a_v4`, `Standard_E64i_v3`, `Standard_E8as_v4+1TB_PS`, `Standard_E8as_v4+2TB_PS`, `Standard_E8a_v4`, `Standard_L16s`, `Standard_L4s` and `Standard_L8s`.

* `capacity` - (Optional) Specifies the node count for the cluster. Boundaries depend on the sku name.

~> **NOTE:** If no `optimized_auto_scale` block is defined, then the capacity is required.
~> **NOTE:** If an `optimized_auto_scale` block is defined and no capacity is set, then the capacity is initially set to the value of `minimum_instances`.

---

A `virtual_network_configuration` block supports the following:

* `subnet_id` - (Required) The subnet resource id.

* `engine_public_ip_id` - (Required) Engine service's public IP address resource id.

* `data_management_public_ip_id` - (Required) Data management's service public IP address resource id.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that is configured on this Kusto Cluster. Possible values are: `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `optimized_auto_scale` block supports the following:

* `minimum_instances` - (Required) The minimum number of allowed instances. Must between `0` and `1000`.

* `maximum_instances` - (Required) The maximum number of allowed instances. Must between `0` and `1000`.

## Attributes Reference

The following attributes are exported:

* `id` - The Kusto Cluster ID.

* `uri` - The FQDN of the Azure Kusto Cluster.

* `data_ingestion_uri` - The Kusto Cluster URI to be used for data ingestion.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this System Assigned Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this System Assigned Managed Service Identity.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Kusto Cluster.
* `update` - (Defaults to 60 minutes) Used when updating the Kusto Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Cluster.
* `delete` - (Defaults to 60 minutes) Used when deleting the Kusto Cluster.

## Import

Kusto Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1
```
