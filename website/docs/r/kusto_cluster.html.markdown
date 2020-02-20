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
  location = "East US"
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

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU. Valid values are: `Dev(No SLA)_Standard_D11_v2`, `Standard_D11_v2`, `Standard_D12_v2`, `Standard_D13_v2`, `Standard_D14_v2`, `Standard_DS13_v2+1TB_PS`, `Standard_DS13_v2+2TB_PS`, `Standard_DS14_v2+3TB_PS`, `Standard_DS14_v2+4TB_PS`, `Standard_L16s`, `Standard_L4s` and `Standard_L8s`

* `capacity` - (Required) Specifies the node count for the cluster. Boundaries depend on the sku name.


## Attributes Reference

The following attributes are exported:

* `id` - The Kusto Cluster ID.

* `uri` - The FQDN of the Azure Kusto Cluster.

* `data_ingestion_uri` - The Kusto Cluster URI to be used for data ingestion.

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
