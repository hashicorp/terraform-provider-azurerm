---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_cluster"
description: |-
  Get information about an existing Kusto (also known as Azure Data Explorer) Cluster
---

# Data Source: azurerm_kusto_cluster

Use this data source to access information about an existing Kusto (also known as Azure Data Explorer) Cluster

## Example Usage

```hcl
data "azurerm_kusto_cluster" "example" {
  name                = "kustocluster"
  resource_group_name = "test_resource_group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Kusto Cluster.

* `resource_group_name` - The name of the Resource Group where the Kusto Cluster exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Kusto Cluster ID.

* `uri` - The FQDN of the Azure Kusto Cluster.

* `data_ingestion_uri` - The Kusto Cluster URI to be used for data ingestion.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Kusto Cluster.

* `identity_ids` - A list of User Assigned Managed Identity IDs to be assigned to this Kusto Cluster.

* `principal_id` - The Principal ID associated with this System Assigned Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this System Assigned Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Cluster.
